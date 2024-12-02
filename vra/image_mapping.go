// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

// imageMappingSchema returns the schema to use for the image_mapping property
func imageMappingSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"cloud_config": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Cloud config for this image. This cloud config will be merged during provisioning with other cloud configurations such as the bootConfig provided in MachineSpecification or vRA cloud templates.",
				},
				"constraints": constraintsSchema(),
				"description": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "A human-friendly description.",
				},
				"external_id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "External entity id on the cloud provider side.",
				},
				"external_region_id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "External region id on the cloud provider side.",
				},
				"image_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The id of this resource instance.",
				},
				"image_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "A human-friendly image name as seen on the cloud provider side.",
				},
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "A human-friendly name of the image mapping.",
				},
				"organization": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "A human-friendly description.",
				},
				"os_family": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Operating system family of the image.",
				},
				"owner": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Email of the user that owns the entity.",
				},
				"private": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Indicates whether this fabric image is private. For vSphere, private images are considered to be templates and snapshots and public are Content Library items.",
				},
			},
		},
	}
}

func expandImageMapping(configImageMappings []interface{}) map[string]models.FabricImageDescription {
	images := make(map[string]models.FabricImageDescription)

	for _, configImageMapping := range configImageMappings {
		image := configImageMapping.(map[string]interface{})

		var i models.FabricImageDescription

		if v, found := image["cloud_config"]; found && v != nil {
			i.CloudConfig = v.(string)
		}

		if v, found := image["constraints"]; found && v != nil {
			i.Constraints = expandConstraints(v.(*schema.Set).List())
		}

		if v, found := image["image_id"]; found && v != nil {
			i.ID = v.(string)
		}

		if v, found := image["image_name"]; found && v != nil {
			i.Name = v.(string)
		}

		images[image["name"].(string)] = i
	}

	return images
}

func flattenImageMappings(list map[string]models.ImageMappingDescription) []interface{} {
	result := make([]interface{}, 0, len(list))
	for mappingName, mappingDescription := range list {
		l := map[string]interface{}{
			"cloud_config":       mappingDescription.CloudConfig,
			"constraints":        flattenImageMappingConstraints(mappingDescription.Constraints),
			"description":        mappingDescription.Description,
			"external_id":        mappingDescription.ExternalID,
			"external_region_id": mappingDescription.ExternalRegionID,
			//"image_id":           *mappingDescription.ID,
			"image_name":   mappingDescription.Name,
			"name":         mappingName,
			"organization": mappingDescription.OrgID,
			"os_family":    mappingDescription.OsFamily,
			"owner":        mappingDescription.Owner,
			"private":      mappingDescription.IsPrivate,
		}

		if mappingDescription.ID != nil {
			l["image_id"] = *mappingDescription.ID
		}

		result = append(result, l)
	}
	return result
}

func flattenImageMappingConstraints(constraints []*models.Constraint) []interface{} {
	if len(constraints) == 0 {
		return make([]interface{}, 0)
	}

	configConstraints := make([]interface{}, 0, len(constraints))

	for _, constraint := range constraints {
		helper := make(map[string]interface{})
		helper["mandatory"] = *constraint.Mandatory
		helper["expression"] = *constraint.Expression

		configConstraints = append(configConstraints, helper)
	}

	return configConstraints
}
