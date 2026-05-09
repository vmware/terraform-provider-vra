// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client"
	"github.com/vmware/vra-sdk-go/pkg/client/blueprint"
	"github.com/vmware/vra-sdk-go/pkg/client/blueprint_requests"
	"github.com/vmware/vra-sdk-go/pkg/client/catalog_items"
	"github.com/vmware/vra-sdk-go/pkg/client/deployment_actions"
	"github.com/vmware/vra-sdk-go/pkg/client/deployments"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

const (
	ChangeOwnerDeploymentActionName = "ChangeOwner"
	ChangeLeaseDeploymentActionName = "ChangeLease"
	EditTagsDeploymentActionName    = "EditTags"
	PowerOffDeploymentActionName    = "PowerOff"
	PowerOnDeploymentActionName     = "PowerOn"
	UpdateDeploymentActionName      = "update"

	// resourceActionModifyKeyword matches Day2 resource-action IDs that represent write/modify intent.
	// Combined with UpdateDeploymentActionName, these filter out read-only display actions.
	resourceActionModifyKeyword = "modify"

	// deploymentUpdateReason is the reason string attached to all Day2 Update actions.
	deploymentUpdateReason = "Updated deployment inputs from vRA provider for Terraform."
)

func resourceDeployment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeploymentCreate,
		ReadContext:   resourceDeploymentRead,
		UpdateContext: resourceDeploymentUpdate,
		DeleteContext: resourceDeploymentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"blueprint_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"blueprint_content", "catalog_item_id"},
				Description:   "The id of the cloud template to be used to request the deployment.",
			},
			"blueprint_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The version of the cloud template to be used to request the deployment.",
			},
			"blueprint_content": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"blueprint_id", "catalog_item_id"},
				Description:   "The content of the the cloud template to be used to request the deployment.",
			},
			"catalog_item_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"blueprint_id", "blueprint_content"},
				ForceNew:      true,
				Description:   "The id of the catalog item to be used to request the deployment.",
			},
			"catalog_item_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The version of the catalog item to be used to request the deployment.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was created. The date is in ISO 6801 and UTC.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user the entity was created by.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A human-friendly description.",
			},
			"expand_last_request": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Deprecated. True by default even if not provided.",
			},
			"expand_project": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Flag to indicate whether to expand project information.",
			},
			"expand_resources": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Deprecated. True by default even if not provided.",
			},
			"expense": expenseSchema(),
			"inputs": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Inputs provided by the user. For inputs including those with default values, refer to inputs_including_defaults.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"inputs_including_defaults": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "All the inputs applied during last create/update operation, including those with default values.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"last_request": deploymentRequestSchema(),
			"last_updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is in ISO 6801 and UTC.",
			},
			"last_updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user that last updated the deployment.",
			},
			"lease_expire_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the deployment lease expire. The date is in ISO 6801 and UTC.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the deployment.",
			},
			"org_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Id of the organization this deployment belongs to.",
			},
			"owner": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The user this deployment belongs to.",
			},
			"project": resourceReferenceSchema(),
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the project this deployment belongs to.",
			},
			"reason": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Reason for requesting/updating a blueprint.",
			},
			"resources": resourcesSchema(),
			// TODO: Add plan / simulate feature
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the deployment with respect to its life cycle operations.",
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceDeploymentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to create vra_deployment resource")
	apiClient := m.(*Client).apiClient

	blueprintID, catalogItemID, blueprintContent := "", "", ""
	if v, ok := d.GetOk("blueprint_id"); ok {
		blueprintID = v.(string)
	}

	if v, ok := d.GetOk("catalog_item_id"); ok {
		catalogItemID = v.(string)
	}

	if v, ok := d.GetOk("blueprint_content"); ok {
		blueprintContent = v.(string)
	}

	if blueprintID != "" && catalogItemID != "" {
		return diag.FromErr(errors.New("only one of (blueprint_id, catalog_item_id) required"))
	}

	deploymentName := d.Get("name").(string)
	projectID := d.Get("project_id").(string)

	getResp, err := apiClient.Deployments.CheckDeploymentNameUsingGET2(deployments.NewCheckDeploymentNameUsingGET2Params().WithName(deploymentName))
	log.Printf("getResp: %v, err: %v", getResp, err)

	if err != nil {
		switch err.(type) {
		case *deployments.CheckDeploymentNameUsingGET2NotFound:
			log.Printf("Deployment '%v' doesn't exist already and hence can be created", deploymentName)
		}
	} else {
		return diag.Errorf("a deployment with name '%v' exists already. Try with a differnet name", deploymentName)
	}

	inputs := make(map[string]interface{})

	// If catalog_item_id is provided, request deployment with the catalog item
	if catalogItemID != "" {
		log.Printf("Requesting vra_deployment '%s' from catalog item", d.Get("name"))
		catalogItemVersion := ""
		if v, ok := d.GetOk("catalog_item_version"); ok {
			catalogItemVersion = v.(string)
		}

		catalogItemRequest := models.CatalogItemRequest{
			DeploymentName: deploymentName,
			ProjectID:      projectID,
			Version:        catalogItemVersion,
		}

		if v, ok := d.GetOk("inputs"); ok {
			inputs, err = getCatalogItemInputsByType(apiClient, catalogItemID, catalogItemVersion, v)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		catalogItemRequest.Inputs = inputs

		if v, ok := d.GetOk("description"); ok {
			catalogItemRequest.Reason = v.(string)
		}

		if v, ok := d.GetOk("reason"); ok {
			catalogItemRequest.Reason = v.(string)
		}

		log.Printf("[DEBUG] Create deployment: %#v", catalogItemRequest)
		postOk, err := apiClient.CatalogItems.RequestCatalogItemInstancesUsingPOST1(
			catalog_items.NewRequestCatalogItemInstancesUsingPOST1Params().WithID(strfmt.UUID(catalogItemID)).
				WithAPIVersion(withString(CatalogAPIVersion)).WithRequest(&catalogItemRequest))

		if err != nil {
			return diag.FromErr(err)
		}

		payload := postOk.GetPayload()
		if len(payload) == 0 {
			return diag.Errorf("failed to request vra_deployment '%s' from catalog item", d.Get("name"))
		}
		d.SetId(payload[0].DeploymentID)
		log.Printf("Finished requesting vra_deployment '%s' from catalog item", d.Get("name"))
	} else {
		blueprintVersion := ""
		if v, ok := d.GetOk("blueprint_version"); ok {
			blueprintVersion = v.(string)
		}

		blueprintRequest := models.BlueprintRequest{
			BlueprintVersion: blueprintVersion,
			DeploymentName:   deploymentName,
			ProjectID:        projectID,
		}

		if blueprintID != "" {
			blueprintRequest.BlueprintID = strfmt.UUID(blueprintID)
		} else {
			// Create empty content in the blueprint
			blueprintRequest.Content = " "
		}

		if blueprintContent != "" {
			blueprintRequest.Content = blueprintContent
		}

		if v, ok := d.GetOk("description"); ok {
			blueprintRequest.Description = v.(string)
		}

		if v, ok := d.GetOk("inputs"); ok {
			if blueprintContent != "" && blueprintID == "" {
				inputs = expandInputs(v)
			} else {
				// If the inputs are provided, get the schema from blueprint to convert the provided input values
				// to the type defined in the schema.
				inputs, err = getBlueprintInputsByType(apiClient, blueprintID, blueprintVersion, v)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
		blueprintRequest.Inputs = inputs

		if v, ok := d.GetOk("reason"); ok {
			blueprintRequest.Reason = v.(string)
		}

		bpRequestCreated, bpRequestAccepted, err := apiClient.BlueprintRequests.CreateBlueprintRequestUsingPOST1(
			blueprint_requests.NewCreateBlueprintRequestUsingPOST1Params().WithRequest(&blueprintRequest))

		if err != nil {
			log.Printf("Received error. err=%s, bpRequestCreated=%v, bpRequestAccepted=%v", err, bpRequestCreated, bpRequestAccepted)
			return diag.FromErr(err)
		}

		// blueprint_requests service may return either 201 or 202 depending on whether the request is in terminal state vs or in-progress
		log.Printf("Requested deployment from blueprint. bpRequestCreated=%v, bpRequestAccepted=%v", bpRequestCreated, bpRequestAccepted)
		deploymentID, status, failureMessage := "", "", ""
		var bpRequest *models.BlueprintRequest
		if bpRequestAccepted != nil {
			bpRequest = bpRequestAccepted.GetPayload()
		} else {
			bpRequest = bpRequestCreated.GetPayload()
		}

		if bpRequest != nil {
			deploymentID = bpRequest.DeploymentID
			status = bpRequest.Status
			failureMessage = bpRequest.FailureMessage
		}

		if deploymentID != "" {
			d.SetId(deploymentID)
		} else {
			return diag.Errorf("failed to request for a deployment. status: %v, message: %v", status, failureMessage)
		}

		log.Printf("Finished requesting vra_deployment '%s' from blueprint %s", d.Get("name"), blueprintID)
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.DeploymentStatusCREATEINPROGRESS, models.DeploymentStatusUPDATEINPROGRESS},
		Refresh:    deploymentStatusRefreshFunc(*apiClient, d.Id()),
		Target:     []string{models.DeploymentStatusCREATESUCCESSFUL, models.DeploymentStatusUPDATESUCCESSFUL},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
	}

	deploymentID, err := stateChangeFunc.WaitForStateContext(ctx)
	if err != nil {
		readErrors := resourceDeploymentRead(ctx, d, m)
		if readErrors.HasError() {
			return append(readErrors, diag.Errorf("failed to create deployment: %v", err.Error())...)
		}
		return diag.FromErr(err)
	}

	d.SetId(deploymentID.(string))
	log.Printf("Finished to create vra_deployment resource with name %s", d.Get("name"))

	return resourceDeploymentRead(ctx, d, m)
}

func resourceDeploymentRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Reading the vra_deployment resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	// Getting the input types map
	inputTypesMap := getInputTypesMap(d, apiClient)

	id := d.Id()
	expandProject := d.Get("expand_project").(bool)

	expand := []string{"resources", "lastRequest"}
	if expandProject {
		expand = append(expand, "project")
	}

	resp, err := apiClient.Deployments.GetDeploymentByIDV3UsingGET(
		deployments.NewGetDeploymentByIDV3UsingGETParams().
			WithDeploymentID(strfmt.UUID(id)).
			WithExpand(expand).
			WithAPIVersion(withString(DeploymentsAPIVersion)))
	if err != nil {
		switch err.(type) {
		case *deployments.GetDeploymentByIDV3UsingGETNotFound:
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	deployment := *resp.Payload
	d.Set("blueprint_id", deployment.BlueprintID)
	d.Set("blueprint_version", deployment.BlueprintVersion)
	d.Set("catalog_item_id", deployment.CatalogItemID)
	d.Set("catalog_item_version", deployment.CatalogItemVersion)
	d.Set("created_at", deployment.CreatedAt.String())
	d.Set("created_by", deployment.CreatedBy)
	d.Set("description", deployment.Description)

	if err := d.Set("expense", flattenExpense(deployment.Expense)); err != nil {
		return diag.Errorf("error setting deployment expense - error: %#v", err)
	}

	if err := d.Set("inputs_including_defaults", expandInputsToString(deployment.Inputs)); err != nil {
		return diag.Errorf("error setting deployment inputs_including_defaults - error: %#v", err)
	}

	allInputs := expandInputs(deployment.Inputs)

	// Collect user-managed input keys so overlayLastRequestInputs can update them.
	var userManagedKeys map[string]bool
	if v, ok := d.GetOk("inputs"); ok {
		if um, ok2 := v.(map[string]interface{}); ok2 {
			userManagedKeys = make(map[string]bool, len(um))
			for k := range um {
				userManagedKeys[k] = true
			}
		}
	}

	lastRequestWasResourceLevel := false
	if deployment.LastRequest != nil && deployment.LastRequest.Status == models.RequestStatusSUCCESSFUL {
		if lastReqInputs, ok := deployment.LastRequest.Inputs.(map[string]interface{}); ok {
			overlayLastRequestInputs(allInputs, lastReqInputs, userManagedKeys)
		}
		if len(deployment.LastRequest.ResourceIds) > 0 && deployment.LastRequest.ActionID != "" {
			lastRequestWasResourceLevel = true
			log.Printf("[DEBUG] last request %s was resource-level (actionId=%s, resourceIds=%v); deployment.Inputs may be stale",
				deployment.LastRequest.ID, deployment.LastRequest.ActionID, deployment.LastRequest.ResourceIds)
		}
	}

	if v, ok := d.GetOk("inputs"); ok {
		userInputs := v.(map[string]interface{})
		updatedInputs := updateUserInputs(allInputs, userInputs, inputTypesMap)

		// After a resource-level action, deployment.Inputs is stale because the
		// platform only updates the individual resource, not deployment-level inputs.
		// We preserve Terraform state values when the last request was resource-level
		// because:
		//   1. Resource-level actions are only submitted by this provider (not the UI)
		//   2. External changes via the vRA UI create deployment-level requests,
		//      which would set lastRequestWasResourceLevel=false, allowing drift
		//      to be reported normally on the next plan.
		if lastRequestWasResourceLevel && updatedInputs != nil {
			for k, stateVal := range userInputs {
				if platformVal, exists := updatedInputs[k]; exists && stateVal != nil && fmt.Sprint(platformVal) != fmt.Sprint(stateVal) {
					log.Printf("[DEBUG] preserving Terraform state value for %q after resource-level action (platform=%v is stale, state=%v)", k, platformVal, stateVal)
					updatedInputs[k] = stateVal
				}
			}
		}

		if err := d.Set("inputs", updatedInputs); err != nil {
			return diag.Errorf("error setting deployment inputs - error: %#v", err)
		}
	}

	if err := d.Set("last_request", flattenDeploymentRequest(deployment.LastRequest)); err != nil {
		return diag.Errorf("error setting deployment last_request - error: %#v", err)
	}

	d.Set("last_updated_at", deployment.LastUpdatedAt.String())
	d.Set("last_updated_by", deployment.LastUpdatedBy)
	d.Set("lease_expire_at", deployment.LeaseExpireAt.String())
	d.Set("name", deployment.Name)
	d.Set("org_id", deployment.OrgID)
	d.Set("owner", deployment.OwnedBy)

	if err := d.Set("project", flattenResourceReference(deployment.Project)); err != nil {
		return diag.Errorf("error setting project in deployment - error: %#v", err)
	}

	d.Set("project_id", deployment.ProjectID)

	getResourcesResp, err := apiClient.Deployments.GetDeploymentResourcesUsingGET2(
		deployments.NewGetDeploymentResourcesUsingGET2Params().
			WithDeploymentID(strfmt.UUID(id)).
			WithExpand([]string{"currentRequest"}).
			WithAPIVersion(withString(DeploymentsAPIVersion)).
			WithDollarTop(withInt32(DefaultDollarTop)))
	if err != nil {
		return diag.Errorf("error retrieving deployment resources - error: %#v", err)
	}

	if err := d.Set("resources", flattenResources(getResourcesResp.GetPayload())); err != nil {
		return diag.Errorf("error setting resources in deployment - error: %#v", err)
	}

	d.Set("status", deployment.Status)

	log.Printf("Finished reading the vra_deployment resource with name '%s'. Current status: '%s'", d.Get("name"), d.Get("status"))
	return nil
}

func resourceDeploymentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to update the vra_deployment resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	_, blueprintContentExists := d.GetOk("blueprint_content")

	if d.HasChange("blueprint_id") || d.HasChange("blueprint_version") || blueprintContentExists {
		err := updateDeploymentWithNewBlueprint(ctx, d, m, apiClient)
		if err.HasError() {
			return err
		}
	} else {
		id := d.Id()
		deploymentUUID := strfmt.UUID(id)
		if d.HasChange("name") || d.HasChange("description") {
			err := updateDeploymentMetadata(d, apiClient, deploymentUUID)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if d.HasChange("inputs") {
			_, updateErr := runDeploymentUpdateAction(ctx, d, apiClient, deploymentUUID)
			if updateErr != nil {
				// Restore prior inputs so a failed update doesn't persist planned values.
				if oldInputs, _ := d.GetChange("inputs"); oldInputs != nil {
					_ = d.Set("inputs", oldInputs)
				}
				return diag.FromErr(updateErr)
			}
		}

		stateChangeFunc := retry.StateChangeConf{
			Delay:      5 * time.Second,
			Pending:    []string{models.DeploymentStatusCREATEINPROGRESS, models.DeploymentStatusUPDATEINPROGRESS},
			Refresh:    deploymentStatusRefreshFunc(*apiClient, d.Id()),
			Target:     []string{models.DeploymentStatusCREATESUCCESSFUL, models.DeploymentStatusUPDATESUCCESSFUL},
			Timeout:    d.Timeout(schema.TimeoutCreate),
			MinTimeout: 5 * time.Second,
		}

		if _, err := stateChangeFunc.WaitForStateContext(ctx); err != nil {
			readErrors := resourceDeploymentRead(ctx, d, m)
			if readErrors.HasError() {
				return append(readErrors, diag.Errorf("failed to create deployment: %v", err.Error())...)
			}
			return diag.FromErr(err)
		}
	}

	if d.HasChange("owner") {
		deploymentUUID := strfmt.UUID(d.Id())
		err := runChangeOwnerDeploymentAction(ctx, d, apiClient, deploymentUUID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("Finished updating the vra_deployment resource with name %s", d.Get("name"))

	return resourceDeploymentRead(ctx, d, m)
}

func resourceDeploymentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to delete the vra_deployment resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	_, err := apiClient.Deployments.DeleteDeploymentUsingDELETE2(deployments.NewDeleteDeploymentUsingDELETE2Params().WithDeploymentID(strfmt.UUID(id)))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("Requested for deleting the vra_deployment resource with name %s", d.Get("name"))

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{reflect.TypeOf((*deployments.GetDeploymentByIDV3UsingGETOK)(nil)).String()},
		Refresh:    deploymentDeleteStatusRefreshFunc(*apiClient, d.Id()),
		Target:     []string{reflect.TypeOf((*deployments.GetDeploymentByIDV3UsingGETNotFound)(nil)).String()},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 5 * time.Second,
	}

	if _, err = stateChangeFunc.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("Finished deleting the vra_deployment resource with name %s", d.Get("name"))
	return nil
}

// Gets the inputs and their types as map[string]string
func getInputTypesMap(d *schema.ResourceData, apiClient *client.API) map[string]string {
	inputTypesMap := make(map[string]string)

	if _, ok := d.GetOk("inputs"); !ok {
		return inputTypesMap
	}

	// Get blueprint_id and catalog_item_id
	blueprintID, catalogItemID := "", ""
	if v, ok := d.GetOk("blueprint_id"); ok {
		blueprintID = v.(string)
	}

	if v, ok := d.GetOk("catalog_item_id"); ok {
		catalogItemID = v.(string)
	}

	if catalogItemID != "" {
		// Get the catalog item inputs and their types
		catalogItemVersion := ""
		if v, ok := d.GetOk("catalog_item_version"); ok {
			catalogItemVersion = v.(string)
		}

		inputTypesMap, _ = getCatalogItemInputTypesMap(apiClient, catalogItemID, catalogItemVersion)
	} else if blueprintID != "" {
		// Get the blueprint inputs and their types
		blueprintVersion := ""
		if v, ok := d.GetOk("blueprint_version"); ok {
			blueprintVersion = v.(string)
		}

		inputTypesMap, _ = getBlueprintInputTypesMap(apiClient, blueprintID, blueprintVersion)
	}

	log.Printf("InputTypesMap: %v", inputTypesMap)
	return inputTypesMap
}

func getCatalogItemInputsByType(apiClient *client.API, catalogItemID string, catalogItemVersion string, inputValues interface{}) (map[string]interface{}, error) {
	inputTypesMap, err := getCatalogItemInputTypesMap(apiClient, catalogItemID, catalogItemVersion)
	if err != nil {
		return nil, err
	}

	log.Printf("InputTypesMap: %v", inputTypesMap)
	inputs, err := getInputsByType(inputValues.(map[string]interface{}), inputTypesMap)
	if err != nil {
		return nil, err
	}

	return inputs, nil
}

func getCatalogItemInputTypesMap(apiClient *client.API, catalogItemID string, catalogItemVersion string) (map[string]string, error) {
	log.Printf("Getting Catalog Item Schema for catalog_item_id: [%v], catalog_item_version: [%v]", catalogItemID, catalogItemVersion)
	inputsSchemaMap, err := getCatalogItemSchema(apiClient, catalogItemID, catalogItemVersion)
	if err != nil {
		return nil, err
	}

	log.Printf("Inputs Schema: %v", inputsSchemaMap)
	inputTypesMap, err := getInputTypesMapFromSchema(inputsSchemaMap)
	if err != nil {
		return nil, err
	}
	return inputTypesMap, nil
}

func getBlueprintInputsByType(apiClient *client.API, blueprintID string, blueprintVersion string, inputValues interface{}) (map[string]interface{}, error) {
	inputTypesMap, err := getBlueprintInputTypesMap(apiClient, blueprintID, blueprintVersion)
	if err != nil {
		return nil, err
	}

	log.Printf("InputTypesMap: %v", inputTypesMap)
	inputs, err := getInputsByType(inputValues.(map[string]interface{}), inputTypesMap)
	if err != nil {
		return nil, err
	}

	return inputs, nil
}

func getBlueprintInputTypesMap(apiClient *client.API, blueprintID string, blueprintVersion string) (map[string]string, error) {
	log.Printf("Getting Blueprint Schema for blueprint_id: [%v], blueprint_version: [%v]", blueprintID, blueprintVersion)
	inputsSchemaMap, err := getBlueprintSchema(apiClient, blueprintID, blueprintVersion)
	if err != nil {
		return nil, err
	}

	log.Printf("Inputs Schema: %v", inputsSchemaMap)
	inputTypesMap, err := getInputTypesMapFromBlueprintInputsSchema(inputsSchemaMap)
	if err != nil {
		return nil, err
	}
	return inputTypesMap, nil
}

func getCatalogItemSchema(apiClient *client.API, catalogItemID string, catalogItemVersion string) (map[string]interface{}, error) {
	// Getting the catalog item schema
	log.Printf("Getting the schema for catalog item: %v version: %v", catalogItemID, catalogItemVersion)
	var catalogItemSchema interface{}
	if catalogItemVersion == "" {
		getItemResp, err := apiClient.CatalogItems.GetCatalogItemUsingGET5(catalog_items.NewGetCatalogItemUsingGET5Params().WithID(strfmt.UUID(catalogItemID)))
		if err != nil {
			return nil, err
		}
		catalogItemSchema = getItemResp.GetPayload().Schema
	} else {
		getVersionResp, err := apiClient.CatalogItems.GetVersionByIDUsingGET2(catalog_items.NewGetVersionByIDUsingGET2Params().WithID(strfmt.UUID(catalogItemID)).WithVersionID(catalogItemVersion))
		if err != nil {
			return nil, err
		}
		catalogItemSchema = getVersionResp.GetPayload().Schema
	}

	if catalogItemSchema != nil && (catalogItemSchema.(map[string]interface{}))["properties"] != nil {
		inputsSchemaMap := (catalogItemSchema.(map[string]interface{}))["properties"].(map[string]interface{})
		return inputsSchemaMap, nil
	}
	return make(map[string]interface{}), nil
}

func getBlueprintSchema(apiClient *client.API, blueprintID string, blueprintVersion string) (map[string]models.PropertyDefinition, error) {
	// Getting the blueprint inputs schema
	log.Printf("Getting the schema for catalog item: %v version: %v", blueprintID, blueprintVersion)
	var blueprintInputsSchema map[string]models.PropertyDefinition
	if blueprintVersion == "" {
		getItemResp, err := apiClient.Blueprint.GetBlueprintInputsSchemaUsingGET1(blueprint.NewGetBlueprintInputsSchemaUsingGET1Params().WithBlueprintID(blueprintID))
		if err != nil {
			return nil, err
		}
		blueprintInputsSchema = getItemResp.GetPayload().Properties
	} else {
		getVersionResp, err := apiClient.Blueprint.GetBlueprintVersionInputsSchemaUsingGET1(
			blueprint.NewGetBlueprintVersionInputsSchemaUsingGET1Params().WithBlueprintID(blueprintID).
				WithVersion(blueprintVersion))
		if err != nil {
			return nil, err
		}
		blueprintInputsSchema = getVersionResp.GetPayload().Properties
	}
	return blueprintInputsSchema, nil
}

func deploymentStatusRefreshFunc(apiClient client.API, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ret, err := apiClient.Deployments.GetDeploymentByIDV3UsingGET(
			deployments.NewGetDeploymentByIDV3UsingGETParams().
				WithDeploymentID(strfmt.UUID(id)).
				WithExpand([]string{"lastRequest"}).
				WithAPIVersion(withString(DeploymentsAPIVersion)))
		if err != nil {
			return id, models.DeploymentStatusCREATEFAILED, err
		}

		status := ret.Payload.Status
		switch status {
		case models.DeploymentStatusCREATEINPROGRESS, models.DeploymentStatusUPDATEINPROGRESS:
			return ret.Payload.ID.String(), status, nil
		case models.DeploymentStatusCREATESUCCESSFUL, models.DeploymentStatusUPDATESUCCESSFUL:
			deploymentID := ret.Payload.ID
			return deploymentID.String(), status, nil
		case models.DeploymentStatusCREATEFAILED, models.DeploymentStatusUPDATEFAILED:
			return ret.Payload.ID.String(), status, errors.New(ret.Payload.LastRequest.Details)
		default:
			return [...]string{id}, ret.Error(), fmt.Errorf("deploymentStatusRefreshFunc: unknown status %v", status)
		}
	}
}

func updateDeploymentWithNewBlueprint(ctx context.Context, d *schema.ResourceData, m interface{}, apiClient *client.API) diag.Diagnostics {
	log.Printf("Noticed changes to blueprint_id/blueprint_version/blueprint_content. Starting to update existing deployment...")

	blueprintID, blueprintVersion, blueprintContent := "", "", ""
	if v, ok := d.GetOk("blueprint_id"); ok {
		blueprintID = v.(string)
	}

	if v, ok := d.GetOk("blueprint_content"); ok {
		blueprintContent = v.(string)
	}

	// Use GetRawConfig to distinguish between values in the user's configuration of the resource vs state.
	// This allows for migration between using blueprint_id and blueprint_content when values are already present (user specified config takes precedence).
	configValue := d.GetRawConfig()

	configHasBlueprintID := !configValue.GetAttr("blueprint_id").IsNull()
	configHasBlueprintContent := !configValue.GetAttr("blueprint_content").IsNull()

	// Empty blueprintContent if the user has specified blueprintID and not blueprintContent.
	if configHasBlueprintID && !configHasBlueprintContent {
		blueprintContent = ""
	}

	// Empty blueprintID if the user has specified blueprintContent and not blueprintID.
	if configHasBlueprintContent && !configHasBlueprintID {
		blueprintID = ""
	}

	if blueprintID != "" && blueprintContent != "" {
		if blueprintID == "inline-blueprint" {
			blueprintID = ""
		} else {
			return diag.FromErr(errors.New("only one of (blueprint_id, blueprintContent) required"))
		}
	}

	deploymentName := d.Get("name").(string)
	projectID := d.Get("project_id").(string)
	deploymentID := d.Id()

	blueprintRequest := models.BlueprintRequest{
		DeploymentID:   deploymentID,
		DeploymentName: deploymentName,
		ProjectID:      projectID,
	}

	if blueprintID != "" {
		blueprintRequest.BlueprintID = strfmt.UUID(blueprintID)
	}

	if v, ok := d.GetOk("blueprint_version"); ok {
		blueprintVersion = v.(string)
		blueprintRequest.BlueprintVersion = blueprintVersion
	}

	if blueprintContent != "" {
		blueprintRequest.Content = blueprintContent
	}

	if v, ok := d.GetOk("description"); ok {
		blueprintRequest.Description = v.(string)
	}

	if v, ok := d.GetOk("inputs"); ok {
		if blueprintContent != "" {
			blueprintRequest.Inputs = expandInputs(v)
		} else {
			// If the inputs are provided, get the schema from blueprint to convert the provided input values
			// to the type defined in the schema.
			inputs, err := getBlueprintInputsByType(apiClient, blueprintID, blueprintVersion, v)
			if err != nil {
				return diag.FromErr(err)
			}
			blueprintRequest.Inputs = inputs
		}
	} else {
		blueprintRequest.Inputs = make(map[string]interface{})
	}

	if v, ok := d.GetOk("reason"); ok {
		blueprintRequest.Reason = v.(string)
	}

	bpRequestCreated, bpRequestAccepted, err := apiClient.BlueprintRequests.CreateBlueprintRequestUsingPOST1(
		blueprint_requests.NewCreateBlueprintRequestUsingPOST1Params().WithRequest(&blueprintRequest))

	if err != nil {
		log.Printf("Received error. err=%s, bpRequestCreated=%v, bpRequestAccepted=%v", err, bpRequestCreated, bpRequestAccepted)
		return diag.FromErr(err)
	}

	// blueprint_requests service may return either 201 or 202 depending on whether the request is in terminal state vs or in-progress
	log.Printf("Requested deployment '%s' update from blueprint '%s'. bpRequestCreated=%v, bpRequestAccepted=%v", deploymentName, blueprintID, bpRequestCreated, bpRequestAccepted)
	deploymentID, status, failureMessage := "", "", ""
	var bpRequest *models.BlueprintRequest
	if bpRequestAccepted != nil {
		bpRequest = bpRequestAccepted.GetPayload()
	} else {
		bpRequest = bpRequestCreated.GetPayload()
	}

	if bpRequest != nil {
		deploymentID = bpRequest.DeploymentID
		status = bpRequest.Status
		failureMessage = bpRequest.FailureMessage
	}

	if deploymentID != "" {
		d.SetId(deploymentID)
	} else {
		return diag.Errorf("failed to request update to existing deployment. status: %v, message: %v", status, failureMessage)
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.DeploymentStatusCREATEINPROGRESS, models.DeploymentStatusUPDATEINPROGRESS},
		Refresh:    deploymentStatusRefreshFunc(*apiClient, deploymentID),
		Target:     []string{models.DeploymentStatusCREATESUCCESSFUL, models.DeploymentStatusUPDATESUCCESSFUL},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
	}

	if _, err = stateChangeFunc.WaitForStateContext(ctx); err != nil {
		readErrors := resourceDeploymentRead(ctx, d, m)
		if readErrors.HasError() {
			return append(readErrors, diag.Errorf("failed to update deployment: %v", err.Error())...)
		}
		return diag.FromErr(err)
	}

	log.Printf("Finished to update vra_deployment '%s' with blueprint '%s'", deploymentName, blueprintID)
	return nil
}

