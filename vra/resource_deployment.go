package vra

import (
	"fmt"
	"strconv"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client"
	"github.com/vmware/vra-sdk-go/pkg/client/blueprint"
	"github.com/vmware/vra-sdk-go/pkg/client/blueprint_requests"
	"github.com/vmware/vra-sdk-go/pkg/client/catalog_items"
	"github.com/vmware/vra-sdk-go/pkg/client/deployment_actions"
	"github.com/vmware/vra-sdk-go/pkg/client/deployments"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"log"
	"reflect"
	"strings"
	"time"
)

func resourceDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceDeploymentCreate,
		Read:   resourceDeploymentRead,
		Update: resourceDeploymentUpdate,
		Delete: resourceDeploymentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"blueprint_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"blueprint_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"blueprint_content": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"catalog_item_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"catalog_item_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expand_last_request": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Deprecated. True by default even if not provided.",
			},
			"expand_project": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"expand_resources": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Deprecated. True by default even if not provided.",
			},
			"expense": expenseSchema(),
			"inputs": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			//TODO: add last_request
			"last_updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lease_expire_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project": resourceReferenceSchema(),
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resources": resourcesSchema(),
			// TODO: Add plan / simulate feature
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceDeploymentCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("Starting to create vra_deployment resource")
	apiClient := m.(*Client).apiClient

	blueprintId, catalogItemId, blueprintContent := "", "", ""
	if v, ok := d.GetOk("blueprint_id"); ok {
		blueprintId = v.(string)
	}

	if v, ok := d.GetOk("catalog_item_id"); ok {
		catalogItemId = v.(string)
	}

	if v, ok := d.GetOk("blueprint_content"); ok {
		blueprintContent = v.(string)
	}

	if blueprintId != "" && catalogItemId != "" {
		return fmt.Errorf("only one of (blueprint_id, catalog_item_id) required")
	}

	deploymentName := d.Get("name").(string)
	projectId := d.Get("project_id").(string)

	getResp, err := apiClient.Deployments.CheckDeploymentNameUsingGET(deployments.NewCheckDeploymentNameUsingGETParams().WithName(deploymentName))
	log.Printf("getResp: %v, err: %v", getResp, err)

	if err != nil {
		switch err.(type) {
		case *deployments.CheckDeploymentNameUsingGETNotFound:
			log.Printf("Deployment '%v' doesn't exist already and hence can be created", deploymentName)
		}
	} else {
		return fmt.Errorf("a deployment with name '%v' exists already. Try with a differnet name", deploymentName)
	}

	// If catalog_item_id is provided, request deployment with the catalog item
	if catalogItemId != "" {
		log.Printf("Requesting vra_deployment '%s' from catalog item", d.Get("name"))
		catalogItemVersion := ""
		if v, ok := d.GetOk("catalog_item_version"); ok {
			catalogItemVersion = v.(string)
		}

		catalogItemRequest := models.CatalogItemRequest{
			DeploymentName: deploymentName,
			ProjectID:      projectId,
			Version:        catalogItemVersion,
		}

		if v, ok := d.GetOk("inputs"); ok {
			inputsSchemaMap, err := getCatalogItemSchema(apiClient, catalogItemId, catalogItemVersion)
			if err != nil {
				return err
			}

			log.Printf("Inputs Schema: %v", inputsSchemaMap)

			inputTypesMap, err := getInputTypesMapFromCatalogItemSchema(inputsSchemaMap)
			if err != nil {
				return err
			}

			log.Printf("InputTypesMap: %v", inputTypesMap)

			catalogItemRequest.Inputs, err = getInputsByType(v.(map[string]interface{}), inputTypesMap)
			if err != nil {
				return err
			}
		} else {
			catalogItemRequest.Inputs = make(map[string]interface{})
		}

		if v, ok := d.GetOk("description"); ok {
			catalogItemRequest.Reason = v.(string)
		}

		log.Printf("[DEBUG] Create deployment: %#v", catalogItemRequest)
		postOk, err := apiClient.CatalogItems.RequestCatalogItemUsingPOST(
			catalog_items.NewRequestCatalogItemUsingPOSTParams().WithID(strfmt.UUID(catalogItemId)).
				WithRequest(&catalogItemRequest))

		if err != nil {
			return err
		}

		d.SetId(postOk.GetPayload().DeploymentID)
		log.Printf("Finished requesting vra_deployment '%s' from catalog item", d.Get("name"))
	} else {
		blueprintRequest := models.BlueprintRequest{
			DeploymentName: deploymentName,
			ProjectID:      projectId,
		}

		if blueprintId != "" {
			blueprintRequest.BlueprintID = strfmt.UUID(blueprintId)
		} else {
			// Create empty content in the blueprint
			blueprintRequest.Content = " "
		}

		if v, ok := d.GetOk("blueprint_version"); ok {
			blueprintRequest.BlueprintVersion = v.(string)
		}

		if blueprintContent != "" {
			blueprintRequest.Content = blueprintContent
		}

		if v, ok := d.GetOk("description"); ok {
			blueprintRequest.Description = v.(string)
		}

		if v, ok := d.GetOk("inputs"); ok {
			blueprintRequest.Inputs = v
		} else {
			blueprintRequest.Inputs = make(map[string]interface{})
		}

		bpRequestCreated, bpRequestAccepted, err := apiClient.BlueprintRequests.CreateBlueprintRequestUsingPOST1(
			blueprint_requests.NewCreateBlueprintRequestUsingPOST1Params().WithRequest(&blueprintRequest))

		if err != nil {
			log.Printf("Received error. err=%s, bpRequestCreated=%v, bpRequestAccepted=%v", err, bpRequestCreated, bpRequestAccepted)
			return err
		}

		// blueprint_requests service may return either 201 or 202 depending on whether the request is in terminal state vs or in-progress
		log.Printf("Requested deployment from blueprint. bpRequestCreated=%v, bpRequestAccepted=%v", bpRequestCreated, bpRequestAccepted)
		deploymentId, status, failureMessage := "", "", ""
		var bpRequest *models.BlueprintRequest
		if bpRequestAccepted != nil {
			bpRequest = bpRequestAccepted.GetPayload()
		} else {
			bpRequest = bpRequestCreated.GetPayload()
		}

		if bpRequest != nil {
			deploymentId = bpRequest.DeploymentID
			status = bpRequest.Status
			failureMessage = bpRequest.FailureMessage
		}

		if deploymentId != "" {
			d.SetId(deploymentId)
		} else {
			return fmt.Errorf("failed to request for a deployment. status: %v, message: %v", status, failureMessage)
		}

		log.Printf("Finished requesting vra_deployment '%s' from blueprint %s", d.Get("name"), blueprintId)
	}

	stateChangeFunc := resource.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.DeploymentStatusCREATEINPROGRESS},
		Refresh:    deploymentCreateStatusRefreshFunc(*apiClient, d.Id()),
		Target:     []string{models.DeploymentStatusCREATESUCCESSFUL},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
	}

	deploymentId, err := stateChangeFunc.WaitForState()
	if err != nil {
		return err
	}

	d.SetId(deploymentId.(string))
	log.Printf("Finished to create vra_deployment resource with name %s", d.Get("name"))

	return resourceDeploymentRead(d, m)
}

