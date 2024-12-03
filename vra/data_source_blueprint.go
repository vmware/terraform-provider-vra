// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"

	"github.com/go-openapi/strfmt"

	"github.com/vmware/vra-sdk-go/pkg/client/blueprint"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBlueprint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBlueprintRead,

		Schema: map[string]*schema.Schema{
			"content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud template YAML content.",
			},
			"content_source_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the content source.",
			},
			"content_source_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The path of the content source.",
			},
			"content_source_sync_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Content source last sync at.",
			},
			"content_source_sync_messages": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Content source last sync status.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"content_source_sync_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Content source last sync status.",
			},
			"content_source_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the content source.",
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
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
				Description:   "The id of the cloud template.",
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"id"},
				Description:   "The name of the cloud template.",
			},
			"org_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the organization this entity belongs to.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The id of the project to narrow the search while looking for cloud templates.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the project the entity belongs to.",
			},
			"request_scope_org": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag to indicate whether this cloud template can be requested from any project in the organization this entity belongs to.",
			},
			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "HATEOAS of the entity.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the cloud template.",
			},
			"total_released_versions": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of released versions.",
			},
			"total_versions": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of versions.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC.",
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user the entity was last updated by.",
			},
			"valid": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag to indicate if the current content of the cloud template is valid.",
			},
			"validation_messages": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Validation message.",
						},
						"metadata": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Validation metadata.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Validation path.",
						},
						"resource_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the resource.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Message type.",
						},
					},
				},
			},
		},
	}
}

func dataSourceBlueprintRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if !idOk && !nameOk {
		return fmt.Errorf("one of id or name is required")
	}

	var blueprintResource *models.Blueprint
	if !idOk {
		listBlueprintsUsingGET1Params := blueprint.NewListBlueprintsUsingGET1Params().WithName(withString(name.(string)))
		if projectID, projectIDOk := d.GetOk("project_id"); projectIDOk {
			listBlueprintsUsingGET1Params = listBlueprintsUsingGET1Params.WithProjects([]string{projectID.(string)})
		}
		getResp, err := apiClient.Blueprint.ListBlueprintsUsingGET1(listBlueprintsUsingGET1Params)
		if err != nil {
			return err
		}

		blueprints := getResp.GetPayload()
		if len(blueprints.Content) > 1 {
			return fmt.Errorf("more than one blueprint found with the same name, try to narrow filter by project_id")
		}
		if len(blueprints.Content) == 0 {
			return fmt.Errorf("blueprint %s not found", name)
		}

		// ListBlueprintsUsingGET1Params does not return the blueprint content, so we need to use the GetBlueprintUsingGET1
		// call in order to retrieve it.
		id = blueprints.Content[0].ID
	}

	getResp, err := apiClient.Blueprint.GetBlueprintUsingGET1(blueprint.NewGetBlueprintUsingGET1Params().WithBlueprintID(strfmt.UUID(id.(string))))
	if err != nil {
		switch err.(type) {
		case *blueprint.GetBlueprintUsingGET1NotFound:
			return fmt.Errorf("blueprint '%s' not found", id)
		default:
			// nop
		}
		return err
	}

	blueprintResource = getResp.GetPayload()

	d.SetId(blueprintResource.ID)
	d.Set("content", blueprintResource.Content)
	d.Set("content_source_id", blueprintResource.ContentSourceID)
	d.Set("content_source_pat", blueprintResource.ContentSourcePath)
	d.Set("content_source_sync_at", blueprintResource.ContentSourceSyncAt)
	d.Set("content_source_sync_messages", blueprintResource.ContentSourceSyncMessages)
	d.Set("content_source_sync_status", blueprintResource.ContentSourceSyncStatus)
	d.Set("content_source_type", blueprintResource.ContentSourceType)
	d.Set("created_at", blueprintResource.CreatedAt)
	d.Set("created_by", blueprintResource.CreatedBy)
	d.Set("description", blueprintResource.Description)
	d.Set("name", blueprintResource.Name)
	d.Set("org_id", blueprintResource.OrgID)
	d.Set("project_id", blueprintResource.ProjectID)
	d.Set("project_name", blueprintResource.ProjectName)
	d.Set("request_scope_org", blueprintResource.RequestScopeOrg)
	d.Set("self_link", blueprintResource.SelfLink)
	d.Set("status", blueprintResource.Status)
	d.Set("total_released_versions", blueprintResource.TotalReleasedVersions)
	d.Set("total_versions", blueprintResource.TotalVersions)
	d.Set("updated_at", blueprintResource.UpdatedAt)
	d.Set("updated_by", blueprintResource.UpdatedBy)
	d.Set("valid", blueprintResource.Valid)

	return nil
}
