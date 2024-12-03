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

	"log"
)

func resourceCatalogSourceEntitlement() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "'vra_catalog_source_entitlement' is deprecated. Use 'vra_content_sharing_policy' instead.",
		CreateContext:      resourceCatalogSourceEntitlementCreate,
		DeleteContext:      resourceCatalogSourceEntitlementDelete,
		ReadContext:        resourceCatalogSourceEntitlementRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"catalog_source_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Catalog source id.",
			},
			"definition": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the catalog source.",
						},
						"icon_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Icon id of associated catalog source.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the catalog source.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the catalog source.",
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

func resourceCatalogSourceEntitlementCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("starting to create vra_catalog_source_entitlement resource")

	apiClient := m.(*Client).apiClient

	catalogSourceID := strfmt.UUID(d.Get("catalog_source_id").(string))

	contentDefinition := models.ContentDefinition{
		ID:   &catalogSourceID,
		Type: withString("CatalogSourceIdentifier"),
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

	d.SetId(createResp.GetPayload().ID.String())
	log.Printf("Finished creating vra_catalog_source_entitlement resource with name %s", d.Get("name"))

	return resourceCatalogSourceEntitlementRead(ctx, d, m)
}

func resourceCatalogSourceEntitlementRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Reading the vra_catalog_source_entitlement resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	resp, err := apiClient.CatalogEntitlements.GetEntitlementsUsingGET2(
		catalog_entitlements.NewGetEntitlementsUsingGET2Params().WithProjectID(withString(d.Get("project_id").(string))))

	if err != nil {
		return diag.FromErr(err)
	}

	setFields := func(entitlement *models.Entitlement) {
		d.SetId(entitlement.ID.String())
		d.Set("project_id", entitlement.ProjectID)
		d.Set("definition", flattenContentDefinition(entitlement.Definition))
	}

	if len(resp.Payload) > 0 {
		for _, entitlement := range resp.Payload {
			if entitlement.Definition.ID.String() == d.Get("catalog_source_id").(string) {
				setFields(entitlement)
				log.Printf("Finished reading the vra_catalog_source_entitlement resource with name %s", d.Get("name"))
				return nil
			}
		}
	}

	d.SetId("")
	return nil
}

func resourceCatalogSourceEntitlementDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to delete the vra_catalog_source_entitlement resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	_, err := apiClient.CatalogEntitlements.DeleteEntitlementUsingDELETE2(
		catalog_entitlements.NewDeleteEntitlementUsingDELETE2Params().WithID(strfmt.UUID(d.Id())))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("Finished deleting the vra_catalog_source_entitlement resource with name %s", d.Get("name"))
	return nil
}