func getCatalogItemSchema(apiClient *client.MulticloudIaaS, catalogItemId string, catalogItemVersion string) (map[string]interface{}, error) {
	// Getting the catalog item schema
	log.Printf("Getting the schema for catalog item: %v version: %v", catalogItemId, catalogItemVersion)
	var catalogItemSchema interface{}
	if catalogItemVersion == "" {
		getItemResp, err := apiClient.CatalogItems.GetCatalogItemUsingGET1(catalog_items.NewGetCatalogItemUsingGET1Params().WithID(strfmt.UUID(catalogItemId)))
		if err != nil {
			return nil, err
		}
		catalogItemSchema = getItemResp.GetPayload().Schema
	} else {
		getVersionResp, err := apiClient.CatalogItems.GetVersionByIDUsingGET(catalog_items.NewGetVersionByIDUsingGETParams().WithID(strfmt.UUID(catalogItemId)).WithVersionID(catalogItemVersion))
		if err != nil {
			return nil, err
		}
		catalogItemSchema = getVersionResp.GetPayload().Schema
	}
	inputsSchemaMap := (catalogItemSchema.(map[string]interface{}))["properties"].(map[string]interface{})
	return inputsSchemaMap, nil
}

func getBlueprintSchema(apiClient *client.MulticloudIaaS, blueprintId string, blueprintVersion string) (map[string]models.PropertyDefinition, error) {
	// Getting the blueprint inputs schema
	log.Printf("Getting the schema for catalog item: %v version: %v", blueprintId, blueprintVersion)
	var blueprintInputsSchema map[string]models.PropertyDefinition
	if blueprintVersion == "" {
		getItemResp, err := apiClient.Blueprint.GetBlueprintInputsSchemaUsingGET1(blueprint.NewGetBlueprintInputsSchemaUsingGET1Params().WithBlueprintID(blueprintId))
		if err != nil {
			return nil, err
		}
		blueprintInputsSchema = getItemResp.GetPayload().Properties
	} else {
		getVersionResp, err := apiClient.Blueprint.GetBlueprintVersionInputsSchemaUsingGET1(
			blueprint.NewGetBlueprintVersionInputsSchemaUsingGET1Params().WithBlueprintID(blueprintId).
				WithVersion(blueprintVersion))
		if err != nil {
			return nil, err
		}
		blueprintInputsSchema = getVersionResp.GetPayload().Properties
	}
	return blueprintInputsSchema, nil
}

