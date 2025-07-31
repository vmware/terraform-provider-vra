// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"

	"github.com/vmware/vra-sdk-go/pkg/client/catalog_sources"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func resourceCatalogSourceBlueprint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCatalogSourceBlueprintCreate,
		DeleteContext: resourceCatalogSourceBlueprintDelete,
		ReadContext:   resourceCatalogSourceBlueprintRead,
		UpdateContext: resourceCatalogSourceBlueprintUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required arguments
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the blueprint content source instance.",
				Required:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "The id of the project the blueprint content source instance belongs to.",
				Required:    true,
			},

			// Optional arguments
			"description": {
				Type:        schema.TypeString,
				Description: "A human-friendly description for the blueprint content source instance.",
				Optional:    true,
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

func resourceCatalogSourceBlueprintCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	config := make(map[string]interface{})
	config["sourceProjectId"] = d.Get("project_id").(string)

	catalogSource := models.CatalogSource{
		Config:      config,
		Description: d.Get("description").(string),
		Name:        withString(d.Get("name").(string)),
		TypeID:      withString("com.vmw.blueprint"),
	}

	_, createResp, err := apiClient.CatalogSources.PostUsingPOST2(catalog_sources.NewPostUsingPOST2Params().WithSource(&catalogSource))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResp.GetPayload().ID.String())

	return resourceCatalogSourceBlueprintRead(ctx, d, m)
}

func resourceCatalogSourceBlueprintRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	resp, err := apiClient.CatalogSources.GetUsingGET2(catalog_sources.NewGetUsingGET2Params().WithSourceID(strfmt.UUID(id)))
	if err != nil {
		switch err.(type) {
		case *catalog_sources.GetUsingGET2NotFound:
			return diag.Errorf("blueprint catalog source with id `%s` is not found", d.Id())
		}
		return diag.FromErr(err)
	}

	catalogSource := *resp.Payload
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

	return nil
}

func resourceCatalogSourceBlueprintUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := strfmt.UUID(d.Id())
	config := make(map[string]interface{})
	config["sourceProjectId"] = d.Get("project_id").(string)

	catalogSource := models.CatalogSource{
		Config:      config,
		Description: d.Get("description").(string),
		ID:          &id,
		Name:        withString(d.Get("name").(string)),
		TypeID:      withString("com.vmw.blueprint"),
	}

	_, updateResp, err := apiClient.CatalogSources.PostUsingPOST2(catalog_sources.NewPostUsingPOST2Params().WithSource(&catalogSource))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(updateResp.GetPayload().ID.String())

	return resourceCatalogSourceBlueprintRead(ctx, d, m)
}

func resourceCatalogSourceBlueprintDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	if _, err := apiClient.CatalogSources.DeleteUsingDELETE4(catalog_sources.NewDeleteUsingDELETE4Params().WithSourceID(strfmt.UUID(id))); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
