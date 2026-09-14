// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/deployment_actions"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceDeploymentResourceActions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDeploymentResourceActionsRead,

		Schema: map[string]*schema.Schema{
			"actions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The day 2 actions that are available on the deployment resource.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the action.",
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The display name of the action.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the action, to be used as the action_id of a vra_deployment_resource_action resource.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the action.",
						},
						"schema_json": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The JSON encoded schema of the inputs of the action.",
						},
						"valid": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the action can be run against the resource in its current state.",
						},
					},
				},
			},
			"deployment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the deployment that owns the resource.",
			},
			"resource_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"resource_id", "resource_name"},
				Description:  "The id of the deployment resource. Conflicts with resource_name.",
			},
			"resource_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"resource_id", "resource_name"},
				Description:  "The name of the deployment resource. The name must be unique within the deployment. Conflicts with resource_id.",
			},
			"valid_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to return only the actions that can be run against the resource in its current state.",
			},
		},
	}
}

func dataSourceDeploymentResourceActionsRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	deploymentUUID := strfmt.UUID(d.Get("deployment_id").(string))
	resourceUUID, err := resolveDeploymentResourceID(apiClient, deploymentUUID, d)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := apiClient.DeploymentActions.GetResourceActionsUsingGET4(
		deployment_actions.NewGetResourceActionsUsingGET4Params().
			WithAPIVersion(withString(DeploymentsAPIVersion)).
			WithDeploymentID(deploymentUUID).
			WithResourceID(resourceUUID))
	if err != nil {
		return diag.FromErr(err)
	}
	if resp == nil || resp.GetPayload() == nil {
		return diag.Errorf("vRA returned an empty response when listing actions of resource %s", resourceUUID)
	}

	validOnly := d.Get("valid_only").(bool)
	actions := make([]map[string]interface{}, 0, len(resp.GetPayload()))
	for _, action := range resp.GetPayload() {
		if action == nil {
			return diag.Errorf("vRA returned an empty action for resource %s", resourceUUID)
		}
		if validOnly && !action.Valid {
			continue
		}
		actions = append(actions, flattenResourceAction(action))
	}

	d.SetId(deploymentUUID.String() + "/" + resourceUUID.String())
	if err := d.Set("actions", actions); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("resource_id", resourceUUID.String()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func flattenResourceAction(action *models.ResourceAction) map[string]interface{} {
	schemaJSON := ""
	if action.Schema != nil {
		encoded, err := json.Marshal(action.Schema)
		if err != nil {
			log.Printf("[WARN] unable to encode the schema of the action %s: %v", action.ID, err)
		} else {
			schemaJSON = string(encoded)
		}
	}

	return map[string]interface{}{
		"description":  action.Description,
		"display_name": action.DisplayName,
		"id":           action.ID,
		"name":         action.Name,
		"schema_json":  schemaJSON,
		"valid":        action.Valid,
	}
}
