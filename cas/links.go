package cas

import (
	"github.com/vmware/cas-sdk-go/pkg/models"
	"github.com/vmware/terraform-provider-cas/sdk"

	"github.com/hashicorp/terraform/helper/schema"
)

// linksSchema returns the schema to use for the links property
func linksSchema() *schema.Schema {
	return &schema.Schema{
		// List is used instead of Map because of: https://github.com/hashicorp/terraform/issues/621
		Type:     schema.TypeList,
		MinItems: 1,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"rel": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"href": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
				"hrefs": &schema.Schema{
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func flattenLinks(links map[string]tango.HypertextReference) []map[string]interface{} {
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

func getSelfLink(configLinks []interface{}) string {
	for _, configLink := range configLinks {
		linkMap := configLink.(map[string]interface{})

		if v, ok := linkMap["rel"].(string); ok && v == "self" {
			return linkMap["href"].(string)
		}
	}

	return ""
}

// linksSDKSchema returns the schema to use for the links property
func linksSDKSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		MinItems: 1,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"rel": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"href": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
				"hrefs": &schema.Schema{
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func flattenSDKLinks(links map[string]models.Href) []map[string]interface{} {
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