func updateDeploymentMetadata(d *schema.ResourceData, apiClient *client.API, deploymentUUID strfmt.UUID) error {
	log.Printf("Starting to update deployment name and description")
	description := d.Get("description").(string)
	name := d.Get("name").(string)

	updateDeploymentSpecification := models.DeploymentUpdate{
		Description: description,
		Name:        name,
	}

	log.Printf("[DEBUG] update deployment: %#v", updateDeploymentSpecification)
	_, err := apiClient.Deployments.PatchDeploymentUsingPATCH2(
		deployments.NewPatchDeploymentUsingPATCH2Params().WithDeploymentID(deploymentUUID).WithUpdate(&updateDeploymentSpecification))
	if err != nil {
		return err
	}

	log.Printf("Finished updating deployment name and description")
	return nil
}

// runDeploymentUpdateAction attempts a deployment-level Update Day2 action.
// If no valid deployment-level action is found it falls back to running Update
// actions on individual deployment resources.
// Returns usedResourceActions=true when the fallback resource-level path was used.
func runDeploymentUpdateAction(ctx context.Context, d *schema.ResourceData, apiClient *client.API, deploymentUUID strfmt.UUID) (usedResourceActions bool, err error) {
	log.Printf("[DEBUG] Noticed changes to inputs. Starting to update deployment with inputs")
	// Get the deployment actions
	deploymentActions, err := apiClient.DeploymentActions.GetDeploymentActionsUsingGET2(deployment_actions.
		NewGetDeploymentActionsUsingGET2Params().WithDeploymentID(deploymentUUID))
	if err != nil {
		return false, err
	}

	actionID := ""
	for _, action := range deploymentActions.Payload {
		if strings.Contains(strings.ToLower(action.ID), UpdateDeploymentActionName) {
			if action.Valid {
				actionID = action.ID
			}
			break
		}
	}

	if actionID == "" {
		// No valid deployment-level Update action; fall back to running Day2 actions
		// on individual resources, which some catalog items use instead.
		log.Printf("[DEBUG] no valid deployment-level Update action found, falling back to resource-level")
		oldInputsRaw, newInputsRaw := d.GetChange("inputs")
		changedKeys := computeChangedKeys(oldInputsRaw, newInputsRaw)
		return true, runResourceLevelUpdateActions(ctx, d, apiClient, deploymentUUID, changedKeys)
	}
	log.Printf("[DEBUG] running deployment-level Update action %q", actionID)

	// Continue if update action is available.
	name := d.Get("name")
	var inputs = make(map[string]interface{})
	blueprintID, catalogItemID := "", ""
	if v, ok := d.GetOk("catalog_item_id"); ok {
		catalogItemID = v.(string)
	}

	if v, ok := d.GetOk("blueprint_id"); ok {
		blueprintID = v.(string)
	}

	log.Printf("Checking values before any update [catalog_item_id]: %s, [blueprint_id]: %s.", catalogItemID, blueprintID)

	// If catalog_item_id is provided, get the catalog item schema deployment with the catalog item
	if catalogItemID != "" {
		catalogItemVersion := ""
		if v, ok := d.GetOk("catalog_item_version"); ok {
			catalogItemVersion = v.(string)
		}

		if v, ok := d.GetOk("inputs"); ok {
			// If the inputs are provided, get the schema from catalog item to convert the provided input values
			// to the type defined in the schema.
			inputs, err = getCatalogItemInputsByType(apiClient, catalogItemID, catalogItemVersion, v)
			if err != nil {

				// If the catalog item version is no longer available,
				// check the inputs from the update action
				log.Printf("Error while getting catalog item inputs. Checking with update action instead")

				if deploymentUUID != "" && actionID != "" {
					inputs, err = getDeploymentActionInputsByType(apiClient, deploymentUUID, actionID, v)
					if err != nil {
						return false, err
					}
				} else {
					return false, err
				}
			}
		}
	} else if blueprintID != "" {
		blueprintVersion := ""
		if v, ok := d.GetOk("blueprint_version"); ok {
			blueprintVersion = v.(string)
		}

		if v, ok := d.GetOk("inputs"); ok {
			// If the inputs are provided, get the schema from blueprint to convert the provided input values
			// to the type defined in the schema.
			inputs, err = getBlueprintInputsByType(apiClient, blueprintID, blueprintVersion, v)
			if err != nil {
				return false, err
			}
		}
	}

	reason := deploymentUpdateReason
	err = runAction(ctx, d, apiClient, deploymentUUID, actionID, inputs, reason)
	if err != nil {
		return false, err
	}

	log.Printf("Finished updating vra_deployment %s with inputs", name)
	return false, nil
}

