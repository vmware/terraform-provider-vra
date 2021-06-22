package vra

import (
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/deployments"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"log"
)

func dataSourceDeployment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDeploymentRead,

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
				Computed: true,
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
				Optional: true,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"inputs": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"last_request": deploymentRequestSchema(),
			"last_updated_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"last_updated_by": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"lease_expire_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project": resourceReferenceSchema(),
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resources": resourcesSchema(),
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceDeploymentRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("Reading the vra_deployment resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if !idOk && !nameOk {
		return fmt.Errorf("one of id or name must be assigned")
	}

	expandLastRequest := d.Get("expand_last_request").(bool)
	expandProject := d.Get("expand_project").(bool)
	expandResources := d.Get("expand_resources").(bool)

	setFields := func(deployment *models.Deployment) error {
		d.SetId(deployment.ID.String())
		d.Set("blueprint_id", deployment.BlueprintID)
		d.Set("blueprint_version", deployment.BlueprintVersion)
		d.Set("catalog_item_id", deployment.CatalogItemID)
		d.Set("catalog_item_version", deployment.CatalogItemVersion)
		d.Set("created_at", deployment.CreatedAt)
		d.Set("created_by", deployment.CreatedBy)
		d.Set("description", deployment.Description)
		d.Set("last_updated_at", deployment.LastUpdatedAt)
		d.Set("last_updated_by", deployment.LastUpdatedBy)
		d.Set("lease_expire_at", deployment.LeaseExpireAt)
		d.Set("name", deployment.Name)
		d.Set("org_id", deployment.OrgID)
		d.Set("owner", deployment.OwnedBy)
		d.Set("project_id", deployment.ProjectID)
		d.Set("status", deployment.Status)

		if err := d.Set("expense", flattenExpense(deployment.Expense)); err != nil {
			return fmt.Errorf("error setting deployment expense - error: %#v", err)
		}

		if err := d.Set("inputs", expandInputs(deployment.Inputs)); err != nil {
			return fmt.Errorf("error setting deployment inputs - error: %#v", err)
		}

		if err := d.Set("last_request", flattenDeploymentRequest(deployment.LastRequest)); err != nil {
			return fmt.Errorf("error setting deployment last_request - error: %#v", err)
		}

		if err := d.Set("project", flattenResourceReference(deployment.Project)); err != nil {
			return fmt.Errorf("error setting project in deployment - error: %#v", err)
		}

		if err := d.Set("resources", flattenResources(deployment.Resources)); err != nil {
			return fmt.Errorf("error setting resources in deployment - error: %#v", err)
		}
		return nil
	}

	if nameOk {
		getAllResp, err := apiClient.Deployments.GetDeploymentsUsingGET(
			deployments.NewGetDeploymentsUsingGETParams().WithName(withString(name.(string))))

		if err != nil {
			return err
		}

		if getAllResp.Payload.NumberOfElements == 1 {
			deployment := getAllResp.Payload.Content[0]
			id = deployment.ID.String()
		} else {
			return fmt.Errorf("deployment %s not found", name)
		}
	}

	// Get the deployment details with all the user provided flags
	getResp, err := apiClient.Deployments.GetDeploymentByIDUsingGET(
		deployments.NewGetDeploymentByIDUsingGETParams().
			WithDeploymentID(strfmt.UUID(id.(string))).
			WithExpandProject(withBool(expandProject)).
			WithExpandResources(withBool(expandResources)).
			WithExpandLastRequest(withBool(expandLastRequest)).
			WithAPIVersion(withString(DeploymentsAPIVersion)).
			WithTimeout(IncreasedTimeOut))

	if err != nil {
		return err
	}

	deployment := getResp.Payload
	return setFields(deployment)
}
