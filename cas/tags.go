package cas

import (
	"github.com/vmware/terraform-provider-cas/sdk"

	"github.com/hashicorp/terraform/helper/schema"
)

// tagsSchema returns the schema to use for the tags property
func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
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

func expandTags(configTags []interface{}) []tango.Tag {
	tags := make([]tango.Tag, 0, len(configTags))

	for _, configTag := range configTags {
		tagMap := configTag.(map[string]interface{})

		tag := tango.Tag{
			Key:   tagMap["key"].(string),
			Value: tagMap["value"].(string),
		}

		tags = append(tags, tag)
	}

	return tags
}

func flattenTags(tags []tango.Tag) []interface{} {
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