func deploymentCreateStatusRefreshFunc(apiClient client.MulticloudIaaS, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ret, err := apiClient.Deployments.GetDeploymentByIDUsingGET(
			deployments.NewGetDeploymentByIDUsingGETParams().WithDepID(strfmt.UUID(id)))
		if err != nil {
			return "", models.DeploymentStatusCREATEFAILED, err
		}

		status := ret.Payload.Status
		switch status {
		case models.DeploymentStatusCREATEFAILED:
			return []string{""}, status, fmt.Errorf(ret.Error())
		case models.DeploymentStatusCREATEINPROGRESS:
			return [...]string{id}, status, nil
		case models.DeploymentStatusCREATESUCCESSFUL:
			deploymentId := ret.Payload.ID
			return deploymentId.String(), status, nil
		default:
			return [...]string{id}, ret.Error(), fmt.Errorf("deploymentCreateStatusRefreshFunc: unknown status %v", status)
		}
	}
}

func resourceDeploymentRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("Reading the vra_deployment resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	expandProject := d.Get("expand_project").(bool)

	resp, err := apiClient.Deployments.GetDeploymentByIDUsingGET(
		deployments.NewGetDeploymentByIDUsingGETParams().WithDepID(strfmt.UUID(id)).
			WithExpandResources(withBool(true)).WithExpandLastRequest(withBool(true)).
			WithExpandProject(withBool(expandProject)))
	if err != nil {
		switch err.(type) {
		case *deployments.GetDeploymentByIDUsingGETNotFound:
			d.SetId("")
			return nil
		}
		return err
	}

	deployment := *resp.Payload
	d.Set("name", deployment.Name)
	d.Set("description", deployment.Description)
	d.Set("blueprint_id", deployment.BlueprintID)
	d.Set("blueprint_version", deployment.BlueprintVersion)
	d.Set("catalog_item_id", deployment.CatalogItemID)
	d.Set("catalog_item_version", deployment.CatalogItemVersion)
	d.Set("created_at", deployment.CreatedAt)
	d.Set("created_by", deployment.CreatedBy)
	//TODO: Set last_request
	d.Set("last_updated_at", deployment.LastUpdatedAt)
	d.Set("last_updated_by", deployment.LastUpdatedBy)
	d.Set("lease_expire_at", deployment.LeaseExpireAt)
	d.Set("project_id", deployment.ProjectID)
	d.Set("status", deployment.Status)

	if err := d.Set("project", flattenResourceReference(deployment.Project)); err != nil {
		return fmt.Errorf("error setting project in deployment - error: %#v", err)
	}

	if err := d.Set("resources", flattenResources(deployment.Resources)); err != nil {
		return fmt.Errorf("error setting resources in deployment - error: %#v", err)
	}

	if err := d.Set("expense", flattenExpense(deployment.Expense)); err != nil {
		return fmt.Errorf("error setting deployment expense - error: %#v", err)
	}

	if err := d.Set("inputs", expandInputs(deployment.Inputs)); err != nil {
		return fmt.Errorf("error setting deployment inputs - error: %#v", err)
	}

	log.Printf("Finished reading the vra_deployment resource with name %s", d.Get("name"))
	return nil
}

func resourceDeploymentUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("Starting to update the vra_deployment resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	deploymentUuid := strfmt.UUID(id)
	if d.HasChange("name") || d.HasChange("description") {
		err := updateDeploymentMetadata(d, apiClient, deploymentUuid)
		if err != nil {
			return err
		}
	}

	if d.HasChange("inputs") {
		err := runDeploymentUpdateAction(d, apiClient, deploymentUuid)
		if err != nil {
			return err
		}
	}

	d.Partial(false)
	log.Printf("Finished updating the vra_deployment resource with name %s", d.Get("name"))
	return resourceDeploymentRead(d, m)
}

