// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"strings"

	"github.com/vmware/vra-sdk-go/pkg/models"
)

// withString will return a string pointer of the passed in string value
func withString(s string) *string {
	return &s
}

func withBool(b bool) *bool {
	return &b
}

func withInt32(i int32) *int32 {
	return &i
}

// expandStringList will convert the interface list into a list of strings
func expandStringList(slist []interface{}) []string {
	vs := make([]string, 0, len(slist))
	for _, v := range slist {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, val)
		}
	}
	return vs
}

/*
func flattenStringList(list []*string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, *v)
	}
	return vs
}
*/

// compareUnique will determine if all of the items passed in are unique
func compareUnique(s []interface{}) bool {
	seen := make(map[string]struct{}, len(s))
	j := 0
	for _, v := range s {
		vs := v.(string)
		if _, ok := seen[vs]; ok {
			continue
		}
		seen[vs] = struct{}{}
		s[j] = vs
		j++
	}
	return j == len(s)
}

// indexOf will lookup and return the index of value in the list of items
func indexOf(value string, items []string) (int, error) {
	for i, v := range items {
		if v == value {
			return i, nil
		}
	}
	return -1, fmt.Errorf("Could not find %s in item list %v", value, items)
}

// flattenAssociatedCloudAccountIDs will return associated cloud account ids from the Href links in the order received
func flattenAssociatedCloudAccountIDs(links map[string]models.Href) []string {
	refStrings := links["associated-cloud-accounts"].Hrefs
	m := make([]string, len(refStrings))
	for i, r := range refStrings {
		m[i] = strings.TrimPrefix(r, "/iaas/api/cloud-accounts/")
	}
	return m
}

// expandInputs will convert the interface  into a map of [string:interface]
func expandInputs(configInputs interface{}) map[string]interface{} {
	if configInputs == nil {
		return nil
	}

	inputs := make(map[string]interface{})
	for key, value := range configInputs.(map[string]interface{}) {
		if value != nil {
			//inputs[key] = fmt.Sprint(value)
			inputs[key] = value
		}
	}

	return inputs
}

// expandInputsToString will convert the interface  into a map of string:string
func expandInputsToString(configInputs interface{}) map[string]string {
	if configInputs == nil {
		return nil
	}

	inputs := make(map[string]string)
	for key, value := range configInputs.(map[string]interface{}) {
		if value != nil {
			inputs[key] = fmt.Sprint(value)
		}
	}

	return inputs
}

// expandCatalogSourceConfig will convert the interface into a map of interface
func expandCatalogSourceConfig(catalogSourceConfig interface{}) map[string]interface{} {
	config := make(map[string]interface{})
	for key, value := range catalogSourceConfig.(map[string]interface{}) {
		if value != nil {
			config[key] = fmt.Sprint(value)
		}
	}

	return config
}

// flattenContentDefinition will convert the ContentDefinition to map of interface
func flattenContentDefinition(contentDefinition *models.ContentDefinition) interface{} {
	helper := make(map[string]interface{})

	helper["description"] = contentDefinition.Description
	helper["icon_id"] = contentDefinition.IconID
	helper["id"] = contentDefinition.ID
	helper["name"] = contentDefinition.Name
	helper["number_of_items"] = contentDefinition.NumItems
	helper["source_name"] = contentDefinition.SourceName
	helper["source_type"] = contentDefinition.SourceType
	helper["type"] = contentDefinition.Type

	definition := make([]interface{}, 0)
	definition = append(definition, helper)
	return definition
}
