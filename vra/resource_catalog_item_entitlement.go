// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/catalog_entitlements"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func resourceCatalogItemEntitlement() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "'vra_catalog_item_entitlement' is deprecated. Use 'vra_content_sharing_policy' instead.",
		CreateContext:      resourceCatalogItemEntitlementCreate,
		DeleteContext:      resourceCatalogItemEntitlementDelete,
		ReadContext:        resourceCatalogItemEntitlementRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"catalog_item_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Catalog item id.",
			},
			"definition": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the catalog item.",
						},
						"icon_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Icon id of associated catalog item.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the catalog item.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the catalog item.",
						},
						"number_of_items": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of items in the associated catalog source.",
						},
						"source_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Catalog source name.",
						},
						"source_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Catalog source type.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Content definition type.",
						},
					},
				},
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Project id.",
			},
		},
	}
}

func resourceCatalogItemEntitlementCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	catalogItemID := strfmt.UUID(d.Get("catalog_item_id").(string))
	contentDefinition := models.ContentDefinition{
		ID:   &catalogItemID,
		Type: withString("CatalogItemIdentifier"),
	}
	entitlement := models.Entitlement{
		Definition: &contentDefinition,
		ProjectID:  withString(d.Get("project_id").(string)),
	}

	_, createResp, err := apiClient.CatalogEntitlements.CreateEntitlementUsingPOST2(
		catalog_entitlements.NewCreateEntitlementUsingPOST2Params().WithEntitlement(&entitlement))
	if err != nil {
		return diag.FromErr(err)
	}

	id := createResp.GetPayload().ID.String()
	d.SetId(id)

	return resourceCatalogItemEntitlementRead(ctx, d, m)
}

func resourceCatalogItemEntitlementRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	resp, err := apiClient.CatalogEntitlements.GetEntitlementsUsingGET2(
		catalog_entitlements.NewGetEntitlementsUsingGET2Params().WithProjectID(withString(d.Get("project_id").(string))))
	if err != nil {
		return diag.FromErr(err)
	}

	if len(resp.Payload) > 0 {
		for _, entitlement := range resp.Payload {
			if entitlement.ID.String() == id {
				d.Set("project_id", entitlement.ProjectID)
				d.Set("definition", flattenContentDefinition(entitlement.Definition))
				return nil
			}
		}
	}

	d.SetId("")
	return nil
}

func resourceCatalogItemEntitlementDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	_, err := apiClient.CatalogEntitlements.DeleteEntitlementUsingDELETE2(
		catalog_entitlements.NewDeleteEntitlementUsingDELETE2Params().WithID(strfmt.UUID(d.Id())))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