func getInputsByType(inputs map[string]interface{}, inputTypesMap map[string]string) (map[string]interface{}, error) {
	log.Printf("Converting the input values to their types")
	inputsByType := make(map[string]interface{})
	var err error
	for k, v := range inputs {
		if t, ok := inputTypesMap[k]; ok {
			log.Printf("input_key: %s, type: %#v, value provided: %#v", k, t, v)
			switch strings.ToLower(t) {
			case "array":
				value := make([]interface{}, 0)
				err := json.Unmarshal([]byte(v.(string)), &value)
				if err != nil {
					return nil, fmt.Errorf("cannot convert input '%v' value '%v' into type '%s'. %#v", k, v, t, err)
				}
				inputsByType[k] = value
			case "boolean":
				inputsByType[k], err = strconv.ParseBool(v.(string))
			case "integer":
				inputsByType[k], err = strconv.Atoi(v.(string))
			case "number":
				inputsByType[k], err = strconv.ParseFloat(v.(string), 32)
			case "object":
				var value interface{}
				err := json.Unmarshal([]byte(v.(string)), &value)
				if err != nil {
					return nil, fmt.Errorf("cannot convert input '%v' value '%v' into type '%s'. %#v", k, v, t, err)
				}
				inputsByType[k] = value
			default:
				inputsByType[k] = v
			}
			if err != nil {
				return nil, fmt.Errorf("cannot convert %v into %s", v, t)
			}
		} else {
			inputsByType[k] = v
		}
	}
	return inputsByType, nil
}

