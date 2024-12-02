// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// nicsSchema returns the schema to use for the nics property
func nicsSchema(isRequired bool) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Required: isRequired,
		Optional: !isRequired,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"description": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"device_index": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"network_id": {
					Type:     schema.TypeString,
					Required: true,
				},
				"addresses": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"security_group_ids": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"custom_properties": {
					Type:     schema.TypeMap,
					Optional: true,
				},
			},
		},
	}
}

func expandNics(configNics []interface{}) []*models.NetworkInterfaceSpecification {
	nics := make([]*models.NetworkInterfaceSpecification, 0, len(configNics))

	for _, configNic := range configNics {
		nicMap := configNic.(map[string]interface{})

		nic := models.NetworkInterfaceSpecification{
			NetworkID: nicMap["network_id"].(string),
		}

		if v, ok := nicMap["name"].(string); ok && v != "" {
			nic.Name = v
		}

		if v, ok := nicMap["description"].(string); ok && v != "" {
			nic.Description = v
		}

		if v, ok := nicMap["device_index"].(int32); ok && v != 0 {
			nic.DeviceIndex = v
		}

		if v, ok := nicMap["addresses"].([]interface{}); ok && len(v) != 0 {
			addresses := make([]string, 0)

			for _, value := range v {
				addresses = append(addresses, value.(string))
			}

			nic.Addresses = addresses
		}

		if v, ok := nicMap["security_group_ids"].([]interface{}); ok && len(v) != 0 {
			securityGroupIDs := make([]string, 0)

			for _, value := range v {
				securityGroupIDs = append(securityGroupIDs, value.(string))
			}

			nic.SecurityGroupIds = securityGroupIDs
		}

		nic.CustomProperties = expandCustomProperties(nicMap["custom_properties"].(map[string]interface{}))

		nics = append(nics, &nic)
	}

	return nics
}
