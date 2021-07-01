package vra

import (
	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/catalog_entitlements"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"log"
)

func resourceCatalogSourceEntitlement() *schema.Resource {
	return &schema.Resource{
		Create: resourceCatalogSourceEntitlementCreate,
		Delete: resourceCatalogSourceEntitlementDelete,
		Read:   resourceCatalogSourceEntitlementRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"catalog_source_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCatalogSourceEntitlementCreate(d *schema.ResourceData, m interface{}) error {
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

	_, createResp, err := apiClient.CatalogEntitlements.CreateEntitlementUsingPOST(
		catalog_entitlements.NewCreateEntitlementUsingPOSTParams().WithEntitlement(&entitlement))

	if err != nil {
		return err
	}

	d.SetId(createResp.GetPayload().ID.String())
	log.Printf("Finished creating vra_catalog_source_entitlement resource with name %s", d.Get("name"))

	return resourceCatalogSourceEntitlementRead(d, m)
}

func resourceCatalogSourceEntitlementRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("Reading the vra_catalog_source_entitlement resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	resp, err := apiClient.CatalogEntitlements.GetEntitlementsUsingGET(
		catalog_entitlements.NewGetEntitlementsUsingGETParams().WithProjectID(withString(d.Get("project_id").(string))))

	if err != nil {
		return err
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

func resourceCatalogSourceEntitlementDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("Starting to delete the vra_catalog_source_entitlement resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	_, err := apiClient.CatalogEntitlements.DeleteEntitlementUsingDELETE(
		catalog_entitlements.NewDeleteEntitlementUsingDELETEParams().WithID(strfmt.UUID(d.Id())))

	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("Finished deleting the vra_catalog_source_entitlement resource with name %s", d.Get("name"))
	return nil
}
