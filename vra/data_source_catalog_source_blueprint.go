// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"log"

	"github.com/vmware/vra-sdk-go/pkg/client/catalog_sources"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCatalogSourceBlueprint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCatalogSourceBlueprintRead,

		Schema: map[string]*schema.Schema{
			"config": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
			"global": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"items_found": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"items_imported": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_import_completed_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_import_errors": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"last_import_started_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCatalogSourceBlueprintRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("Looking up the vra_catalog_source_blueprint data resource with name")
	apiClient := m.(*Client).apiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")
	projectID, projectIDOk := d.GetOk("project_id")

	if !idOk && !nameOk && !projectIDOk {
		return fmt.Errorf("one of id or name or project_id must be provided")
	}

	setFields := func(catalogSource *models.CatalogSource) {
		d.SetId(catalogSource.ID.String())

		d.Set("config", expandCatalogSourceConfig(catalogSource.Config))
		d.Set("created_at", catalogSource.CreatedAt)
		d.Set("created_by", catalogSource.CreatedBy)
		d.Set("description", catalogSource.Description)
		d.Set("global", catalogSource.Global)
		d.Set("items_found", catalogSource.ItemsFound)
		d.Set("items_imported", catalogSource.ItemsImported)
		d.Set("last_import_completed_at", catalogSource.LastImportCompletedAt)
		d.Set("last_import_errors", catalogSource.LastImportErrors)
		d.Set("last_import_started_at", catalogSource.LastImportStartedAt)
		d.Set("last_updated_at", catalogSource.LastUpdatedAt)
		d.Set("last_updated_by", catalogSource.LastUpdatedBy)
		d.Set("name", catalogSource.Name)
		d.Set("project_id", catalogSource.ProjectID)
		d.Set("type_id", catalogSource.TypeID)
	}

	// Get catalog source by id if id is provided
	if idOk {
		resp, err := apiClient.CatalogSources.GetUsingGET2(
			catalog_sources.NewGetUsingGET2Params().WithSourceID(strfmt.UUID(id.(string))))

		if err != nil {
			switch err.(type) {
			case *catalog_sources.GetUsingGET2NotFound:
				return fmt.Errorf("blueprint catalog source with id '%v' is not found", id)
			}
			return err
		}

		setFields(resp.Payload)
		log.Printf("Finished reading the vra_catalog_source_blueprint resource with id  '%v'", id)
		return nil
	}

	// Search catalog sources if name is provided and look for exact match with type com.vmw.blueprint
	if nameOk {
		resp, err := apiClient.CatalogSources.GetPageUsingGET2(
			catalog_sources.NewGetPageUsingGET2Params().WithSearch(withString(name.(string))))

		if err != nil {
			return err
		}

		if resp.Payload.NumberOfElements > 0 {
			for _, catalogSource := range resp.Payload.Content {
				if *catalogSource.Name == name.(string) && *catalogSource.TypeID == "com.vmw.blueprint" {
					setFields(catalogSource)
					log.Printf("Finished reading the vra_catalog_source_blueprint resource with name  '%v'", name)
					return nil
				}
			}
		}

		return fmt.Errorf("blueprint catalog source with name '%v' is not found", name)
	}

	// Filter catalog sources if projectId is provided and look for exact match with type com.vmw.blueprint and projectId as global catalog sources are returned as well
	if projectIDOk {
		resp, err := apiClient.CatalogSources.GetPageUsingGET2(
			catalog_sources.NewGetPageUsingGET2Params().WithProjectID(withString(projectID.(string))))

		if err != nil {
			return err
		}

		if resp.Payload.NumberOfElements > 0 {
			for _, catalogSource := range resp.Payload.Content {
				if catalogSource.ProjectID == projectID.(string) && *catalogSource.TypeID == "com.vmw.blueprint" {
					setFields(catalogSource)
					log.Printf("Finished reading the vra_catalog_source_blueprint resource with project_id '%v'", projectID)
					return nil
				}
			}
		}

		return fmt.Errorf("blueprint catalog source with project_id '%v' is not found", projectID)
	}

	return nil
}
