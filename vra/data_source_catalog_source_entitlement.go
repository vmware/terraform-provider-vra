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
		Read: dataSourceCatalogSourceEntitlementRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"catalog_source_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"definition": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"number_of_items": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"source_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
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

	resp, err := apiClient.CatalogEntitlements.GetEntitlementsUsingGET(
		catalog_entitlements.NewGetEntitlementsUsingGETParams().WithProjectID(withString(projectID)))

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
