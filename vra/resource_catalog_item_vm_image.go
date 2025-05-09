// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/catalog_items"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func resourceCatalogItemVMImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCatalogItemVMImageCreate,
		ReadContext:   resourceCatalogItemVMImageRead,
		UpdateContext: resourceCatalogItemVMImageUpdate,
		DeleteContext: resourceCatalogItemVMImageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required arguments
			"image_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the VM image to publish.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the catalog item.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project to share this catalog item with.",
			},

			// Optional arguments
			"cloud_config": {
				Type:        schema.TypeString,
				Description: "Cloud config script to be applied to VMs provisioned from this image.",
				Optional:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "A human-friendly description for the catalog item.",
				Optional:    true,
			},
			"global": {
				Type:        schema.TypeBool,
				Description: "Whether to allow this catalog to be shared with multiple projects or to restrict it to the specified project.",
				Default:     false,
				Optional:    true,
			},
			"icon_id": {
				Type:        schema.TypeString,
				Description: "ID of the icon to associate with this catalog item.",
				Optional:    true,
			},
			"select_zone": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to create a zone input for the published catalog item.",
				Optional:    true,
			},

			// Computed attributes
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
		},
	}
}

func resourceCatalogItemVMImageCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	spec := &CatalogItemVMImagePublishSpec{
		ImageName: d.Get("image_name").(string),
	}
	if cloudConfig, ok := d.GetOk("cloud_config"); ok {
		spec.CloudConfig = withString(cloudConfig.(string))
	}
	if selectZone, ok := d.GetOk("select_zone"); ok {
		spec.SelectZone = withBool(selectZone.(bool))
	}

	createResp, err := apiClient.CatalogItems.PublishCatalogItem(catalog_items.NewPublishCatalogItemParams().WithRequest(&models.CatalogItemPublishRequest{
		Description: d.Get("description").(string),
		Global:      d.Get("global").(bool),
		IconID:      strfmt.UUID(d.Get("icon_id").(string)),
		Name:        d.Get("name").(string),
		ProjectID:   strfmt.UUID(d.Get("project_id").(string)),
		Spec:        spec,
		TypeID:      CatalogItemVMImageTypeID,
	}))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResp.GetPayload().ID.String())

	return resourceCatalogItemVMImageRead(ctx, d, m)
}

func resourceCatalogItemVMImageRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	getResp, err := apiClient.CatalogItems.GetCatalogItemUsingGET5(catalog_items.NewGetCatalogItemUsingGET5Params().WithID(strfmt.UUID(id)).WithExpand(withString("spec")))
	if err != nil {
		switch err.(type) {
		case *catalog_items.GetCatalogItemUsingGET5NotFound:
			return diag.Errorf("catalog item with id `%s` not found", id)
		default:
			// nop
		}
		return diag.FromErr(err)
	}

	catalogItem := getResp.GetPayload()

	if catalogItem.Type.ID != CatalogItemVMImageTypeID {
		return diag.Errorf("catalog item with id `%s` is not a VM Image", id)
	}

	d.SetId(catalogItem.ID.String())
	d.Set("created_at", catalogItem.CreatedAt.String())
	d.Set("created_by", catalogItem.CreatedBy)
	d.Set("description", catalogItem.Description)
	d.Set("icon_id", catalogItem.IconID)
	d.Set("global", catalogItem.Global)
	d.Set("last_updated_at", catalogItem.LastUpdatedAt.String())
	d.Set("last_updated_by", catalogItem.LastUpdatedBy)
	d.Set("name", catalogItem.Name)
	d.Set("project_id", catalogItem.SourceProjectID)

	var spec CatalogItemVMImagePublishSpec
	if err := catalogItemSpecConvert(catalogItem.Spec, &spec); err != nil {
		return diag.FromErr(err)
	}

	if spec.CloudConfig != nil {
		d.Set("cloud_config", spec.CloudConfig)
	}
	d.Set("image_name", spec.ImageName)
	if spec.SelectZone != nil {
		d.Set("select_zone", spec.SelectZone)
	}

	return nil
}

func resourceCatalogItemVMImageUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	spec := &CatalogItemVMImagePublishSpec{
		ImageName: d.Get("image_name").(string),
	}
	if cloudConfig, ok := d.GetOk("cloud_config"); ok {
		spec.CloudConfig = withString(cloudConfig.(string))
	}
	if selectZone, ok := d.GetOk("select_zone"); ok {
		spec.SelectZone = withBool(selectZone.(bool))
	}

	_, err := apiClient.CatalogItems.RepublishCatalogItem(catalog_items.NewRepublishCatalogItemParams().WithCatalogItemID(strfmt.UUID(id)).WithRequest(&models.CatalogItemPublishRequest{
		Description: d.Get("description").(string),
		Global:      d.Get("global").(bool),
		IconID:      strfmt.UUID(d.Get("icon_id").(string)),
		Name:        d.Get("name").(string),
		ProjectID:   strfmt.UUID(d.Get("project_id").(string)),
		Spec:        spec,
		TypeID:      CatalogItemVMImageTypeID,
	}))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCatalogItemVMImageRead(ctx, d, m)
}

func resourceCatalogItemVMImageDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	if _, err := apiClient.CatalogItems.UnpublishCatalogItem(catalog_items.NewUnpublishCatalogItemParams().WithCatalogItemID(strfmt.UUID(id))); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
