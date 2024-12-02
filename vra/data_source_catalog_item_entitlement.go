// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/catalog_entitlements"
)

func dataSourceCatalogItemEntitlement() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "'vra_catalog_item_entitlement' is deprecated. Use 'vra_content_sharing_policy' instead.",
		Read:               dataSourceCatalogItemEntitlementRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"catalog_item_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Catalog item id.",
				ConflictsWith: []string{"id"},
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
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Entitlement id.",
				ConflictsWith: []string{"catalog_item_id"},
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project id.",
			},
		},
	}
}

func dataSourceCatalogItemEntitlementRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	id, idOk := d.GetOk("id")
	catalogItemID, catalogItemIDOk := d.GetOk("catalog_item_id")
	projectID := d.Get("project_id").(string)

	if !idOk && !catalogItemIDOk {
		return fmt.Errorf("one of id or catalog_item_id must be provided with project_id")
	}

	resp, err := apiClient.CatalogEntitlements.GetEntitlementsUsingGET2(
		catalog_entitlements.NewGetEntitlementsUsingGET2Params().WithProjectID(withString(projectID)))
	if err != nil {
		return err
	}

	if len(resp.Payload) > 0 {
		for _, entitlement := range resp.Payload {
			if (idOk && entitlement.ID.String() == id.(string)) || (catalogItemIDOk && entitlement.Definition.ID.String() == catalogItemID.(string)) {
				d.SetId(entitlement.ID.String())
				d.Set("catalog_item_id", entitlement.Definition.ID)
				d.Set("definition", flattenContentDefinition(entitlement.Definition))
				d.Set("project_id", entitlement.ProjectID)
				return nil
			}
		}
	}

	return fmt.Errorf("no catalog item entitlements found for the project_id '%s'", projectID)
}
