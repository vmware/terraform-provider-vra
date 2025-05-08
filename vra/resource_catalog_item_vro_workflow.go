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

func resourceCatalogItemVroWorkflow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCatalogItemVroWorkflowCreate,
		ReadContext:   resourceCatalogItemVroWorkflowRead,
		UpdateContext: resourceCatalogItemVroWorkflowUpdate,
		DeleteContext: resourceCatalogItemVroWorkflowDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required arguments
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
			"workflow_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the vRO workflow to publish.",
			},

			// Optional arguments
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

func resourceCatalogItemVroWorkflowCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	spec := &CatalogItemVroWorkflowPublishSpec{
		WorkflowID: d.Get("workflow_id").(string),
	}

	createResp, err := apiClient.CatalogItems.PublishCatalogItem(catalog_items.NewPublishCatalogItemParams().WithRequest(&models.CatalogItemPublishRequest{
		Description: d.Get("description").(string),
		Global:      d.Get("global").(bool),
		IconID:      strfmt.UUID(d.Get("icon_id").(string)),
		Name:        d.Get("name").(string),
		ProjectID:   strfmt.UUID(d.Get("project_id").(string)),
		Spec:        spec,
		TypeID:      CatalogItemVroWorkflowTypeID,
	}))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResp.GetPayload().ID.String())

	return resourceCatalogItemVroWorkflowRead(ctx, d, m)
}

func resourceCatalogItemVroWorkflowRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	if catalogItem.Type.ID != CatalogItemVroWorkflowTypeID {
		return diag.Errorf("catalog item with id `%s` is not a vRO Workflow", id)
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

	var spec CatalogItemVroWorkflowPublishSpec
	if err := catalogItemSpecConvert(catalogItem.Spec, &spec); err != nil {
		return diag.FromErr(err)
	}

	d.Set("workflow_id", spec.WorkflowID)

	return nil
}

func resourceCatalogItemVroWorkflowUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	spec := &CatalogItemVroWorkflowPublishSpec{
		WorkflowID: d.Get("workflow_id").(string),
	}

	_, err := apiClient.CatalogItems.RepublishCatalogItem(catalog_items.NewRepublishCatalogItemParams().WithCatalogItemID(strfmt.UUID(id)).WithRequest(&models.CatalogItemPublishRequest{
		Description: d.Get("description").(string),
		Global:      d.Get("global").(bool),
		IconID:      strfmt.UUID(d.Get("icon_id").(string)),
		Name:        d.Get("name").(string),
		ProjectID:   strfmt.UUID(d.Get("project_id").(string)),
		Spec:        spec,
		TypeID:      CatalogItemVroWorkflowTypeID,
	}))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCatalogItemVroWorkflowRead(ctx, d, m)
}

func resourceCatalogItemVroWorkflowDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	if _, err := apiClient.CatalogItems.UnpublishCatalogItem(catalog_items.NewUnpublishCatalogItemParams().WithCatalogItemID(strfmt.UUID(id))); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