// Returns a map of string, string with input variables and their types defined in the given schema map.
// Used for getting map of inputs and their types for Catalog item and Deployment actions.
// Properties without a simple "type" string are skipped rather than panicking.
func getInputTypesMapFromSchema(schema map[string]interface{}) (map[string]string, error) {
	log.Printf("Getting the map of inputs and their types")
	inputTypesMap := make(map[string]string, len(schema))
	for k, v := range schema {
		propMap, ok := v.(map[string]interface{})
		if !ok {
			log.Printf("[DEBUG] schema property %q is not a map; skipping", k)
			continue
		}
		rawType, hasType := propMap["type"]
		if !hasType {
			log.Printf("[DEBUG] schema property %q has no 'type' field; skipping", k)
			continue
		}
		typeStr, ok := rawType.(string)
		if !ok {
			log.Printf("[DEBUG] schema property %q 'type' is %T, not a string; skipping", k, rawType)
			continue
		}
		inputTypesMap[k] = typeStr
	}
	return inputTypesMap, nil
}

// Returns a map of string, string with input variables and their types defined in the vRA blueprint
func getInputTypesMapFromBlueprintInputsSchema(schema map[string]models.PropertyDefinition) (map[string]string, error) {
	log.Printf("Getting the map of inputs and their types")
	inputTypesMap := make(map[string]string, len(schema))
	for k, v := range schema {
		inputTypesMap[k] = v.Type
	}
	return inputTypesMap, nil
}

