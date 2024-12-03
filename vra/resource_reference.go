// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

// resourceReferenceSchema returns the schema to use for the catalog item type property
func resourceReferenceSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"description": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"id": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"version": {
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
		helper["description"] = resourceRef.Description
		helper["id"] = resourceRef.ID
		helper["name"] = resourceRef.Name
		helper["version"] = resourceRef.Version

		resourceRefs = append(resourceRefs, helper)
	}

	return resourceRefs
}

func flattenResourceReference(resourceReference *models.ResourceReference) []interface{} {
	if resourceReference == nil {
		return make([]interface{}, 0)
	}
	helper := make(map[string]interface{})
	helper["description"] = resourceReference.Description
	helper["id"] = resourceReference.ID
	helper["name"] = resourceReference.Name
	helper["version"] = resourceReference.Version

	return []interface{}{helper}
}
