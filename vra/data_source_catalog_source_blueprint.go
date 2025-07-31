// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"

	"github.com/vmware/vra-sdk-go/pkg/client/catalog_sources"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCatalogSourceBlueprint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCatalogSourceBlueprintRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Computed:      true,
				ConflictsWith: []string{"name", "project_id"},
				Description:   "The id of the blueprint content source instance.",
				Optional:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Computed:      true,
				ConflictsWith: []string{"id", "project_id"},
				Description:   "The name of the blueprint content source instance.",
				Optional:      true,
			},
			"project_id": {
				Type:          schema.TypeString,
				Computed:      true,
				ConflictsWith: []string{"id", "name"},
				Description:   "The id of the project the blueprint content source instance belongs to.",
				Optional:      true,
			},

			// Computed attributes
			"config": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The content source custom configuration.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
				Description: "A human-friendly description for the blueprint content source instance.",
			},
			"global": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Global flag indicating that all the items can be requested across all projects.",
			},
			"icon_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default Icon Identifier.",
			},
			"items_found": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of items found.",
			},
			"items_imported": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of items imported.",
			},
			"last_import_completed_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the last import was completed. The date is in ISO 8601 and UTC.",
			},
			"last_import_errors": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "A list of errors seen at last time the content source is imported.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"last_import_started_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the last import was started. The date is in ISO 8601 and UTC.",
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
			"type_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of this content source.",
			},
		},
	}
}

func dataSourceCatalogSourceBlueprintRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")
	projectID, projectIDOk := d.GetOk("project_id")

	if !idOk && !nameOk && !projectIDOk {
		return fmt.Errorf("one of `id` or `name` or `project_id` must be provided")
	}

	setFields := func(catalogSource *models.CatalogSource) {
		d.SetId(catalogSource.ID.String())

		d.Set("config", expandCatalogSourceConfig(catalogSource.Config))
		d.Set("created_at", catalogSource.CreatedAt.String())
		d.Set("created_by", catalogSource.CreatedBy)
		d.Set("description", catalogSource.Description)
		d.Set("global", catalogSource.Global)
		d.Set("icon_id", catalogSource.IconID)
		d.Set("items_found", catalogSource.ItemsFound)
		d.Set("items_imported", catalogSource.ItemsImported)
		d.Set("last_import_completed_at", catalogSource.LastImportCompletedAt.String())
		d.Set("last_import_errors", catalogSource.LastImportErrors)
		d.Set("last_import_started_at", catalogSource.LastImportStartedAt.String())
		d.Set("last_updated_at", catalogSource.LastUpdatedAt.String())
		d.Set("last_updated_by", catalogSource.LastUpdatedBy)
		d.Set("name", catalogSource.Name)
		d.Set("project_id", catalogSource.ProjectID)
		d.Set("type_id", catalogSource.TypeID)
	}

	// Get catalog source by id if id is provided
	if idOk {
		resp, err := apiClient.CatalogSources.GetUsingGET2(catalog_sources.NewGetUsingGET2Params().WithSourceID(strfmt.UUID(id.(string))))
		if err != nil {
			switch err.(type) {
			case *catalog_sources.GetUsingGET2NotFound:
				return fmt.Errorf("blueprint catalog source with id `%s` is not found", id)
			}
			return err
		}

		setFields(resp.Payload)
		return nil
	}

	// Search catalog sources if name is provided and look for exact match with type com.vmw.blueprint
	if nameOk {
		resp, err := apiClient.CatalogSources.GetPageUsingGET2(catalog_sources.NewGetPageUsingGET2Params().WithSearch(withString(name.(string))))
		if err != nil {
			return err
		}

		if resp.Payload.NumberOfElements > 0 {
			for _, catalogSource := range resp.Payload.Content {
				if *catalogSource.Name == name.(string) && *catalogSource.TypeID == "com.vmw.blueprint" {
					setFields(catalogSource)
					return nil
				}
			}
		}

		return fmt.Errorf("blueprint catalog source with name `%s` is not found", name)

	}

	// Filter catalog sources if projectId is provided and look for exact match with type com.vmw.blueprint and projectId as global catalog sources are returned as well
	if projectIDOk {
		resp, err := apiClient.CatalogSources.GetPageUsingGET2(catalog_sources.NewGetPageUsingGET2Params().WithProjectID(withString(projectID.(string))))
		if err != nil {
			return err
		}

		if resp.Payload.NumberOfElements > 0 {
			for _, catalogSource := range resp.Payload.Content {
				if catalogSource.ProjectID == projectID.(string) && *catalogSource.TypeID == "com.vmw.blueprint" {
					setFields(catalogSource)
					return nil
				}
			}
		}

		return fmt.Errorf("blueprint catalog source with project_id `%s` is not found", projectID)
	}

	return nil
}
