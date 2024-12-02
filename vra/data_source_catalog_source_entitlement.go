// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/catalog_entitlements"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceCatalogSourceEntitlement() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "'vra_catalog_source_entitlement' is deprecated. Use 'vra_content_sharing_policy' instead.",
		Read:               dataSourceCatalogSourceEntitlementRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"catalog_source_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Catalog source id.",
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
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Entitlement id.",
				ConflictsWith: []string{"catalog_source_id"},
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project id.",
			},
		},
	}
}

func dataSourceCatalogSourceEntitlementRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("Reading the vra_catalog_source_entitlement data source")
	apiClient := m.(*Client).apiClient

	id, idOk := d.GetOk("id")
	catalogSourceID, catalogSourceIDOk := d.GetOk("catalog_source_id")
	projectID := d.Get("project_id").(string)

	if !idOk && !catalogSourceIDOk {
		return fmt.Errorf("one of id or catalog_source_id must be provided with project_id")
	}

	resp, err := apiClient.CatalogEntitlements.GetEntitlementsUsingGET2(
		catalog_entitlements.NewGetEntitlementsUsingGET2Params().WithProjectID(withString(projectID)))

	if err != nil {
		return err
	}

	setFields := func(entitlement *models.Entitlement) {
		d.SetId(entitlement.ID.String())
		d.Set("project_id", entitlement.ProjectID)
		d.Set("catalog_source_id", entitlement.Definition.ID)
		d.Set("definition", flattenContentDefinition(entitlement.Definition))
	}

	if len(resp.Payload) > 0 {
		for _, entitlement := range resp.Payload {
			if idOk && entitlement.ID.String() == id.(string) {
				setFields(entitlement)
				log.Printf("Finished reading the vra_catalog_source_entitlement data source")
				return nil
			}

			if catalogSourceIDOk && entitlement.Definition.ID.String() == catalogSourceID.(string) {
				setFields(entitlement)
				log.Printf("Finished reading the vra_catalog_source_entitlement data source")
				return nil
			}
		}
	}

	return fmt.Errorf("no catalog source entitlements found for the project_id '%v'", projectID)

}