// Returns whether the day2 action is valid currently, exact action ID for a given action string
func getDeploymentDay2ActionID(apiClient *client.API, deploymentUUID strfmt.UUID, actionName string) (bool, string, error) {
	// Get the deployment actions
	deploymentActions, err := apiClient.DeploymentActions.GetDeploymentActionsUsingGET2(deployment_actions.
		NewGetDeploymentActionsUsingGET2Params().WithDeploymentID(deploymentUUID))
	if err != nil {
		return false, "", err
	}

	actionAvailable := false
	actionID := ""

	for _, action := range deploymentActions.Payload {
		if strings.Contains(strings.ToLower(action.ID), strings.ToLower(actionName)) {
			actionID = action.ID
			if action.Valid {
				log.Printf("[DEBUG] %s action is available on the deployment", actionName)
				actionAvailable = true
				return actionAvailable, actionID, nil
			}

			// Day-2 action id is not valid
			log.Printf("[DEBUG] %s action is not valid based on current state of the deployment", actionID)
			return actionAvailable, actionID, fmt.Errorf("%s action is not valid based on current state of the deployment", actionID)
		}
	}
	return actionAvailable, actionID, fmt.Errorf("%s action is not found in the list of day2 actions allowed on the deployment", actionName)
}

func getDeploymentActionInputsByType(apiClient *client.API, deploymentUUID strfmt.UUID, actionID string, inputValues interface{}) (map[string]interface{}, error) {
	inputTypesMap, err := getDeploymentActionInputTypesMap(apiClient, deploymentUUID, actionID)
	if err != nil {
		return nil, err
	}

	log.Printf("InputTypesMap: %v", inputTypesMap)
	inputs, err := getInputsByType(inputValues.(map[string]interface{}), inputTypesMap)
	if err != nil {
		return nil, err
	}

	return inputs, nil
}

// Gets the schema for a given deployment action id
func getDeploymentActionSchema(apiClient *client.API, deploymentUUID strfmt.UUID, actionID string) (map[string]interface{}, error) {
	// Getting the catalog item schema
	log.Printf("Getting the schema for deploymentID: %v, actionID: %v", deploymentUUID, actionID)
	var actionSchema interface{}

	deploymentAction, err := apiClient.DeploymentActions.GetDeploymentActionUsingGET2(deployment_actions.
		NewGetDeploymentActionUsingGET2Params().WithDeploymentID(deploymentUUID).WithActionID(actionID))
	if err != nil {
		return nil, err
	}

	actionSchema = deploymentAction.GetPayload().Schema

	if actionSchema != nil && (actionSchema.(map[string]interface{}))["properties"] != nil {
		actionInputsSchemaMap := (actionSchema.(map[string]interface{}))["properties"].(map[string]interface{})
		return actionInputsSchemaMap, nil
	}
	return make(map[string]interface{}), nil
}

func getDeploymentActionInputTypesMap(apiClient *client.API, deploymentUUID strfmt.UUID, actionID string) (map[string]string, error) {
	inputsSchemaMap, err := getDeploymentActionSchema(apiClient, deploymentUUID, actionID)
	if err != nil {
		return nil, err
	}

	inputTypesMap, err := getInputTypesMapFromSchema(inputsSchemaMap)
	if err != nil {
		return nil, err
	}
	return inputTypesMap, nil
}

// overlayLastRequestInputs merges lastRequest.Inputs into allInputs for keys already
// known to Terraform, so that stale deployment.Inputs are updated after Day2 actions.
func overlayLastRequestInputs(allInputs map[string]interface{}, lastReqInputs map[string]interface{}, userManagedKeys map[string]bool) {
	for k, v := range lastReqInputs {
		if _, exists := allInputs[k]; exists || userManagedKeys[k] {
			allInputs[k] = v
		}
	}
}

// selectBestWriteAction returns the valid write-intent Day2 action whose schema
// overlaps the most with changedKeys, or ("", nil) if none match.
func selectBestWriteAction(apiClient *client.API, deploymentUUID strfmt.UUID, resourceID strfmt.UUID, actions []*models.ResourceAction, changedKeys map[string]bool) (string, map[string]interface{}) {
	bestID := ""
	bestScore := 0
	var bestSchema map[string]interface{}
	for _, action := range actions {
		if !action.Valid {
			continue
		}
		id := strings.ToLower(action.ID)
		if !strings.Contains(id, UpdateDeploymentActionName) && !strings.Contains(id, resourceActionModifyKeyword) {
			continue
		}
		actionSchema, schemaErr := getResourceActionSchema(apiClient, deploymentUUID, resourceID, action.ID)
		if schemaErr != nil {
			log.Printf("[DEBUG] could not fetch schema for action %s: %v", action.ID, schemaErr)
			continue
		}
		score := 0
		for k := range changedKeys {
			if _, ok := actionSchema[k]; ok {
				score++
			}
		}
		if score > bestScore {
			bestID = action.ID
			bestScore = score
			bestSchema = actionSchema
		}
	}
	return bestID, bestSchema
}

// runUpdateActionsForResource runs one or more best-matching write actions on a
// single resource until all changed keys are handled (or no more matching actions
// remain). This handles cases where ports and IPs require separate Day2 actions.
// Returns whether any action ran, which keys were handled, and any error.
func runUpdateActionsForResource(ctx context.Context, d *schema.ResourceData, apiClient *client.API, deploymentUUID strfmt.UUID, resource *models.DeploymentResource, enrichedInputs map[string]interface{}, changedKeys map[string]bool) (bool, map[string]bool, error) {
	actionsResp, err := apiClient.DeploymentActions.GetResourceActionsUsingGET4(
		deployment_actions.NewGetResourceActionsUsingGET4Params().
			WithAPIVersion(withString(DeploymentsAPIVersion)).
			WithDeploymentID(deploymentUUID).
			WithResourceID(resource.ID))
	if err != nil {
		log.Printf("[WARN] failed to list actions for resource %s (%s): %v; skipping", *resource.Name, resource.ID, err)
		return false, nil, nil
	}

	allHandledKeys := make(map[string]bool)
	localRemaining := make(map[string]bool, len(changedKeys))
	for k := range changedKeys {
		localRemaining[k] = true
	}

	anyRan := false
	for len(localRemaining) > 0 {
		bestActionID, bestSchema := selectBestWriteAction(apiClient, deploymentUUID, resource.ID, actionsResp.Payload, localRemaining)
		if bestActionID == "" {
			break
		}

		handledKeys := make(map[string]bool)
		for k := range localRemaining {
			if _, ok := bestSchema[k]; ok {
				handledKeys[k] = true
			}
		}

		inputs, buildErr := buildResourceActionInputs(bestSchema, resource, bestActionID, enrichedInputs)
		if buildErr != nil {
			return anyRan, allHandledKeys, fmt.Errorf("failed to build inputs for resource action %s on %s: %w", bestActionID, *resource.Name, buildErr)
		}
		log.Printf("[DEBUG] running resource-level action %q on resource %s (%s) for keys %v", bestActionID, *resource.Name, resource.ID, handledKeys)
		if runErr := runResourceAction(ctx, d, apiClient, deploymentUUID, resource.ID, bestActionID, inputs); runErr != nil {
			return anyRan, allHandledKeys, runErr
		}

		anyRan = true
		for k := range handledKeys {
			allHandledKeys[k] = true
			delete(localRemaining, k)
		}
	}

	return anyRan, allHandledKeys, nil
}

