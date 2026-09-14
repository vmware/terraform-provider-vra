// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client"
	"github.com/vmware/vra-sdk-go/pkg/client/deployment_actions"
	"github.com/vmware/vra-sdk-go/pkg/client/deployments"
	"github.com/vmware/vra-sdk-go/pkg/client/requests"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

// resourceActionReason is the reason recorded against day 2 action requests
// submitted by the vra_deployment_resource_action resource.
const resourceActionReason = "Submitted by the vRA provider for Terraform."

type deploymentResourceActionInvocation struct {
	deploymentUUID strfmt.UUID
	resourceUUID   strfmt.UUID
	actionID       string
	inputs         interface{}
	timeoutKey     string
}

// deploymentActionLocks serializes the day 2 actions that Terraform runs against a single
// deployment. vRA accepts one request at a time per deployment and rejects the others, so
// without this lock a configuration that changes the inputs of two actions of the same
// deployment in one apply fails depending on the order in which Terraform schedules them.
var deploymentActionLocks sync.Map

// lockDeployment waits until no other action of the given deployment is running and returns
// the function that releases the deployment. Waiting stops when Terraform cancels the context.
func lockDeployment(ctx context.Context, deploymentUUID strfmt.UUID) (func(), error) {
	newLock := make(chan struct{}, 1)
	newLock <- struct{}{}
	value, _ := deploymentActionLocks.LoadOrStore(deploymentUUID.String(), newLock)
	lock := value.(chan struct{})

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-lock:
		if err := ctx.Err(); err != nil {
			lock <- struct{}{}
			return nil, err
		}
		return func() { lock <- struct{}{} }, nil
	}
}

func deploymentResourceActionPendingStatuses() []string {
	return []string{
		models.RequestStatusCREATED,
		models.RequestStatusPENDING,
		models.RequestStatusINITIALIZATION,
		models.RequestStatusCHECKINGAPPROVAL,
		models.RequestStatusAPPROVALPENDING,
		models.RequestStatusUSERINTERACTIONPENDING,
		models.RequestStatusINPROGRESS,
		models.RequestStatusCOMPLETION,
	}
}

func isDeploymentResourceActionPending(status string) bool {
	for _, pending := range deploymentResourceActionPendingStatuses() {
		if status == pending {
			return true
		}
	}
	return false
}

func resourceDeploymentResourceAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeploymentResourceActionCreate,
		ReadContext:   resourceDeploymentResourceActionRead,
		UpdateContext: resourceDeploymentResourceActionUpdate,
		DeleteContext: resourceDeploymentResourceActionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"action_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the day 2 action to run on the resource, for example 'Cloud.vSphere.Machine.Snapshot.Create'. Use the vra_deployment_resource_actions data source to discover the available action ids.",
			},
			"deployment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the deployment that owns the resource.",
			},
			"destroy_action_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The id of an optional day 2 action to run when this resource is destroyed. When omitted, destroying this resource only removes it from the Terraform state and does not call vRA.",
			},
			"destroy_inputs": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The inputs for the action identified by destroy_action_id.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"inputs": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The inputs for the action. The values are converted to the types declared by the action schema. Values for the array and object types must be supplied as JSON encoded strings.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"last_request_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the most recent action request submitted by this resource.",
			},
			"last_request_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the most recent action request submitted by this resource.",
			},
			"resource_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"resource_id", "resource_name"},
				Description:  "The id of the deployment resource to run the action on. Conflicts with resource_name.",
			},
			"resource_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"resource_id", "resource_name"},
				Description:  "The name of the deployment resource to run the action on. The name must be unique within the deployment. Conflicts with resource_id.",
			},
			"run_on_create": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to run the action when this resource is created. Set to false when the deployment request has already applied the desired values and the action should only run on subsequent changes.",
			},
			"triggers": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "An arbitrary map of values that, when changed, causes the action to run again.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"wait_for_completion": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to wait for the action request to reach a terminal state before continuing.",
			},
		},
	}
}

func resourceDeploymentResourceActionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	deploymentUUID := strfmt.UUID(d.Get("deployment_id").(string))
	resourceUUID, err := resolveDeploymentResourceID(apiClient, deploymentUUID, d)
	if err != nil {
		return diag.FromErr(err)
	}

	actionID := d.Get("action_id").(string)
	d.SetId(deploymentResourceActionID(deploymentUUID, resourceUUID, actionID))
	if err := d.Set("resource_id", resourceUUID.String()); err != nil {
		return diag.FromErr(err)
	}

	if !d.Get("run_on_create").(bool) {
		log.Printf("[DEBUG] run_on_create is false, recording the desired inputs for action %s without calling vRA", actionID)
		return resourceDeploymentResourceActionRead(ctx, d, m)
	}

	invocation := deploymentResourceActionInvocation{
		deploymentUUID: deploymentUUID,
		resourceUUID:   resourceUUID,
		actionID:       actionID,
		inputs:         d.Get("inputs"),
		timeoutKey:     schema.TimeoutCreate,
	}
	if err := submitDeploymentResourceAction(ctx, d, apiClient, invocation); err != nil {
		return diag.FromErr(err)
	}

	return resourceDeploymentResourceActionRead(ctx, d, m)
}

func resourceDeploymentResourceActionRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	deploymentUUID, resourceUUID, actionID, err := parseDeploymentResourceActionID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// The resource records that an action was submitted with a given set of inputs. vRA does not
	// report the inputs of a past action as resource state, so inputs and triggers are not
	// refreshed here. The resource is removed from the state only when the deployment or the
	// deployment resource no longer exists.
	resourceResp, err := apiClient.Deployments.GetResourceByIDUsingGET4(
		deployments.NewGetResourceByIDUsingGET4Params().
			WithAPIVersion(withString(DeploymentsAPIVersion)).
			WithDeploymentID(deploymentUUID).
			WithResourceID(resourceUUID))
	if err != nil {
		switch err.(type) {
		case *deployments.GetResourceByIDUsingGET4NotFound:
			log.Printf("[DEBUG] deployment resource %s no longer exists, removing it from the state", resourceUUID)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	if resourceResp == nil || resourceResp.GetPayload() == nil {
		return diag.Errorf("vRA returned an empty response for deployment resource %s", resourceUUID)
	}

	if err := d.Set("action_id", actionID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("deployment_id", deploymentUUID.String()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("resource_id", resourceUUID.String()); err != nil {
		return diag.FromErr(err)
	}

	if requestID, ok := d.GetOk("last_request_id"); ok {
		requestResp, err := apiClient.Requests.GetRequestUsingGET2(
			requests.NewGetRequestUsingGET2Params().
				WithAPIVersion(withString(DeploymentsAPIVersion)).
				WithRequestID(strfmt.UUID(requestID.(string))))
		if err != nil {
			log.Printf("[WARN] unable to read the status of the action request %s: %v", requestID, err)
		} else if requestResp == nil || requestResp.GetPayload() == nil {
			log.Printf("[WARN] vRA returned an empty response for action request %s", requestID)
		} else if err := d.Set("last_request_status", requestResp.GetPayload().Status); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceDeploymentResourceActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	deploymentUUID, resourceUUID, _, err := parseDeploymentResourceActionID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// destroy_action_id and destroy_inputs are only read when the resource is destroyed, so a
	// change to either of them must not run the action again.
	if !d.HasChanges("action_id", "inputs", "triggers") {
		return resourceDeploymentResourceActionRead(ctx, d, m)
	}

	actionID := d.Get("action_id").(string)
	invocation := deploymentResourceActionInvocation{
		deploymentUUID: deploymentUUID,
		resourceUUID:   resourceUUID,
		actionID:       actionID,
		inputs:         d.Get("inputs"),
		timeoutKey:     schema.TimeoutUpdate,
	}
	if err := submitDeploymentResourceAction(ctx, d, apiClient, invocation); err != nil {
		// Terraform writes the state back even when the update fails, which would record the
		// new inputs although the action never ran. Restoring the values that were read from
		// vRA keeps the change in the next plan so that the failed action is retried.
		return append(revertDeploymentResourceActionInputs(d), diag.FromErr(err)...)
	}

	d.SetId(deploymentResourceActionID(deploymentUUID, resourceUUID, actionID))

	return resourceDeploymentResourceActionRead(ctx, d, m)
}

// revertDeploymentResourceActionInputs restores the attributes that describe the last action
// that was run to the values they had before the failed update.
func revertDeploymentResourceActionInputs(d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	for _, name := range []string{"action_id", "inputs", "triggers"} {
		if !d.HasChange(name) {
			continue
		}
		previous, _ := d.GetChange(name)
		if err := d.Set(name, previous); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}
	return diags
}

func resourceDeploymentResourceActionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	destroyActionID, ok := d.GetOk("destroy_action_id")
	if !ok {
		log.Printf("[DEBUG] destroy_action_id is not set, removing %s from the state without calling vRA", d.Id())
		d.SetId("")
		return nil
	}

	deploymentUUID, resourceUUID, _, err := parseDeploymentResourceActionID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	invocation := deploymentResourceActionInvocation{
		deploymentUUID: deploymentUUID,
		resourceUUID:   resourceUUID,
		actionID:       destroyActionID.(string),
		inputs:         d.Get("destroy_inputs"),
		timeoutKey:     schema.TimeoutDelete,
	}
	if err := submitDeploymentResourceAction(ctx, d, apiClient, invocation); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

// submitDeploymentResourceAction converts the given inputs to the types declared by the action
// schema, submits the action, and optionally waits for the request to reach a terminal state.
func submitDeploymentResourceAction(ctx context.Context, d *schema.ResourceData, apiClient *client.API, invocation deploymentResourceActionInvocation) error {
	// vRA runs one request at a time per deployment. Synchronous actions hold the deployment
	// until the request reaches a terminal state; asynchronous actions release it after submission.
	unlock, err := lockDeployment(ctx, invocation.deploymentUUID)
	if err != nil {
		return fmt.Errorf("unable to wait to run an action on deployment %s: %w", invocation.deploymentUUID, err)
	}
	defer unlock()

	inputs, err := buildDeploymentResourceActionInputs(apiClient, invocation.deploymentUUID, invocation.resourceUUID, invocation.actionID, invocation.inputs)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Running the action %s on the resource %s of the deployment %s", invocation.actionID, invocation.resourceUUID, invocation.deploymentUUID)
	resp, err := apiClient.DeploymentActions.SubmitResourceActionRequestUsingPOST4(
		deployment_actions.NewSubmitResourceActionRequestUsingPOST4Params().
			WithContext(ctx).
			WithAPIVersion(withString(DeploymentsAPIVersion)).
			WithDeploymentID(invocation.deploymentUUID).
			WithResourceID(invocation.resourceUUID).
			WithActionRequest(&models.ResourceActionRequest{
				ActionID: invocation.actionID,
				Reason:   resourceActionReason,
				Inputs:   inputs,
			}))
	if err != nil {
		return fmt.Errorf("unable to run the action %s on the resource %s: %w", invocation.actionID, invocation.resourceUUID, err)
	}
	if resp == nil || resp.GetPayload() == nil || resp.GetPayload().ID == "" {
		return fmt.Errorf("vRA returned an empty response after submitting action %s on resource %s", invocation.actionID, invocation.resourceUUID)
	}

	requestID := resp.GetPayload().ID
	if err := d.Set("last_request_id", requestID.String()); err != nil {
		return err
	}
	if !d.Get("wait_for_completion").(bool) {
		return d.Set("last_request_status", resp.GetPayload().Status)
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    deploymentResourceActionPendingStatuses(),
		Refresh:    deploymentResourceActionStatusRefreshFunc(ctx, apiClient, requestID),
		Target:     []string{models.RequestStatusSUCCESSFUL},
		Timeout:    d.Timeout(invocation.timeoutKey),
		MinTimeout: 5 * time.Second,
	}
	if _, err := stateChangeFunc.WaitForStateContext(ctx); err != nil {
		return err
	}

	return d.Set("last_request_status", models.RequestStatusSUCCESSFUL)
}

// buildDeploymentResourceActionInputs converts the given inputs to the types declared by the
// schema of the action. Inputs that the action schema does not declare are rejected so that a
// typo in the configuration fails the apply instead of being silently dropped by vRA.
func buildDeploymentResourceActionInputs(apiClient *client.API, deploymentUUID, resourceUUID strfmt.UUID, actionID string, rawInputs interface{}) (map[string]interface{}, error) {
	inputs, _ := rawInputs.(map[string]interface{})
	if len(inputs) == 0 {
		return make(map[string]interface{}), nil
	}

	inputTypesMap, err := getResourceActionInputTypesMap(apiClient, deploymentUUID, resourceUUID, actionID)
	if err != nil {
		return nil, err
	}

	unknown := make([]string, 0)
	for name := range inputs {
		if _, ok := inputTypesMap[name]; !ok {
			unknown = append(unknown, name)
		}
	}
	if len(unknown) > 0 {
		known := make([]string, 0, len(inputTypesMap))
		for name := range inputTypesMap {
			known = append(known, name)
		}
		sort.Strings(unknown)
		sort.Strings(known)
		return nil, fmt.Errorf("the action %s does not declare the input(s) %s. The action declares the input(s) %s",
			actionID, strings.Join(unknown, ", "), strings.Join(known, ", "))
	}

	typedInputs, err := getInputsByType(inputs, inputTypesMap)
	if err != nil {
		return nil, fmt.Errorf("unable to create the inputs for the action %s: %w", actionID, err)
	}

	return typedInputs, nil
}

// getResourceActionInputTypesMap returns the input names and types declared by the schema of a
// day 2 action of a deployment resource.
func getResourceActionInputTypesMap(apiClient *client.API, deploymentUUID, resourceUUID strfmt.UUID, actionID string) (map[string]string, error) {
	action, err := apiClient.DeploymentActions.GetResourceActionUsingGET4(
		deployment_actions.NewGetResourceActionUsingGET4Params().
			WithAPIVersion(withString(DeploymentsAPIVersion)).
			WithDeploymentID(deploymentUUID).
			WithResourceID(resourceUUID).
			WithActionID(actionID))
	if err != nil {
		return nil, fmt.Errorf("unable to read the action %s of the resource %s: %w", actionID, resourceUUID, err)
	}
	if action == nil || action.GetPayload() == nil {
		return nil, fmt.Errorf("vRA returned an empty response for action %s of resource %s", actionID, resourceUUID)
	}

	if !action.GetPayload().Valid {
		return nil, fmt.Errorf("the action %s is not valid for the resource %s in its current state", actionID, resourceUUID)
	}

	return getResourceActionInputTypesMapFromSchema(flattenResourceActionSchemaProperties(action.GetPayload().Schema))
}

func getResourceActionInputTypesMapFromSchema(actionSchema map[string]interface{}) (map[string]string, error) {
	inputTypesMap := make(map[string]string, len(actionSchema))
	for name, rawProperty := range actionSchema {
		property, ok := rawProperty.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("input %q has an invalid schema: expected an object", name)
		}

		rawType, ok := property["type"]
		if !ok {
			return nil, fmt.Errorf("input %q has an invalid schema: missing type", name)
		}

		inputType, ok := rawType.(string)
		if !ok || inputType == "" {
			return nil, fmt.Errorf("input %q has an invalid schema: type must be a non-empty string", name)
		}

		inputTypesMap[name] = inputType
	}
	return inputTypesMap, nil
}

// flattenResourceActionSchemaProperties returns the properties of an action schema, or an empty
// map when the schema is absent or does not have the expected shape.
func flattenResourceActionSchemaProperties(actionSchema interface{}) map[string]interface{} {
	schemaMap, ok := actionSchema.(map[string]interface{})
	if !ok {
		return make(map[string]interface{})
	}

	properties, ok := schemaMap["properties"].(map[string]interface{})
	if !ok {
		return make(map[string]interface{})
	}

	return properties
}

// resolveDeploymentResourceID returns the id of the deployment resource identified by either the
// resource_id or the resource_name attribute.
func resolveDeploymentResourceID(apiClient *client.API, deploymentUUID strfmt.UUID, d *schema.ResourceData) (strfmt.UUID, error) {
	if v, ok := d.GetOk("resource_id"); ok {
		return strfmt.UUID(v.(string)), nil
	}

	name := d.Get("resource_name").(string)
	resp, err := apiClient.Deployments.GetDeploymentResourcesUsingGET2(
		deployments.NewGetDeploymentResourcesUsingGET2Params().
			WithAPIVersion(withString(DeploymentsAPIVersion)).
			WithDeploymentID(deploymentUUID).
			WithDollarTop(withInt32(DefaultDollarTop)))
	if err != nil {
		return "", fmt.Errorf("unable to list the resources of the deployment %s: %w", deploymentUUID, err)
	}
	if resp == nil || resp.GetPayload() == nil {
		return "", fmt.Errorf("vRA returned an empty response when listing resources of deployment %s", deploymentUUID)
	}

	matches := make([]strfmt.UUID, 0, 1)
	names := make([]string, 0, len(resp.GetPayload().Content))
	for _, resource := range resp.GetPayload().Content {
		if resource == nil || resource.Name == nil {
			continue
		}
		names = append(names, *resource.Name)
		if *resource.Name == name {
			matches = append(matches, resource.ID)
		}
	}

	switch len(matches) {
	case 1:
		return matches[0], nil
	case 0:
		return "", fmt.Errorf("no resource named %s was found in the deployment %s. The deployment has the resource(s) %s",
			name, deploymentUUID, strings.Join(names, ", "))
	default:
		return "", fmt.Errorf("%d resources named %s were found in the deployment %s. Use resource_id to select one of them",
			len(matches), name, deploymentUUID)
	}
}

func deploymentResourceActionID(deploymentUUID, resourceUUID strfmt.UUID, actionID string) string {
	return fmt.Sprintf("%s/%s/%s", deploymentUUID, resourceUUID, actionID)
}

func parseDeploymentResourceActionID(id string) (strfmt.UUID, strfmt.UUID, string, error) {
	parts := strings.SplitN(id, "/", 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", fmt.Errorf("invalid id %q, expected the format 'deployment_id/resource_id/action_id'", id)
	}

	return strfmt.UUID(parts[0]), strfmt.UUID(parts[1]), parts[2], nil
}

func deploymentResourceActionStatusRefreshFunc(ctx context.Context, apiClient *client.API, requestID strfmt.UUID) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ret, err := apiClient.Requests.GetRequestUsingGET2(
			requests.NewGetRequestUsingGET2Params().
				WithContext(ctx).
				WithAPIVersion(withString(DeploymentsAPIVersion)).
				WithRequestID(requestID))
		if err != nil {
			return "", models.RequestStatusFAILED, err
		}
		if ret == nil || ret.GetPayload() == nil {
			return "", models.RequestStatusFAILED, fmt.Errorf("vRA returned an empty response for action request %s", requestID)
		}

		request := ret.GetPayload()
		if isDeploymentResourceActionPending(request.Status) {
			return request, request.Status, nil
		}

		switch request.Status {
		case models.RequestStatusSUCCESSFUL:
			return request, request.Status, nil
		case models.RequestStatusFAILED, models.RequestStatusAPPROVALREJECTED, models.RequestStatusABORTED:
			return request, request.Status, fmt.Errorf("the action request %s ended with the status %s: %s", requestID, request.Status, request.Details)
		default:
			return request, request.Status, errors.New("deploymentResourceActionStatusRefreshFunc: unknown status " + request.Status)
		}
	}
}
