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

	"log"
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
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCatalogSourceBlueprintCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("starting to create vra_catalog_source_blueprint resource")

	apiClient := m.(*Client).apiClient

	name := d.Get("name").(string)

	config := make(map[string]interface{})
	config["sourceProjectId"] = d.Get("project_id").(string)

	catalogSource := models.CatalogSource{
		Config: config,
		Name:   withString(name),
		TypeID: withString("com.vmw.blueprint"),
	}

	if v, ok := d.GetOk("description"); ok {
		catalogSource.Description = v.(string)
	}

	_, createResp, err := apiClient.CatalogSources.PostUsingPOST2(
		catalog_sources.NewPostUsingPOST2Params().WithSource(&catalogSource))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResp.GetPayload().ID.String())
	log.Printf("Finished creating vra_catalog_source_blueprint resource with name %s", d.Get("name"))

	return resourceCatalogSourceBlueprintRead(ctx, d, m)
}

func resourceCatalogSourceBlueprintRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Reading the vra_catalog_source_blueprint resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	resp, err := apiClient.CatalogSources.GetUsingGET2(
		catalog_sources.NewGetUsingGET2Params().WithSourceID(strfmt.UUID(d.Id())))

	if err != nil {
		switch err.(type) {
		case *catalog_sources.GetUsingGET2NotFound:
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	catalogSource := *resp.Payload
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

	log.Printf("Finished reading the vra_catalog_source_blueprint resource with name %s", d.Get("name"))
	return nil
}

func resourceCatalogSourceBlueprintUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("starting to update vra_catalog_source_blueprint resource")

	apiClient := m.(*Client).apiClient

	name := d.Get("name").(string)

	config := make(map[string]interface{})
	config["sourceProjectId"] = d.Get("project_id").(string)

	csID := strfmt.UUID(d.Id())

	catalogSource := models.CatalogSource{
		Config: config,
		ID:     &csID,
		Name:   withString(name),
		TypeID: withString("com.vmw.blueprint"),
	}

	if v, ok := d.GetOk("description"); ok {
		catalogSource.Description = v.(string)
	}

	_, createResp, err := apiClient.CatalogSources.PostUsingPOST2(
		catalog_sources.NewPostUsingPOST2Params().WithSource(&catalogSource))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResp.GetPayload().ID.String())
	log.Printf("Finished updating vra_catalog_source_blueprint resource with name %s", d.Get("name"))

	return resourceCatalogSourceBlueprintRead(ctx, d, m)
}

func resourceCatalogSourceBlueprintDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to delete the vra_catalog_source_blueprint resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	_, err := apiClient.CatalogSources.DeleteUsingDELETE4(
		catalog_sources.NewDeleteUsingDELETE4Params().WithSourceID(strfmt.UUID(d.Id())))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("Finished deleting the vra_catalog_source_blueprint resource with name %s", d.Get("name"))
	return nil
}