func updateDeploymentMetadata(d *schema.ResourceData, apiClient *client.MulticloudIaaS, deploymentUuid strfmt.UUID) error {
	log.Printf("Starting to update deployment name and description")
	description := d.Get("description").(string)
	name := d.Get("name").(string)

	updateDeploymentSpecification := models.DeploymentUpdate{
		Description: description,
		Name:        name,
	}

	log.Printf("[DEBUG] update deployment: %#v", updateDeploymentSpecification)
	_, err := apiClient.Deployments.PatchDeploymentUsingPATCH(
		deployments.NewPatchDeploymentUsingPATCHParams().WithDepID(deploymentUuid).WithUpdate(&updateDeploymentSpecification))
	if err != nil {
		return err
	}

	d.SetPartial("name")
	d.SetPartial("description")
	log.Printf("Finished updating deployment name and description")
	return nil
}

func runDeploymentUpdateAction(d *schema.ResourceData, apiClient *client.MulticloudIaaS, deploymentUuid strfmt.UUID) error {
	log.Printf("Noticed changes to inputs. Starting to update deployment with inputs")
	// Get the deployment actions
	deploymentActions, err := apiClient.DeploymentActions.GetDeploymentActionsUsingGET(deployment_actions.
		NewGetDeploymentActionsUsingGETParams().WithDepID(deploymentUuid))
	if err != nil {
		return err
	}

	updateAvailable := false
	actionId := ""
	for _, action := range deploymentActions.Payload {
		if strings.Contains(strings.ToLower(action.ID), strings.ToLower("Update")) {
			// Update action is available on the deployment
			updateAvailable = true
			actionId = action.ID
			break
		}
	}

	name := d.Get("name")
	if !updateAvailable {
		log.Printf("Update action is not available on the vra_deployment %s, hence new inputs are not applied", name)
	} else {
		var inputs = make(map[string]interface{})
		blueprintId, catalogItemId := "", ""
		if v, ok := d.GetOk("catalog_item_id"); ok {
			catalogItemId = v.(string)
		}

		if v, ok := d.GetOk("blueprint_id"); ok {
			blueprintId = v.(string)
		}

		// If catalog_item_id is provided, get the catalog item schema deployment with the catalog item
		if catalogItemId != "" {
			catalogItemVersion := ""
			if v, ok := d.GetOk("catalog_item_version"); ok {
				catalogItemVersion = v.(string)
			}

			if v, ok := d.GetOk("inputs"); ok {
				// If the inputs are provided, get the schema from catalog item to convert the provided input values
				// to the type defined in the schema.
				inputsSchemaMap, err := getCatalogItemSchema(apiClient, catalogItemId, catalogItemVersion)
				if err != nil {
					return err
				}

				log.Printf("Inputs Schema: %v", inputsSchemaMap)
				inputTypesMap, err := getInputTypesMapFromCatalogItemSchema(inputsSchemaMap)
				if err != nil {
					return err
				}

				log.Printf("InputTypesMap: %v", inputTypesMap)
				inputs, err = getInputsByType(v.(map[string]interface{}), inputTypesMap)
				if err != nil {
					return err
				}
			}
		} else if blueprintId != "" {
			blueprintVersion := ""
			if v, ok := d.GetOk("blueprint_version"); ok {
				blueprintVersion = v.(string)
			}

			if v, ok := d.GetOk("inputs"); ok {
				// If the inputs are provided, get the schema from blueprint to convert the provided input values
				// to the type defined in the schema.
				inputsSchemaMap, err := getBlueprintSchema(apiClient, blueprintId, blueprintVersion)
				if err != nil {
					return err
				}

				log.Printf("Inputs Schema: %v", inputsSchemaMap)
				inputTypesMap, err := getInputTypesMapFromBlueprintInputsSchema(inputsSchemaMap)
				if err != nil {
					return err
				}

				log.Printf("InputTypesMap: %v", inputTypesMap)
				inputs, err = getInputsByType(v.(map[string]interface{}), inputTypesMap)
				if err != nil {
					return err
				}
			}
		}

		reason := "Updated deployment inputs from vRA provider for Terraform."
		err := runAction(d, apiClient, deploymentUuid, actionId, inputs, reason)
		if err != nil {
			return err
		}

		d.SetPartial("inputs")
		log.Printf("Finished updating vra_deployment %s with inputs", name)
	}

	return nil
}

