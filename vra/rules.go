// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// rulesSchema returns the schema to use for the rules property
func rulesSchema(isRequired bool) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Required: isRequired,
		Optional: !isRequired,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"access": {
					Type:     schema.TypeString,
					Required: true,
				},
				"direction": {
					Type:     schema.TypeString,
					Required: true,
				},
				"ip_range_cidr": {
					Type:     schema.TypeInt,
					Required: true,
				},
				"name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"ports": {
					Type:     schema.TypeString,
					Required: true,
				},
				"protocol": {
					Type:     schema.TypeString,
					Required: true,
				},
				"service": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

/*
func expandRules(configRules []interface{}) []*models.Rule {
	rules := make([]*models.Rule, 0, len(configRules))

	for _, configRule := range configRules {
		ruleMap := configRule.(map[string]interface{})

		rule := models.Rule{
			Access:      withString(ruleMap["access"].(string)),
			IPRangeCidr: withString(ruleMap["ip_range_cidr"].(string)),
			Ports:       withString(ruleMap["ports"].(string)),
			Protocol:    withString(ruleMap["protocol"].(string)),
		}

		if v, ok := ruleMap["name"].(string); ok && v != "" {
			rule.Name = v
		}

		rules = append(rules, &rule)
	}

	return rules
}
*/
