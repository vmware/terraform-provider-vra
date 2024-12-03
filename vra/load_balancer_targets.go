// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// LoadBalancerTargetSchema returns the schema to use for the targets property
func LoadBalancerTargetSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"machine_id": {
					Type:     schema.TypeString,
					Required: true,
				},
				"network_interface_id": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func expandLoadBalancerTargets(configLoadBalancerTargets []interface{}) []string {
	targets := make([]string, 0)

	for _, target := range configLoadBalancerTargets {
		targetMap := target.(map[string]interface{})

		machineID := targetMap["machine_id"].(string)
		networkInterfaceID := targetMap["network_interface_id"].(string)

		link := fmt.Sprintf("/iaas/api/machines/%s", machineID)
		if networkInterfaceID != "" {
			link = fmt.Sprintf("%s/network-interfaces/%s", link, networkInterfaceID)
		}
		targets = append(targets, link)
	}
	return targets
}