func getInputsByType(inputs map[string]interface{}, inputTypesMap map[string]string) (map[string]interface{}, error) {
	log.Printf("Converting the input values to their types")
	inputsByType := make(map[string]interface{})
	var err error
	for k, v := range inputs {
		if t, ok := inputTypesMap[k]; ok {
			log.Printf("input_key: %s, type: %#v", k, t)
			switch strings.ToLower(t) {
			case "boolean":
				inputsByType[k], err = strconv.ParseBool(v.(string))
			case "integer":
				inputsByType[k], err = strconv.Atoi(v.(string))
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

// Returns a map of string, string with input variables and their types defined in the vRA catalog item
func getInputTypesMapFromCatalogItemSchema(schema map[string]interface{}) (map[string]string, error) {
	log.Printf("Getting the map of inputs and their types")
	inputTypesMap := make(map[string]string, len(schema))
	for k, v := range schema {
		inputTypesMap[k] = (v.(map[string]interface{}))["type"].(string)
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

func runAction(d *schema.ResourceData, apiClient *client.MulticloudIaaS, deploymentUuid strfmt.UUID, actionId string, inputs map[string]interface{}, reason string) error {
	resourceActionRequest := models.ResourceActionRequest{
		ActionID: actionId,
		Reason:   reason,
		Inputs:   inputs,
	}

	resp, err := apiClient.DeploymentActions.SubmitDeploymentActionRequestUsingPOST(
		deployment_actions.NewSubmitDeploymentActionRequestUsingPOSTParams().WithDepID(deploymentUuid).
			WithActionRequest(&resourceActionRequest))
	if err != nil {
		return err
	}

	requestId := resp.GetPayload().ID

	stateChangeFunc := resource.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.DeploymentRequestStatusPENDING, models.DeploymentRequestStatusINPROGRESS},
		Refresh:    deploymentActionStatusRefreshFunc(*apiClient, deploymentUuid, requestId),
		Target:     []string{models.DeploymentRequestStatusSUCCESSFUL},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
	}
	_, err = stateChangeFunc.WaitForState()
	if err != nil {
		return err
	}
	return nil
}

func deploymentActionStatusRefreshFunc(apiClient client.MulticloudIaaS, deploymentUuid strfmt.UUID, requestId strfmt.UUID) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ret, err := apiClient.Deployments.GetDeploymentByIDUsingGET(
			deployments.NewGetDeploymentByIDUsingGETParams().WithDepID(deploymentUuid).
				WithExpandLastRequest(withBool(true)))
		if err != nil {
			return "", models.DeploymentRequestStatusFAILED, err
		}

		status := ret.Payload.LastRequest.Status
		switch status {
		case models.DeploymentRequestStatusPENDING, models.DeploymentRequestStatusINPROGRESS:
			return [...]string{deploymentUuid.String()}, status, nil
		case models.DeploymentRequestStatusFAILED:
			return []string{""}, status, fmt.Errorf(ret.Error())
		case models.DeploymentRequestStatusSUCCESSFUL:
			deploymentID := ret.Payload.ID
			return deploymentID.String(), status, nil
		default:
			return [...]string{deploymentUuid.String()}, ret.Error(), fmt.Errorf("deploymentActionStatusRefreshFunc: unknown status %v", status)
		}
	}
}

func resourceDeploymentDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("Starting to delete the vra_deployment resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	_, err := apiClient.Deployments.DeleteDeploymentUsingDELETE(deployments.NewDeleteDeploymentUsingDELETEParams().WithDepID(strfmt.UUID(id)))
	if err != nil {
		return err
	}

	log.Printf("Requested for deleting the vra_deployment resource with name %s", d.Get("name"))

	stateChangeFunc := resource.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{reflect.TypeOf((*deployments.GetDeploymentByIDUsingGETOK)(nil)).String()},
		Refresh:    deploymentDeleteStatusRefreshFunc(*apiClient, d.Id()),
		Target:     []string{reflect.TypeOf((*deployments.GetDeploymentByIDUsingGETNotFound)(nil)).String()},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 5 * time.Second,
	}

	_, err = stateChangeFunc.WaitForState()
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("Finished deleting the vra_deployment resource with name %s", d.Get("name"))
	return nil
}

func deploymentDeleteStatusRefreshFunc(apiClient client.MulticloudIaaS, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ret, err := apiClient.Deployments.GetDeploymentByIDUsingGET(
			deployments.NewGetDeploymentByIDUsingGETParams().WithDepID(strfmt.UUID(id)))
		if err != nil {
			switch err.(type) {
			case *deployments.GetDeploymentByIDUsingGETNotFound:
				return "", reflect.TypeOf(err).String(), nil
			default:
				return [...]string{id}, reflect.TypeOf(err).String(), fmt.Errorf(ret.Error())
			}
		}

		return [...]string{id}, reflect.TypeOf(ret).String(), nil
	}
}
