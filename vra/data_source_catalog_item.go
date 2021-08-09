package vra

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/catalog_items"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceCatalogItem() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCatalogItemRead,

		Schema: map[string]*schema.Schema{
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
				Computed: true,
			},
			"expand_projects": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"expand_versions": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"last_updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"project_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"projects": resourceReferenceSchema(),
			"schema": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type":     resourceReferenceSchema(),
			"versions": catalogItemVersionSchema(),
		},
	}
}

func dataSourceCatalogItemRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if !idOk && !nameOk {
		return fmt.Errorf("one of id or name must be assigned")
	}

	expandProjects := d.Get("expand_projects").(bool)

	getItemsResp, err := apiClient.CatalogItems.GetCatalogItemsUsingGET1(
		catalog_items.NewGetCatalogItemsUsingGET1Params().
			WithSearch(withString(name.(string))).
			WithExpandProjects(withBool(expandProjects)))
	if err != nil {
		return err
	}
	fmt.Printf(id.(string))
	fmt.Print(getItemsResp.GetPayload())

	setFields := func(catalogItem *models.CatalogItem, versions []*models.CatalogItemVersion) {
		d.SetId(catalogItem.ID.String())
		d.Set("created_at", catalogItem.CreatedAt)
		d.Set("created_by", catalogItem.CreatedBy)
		d.Set("description", catalogItem.Description)
		d.Set("last_updated_at", catalogItem.LastUpdatedAt)
		d.Set("last_updated_by", catalogItem.LastUpdatedBy)
		d.Set("name", catalogItem.Name)
		d.Set("project_ids", catalogItem.ProjectIds)
		d.Set("projects", flattenResourceReferences(catalogItem.Projects))
		d.Set("source_id", catalogItem.SourceID.String())
		d.Set("source_name", catalogItem.SourceName)
		d.Set("type", flattenResourceReference(catalogItem.Type))
		d.Set("versions", flattenCatalogItemVersions(versions))

		schemaJSON, _ := json.Marshal(catalogItem.Schema)
		d.Set("schema", string(schemaJSON))
	}

	for _, catalogItem := range getItemsResp.Payload.Content {
		if (idOk && catalogItem.ID.String() == id) || (nameOk && *catalogItem.Name == name.(string)) {
			getItemResp, err := apiClient.CatalogItems.GetCatalogItemUsingGET1(catalog_items.NewGetCatalogItemUsingGET1Params().WithID(*catalogItem.ID).WithExpandProjects(withBool(expandProjects)))

			if err != nil {
				return err
			}

			if d.Get("expand_versions").(bool) {
				getVersionsResp, err := apiClient.CatalogItems.GetVersionsUsingGET(catalog_items.NewGetVersionsUsingGETParams().WithID(*catalogItem.ID))

				if err != nil {
					return err
				}

				setFields(getItemResp.Payload, getVersionsResp.Payload.Content)
			} else {
				setFields(getItemResp.Payload, nil)
			}

			return nil
		}
	}

	return fmt.Errorf("catalog item %s not found", name)
}
