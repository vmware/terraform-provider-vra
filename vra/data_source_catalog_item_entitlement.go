package vra

import (
        "fmt"
        "log"

        "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
        "github.com/vmware/vra-sdk-go/pkg/client/catalog_entitlements"
        "github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceCatalogItemEntitlement() *schema.Resource {
        return &schema.Resource{
                Read: dataSourceCatalogItemEntitlementRead,
                Importer: &schema.ResourceImporter{
                        State: schema.ImportStatePassthrough,
                },

                Schema: map[string]*schema.Schema{
                        "catalog_item_id": {
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

func dataSourceCatalogItemEntitlementRead(d *schema.ResourceData, m interface{}) error {
        log.Printf("Reading the vra_catalog_item_entitlement data source")
        apiClient := m.(*Client).apiClient

        id, idOk := d.GetOk("id")
        catalogItemID, catalogItemIDOk := d.GetOk("catalog_item_id")
        projectID := d.Get("project_id").(string)

        if !idOk && !catalogItemIDOk {
                return fmt.Errorf("one of id or catalog_item_id must be provided with project_id")
        }

        resp, err := apiClient.CatalogEntitlements.GetEntitlementsUsingGET(
                catalog_entitlements.NewGetEntitlementsUsingGETParams().WithProjectID(withString(projectID)))

        if err != nil {
                return err
        }

        setFields := func(entitlement *models.Entitlement) {
                d.SetId(entitlement.ID.String())
                d.Set("project_id", entitlement.ProjectID)
                d.Set("catalog_item_id", entitlement.Definition.ID)
                d.Set("definition", flattenContentDefinition(entitlement.Definition))
        }

        if len(resp.Payload) > 0 {
                for _, entitlement := range resp.Payload {
                        if idOk && entitlement.ID.String() == id.(string) {
                                setFields(entitlement)
                                log.Printf("Finished reading the vra_catalog_item_entitlement data item")
                                return nil
                        }

                        if catalogItemIDOk && entitlement.Definition.ID.String() == catalogItemID.(string) {
                                setFields(entitlement)
                                log.Printf("Finished reading the vra_catalog_item_entitlement data item")
                                return nil
                        }
                }
        }

        return fmt.Errorf("no catalog item entitlements found for the project_id '%v'", projectID)

}
