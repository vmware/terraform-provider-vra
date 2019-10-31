package vra

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

// resourceReferenceSchema returns the schema to use for the catalog item type property
func resourceReferenceSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"link": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"name": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func flattenResourceReferences(resourceReferences []*models.ResourceReference) []map[string]interface{} {
	if len(resourceReferences) == 0 {
		return make([]map[string]interface{}, 0)
	}

	resourceRefs := make([]map[string]interface{}, 0, len(resourceReferences))

	for _, resourceRef := range resourceReferences {
		helper := make(map[string]interface{})
		helper["id"] = resourceRef.ID
		helper["link"] = resourceRef.Link
		helper["name"] = resourceRef.Name

		resourceRefs = append(resourceRefs, helper)
	}

	return resourceRefs
}

func flattenResourceReference(resourceReference *models.ResourceReference) []interface{} {
	if resourceReference == nil {
		return make([]interface{}, 0)
	}
	helper := make(map[string]interface{})
	helper["id"] = resourceReference.ID
	helper["link"] = resourceReference.Link
	helper["name"] = resourceReference.Name

	return []interface{}{helper}
}