// runResourceLevelUpdateActions runs Update Day2 actions on individual resources
// as a fallback when no deployment-level Update action is available.
func runResourceLevelUpdateActions(ctx context.Context, d *schema.ResourceData, apiClient *client.API, deploymentUUID strfmt.UUID, changedKeys map[string]bool) error {
	name := d.Get("name")
	log.Printf("[DEBUG] attempting resource-level Update actions for deployment %s", name)

	resourcesResp, err := apiClient.Deployments.GetDeploymentResourcesUsingGET2(
		deployments.NewGetDeploymentResourcesUsingGET2Params().
			WithAPIVersion(withString(DeploymentsAPIVersion)).
			WithDeploymentID(deploymentUUID).
			WithDollarTop(withInt32(DefaultDollarTop)))
	if err != nil {
		return fmt.Errorf("failed to list resources for deployment %s: %w", name, err)
	}

	// Build enriched inputs once: merge Terraform inputs with API-stored deployment inputs.
	terraformInputs, _ := d.GetOk("inputs")
	rawInputs, _ := terraformInputs.(map[string]interface{})
	depResp, depErr := apiClient.Deployments.GetDeploymentByIDV3UsingGET(
		deployments.NewGetDeploymentByIDV3UsingGETParams().
			WithDeploymentID(deploymentUUID).
			WithAPIVersion(withString(DeploymentsAPIVersion)))
	if depErr == nil && depResp.Payload != nil {
		if apiInputs, ok := depResp.Payload.Inputs.(map[string]interface{}); ok {
			enriched := make(map[string]interface{}, len(apiInputs)+len(rawInputs))
			for k, v := range apiInputs {
				enriched[k] = v
			}
			for k, v := range rawInputs {
				enriched[k] = v
			}
			rawInputs = enriched
		}
	}

	remainingKeys := make(map[string]bool, len(changedKeys))
	for k := range changedKeys {
		remainingKeys[k] = true
	}
	updatedCount := 0
	for _, resource := range resourcesResp.Payload.Content {
		ran, handledKeys, runErr := runUpdateActionsForResource(ctx, d, apiClient, deploymentUUID, resource, rawInputs, remainingKeys)
		if runErr != nil {
			return runErr
		}
		if ran {
			updatedCount++
			for k := range handledKeys {
				delete(remainingKeys, k)
			}
		}
	}

	if updatedCount == 0 {
		return fmt.Errorf("'Update' action is not supported at deployment or resource level for deployment %s", name)
	}

	log.Printf("[DEBUG] submitted Update action for %d resource(s) in deployment %s", updatedCount, name)
	return nil
}

func getResourceActionSchema(apiClient *client.API, deploymentUUID strfmt.UUID, resourceID strfmt.UUID, actionID string) (map[string]interface{}, error) {
	log.Printf("[DEBUG] getting schema for resource action %s on resource %s in deployment %s", actionID, resourceID, deploymentUUID)
	action, err := apiClient.DeploymentActions.GetResourceActionUsingGET4(
		deployment_actions.NewGetResourceActionUsingGET4Params().
			WithAPIVersion(withString(DeploymentsAPIVersion)).
			WithDeploymentID(deploymentUUID).
			WithResourceID(resourceID).
			WithActionID(actionID))
	if err != nil {
		return nil, err
	}

	actionSchema := action.GetPayload().Schema
	if actionSchema != nil {
		schemaMap, ok := actionSchema.(map[string]interface{})
		if !ok {
			log.Printf("[DEBUG] action %q schema is %T, not map; treating as empty", actionID, actionSchema)
			return make(map[string]interface{}), nil
		}
		if props, ok := schemaMap["properties"]; ok && props != nil {
			schemaProps, ok := props.(map[string]interface{})
			if !ok {
				log.Printf("[DEBUG] action %q schema properties is %T, not map; treating as empty", actionID, props)
				return make(map[string]interface{}), nil
			}
			if b, err := json.Marshal(schemaProps); err == nil {
				log.Printf("[DEBUG] action %q schema properties for resource %s: %s", actionID, resourceID, b)
			}
			return schemaProps, nil
		}
	}
	log.Printf("[DEBUG] action %q returned empty/nil schema for resource %s", actionID, resourceID)
	return make(map[string]interface{}), nil
}

// computeChangedKeys returns keys whose values differ between old and new input maps.
func computeChangedKeys(oldRaw, newRaw interface{}) map[string]bool {
	changed := make(map[string]bool)
	oldMap, _ := oldRaw.(map[string]interface{})
	newMap, _ := newRaw.(map[string]interface{})
	for k, newV := range newMap {
		if !reflect.DeepEqual(oldMap[k], newV) {
			changed[k] = true
		}
	}
	for k := range oldMap {
		if _, ok := newMap[k]; !ok {
			changed[k] = true
		}
	}
	return changed
}

// isLiteralDefault returns true if v is a usable literal default (not nil, not a map).
func isLiteralDefault(v interface{}) bool {
	if v == nil {
		return false
	}
	_, isMap := v.(map[string]interface{})
	return !isMap
}

// nativePlatformID returns the platform identifier for a resource (providerId or uuid).
func nativePlatformID(props map[string]interface{}) string {
	if id, _ := props["providerId"].(string); id != "" {
		return id
	}
	if id, _ := props["uuid"].(string); id != "" {
		return id
	}
	return ""
}

// pickFieldValue returns the best value for a schema field, checking in order:
// deployment inputs, resource properties, literal schema default, then self-referential
// resource ID (derived from the resource name, e.g. "Pool" → "poolId").
func pickFieldValue(k, fieldType string, typedInputs, currentProps map[string]interface{}, propDefault interface{}, platformID, selfRefIDKey string) interface{} {
	if v := typedInputs[k]; v != nil {
		return v
	}
	if v := currentProps[k]; v != nil {
		return v
	}
	if isLiteralDefault(propDefault) {
		return propDefault
	}
	if fieldType == "string" && platformID != "" && selfRefIDKey != "" && strings.EqualFold(k, selfRefIDKey) {
		return platformID
	}
	return nil
}

// resolveUnmatchedArrayFields matches unresolved array-of-object schema fields to
// deployment inputs by comparing the schema's items.properties keys against the
// object keys in JSON-encoded input values (structural matching for renamed fields).
func resolveUnmatchedArrayFields(schemaMap map[string]interface{}, inputTypesMap map[string]string, rawInputs map[string]interface{}, result map[string]interface{}) {
	for k := range inputTypesMap {
		if _, matched := result[k]; matched {
			continue
		}
		if inputTypesMap[k] != "array" {
			continue
		}
		propMap, _ := schemaMap[k].(map[string]interface{})
		if propMap == nil {
			continue
		}
		itemsMap, _ := propMap["items"].(map[string]interface{})
		if itemsMap == nil {
			continue
		}
		itemType, _ := itemsMap["type"].(string)
		if itemType != "object" {
			continue
		}
		itemProps, _ := itemsMap["properties"].(map[string]interface{})
		if len(itemProps) == 0 {
			continue
		}
		schemaKeys := make(map[string]bool, len(itemProps))
		for pk := range itemProps {
			schemaKeys[pk] = true
		}
		// Sort input names for deterministic matching when multiple inputs
		// have identically-structured JSON arrays.
		sortedInputNames := make([]string, 0, len(rawInputs))
		for inputName := range rawInputs {
			sortedInputNames = append(sortedInputNames, inputName)
		}
		sort.Strings(sortedInputNames)
		for _, inputName := range sortedInputNames {
			inputVal := rawInputs[inputName]
			if _, used := result[inputName]; used {
				continue
			}
			str, isStr := inputVal.(string)
			if !isStr || len(str) == 0 || str[0] != '[' {
				continue
			}
			var parsed []interface{}
			if err := json.Unmarshal([]byte(str), &parsed); err != nil || len(parsed) == 0 {
				continue
			}
			firstObj, ok := parsed[0].(map[string]interface{})
			if !ok {
				continue
			}
			if len(firstObj) != len(schemaKeys) {
				continue
			}
			allMatch := true
			for pk := range schemaKeys {
				if _, ok := firstObj[pk]; !ok {
					allMatch = false
					break
				}
			}
			if allMatch {
				result[k] = parsed
				log.Printf("[DEBUG] structurally matched deployment input %q to schema field %q", inputName, k)
				break
			}
		}
	}
}

