package vra

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

// resourcesSchema returns the schema to use for the resource property
func resourcesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"created_at": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
				"depends_on": &schema.Schema{
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"description": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
				"expense": expenseSchema(),
				// TODO:  Add metadata
				"id": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"name": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"properties_json": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
				"state": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
				"sync_status": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
				"type": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func flattenResources(resources []*models.Resource) []map[string]interface{} {
	if len(resources) == 0 {
		return make([]map[string]interface{}, 0)
	}

	configResources := make([]map[string]interface{}, 0, len(resources))

	for _, value := range resources {
		helper := make(map[string]interface{})

		helper["created_at"] = value.CreatedAt.String()
		helper["depends_on"] = value.DependsOn
		helper["description"] = value.Description
		helper["id"] = value.ID
		helper["name"] = value.Name
		helper["state"] = value.State
		helper["sync_status"] = value.SyncStatus
		helper["type"] = value.Type
		helper["expense"] = flattenExpense(value.Expense)

		propertiesSlice, _ := json.Marshal(value.Properties)
		helper["properties_json"] = string(propertiesSlice)

		configResources = append(configResources, helper)
	}

	return configResources
}

//func expandResources(configResources []interface{}) []*models.Resource {
//	resources := make([]*models.Resource, 0, len(configResources))
//
//	for _, configResource := range configResources {
//		resourceMap := configResource.(map[string]interface{})
//
//		resource := models.Resource{
//			ID: strfmt.UUID(resourceMap["id"].(string)),
//		}
//
//		if v, ok := resourceMap["created_at"].(string); ok && v != "" {
//			resource.CreatedAt, _ = strfmt.ParseDateTime(v)
//		}
//
//		if v, ok := resourceMap["depends_on"].([]interface{}); ok && len(v) != 0 {
//			dependsOn := make([]string, 0)
//
//			for _, value := range v {
//				dependsOn = append(dependsOn, value.(string))
//			}
//
//			resource.DependsOn = dependsOn
//		}
//
//		if v, ok := resourceMap["description"].(string); ok && v != "" {
//			resource.Description = v
//		}
//
//		if v, ok := resourceMap["name"].(string); ok && v != "" {
//			resource.Name = withString(v)
//		}
//
//		resource.Properties = expandCustomProperties(resourceMap["properties"].(map[string]interface{}))
//
//		if v, ok := resourceMap["state"].(string); ok && v != "" {
//			resource.State = v
//		}
//
//		if v, ok := resourceMap["sync_status"].(string); ok && v != "" {
//			resource.SyncStatus = v
//		}
//
//		if v, ok := resourceMap["type"].(string); ok && v != "" {
//			resource.Type = &v
//		}
//
//		resources = append(resources, &resource)
//	}
//
//	return resources
//}
