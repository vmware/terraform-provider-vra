package vra

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

// resourceReferenceSchema returns the schema to use for the catalog item type property
func catalogItemVersionSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"created_at": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"description": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"id": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func flattenCatalogItemVersions(catalogItemVersions []*models.CatalogItemVersion) []map[string]interface{} {
	if len(catalogItemVersions) == 0 {
		return make([]map[string]interface{}, 0)
	}

	versions := make([]map[string]interface{}, 0, len(catalogItemVersions))

	for _, version := range catalogItemVersions {
		helper := make(map[string]interface{})
		helper["created_at"] = version.CreatedAt.String()
		helper["description"] = version.Description
		helper["id"] = version.ID

		versions = append(versions, helper)
	}

	return versions
}
