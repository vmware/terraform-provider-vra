package vra

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

// deploymentRequest returns the schema to use for the last_request property
func deploymentRequestSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"action_id": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"approved_at": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"blueprint_id": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"cancelable": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
				"catalog_item_id": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"completed_at": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"completed_tasks": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"created_at": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"deployment_id": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"details": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"dismissed": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
				"id": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"initialized_at": {
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
				"name": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"parent_id": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"requested_by": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"resource_name": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"resource_type": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"status": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"total_tasks": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"updated_at": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func flattenDeploymentRequest(deploymentRequest *models.DeploymentRequest) interface{} {
	if deploymentRequest == nil {
		return make([]interface{}, 0)
	}

	helper := make(map[string]interface{})

	helper["action_id"] = deploymentRequest.ActionID
	helper["approved_at"] = deploymentRequest.ApprovedAt.String()
	helper["blueprint_id"] = deploymentRequest.BlueprintID
	helper["cancelable"] = deploymentRequest.Cancelable
	helper["catalog_item_id"] = deploymentRequest.CatalogItemID
	helper["completed_at"] = deploymentRequest.CompletedAt.String()
	helper["completed_tasks"] = deploymentRequest.CompletedTasks
	helper["created_at"] = deploymentRequest.CreatedAt.String()
	helper["deployment_id"] = deploymentRequest.DeploymentID.String()
	helper["details"] = deploymentRequest.Details
	helper["dismissed"] = deploymentRequest.Dismissed
	helper["id"] = deploymentRequest.ID.String()
	helper["initialized_at"] = deploymentRequest.InitializedAt.String()
	helper["inputs"] = expandInputs(deploymentRequest.Inputs)
	helper["name"] = deploymentRequest.Name
	helper["parent_id"] = deploymentRequest.ParentID.String()
	helper["requested_by"] = deploymentRequest.RequestedBy
	helper["resource_name"] = deploymentRequest.ResourceName
	helper["resource_type"] = deploymentRequest.ResourceType
	helper["status"] = deploymentRequest.Status
	helper["total_tasks"] = deploymentRequest.TotalTasks
	helper["updated_at"] = deploymentRequest.UpdatedAt.String()

	return []interface{}{helper}
}