// resolveByTitle matches unresolved schema fields that have $dynamicDefault or $data
// to deployment inputs by normalizing the field's title (lowercase, no spaces) and
// looking up a matching input name.
func resolveByTitle(schemaMap map[string]interface{}, rawInputs map[string]interface{}, inputTypesMap map[string]string, result map[string]interface{}) {
	inputsByLower := make(map[string]interface{}, len(rawInputs))
	for name, val := range rawInputs {
		inputsByLower[strings.ToLower(name)] = val
	}

	for k := range inputTypesMap {
		if _, matched := result[k]; matched {
			continue
		}
		propMap, _ := schemaMap[k].(map[string]interface{})
		if propMap == nil {
			continue
		}
		if ro, _ := propMap["readOnly"].(bool); ro {
			continue
		}
		_, hasDynamic := propMap["$dynamicDefault"]
		_, hasData := propMap["$data"]
		if !hasDynamic && !hasData {
			continue
		}
		title, _ := propMap["title"].(string)
		if title == "" {
			continue
		}
		normalizedTitle := strings.ToLower(strings.ReplaceAll(title, " ", ""))
		v, ok := inputsByLower[normalizedTitle]
		if !ok {
			continue
		}
		fieldType := inputTypesMap[k]
		switch strings.ToLower(fieldType) {
		case "array":
			if s, isStr := v.(string); isStr && len(s) > 0 && s[0] == '[' {
				var parsed interface{}
				if err := json.Unmarshal([]byte(s), &parsed); err == nil {
					result[k] = parsed
					log.Printf("[DEBUG] title-matched deployment input %q to schema field %q (title %q)", normalizedTitle, k, title)
					continue
				}
			}
		case "boolean":
			if s, isStr := v.(string); isStr {
				if b, err := strconv.ParseBool(s); err == nil {
					result[k] = b
					log.Printf("[DEBUG] title-matched deployment input %q to schema field %q (title %q)", normalizedTitle, k, title)
					continue
				}
			}
		case "integer":
			if s, isStr := v.(string); isStr {
				if n, err := strconv.Atoi(s); err == nil {
					result[k] = n
					log.Printf("[DEBUG] title-matched deployment input %q to schema field %q (title %q)", normalizedTitle, k, title)
					continue
				}
			}
		case "number":
			if s, isStr := v.(string); isStr {
				if f, err := strconv.ParseFloat(s, 64); err == nil {
					result[k] = f
					log.Printf("[DEBUG] title-matched deployment input %q to schema field %q (title %q)", normalizedTitle, k, title)
					continue
				}
			}
		default:
			result[k] = v
			log.Printf("[DEBUG] title-matched deployment input %q to schema field %q (title %q)", normalizedTitle, k, title)
		}
	}
}

// resolveFromResourceProperties populates unresolved array-of-string schema fields
// from the resource's current properties. It matches fields (without $dynamicDefault)
// whose title keywords appear as keys in resource property arrays.
func resolveFromResourceProperties(schemaMap map[string]interface{}, inputTypesMap map[string]string, currentProps map[string]interface{}, result map[string]interface{}) {
	for k := range inputTypesMap {
		if _, matched := result[k]; matched {
			continue
		}
		if inputTypesMap[k] != "array" {
			continue
		}
		propMap, _ := schemaMap[k].(map[string]interface{})
		if propMap == nil {
			continue
		}
		if ro, _ := propMap["readOnly"].(bool); ro {
			continue
		}
		itemsMap, _ := propMap["items"].(map[string]interface{})
		if itemsMap != nil {
			if itemType, _ := itemsMap["type"].(string); itemType == "object" {
				continue
			}
		}
		if _, hasDyn := propMap["$dynamicDefault"]; hasDyn {
			continue
		}
		if _, hasData := propMap["$data"]; hasData {
			continue
		}
		title, _ := propMap["title"].(string)
		if title == "" {
			continue
		}
		values := extractValuesFromProperties(title, currentProps)
		if len(values) > 0 {
			result[k] = values
			log.Printf("[DEBUG] populated schema field %q (title %q) from resource properties: %v", k, title, values)
		}
	}
}

// fieldNameMatchesTitle returns true if fieldName appears as a word-prefix in
// any word of the title. This handles plurals (field "port" matches title word
// "Ports") while rejecting substring matches ("port" does not match "transport"
// because "transport" does not start with "port").
func fieldNameMatchesTitle(fieldName, title string) bool {
	fieldLower := strings.ToLower(fieldName)
	for _, word := range strings.Fields(title) {
		wordLower := strings.ToLower(word)
		if strings.HasPrefix(wordLower, fieldLower) {
			return true
		}
	}
	return false
}

// extractValuesFromProperties scans resource property arrays for fields whose name
// appears in the given title, and returns their string values.
func extractValuesFromProperties(title string, props map[string]interface{}) []interface{} {
	// Sort property keys for deterministic matching.
	sortedPropKeys := make([]string, 0, len(props))
	for k := range props {
		sortedPropKeys = append(sortedPropKeys, k)
	}
	sort.Strings(sortedPropKeys)
	for _, propKey := range sortedPropKeys {
		propVal := props[propKey]
		arr, ok := propVal.([]interface{})
		if !ok || len(arr) == 0 {
			continue
		}
		firstObj, ok := arr[0].(map[string]interface{})
		if !ok {
			continue
		}
		// Sort field names for deterministic matching.
		sortedFields := make([]string, 0, len(firstObj))
		for fn := range firstObj {
			sortedFields = append(sortedFields, fn)
		}
		sort.Strings(sortedFields)
		for _, fieldName := range sortedFields {
			// Match field names as word prefixes in the title to handle plurals
			// (e.g. field "port" matches title word "Ports") while avoiding false
			// positives from substring matches (e.g. "port" must not match "transport").
			if !fieldNameMatchesTitle(fieldName, title) {
				continue
			}
			var values []interface{}
			for _, item := range arr {
				if obj, ok := item.(map[string]interface{}); ok {
					if v, exists := obj[fieldName]; exists {
						values = append(values, fmt.Sprint(v))
					}
				}
			}
			if len(values) > 0 {
				return values
			}
		}
	}
	return nil
}

// buildResourceActionInputs constructs the input map for a resource-level Day2 action
// by merging deployment inputs, resource properties, schema defaults, and heuristics
// for fields with auto-generated IDs. The schemaMap parameter should be the pre-fetched
// action schema from selectBestWriteAction to avoid a redundant API call.
// The enrichedInputs parameter should already contain both Terraform and API-stored
// deployment inputs, merged by the caller.
func buildResourceActionInputs(schemaMap map[string]interface{}, resource *models.DeploymentResource, actionID string, enrichedInputs map[string]interface{}) (map[string]interface{}, error) {

	inputTypesMap, err := getInputTypesMapFromSchema(schemaMap)
	if err != nil {
		return nil, err
	}
	if len(inputTypesMap) == 0 {
		return make(map[string]interface{}), nil
	}

	currentProps, _ := resource.Properties.(map[string]interface{})
	if currentProps == nil {
		currentProps = make(map[string]interface{})
	}

	rawInputs := enrichedInputs
	log.Printf("[DEBUG] buildResourceActionInputs: %d enriched input keys", len(rawInputs))
	filtered := make(map[string]interface{}, len(inputTypesMap))
	for k := range inputTypesMap {
		if v, ok := rawInputs[k]; ok {
			filtered[k] = v
		}
	}
	typedInputs, castErr := getInputsByType(filtered, inputTypesMap)
	if castErr != nil {
		log.Printf("[WARN] type conversion failed for resource action %q inputs: %v; using raw values", actionID, castErr)
		typedInputs = filtered
	}

	platformID := nativePlatformID(currentProps)
	selfRefIDKey := ""
	if resource.Name != nil && len(*resource.Name) > 0 {
		n := *resource.Name
		candidate := strings.ToLower(n[:1]) + n[1:] + "Id"
		// Only use the heuristic if the schema actually contains this key.
		if _, exists := schemaMap[candidate]; exists {
			selfRefIDKey = candidate
		}
	}
	result := make(map[string]interface{}, len(schemaMap))
	for k, fieldType := range inputTypesMap {
		propMap, _ := schemaMap[k].(map[string]interface{})
		if ro, _ := propMap["readOnly"].(bool); ro {
			continue
		}
		if v := pickFieldValue(k, fieldType, typedInputs, currentProps, propMap["default"], platformID, selfRefIDKey); v != nil {
			result[k] = v
		}
	}

	resolveUnmatchedArrayFields(schemaMap, inputTypesMap, rawInputs, result)

	resolveByTitle(schemaMap, rawInputs, inputTypesMap, result)

	resolveFromResourceProperties(schemaMap, inputTypesMap, currentProps, result)

	// Handle properties without a simple "type" that getInputTypesMapFromSchema skipped.
	for k, propRaw := range schemaMap {
		if _, handled := inputTypesMap[k]; handled {
			continue
		}
		if _, exists := result[k]; exists {
			continue
		}
		propMap, _ := propRaw.(map[string]interface{})
		if propMap != nil {
			if ro, _ := propMap["readOnly"].(bool); ro {
				continue
			}
		}
		if v, ok := rawInputs[k]; ok {
			if str, isStr := v.(string); isStr && len(str) > 0 && (str[0] == '[' || str[0] == '{') {
				var parsed interface{}
				if err := json.Unmarshal([]byte(str), &parsed); err == nil {
					result[k] = parsed
					continue
				}
			}
			result[k] = v
			continue
		}
		if v := currentProps[k]; v != nil {
			result[k] = v
		}
	}

	if b, err := json.Marshal(result); err == nil {
		log.Printf("[DEBUG] resource action %q inputs for %q: %s", actionID, *resource.Name, b)
	}
	return result, nil
}

