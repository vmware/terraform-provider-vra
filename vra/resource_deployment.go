package vra

import (
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client"
	"github.com/vmware/vra-sdk-go/pkg/client/blueprint_requests"
	"github.com/vmware/vra-sdk-go/pkg/client/catalog_items"
	"github.com/vmware/vra-sdk-go/pkg/client/deployments"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"log"
	"reflect"
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
				Type:     schema.TypeBool,
				Optional: true,
			},
			"expand_project": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"expand_resources": {
				Type:     schema.TypeBool,
				Optional: true,
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
	log.Printf("starting to create vra_deployment resource")
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
			log.Printf("deployment '%v' doesn't exist already and hence can be created", deploymentName)
		}
	} else {
		return fmt.Errorf("a deployment with name '%v' exists already. Try with a differnet name", deploymentName)
	}

	// If catalog_item_id is provided, request deployment with the catalog item
	if catalogItemId != "" {
		log.Printf("requesting vra_deployment '%s' from catalog item", d.Get("name"))
		catalogItemRequest := models.CatalogItemRequest{
			DeploymentName: deploymentName,
			ProjectID:      projectId,
		}

		if v, ok := d.GetOk("inputs"); ok {
			catalogItemRequest.Inputs = v
		} else {
			catalogItemRequest.Inputs = make(map[string]interface{})
		}

		if v, ok := d.GetOk("description"); ok {
			catalogItemRequest.Reason = v.(string)
		}

		if v, ok := d.GetOk("catalog_item_version"); ok {
			catalogItemRequest.Version = v.(string)
		}

		log.Printf("[DEBUG] create deployment: %#v", catalogItemRequest)
		postOk, err := apiClient.CatalogItems.RequestCatalogItemUsingPOST(
			catalog_items.NewRequestCatalogItemUsingPOSTParams().WithID(strfmt.UUID(catalogItemId)).
				WithRequest(&catalogItemRequest))

		if err != nil {
			return err
		}

		d.SetId(postOk.GetPayload().DeploymentID)
		log.Printf("finished requesting vra_deployment '%s' from catalog item", d.Get("name"))
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
			log.Printf("received error. err=%s, bpRequestCreated=%v, bpRequestAccepted=%v", err, bpRequestCreated, bpRequestAccepted)
			return err
		}

		// blueprint_requests service may return either 201 or 202 depending on whether the request is in terminal state vs or in-progress
		log.Printf("requested deployment from blueprint. bpRequestCreated=%v, bpRequestAccepted=%v", bpRequestCreated, bpRequestAccepted)
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

		log.Printf("finished requesting vra_deployment '%s' from blueprint %s", d.Get("name"), blueprintId)
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
	log.Printf("finished to create vra_deployment resource with name %s", d.Get("name"))

	return resourceDeploymentRead(d, m)
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
	expandLastRequest := d.Get("expand_last_request").(bool)
	expandProject := d.Get("expand_project").(bool)
	expandResources := d.Get("expand_resources").(bool)

	resp, err := apiClient.Deployments.GetDeploymentByIDUsingGET(
		deployments.NewGetDeploymentByIDUsingGETParams().WithDepID(strfmt.UUID(id)).
			WithExpandResources(withBool(expandResources)).WithExpandLastRequest(withBool(expandLastRequest)).
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

	log.Printf("finished reading the vra_deployment resource with name %s", d.Get("name"))
	return nil
}

func resourceDeploymentUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("starting to update the vra_deployment resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	description := d.Get("description").(string)
	name := d.Get("name").(string)

	updateDeploymentSpecification := models.DeploymentUpdate{
		Description: description,
		Name:        name,
	}

	log.Printf("[DEBUG] update deployment: %#v", updateDeploymentSpecification)
	_, err := apiClient.Deployments.PatchDeploymentUsingPATCH(deployments.NewPatchDeploymentUsingPATCHParams().WithDepID(strfmt.UUID(id)).WithUpdate(&updateDeploymentSpecification))
	if err != nil {
		return err
	}

	log.Printf("finished updating the vra_deployment resource with name %s", d.Get("name"))
	return resourceDeploymentRead(d, m)
}

func resourceDeploymentDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("starting to delete the vra_deployment resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	_, err := apiClient.Deployments.DeleteDeploymentUsingDELETE(deployments.NewDeleteDeploymentUsingDELETEParams().WithDepID(strfmt.UUID(id)))
	if err != nil {
		return err
	}

	log.Printf("requested for deleting the vra_deployment resource with name %s", d.Get("name"))

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
	log.Printf("finished deleting the vra_deployment resource with name %s", d.Get("name"))
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
