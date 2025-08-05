// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/vmware/vra-sdk-go/pkg/client/content_source"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceContentSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContentSourceRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Input attributes
			"id": {
				Type:          schema.TypeString,
				Computed:      true,
				ConflictsWith: []string{"name"},
				Description:   "The id of the content source instance.",
				Optional:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Computed:      true,
				ConflictsWith: []string{"id"},
				Description:   "The name of the content source instance.",
				Optional:      true,
			},

			// Computed attributes
			"config": {
				Type:        schema.TypeSet,
				Description: "The content source custom configuration.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"branch": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The content source branch name.",
						},
						"content_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The content source type.",
						},
						"integration_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The content source integration id as seen integrations.",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path to refer to in the content source repository and branch.",
						},
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the project.",
						},
						"repository": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The content source repository.",
						},
					},
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was created. The date is in ISO 8601 and UTC.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user the entity was created by.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly description for the catalog source instance.",
			},
			"last_updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC.",
			},
			"last_updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user the entity was last updated by.",
			},
			"org_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the organization this entity belongs to.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the project this entity belongs to.",
			},
			"sync_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Wether or not sync is enabled for this content source.",
			},
			"type_id": {
				Type:        schema.TypeString,
				Description: "The type of this content source.",
				Computed:    true,
			},
		},
	}
}

func dataSourceContentSourceRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")
	if !idOk && !nameOk {
		return diag.Errorf("one of id or name must be provided")
	}

	var contentSource *models.ContentSource
	if id != "" {
		getResp, err := apiClient.ContentSource.GetContentSourceUsingGET(content_source.NewGetContentSourceUsingGETParams().WithID(strfmt.UUID(id.(string))))
		if err != nil {
			switch err.(type) {
			case *content_source.GetContentSourceUsingGETNotFound:
				return diag.Errorf("content source with id `%s` not found", id)
			default:
				// nop
			}
			return diag.FromErr(err)
		}

		contentSource = getResp.GetPayload()
	} else {
		getResp, err := apiClient.ContentSource.ListContentSourcesUsingGET(content_source.NewListContentSourcesUsingGETParams().WithSearch(withString(name.(string))))
		if err != nil {
			return diag.FromErr(err)
		}

		contentSources := getResp.Payload
		if len(contentSources.Content) == 0 {
			return diag.Errorf("vra_content_source `name` criteria did not match any content source")
		}
		if len(contentSources.Content) > 1 {
			return diag.Errorf("vra_content_source `name` criteria must filter to a single content source")
		}

		contentSource = contentSources.Content[0]
	}

	d.SetId(contentSource.ID.String())
	_ = d.Set("created_at", contentSource.CreatedAt.String())
	_ = d.Set("created_by", contentSource.CreatedBy)
	_ = d.Set("description", contentSource.Description)
	_ = d.Set("name", contentSource.Name)
	_ = d.Set("last_updated_at", contentSource.LastUpdatedAt.String())
	_ = d.Set("last_updated_by", contentSource.LastUpdatedBy)
	_ = d.Set("org_id", contentSource.OrgID)
	_ = d.Set("project_id", contentSource.ProjectID)
	_ = d.Set("sync_enabled", contentSource.SyncEnabled)
	_ = d.Set("type_id", contentSource.TypeID)

	config, err := flattenContentsourceRepositoryConfig(contentSource.Config)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("config", config)

	return nil
}
