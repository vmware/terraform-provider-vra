// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"encoding/json"
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/catalog_items"
)

func dataSourceCatalogItem() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCatalogItemRead,

		Schema: map[string]*schema.Schema{
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date-time when the entity was created.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user the entity was created by.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Catalog item description.",
			},
			"expand_projects": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Flag to indicate whether to expand detailed project data for the catalog item.",
			},
			"expand_versions": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Flag to indicate whether to expand detailed versions of the catalog item.",
			},
			"form_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Form ID.",
			},
			"icon_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Icon ID.",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The id of catalog item. One of `id` or `name` must be provided.",
			},
			"last_updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date-time when the entity was last updated.",
			},
			"last_updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user the entity was last updated by.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the catalog item. One of `id` or `name` must be provided.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The id of the project to narrow the search while looking for catalog items.",
			},
			"project_ids": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of associated project IDs that can be used for requesting this catalog item.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"projects": resourceReferenceSchema(),
			"schema": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Json schema describing request parameters.",
			},
			"source_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "LibraryItem source ID.",
			},
			"source_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "LibraryItem source name.",
			},
			"type":     resourceReferenceSchema(),
			"versions": catalogItemVersionSchema(),
		},
	}
}

func dataSourceCatalogItemRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")
	expandProjects := d.Get("expand_projects").(bool)

	if !idOk && !nameOk {
		return fmt.Errorf("one of id or name is required")
	}

	if !idOk {
		getCatalogItemsUsingGET5Params := catalog_items.NewGetCatalogItemsUsingGET5Params().
			WithSearch(withString(name.(string))).
			WithExpandProjects(withBool(expandProjects))
		if projectID, projectIDOk := d.GetOk("project_id"); projectIDOk {
			getCatalogItemsUsingGET5Params = getCatalogItemsUsingGET5Params.WithProjects([]string{projectID.(string)})
		}
		getResp, err := apiClient.CatalogItems.GetCatalogItemsUsingGET5(getCatalogItemsUsingGET5Params)
		if err != nil {
			return err
		}
		catalogItems := getResp.GetPayload()
		for _, catalogItem := range catalogItems.Content {
			if *catalogItem.Name == name {
				id = catalogItem.ID.String()
				break
			}
		}
		if id == "" {
			return fmt.Errorf("catalog item '%s' not found", name)
		}
	}

	getResp, err := apiClient.CatalogItems.GetCatalogItemUsingGET5(catalog_items.NewGetCatalogItemUsingGET5Params().WithID(strfmt.UUID(id.(string))).WithExpandProjects(withBool(expandProjects)))
	if err != nil {
		switch err.(type) {
		case *catalog_items.GetCatalogItemUsingGET5NotFound:
			return fmt.Errorf("catalog item '%s' not found", id)
		default:
			// nop
		}
		return err
	}

	catalogItem := getResp.GetPayload()

	d.SetId(catalogItem.ID.String())
	d.Set("created_at", catalogItem.CreatedAt)
	d.Set("created_by", catalogItem.CreatedBy)
	d.Set("description", catalogItem.Description)
	d.Set("form_id", catalogItem.FormID)
	d.Set("icon_id", catalogItem.IconID)
	d.Set("last_updated_at", catalogItem.LastUpdatedAt)
	d.Set("last_updated_by", catalogItem.LastUpdatedBy)
	d.Set("name", catalogItem.Name)
	d.Set("project_ids", catalogItem.ProjectIds)
	d.Set("projects", flattenResourceReferences(catalogItem.Projects))
	schemaJSON, _ := json.Marshal(catalogItem.Schema)
	d.Set("schema", string(schemaJSON))
	d.Set("source_id", catalogItem.SourceID.String())
	d.Set("source_name", catalogItem.SourceName)
	d.Set("type", flattenResourceReference(catalogItem.Type))

	if d.Get("expand_versions").(bool) {
		getVersionsResp, err := apiClient.CatalogItems.GetVersionsUsingGET2(catalog_items.NewGetVersionsUsingGET2Params().WithID(*catalogItem.ID))
		if err != nil {
			return err
		}

		d.Set("versions", flattenCatalogItemVersions(getVersionsResp.GetPayload().Content))
	}

	return nil
}