// runResourceAction submits a Day2 action on a specific resource and waits for completion.
func runResourceAction(ctx context.Context, d *schema.ResourceData, apiClient *client.API, deploymentUUID strfmt.UUID, resourceID strfmt.UUID, actionID string, inputs map[string]interface{}) error {
	resourceActionRequest := models.ResourceActionRequest{
		ActionID: actionID,
		Reason:   deploymentUpdateReason,
		Inputs:   inputs,
	}

	resp, err := apiClient.DeploymentActions.SubmitResourceActionRequestUsingPOST4(
		deployment_actions.NewSubmitResourceActionRequestUsingPOST4Params().
			WithAPIVersion(withString(DeploymentsAPIVersion)).
			WithDeploymentID(deploymentUUID).
			WithResourceID(resourceID).
			WithActionRequest(&resourceActionRequest))
	if err != nil {
		return err
	}

	requestID := resp.GetPayload().ID

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestStatusPENDING, models.RequestStatusINITIALIZATION, models.RequestStatusCHECKINGAPPROVAL, models.RequestStatusAPPROVALPENDING, models.RequestStatusINPROGRESS},
		Refresh:    deploymentActionStatusRefreshFunc(*apiClient, deploymentUUID, requestID),
		Target:     []string{models.RequestStatusCOMPLETION, models.RequestStatusAPPROVALREJECTED, models.RequestStatusABORTED, models.RequestStatusSUCCESSFUL, models.RequestStatusFAILED},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 5 * time.Second,
	}
	if _, err = stateChangeFunc.WaitForStateContext(ctx); err != nil {
		return err
	}
	return nil
}

func runChangeOwnerDeploymentAction(ctx context.Context, d *schema.ResourceData, apiClient *client.API, deploymentUUID strfmt.UUID) error {
	oldOwner, newOwner := d.GetChange("owner")
	log.Printf("Noticed changes to owner. Starting to change deployment owner from %s to %s", oldOwner.(string), newOwner.(string))

	// Get the deployment actionID for Change Owner
	isActionValid, actionID, err := getDeploymentDay2ActionID(apiClient, deploymentUUID, ChangeOwnerDeploymentActionName)
	if err != nil {
		return fmt.Errorf("noticed changes to owner. But, %s", err.Error())
	}

	if !isActionValid {
		return fmt.Errorf("noticed changes to owner, but 'Change Owner' action is not found or supported")
	}

	// Continue if 'Change Owner' action is available. Get action inputs for the 'ChangeOwner' action
	actionInputs := make(map[string]interface{})
	actionInputs["New Owner"] = newOwner

	actionInputTypesMap, err := getDeploymentActionInputTypesMap(apiClient, deploymentUUID, actionID)
	if err != nil {
		return err
	}

	inputs, err := getInputsByType(actionInputs, actionInputTypesMap)
	if err != nil {
		return fmt.Errorf("unable to create action inputs for %v. %v", actionID, err.Error())
	}

	reason := "Updated deployment owner from vRA provider for Terraform."
	err = runAction(ctx, d, apiClient, deploymentUUID, actionID, inputs, reason)
	if err != nil {
		return err
	}

	log.Printf("Finished changing owner for vra_deployment %s with new owner %v", d.Get("name").(string), newOwner)
	return nil
}

func runAction(ctx context.Context, d *schema.ResourceData, apiClient *client.API, deploymentUUID strfmt.UUID, actionID string, inputs map[string]interface{}, reason string) error {
	resourceActionRequest := models.ResourceActionRequest{
		ActionID: actionID,
		Reason:   reason,
		Inputs:   inputs,
	}

	resp, err := apiClient.DeploymentActions.SubmitDeploymentActionRequestUsingPOST2(
		deployment_actions.NewSubmitDeploymentActionRequestUsingPOST2Params().
			WithAPIVersion(withString(DeploymentsAPIVersion)).
			WithDeploymentID(deploymentUUID).
			WithActionRequest(&resourceActionRequest))
	if err != nil {
		return err
	}

	requestID := resp.GetPayload().ID

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestStatusPENDING, models.RequestStatusINITIALIZATION, models.RequestStatusCHECKINGAPPROVAL, models.RequestStatusAPPROVALPENDING, models.RequestStatusINPROGRESS},
		Refresh:    deploymentActionStatusRefreshFunc(*apiClient, deploymentUUID, requestID),
		Target:     []string{models.RequestStatusCOMPLETION, models.RequestStatusAPPROVALREJECTED, models.RequestStatusABORTED, models.RequestStatusSUCCESSFUL, models.RequestStatusFAILED},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 5 * time.Second,
	}
	if _, err = stateChangeFunc.WaitForStateContext(ctx); err != nil {
		return err
	}
	return nil
}

func deploymentActionStatusRefreshFunc(apiClient client.API, deploymentUUID strfmt.UUID, _ strfmt.UUID) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ret, err := apiClient.Deployments.GetDeploymentByIDV3UsingGET(
			deployments.NewGetDeploymentByIDV3UsingGETParams().
				WithDeploymentID(deploymentUUID).
				WithExpand([]string{"lastRequest"}).
				WithAPIVersion(withString(DeploymentsAPIVersion)))
		if err != nil {
			return "", models.RequestStatusFAILED, err
		}

		status := ret.Payload.LastRequest.Status
		switch status {
		case models.RequestStatusPENDING, models.RequestStatusINITIALIZATION, models.RequestStatusCHECKINGAPPROVAL, models.RequestStatusAPPROVALPENDING, models.RequestStatusINPROGRESS, models.RequestStatusCOMPLETION:
			return [...]string{deploymentUUID.String()}, status, nil
		case models.RequestStatusAPPROVALREJECTED, models.RequestStatusABORTED:
			return []string{""}, status, errors.New(ret.Error())
		case models.RequestStatusFAILED:
			return [...]string{deploymentUUID.String()}, status, errors.New(ret.Payload.LastRequest.Details)
		case models.RequestStatusSUCCESSFUL:
			deploymentID := ret.Payload.ID
			return deploymentID.String(), status, nil
		default:
			return [...]string{deploymentUUID.String()}, ret.Error(), fmt.Errorf("deploymentActionStatusRefreshFunc: unknown status %v", status)
		}
	}
}

func deploymentDeleteStatusRefreshFunc(apiClient client.API, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ret, err := apiClient.Deployments.GetDeploymentByIDV3UsingGET(
			deployments.NewGetDeploymentByIDV3UsingGETParams().
				WithDeploymentID(strfmt.UUID(id)).
				WithExpand([]string{"lastRequest"}).
				WithAPIVersion(withString(DeploymentsAPIVersion)))
		if err != nil {
			switch err.(type) {
			case *deployments.GetDeploymentByIDV3UsingGETNotFound:
				return "", reflect.TypeOf(err).String(), nil
			default:
				return [...]string{id}, reflect.TypeOf(err).String(), errors.New(ret.Error())
			}
		}

		status := ret.Payload.Status
		switch status {
		case models.DeploymentStatusDELETEINPROGRESS:
			return [...]string{id}, reflect.TypeOf(ret).String(), nil
		case models.DeploymentStatusDELETESUCCESSFUL:
			return [...]string{id}, reflect.TypeOf(ret).String(), nil
		case models.DeploymentStatusDELETEFAILED:
			return [...]string{id}, reflect.TypeOf(ret).String(), errors.New(ret.Payload.LastRequest.Details)
		default:
			return [...]string{id}, reflect.TypeOf(ret).String(), fmt.Errorf("deploymentStatusRefreshFunc: unknown status %v", status)
		}
	}
}

// updateUserInputs reconciles Terraform's known inputs with platform-reported values.
// Platform values are preferred; user values are kept only for keys the platform doesn't expose.
func updateUserInputs(allInputs, userInputs map[string]interface{}, inputTypesMap map[string]string) map[string]interface{} {
	if allInputs == nil || userInputs == nil {
		return nil
	}

	inputs := make(map[string]interface{})
	for name, value := range userInputs {
		if value != nil {
			platformVal := allInputs[name]
			if platformVal == nil {
				inputs[name] = value
			} else {
				inputs[name] = decodeInputValue(name, platformVal, inputTypesMap)
			}
		}
		log.Printf("Converted incoming value to string: Key: %v, Value: %v, Converted value: %#v", name, allInputs[name], inputs[name])
	}

	return inputs
}

// decodeInputValue converts a platform input value to the string representation
// expected by Terraform state, marshaling arrays/objects to JSON.
func decodeInputValue(inputName string, inputValue interface{}, inputTypesMap map[string]string) interface{} {
	if t, ok := inputTypesMap[inputName]; ok {
		log.Printf("[DEBUG] input_key: %s, schema_type: %s, value: %#v", inputName, t, inputValue)
		switch strings.ToLower(t) {
		case "array", "object":
			value, err := json.Marshal(inputValue)
			if err != nil {
				log.Printf("[ERROR] Cannot marshal input '%s' to JSON: %v", inputName, err)
				return fmt.Sprint(inputValue)
			}
			return string(value)
		default:
			return fmt.Sprint(inputValue)
		}
	}

	switch inputValue.(type) {
	case []interface{}, map[string]interface{}:
		value, err := json.Marshal(inputValue)
		if err != nil {
			log.Printf("[ERROR] Cannot auto-marshal input '%s' to JSON: %v", inputName, err)
			return fmt.Sprint(inputValue)
		}
		return string(value)
	default:
		return fmt.Sprint(inputValue)
	}
}
