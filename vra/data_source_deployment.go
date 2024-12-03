// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/deployments"
)

func dataSourceDeployment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDeploymentRead,

		Schema: map[string]*schema.Schema{
			"blueprint_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the cloud template used to request the deployment.",
			},
			"blueprint_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the cloud template used to request the deployment.",
			},
			"catalog_item_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the catalog item used to request the deployment.",
			},
			"catalog_item_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the catalog item used to request the deployment.",
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
				Computed:    true,
				Description: "A human-friendly description.",
			},
			"expand_last_request": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Flag to indicate whether to expand last request on the deployment.",
			},
			"expand_project": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Flag to indicate whether to expand project information.",
			},
			"expand_resources": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Flag to indicate whether to expand resources in the deployment.",
			},
			"expense": expenseSchema(),
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
				Description:   "The id of the deployment. One of `id` or `name` must be provided.",
			},
			"inputs": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Inputs provided by the user while requesting / updating the deployment.",
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
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"id"},
				Description:   "The name of the deployment. One of `id` or `name` must be provided.",
			},
			"org_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Id of the organization this deployment belongs to.",
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user this deployment belongs to.",
			},
			"project": resourceReferenceSchema(),
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the project this deployment belongs to.",
			},
			"resources": resourcesSchema(),
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the deployment with respect to its life cycle operations.",
			},
		},
	}
}

func dataSourceDeploymentRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if !idOk && !nameOk {
		return diag.Errorf("one of id or name is required")
	}

	if nameOk {
		getAllResp, err := apiClient.Deployments.GetDeploymentsV3UsingGET(
			deployments.NewGetDeploymentsV3UsingGETParams().WithName(withString(name.(string))))

		if err != nil {
			return diag.FromErr(err)
		}

		if getAllResp.Payload.NumberOfElements == 1 {
			deployment := getAllResp.Payload.Content[0]
			id = deployment.ID.String()
		} else {
			return diag.Errorf("deployment %s not found", name)
		}
	}

	// Get the deployment details with all the user provided flags
	expand := []string{}
	expandProject := d.Get("expand_project").(bool)
	if expandProject {
		expand = append(expand, "project")
	}
	expandResources := d.Get("expand_resources").(bool)
	if expandResources {
		expand = append(expand, "resources")
	}
	expandLastRequest := d.Get("expand_last_request").(bool)
	if expandLastRequest {
		expand = append(expand, "lastRequest")
	}

	getResp, err := apiClient.Deployments.GetDeploymentByIDV3UsingGET(
		deployments.NewGetDeploymentByIDV3UsingGETParams().
			WithDeploymentID(strfmt.UUID(id.(string))).
			WithExpand(expand).
			WithAPIVersion(withString(DeploymentsAPIVersion)))

	if err != nil {
		return diag.FromErr(err)
	}

	deployment := getResp.Payload
	d.SetId(deployment.ID.String())
	d.Set("blueprint_id", deployment.BlueprintID)
	d.Set("blueprint_version", deployment.BlueprintVersion)
	d.Set("catalog_item_id", deployment.CatalogItemID)
	d.Set("catalog_item_version", deployment.CatalogItemVersion)
	d.Set("created_at", deployment.CreatedAt.String())
	d.Set("created_by", deployment.CreatedBy)
	d.Set("description", deployment.Description)
	d.Set("last_updated_at", deployment.LastUpdatedAt.String())
	d.Set("last_updated_by", deployment.LastUpdatedBy)
	d.Set("lease_expire_at", deployment.LeaseExpireAt.String())
	d.Set("name", deployment.Name)
	d.Set("org_id", deployment.OrgID)
	d.Set("owner", deployment.OwnedBy)
	d.Set("project_id", deployment.ProjectID)
	d.Set("status", deployment.Status)

	if err := d.Set("expense", flattenExpense(deployment.Expense)); err != nil {
		return diag.Errorf("error setting deployment expense - error: %#v", err)
	}

	if err := d.Set("inputs", expandInputs(deployment.Inputs)); err != nil {
		return diag.Errorf("error setting deployment inputs - error: %#v", err)
	}

	if err := d.Set("last_request", flattenDeploymentRequest(deployment.LastRequest)); err != nil {
		return diag.Errorf("error setting deployment last_request - error: %#v", err)
	}

	if err := d.Set("project", flattenResourceReference(deployment.Project)); err != nil {
		return diag.Errorf("error setting project in deployment - error: %#v", err)
	}

	if expandResources {
		getResourcesResp, err := apiClient.Deployments.GetDeploymentResourcesUsingGET2(
			deployments.NewGetDeploymentResourcesUsingGET2Params().
				WithDeploymentID(strfmt.UUID(id.(string))).
				WithExpand([]string{"currentRequest"}).
				WithAPIVersion(withString(DeploymentsAPIVersion)).
				WithDollarTop(withInt32(DefaultDollarTop)))
		if err != nil {
			return diag.Errorf("error retrieving deployment resources - error: %#v", err)
		}

		if err := d.Set("resources", flattenResources(getResourcesResp.GetPayload())); err != nil {
			return diag.Errorf("error setting resources in deployment - error: %#v", err)
		}
	}

	return nil
}
