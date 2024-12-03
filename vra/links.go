// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// linksSchema returns the schema to use for the links property
func linksSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"rel": {
					Type:     schema.TypeString,
					Required: true,
				},
				"href": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"hrefs": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func flattenLinks(links map[string]models.Href) []map[string]interface{} {
	if len(links) == 0 {
		return make([]map[string]interface{}, 0)
	}

	configLinks := make([]map[string]interface{}, 0, len(links))

	for key, value := range links {
		helper := make(map[string]interface{})

		helper["rel"] = key
		helper["href"] = value.Href
		helper["hrefs"] = value.Hrefs

		configLinks = append(configLinks, helper)
	}

	return configLinks
}

/*
func getSelfLink(configLinks []interface{}) string {
	for _, configLink := range configLinks {
		linkMap := configLink.(map[string]interface{})

		if v, ok := linkMap["rel"].(string); ok && v == "self" {
			return linkMap["href"].(string)
		}
	}

	return ""
}
*/
