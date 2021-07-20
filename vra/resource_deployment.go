package vra

import (
	"encoding/json"
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

const (
	ChangeOwnerDeploymentActionName = "ChangeOwner"
	ChangeLeaseDeploymentActionName = "ChangeLease"
	EditTagsDeploymentActionName    = "EditTags"
	PowerOffDeploymentActionName    = "PowerOff"
	PowerOnDeploymentActionName     = "PowerOn"
	UpdateDeploymentActionName      = "update"
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
				ForceNew: true,
			},
			"catalog_item_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceDeploymentCreate(d *schema.ResourceData, m interface{}) error {
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
		return fmt.Errorf("only one of (blueprint_id, catalog_item_id) required")
	}

	deploymentName := d.Get("name").(string)
	projectID := d.Get("project_id").(string)

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
				return err
			}
		}
		catalogItemRequest.Inputs = inputs

		if v, ok := d.GetOk("description"); ok {
			catalogItemRequest.Reason = v.(string)
		}

		log.Printf("[DEBUG] Create deployment: %#v", catalogItemRequest)
		postOk, err := apiClient.CatalogItems.RequestCatalogItemUsingPOST(
			catalog_items.NewRequestCatalogItemUsingPOSTParams().WithID(strfmt.UUID(catalogItemID)).
				WithAPIVersion(withString(CatalogAPIVersion)).WithRequest(&catalogItemRequest))

		if err != nil {
			return err
		}

		d.SetId(postOk.GetPayload().DeploymentID)
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
			// If the inputs are provided, get the schema from blueprint to convert the provided input values
			// to the type defined in the schema.
			inputs, err = getBlueprintInputsByType(apiClient, blueprintID, blueprintVersion, v)
			if err != nil {
				return err
			}
		}
		blueprintRequest.Inputs = inputs

		bpRequestCreated, bpRequestAccepted, err := apiClient.BlueprintRequests.CreateBlueprintRequestUsingPOST1(
			blueprint_requests.NewCreateBlueprintRequestUsingPOST1Params().WithRequest(&blueprintRequest))

		if err != nil {
			log.Printf("Received error. err=%s, bpRequestCreated=%v, bpRequestAccepted=%v", err, bpRequestCreated, bpRequestAccepted)
			return err
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
			return fmt.Errorf("failed to request for a deployment. status: %v, message: %v", status, failureMessage)
		}

		log.Printf("Finished requesting vra_deployment '%s' from blueprint %s", d.Get("name"), blueprintID)
	}

	stateChangeFunc := resource.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.DeploymentStatusCREATEINPROGRESS, models.DeploymentStatusUPDATEINPROGRESS},
		Refresh:    deploymentStatusRefreshFunc(*apiClient, d.Id()),
		Target:     []string{models.DeploymentStatusCREATESUCCESSFUL, models.DeploymentStatusUPDATESUCCESSFUL},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
	}

	deploymentID, err := stateChangeFunc.WaitForState()
	if err != nil {
		readError := resourceDeploymentRead(d, m)
		if readError != nil {
			return fmt.Errorf("failed to create deployment: %v \nfailed to read deployment state: %v", err.Error(), readError.Error())
		}
		return err
	}

	d.SetId(deploymentID.(string))
	log.Printf("Finished to create vra_deployment resource with name %s", d.Get("name"))

	return resourceDeploymentRead(d, m)
}

func resourceDeploymentRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("Reading the vra_deployment resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	// Getting the input types map
	inputTypesMap := getInputTypesMap(d, apiClient)

	id := d.Id()
	expandProject := d.Get("expand_project").(bool)

	resp, err := apiClient.Deployments.GetDeploymentByIDUsingGET(
		deployments.NewGetDeploymentByIDUsingGETParams().
			WithDeploymentID(strfmt.UUID(id)).
			WithExpandResources(withBool(true)).
			WithExpandLastRequest(withBool(true)).
			WithExpandProject(withBool(expandProject)).
			WithAPIVersion(withString(DeploymentsAPIVersion)).
			WithTimeout(IncreasedTimeOut))
	if err != nil {
		switch err.(type) {
		case *deployments.GetDeploymentByIDUsingGETNotFound:
			d.SetId("")
			return nil
		}
		return err
	}

	deployment := *resp.Payload
	d.Set("blueprint_id", deployment.BlueprintID)
	d.Set("blueprint_version", deployment.BlueprintVersion)
	d.Set("catalog_item_id", deployment.CatalogItemID)
	d.Set("catalog_item_version", deployment.CatalogItemVersion)
	d.Set("created_at", deployment.CreatedAt)
	d.Set("created_by", deployment.CreatedBy)
	d.Set("description", deployment.Description)

	if err := d.Set("expense", flattenExpense(deployment.Expense)); err != nil {
		return fmt.Errorf("error setting deployment expense - error: %#v", err)
	}

	if err := d.Set("inputs_including_defaults", expandInputsToString(deployment.Inputs)); err != nil {
		return fmt.Errorf("error setting deployment inputs_including_defaults - error: %#v", err)
	}

	allInputs := expandInputs(deployment.Inputs)
	if v, ok := d.GetOk("inputs"); ok {
		userInputs := v.(map[string]interface{})
		if err := d.Set("inputs", updateUserInputs(allInputs, userInputs, inputTypesMap)); err != nil {
			return fmt.Errorf("error setting deployment inputs - error: %#v", err)
		}
	}

	if err := d.Set("last_request", flattenDeploymentRequest(deployment.LastRequest)); err != nil {
		return fmt.Errorf("error setting deployment last_request - error: %#v", err)
	}

	d.Set("last_updated_at", deployment.LastUpdatedAt)
	d.Set("last_updated_by", deployment.LastUpdatedBy)
	d.Set("lease_expire_at", deployment.LeaseExpireAt)
	d.Set("name", deployment.Name)
	d.Set("org_id", deployment.OrgID)
	d.Set("owner", deployment.OwnedBy)

	if err := d.Set("project", flattenResourceReference(deployment.Project)); err != nil {
		return fmt.Errorf("error setting project in deployment - error: %#v", err)
	}

	d.Set("project_id", deployment.ProjectID)

	if err := d.Set("resources", flattenResources(deployment.Resources)); err != nil {
		return fmt.Errorf("error setting resources in deployment - error: %#v", err)
	}

	d.Set("status", deployment.Status)

	log.Printf("Finished reading the vra_deployment resource with name '%s'. Current status: '%s'", d.Get("name"), d.Get("status"))
	return nil
}

func resourceDeploymentUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("Starting to update the vra_deployment resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	if d.HasChange("blueprint_id") || d.HasChange("blueprint_version") || d.HasChange("blueprint_content") {
		err := updateDeploymentWithNewBlueprint(d, m, apiClient)
		if err != nil {
			return err
		}
	} else {
		id := d.Id()
		deploymentUUID := strfmt.UUID(id)
		if d.HasChange("name") || d.HasChange("description") {
			err := updateDeploymentMetadata(d, apiClient, deploymentUUID)
			if err != nil {
				return err
			}
		}

		if d.HasChange("inputs") {
			err := runDeploymentUpdateAction(d, apiClient, deploymentUUID)
			if err != nil {
				return err
			}
		}

		stateChangeFunc := resource.StateChangeConf{
			Delay:      5 * time.Second,
			Pending:    []string{models.DeploymentStatusCREATEINPROGRESS, models.DeploymentStatusUPDATEINPROGRESS},
			Refresh:    deploymentStatusRefreshFunc(*apiClient, d.Id()),
			Target:     []string{models.DeploymentStatusCREATESUCCESSFUL, models.DeploymentStatusUPDATESUCCESSFUL},
			Timeout:    d.Timeout(schema.TimeoutCreate),
			MinTimeout: 5 * time.Second,
		}

		_, err := stateChangeFunc.WaitForState()
		if err != nil {
			readError := resourceDeploymentRead(d, m)
			if readError != nil {
				return fmt.Errorf("failed to update deployment: %v \nfailed to read deployment state: %v", err.Error(), readError.Error())
			}
			return err
		}
	}

	if d.HasChange("owner") {
		deploymentUUID := strfmt.UUID(d.Id())
		err := runChangeOwnerDeploymentAction(d, apiClient, deploymentUUID)
		if err != nil {
			return err
		}
	}

	log.Printf("Finished updating the vra_deployment resource with name %s", d.Get("name"))
	return resourceDeploymentRead(d, m)
}

func resourceDeploymentDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("Starting to delete the vra_deployment resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	_, err := apiClient.Deployments.DeleteDeploymentUsingDELETE(deployments.NewDeleteDeploymentUsingDELETEParams().WithDeploymentID(strfmt.UUID(id)))
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

// Gets the inputs and their types as map[string]string
func getInputTypesMap(d *schema.ResourceData, apiClient *client.MulticloudIaaS) map[string]string {
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

func getCatalogItemInputsByType(apiClient *client.MulticloudIaaS, catalogItemID string, catalogItemVersion string, inputValues interface{}) (map[string]interface{}, error) {
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

func getCatalogItemInputTypesMap(apiClient *client.MulticloudIaaS, catalogItemID string, catalogItemVersion string) (map[string]string, error) {
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

func getBlueprintInputsByType(apiClient *client.MulticloudIaaS, blueprintID string, blueprintVersion string, inputValues interface{}) (map[string]interface{}, error) {
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

func getBlueprintInputTypesMap(apiClient *client.MulticloudIaaS, blueprintID string, blueprintVersion string) (map[string]string, error) {
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

func getCatalogItemSchema(apiClient *client.MulticloudIaaS, catalogItemID string, catalogItemVersion string) (map[string]interface{}, error) {
	// Getting the catalog item schema
	log.Printf("Getting the schema for catalog item: %v version: %v", catalogItemID, catalogItemVersion)
	var catalogItemSchema interface{}
	if catalogItemVersion == "" {
		getItemResp, err := apiClient.CatalogItems.GetCatalogItemUsingGET1(catalog_items.NewGetCatalogItemUsingGET1Params().WithID(strfmt.UUID(catalogItemID)))
		if err != nil {
			return nil, err
		}
		catalogItemSchema = getItemResp.GetPayload().Schema
	} else {
		getVersionResp, err := apiClient.CatalogItems.GetVersionByIDUsingGET(catalog_items.NewGetVersionByIDUsingGETParams().WithID(strfmt.UUID(catalogItemID)).WithVersionID(catalogItemVersion))
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

func getBlueprintSchema(apiClient *client.MulticloudIaaS, blueprintID string, blueprintVersion string) (map[string]models.PropertyDefinition, error) {
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

func deploymentStatusRefreshFunc(apiClient client.MulticloudIaaS, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ret, err := apiClient.Deployments.GetDeploymentByIDUsingGET(
			deployments.NewGetDeploymentByIDUsingGETParams().
				WithDeploymentID(strfmt.UUID(id)).
				WithExpandLastRequest(withBool(true)).
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
			return ret.Payload.ID.String(), status, fmt.Errorf(ret.Payload.LastRequest.Details)
		default:
			return [...]string{id}, ret.Error(), fmt.Errorf("deploymentStatusRefreshFunc: unknown status %v", status)
		}
	}
}

func updateDeploymentWithNewBlueprint(d *schema.ResourceData, m interface{}, apiClient *client.MulticloudIaaS) error {
	log.Printf("Noticed changes to blueprint_id/blueprint_version/blueprint_content. Starting to update existing deployment...")

	blueprintID, blueprintVersion, blueprintContent := "", "", ""
	if v, ok := d.GetOk("blueprint_id"); ok {
		blueprintID = v.(string)
	}

	if v, ok := d.GetOk("blueprint_content"); ok {
		blueprintContent = v.(string)
	}

	if blueprintID != "" && blueprintContent != "" {
		return fmt.Errorf("only one of (blueprint_id, blueprintContent) required")
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
		// If the inputs are provided, get the schema from blueprint to convert the provided input values
		// to the type defined in the schema.
		inputs, err := getBlueprintInputsByType(apiClient, blueprintID, blueprintVersion, v)
		if err != nil {
			return err
		}
		blueprintRequest.Inputs = inputs
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
		return fmt.Errorf("failed to request update to existing deployment. status: %v, message: %v", status, failureMessage)
	}

	stateChangeFunc := resource.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.DeploymentStatusCREATEINPROGRESS, models.DeploymentStatusUPDATEINPROGRESS},
		Refresh:    deploymentStatusRefreshFunc(*apiClient, deploymentID),
		Target:     []string{models.DeploymentStatusCREATESUCCESSFUL, models.DeploymentStatusUPDATESUCCESSFUL},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
	}

	_, err = stateChangeFunc.WaitForState()
	if err != nil {
		readError := resourceDeploymentRead(d, m)
		if readError != nil {
			return fmt.Errorf("failed to update deployment: %v \nfailed to read deployment state: %v", err.Error(), readError.Error())
		}
		return err
	}

	log.Printf("Finished to update vra_deployment '%s' with blueprint '%s'", deploymentName, blueprintID)
	return nil
}

func updateDeploymentMetadata(d *schema.ResourceData, apiClient *client.MulticloudIaaS, deploymentUUID strfmt.UUID) error {
	log.Printf("Starting to update deployment name and description")
	description := d.Get("description").(string)
	name := d.Get("name").(string)

	updateDeploymentSpecification := models.DeploymentUpdate{
		Description: description,
		Name:        name,
	}

	log.Printf("[DEBUG] update deployment: %#v", updateDeploymentSpecification)
	_, err := apiClient.Deployments.PatchDeploymentUsingPATCH(
		deployments.NewPatchDeploymentUsingPATCHParams().WithDeploymentID(deploymentUUID).WithUpdate(&updateDeploymentSpecification))
	if err != nil {
		return err
	}

	log.Printf("Finished updating deployment name and description")
	return nil
}

func runDeploymentUpdateAction(d *schema.ResourceData, apiClient *client.MulticloudIaaS, deploymentUUID strfmt.UUID) error {
	log.Printf("Noticed changes to inputs. Starting to update deployment with inputs")
	// Get the deployment actions
	deploymentActions, err := apiClient.DeploymentActions.GetDeploymentActionsUsingGET(deployment_actions.
		NewGetDeploymentActionsUsingGETParams().WithDeploymentID(deploymentUUID))
	if err != nil {
		return err
	}

	updateAvailable := false
	actionID := ""
	for _, action := range deploymentActions.Payload {
		if strings.Contains(strings.ToLower(action.ID), strings.ToLower("Update")) {
			if action.Valid {
				log.Printf("[DEBUG] update action is available on the deployment")
				updateAvailable = true
				actionID = action.ID
				break
			} else {
				return fmt.Errorf("noticed changes to inputs, but 'Update' action is not supported based on the current state of the deployment")
			}
		}
	}

	name := d.Get("name")
	if !updateAvailable {
		return fmt.Errorf("noticed changes to inputs, but 'Update' action is not supported based on the current state of the deployment %s", name)
	}

	// Continue if update action is available.
	var inputs = make(map[string]interface{})
	blueprintID, catalogItemID := "", ""
	if v, ok := d.GetOk("catalog_item_id"); ok {
		catalogItemID = v.(string)
	}

	if v, ok := d.GetOk("blueprint_id"); ok {
		blueprintID = v.(string)
	}

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
				return err
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
				return err
			}
		}
	}

	reason := "Updated deployment inputs from vRA provider for Terraform."
	err = runAction(d, apiClient, deploymentUUID, actionID, inputs, reason)
	if err != nil {
		return err
	}

	log.Printf("Finished updating vra_deployment %s with inputs", name)

	return nil
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
//Used for getting map of inputs and their types for Catalog item and Deployment actions
func getInputTypesMapFromSchema(schema map[string]interface{}) (map[string]string, error) {
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

// Returns whether the day2 action is valid currently, exact action ID for a given action string
func getDeploymentDay2ActionID(apiClient *client.MulticloudIaaS, deploymentUUID strfmt.UUID, actionName string) (bool, string, error) {
	// Get the deployment actions
	deploymentActions, err := apiClient.DeploymentActions.GetDeploymentActionsUsingGET(deployment_actions.
		NewGetDeploymentActionsUsingGETParams().WithDeploymentID(deploymentUUID))
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

// Gets the schema for a given deployment action id
func getDeploymentActionSchema(apiClient *client.MulticloudIaaS, deploymentUUID strfmt.UUID, actionID string) (map[string]interface{}, error) {
	// Getting the catalog item schema
	log.Printf("Getting the schema for deploymentID: %v, actionID: %v", deploymentUUID, actionID)
	var actionSchema interface{}

	deploymentAction, err := apiClient.DeploymentActions.GetDeploymentActionUsingGET(deployment_actions.
		NewGetDeploymentActionUsingGETParams().WithDeploymentID(deploymentUUID).WithActionID(actionID))
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

func getDeploymentActionInputTypesMap(apiClient *client.MulticloudIaaS, deploymentUUID strfmt.UUID, actionID string) (map[string]string, error) {
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

func runChangeOwnerDeploymentAction(d *schema.ResourceData, apiClient *client.MulticloudIaaS, deploymentUUID strfmt.UUID) error {
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
	err = runAction(d, apiClient, deploymentUUID, actionID, inputs, reason)
	if err != nil {
		return err
	}

	log.Printf("Finished changing owner for vra_deployment %s with new owner %v", d.Get("name").(string), newOwner)
	return nil
}

func runAction(d *schema.ResourceData, apiClient *client.MulticloudIaaS, deploymentUUID strfmt.UUID, actionID string, inputs map[string]interface{}, reason string) error {
	resourceActionRequest := models.ResourceActionRequest{
		ActionID: actionID,
		Reason:   reason,
		Inputs:   inputs,
	}

	resp, err := apiClient.DeploymentActions.SubmitDeploymentActionRequestUsingPOST(
		deployment_actions.NewSubmitDeploymentActionRequestUsingPOSTParams().
			WithAPIVersion(withString(DeploymentsAPIVersion)).
			WithDeploymentID(deploymentUUID).
			WithActionRequest(&resourceActionRequest))
	if err != nil {
		return err
	}

	requestID := resp.GetPayload().ID

	stateChangeFunc := resource.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestStatusPENDING, models.RequestStatusINITIALIZATION, models.RequestStatusCHECKINGAPPROVAL, models.RequestStatusAPPROVALPENDING, models.RequestStatusINPROGRESS},
		Refresh:    deploymentActionStatusRefreshFunc(*apiClient, deploymentUUID, requestID),
		Target:     []string{models.RequestStatusCOMPLETION, models.RequestStatusAPPROVALREJECTED, models.RequestStatusABORTED, models.RequestStatusSUCCESSFUL, models.RequestStatusFAILED},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 5 * time.Second,
	}
	_, err = stateChangeFunc.WaitForState()
	if err != nil {
		return err
	}
	return nil
}

func deploymentActionStatusRefreshFunc(apiClient client.MulticloudIaaS, deploymentUUID strfmt.UUID, requestID strfmt.UUID) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ret, err := apiClient.Deployments.GetDeploymentByIDUsingGET(
			deployments.NewGetDeploymentByIDUsingGETParams().
				WithDeploymentID(deploymentUUID).
				WithExpandLastRequest(withBool(true)).
				WithAPIVersion(withString(DeploymentsAPIVersion)).
				WithTimeout(IncreasedTimeOut))
		if err != nil {
			return "", models.RequestStatusFAILED, err
		}

		status := ret.Payload.LastRequest.Status
		switch status {
		case models.RequestStatusPENDING, models.RequestStatusINITIALIZATION, models.RequestStatusCHECKINGAPPROVAL, models.RequestStatusAPPROVALPENDING, models.RequestStatusINPROGRESS, models.RequestStatusCOMPLETION:
			return [...]string{deploymentUUID.String()}, status, nil
		case models.RequestStatusAPPROVALREJECTED, models.RequestStatusABORTED, models.RequestStatusFAILED:
			return []string{""}, status, fmt.Errorf(ret.Error())
		case models.RequestStatusSUCCESSFUL:
			deploymentID := ret.Payload.ID
			return deploymentID.String(), status, nil
		default:
			return [...]string{deploymentUUID.String()}, ret.Error(), fmt.Errorf("deploymentActionStatusRefreshFunc: unknown status %v", status)
		}
	}
}

func deploymentDeleteStatusRefreshFunc(apiClient client.MulticloudIaaS, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ret, err := apiClient.Deployments.GetDeploymentByIDUsingGET(
			deployments.NewGetDeploymentByIDUsingGETParams().WithDeploymentID(strfmt.UUID(id)))
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

func updateUserInputs(allInputs, userInputs map[string]interface{}, inputTypesMap map[string]string) map[string]interface{} {
	if allInputs == nil || userInputs == nil {
		return nil
	}

	inputs := make(map[string]interface{})
	for name, value := range userInputs {
		if value != nil {
			inputs[name] = decodeInputValue(name, allInputs[name], inputTypesMap)
		}
		log.Printf("Converted incoming value to string: Key: %v, Value: %v, Converted value: %#v", name, allInputs[name], inputs[name])
	}

	return inputs
}

func decodeInputValue(inputName string, inputValue interface{}, inputTypesMap map[string]string) interface{} {
	if t, ok := inputTypesMap[inputName]; ok {
		log.Printf("input_key: %s, type: %#v, value: %#v", inputName, t, inputValue)
		switch strings.ToLower(t) {
		case "array":
			fallthrough
		case "object":
			log.Printf("Converting input '%v' of type '%v'", inputName, t)
			value, err := json.Marshal(inputValue)
			if err != nil {
				log.Printf("cannot convert input '%v' value '%v' into type '%v'. %#v", inputName, inputValue, t, err)
				return fmt.Sprint(inputValue)
			}
			return string(value)
		default:
			return fmt.Sprint(inputValue)
		}
	}

	return fmt.Sprint(inputValue)
}
