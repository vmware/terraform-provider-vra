// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

// tagsSchema returns the schema to use for the tags property
func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:     schema.TypeString,
					Required: true,
				},
				"value": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

func expandTags(configTags []interface{}) []*models.Tag {
	tags := make([]*models.Tag, 0, len(configTags))

	for _, configTag := range configTags {
		tagMap := configTag.(map[string]interface{})

		tag := models.Tag{
			Key:   withString(tagMap["key"].(string)),
			Value: withString(tagMap["value"].(string)),
		}

		tags = append(tags, &tag)
	}

	return tags
}

func flattenTags(tags []*models.Tag) []interface{} {
	if len(tags) == 0 {
		return make([]interface{}, 0)
	}

	configTags := make([]interface{}, 0, len(tags))

	for _, tag := range tags {
		helper := make(map[string]interface{})
		helper["key"] = tag.Key
		helper["value"] = tag.Value

		configTags = append(configTags, helper)
	}

	return configTags
}
