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

// flattenAndNormalizeCloudAccountRegionIds will return region id's in the same order as regionOrder
func flattenAndNormalizeCloudAccountRegionIds(regionOrder []string, cloudAccount *models.CloudAccount) ([]string, error) {
	returnOrder := cloudAccount.EnabledRegionIds
	refStrings := cloudAccount.Links["regions"].Hrefs
	m := make([]string, len(regionOrder))
	for i, r := range regionOrder {
		index, err := indexOf(r, returnOrder)
		if err != nil {
			return []string{}, err
		}
		m[i] = strings.TrimPrefix(refStrings[index], "/iaas/api/regions/")
	}
	return m, nil
}

// flattenAndNormalizeCLoudAccountAWSRegionIds will return region id's in the same order as regionOrder
func flattenAndNormalizeCLoudAccountAWSRegionIds(regionOrder []string, cloudAccount *models.CloudAccountAws) ([]string, error) {
	returnOrder := cloudAccount.EnabledRegionIds
	refStrings := cloudAccount.Links["regions"].Hrefs
	m := make([]string, len(regionOrder))
	for i, r := range regionOrder {
		index, err := indexOf(r, returnOrder)
		if err != nil {
			return []string{}, err
		}
		m[i] = strings.TrimPrefix(refStrings[index], "/iaas/api/regions/")
	}
	return m, nil
}

// flattenAndNormalizeCLoudAccountAzureRegionIds will return region id's in the same order as regionOrder
func flattenAndNormalizeCLoudAccountAzureRegionIds(regionOrder []string, cloudAccount *models.CloudAccountAzure) ([]string, error) {
	returnOrder := cloudAccount.EnabledRegionIds
	refStrings := cloudAccount.Links["regions"].Hrefs
	m := make([]string, len(regionOrder))
	for i, r := range regionOrder {
		index, err := indexOf(r, returnOrder)
		if err != nil {
			return []string{}, err
		}
		m[i] = strings.TrimPrefix(refStrings[index], "/iaas/api/regions/")
	}
	return m, nil
}

// flattenAndNormalizeCloudAccountVsphereRegionIds will return region id's in the same order as regionOrder
func flattenAndNormalizeCloudAccountVsphereRegionIds(regionOrder []string, cloudAccount *models.CloudAccountVsphere) ([]string, error) {
	returnOrder := cloudAccount.EnabledRegionIds
	refStrings := cloudAccount.Links["regions"].Hrefs
	m := make([]string, len(regionOrder))
	for i, r := range regionOrder {
		index, err := indexOf(r, returnOrder)
		if err != nil {
			return []string{}, err
		}
		m[i] = strings.TrimPrefix(refStrings[index], "/iaas/api/regions/")
	}
	return m, nil
}

// flattenAssociatedCloudAccountIds will return associated cloud account ids from the Href links in the order received
func flattenAssociatedCloudAccountIds(links map[string]models.Href) []string {
	refStrings := links["associated-cloud-accounts"].Hrefs
	m := make([]string, len(refStrings))
	for i, r := range refStrings {
		m[i] = strings.TrimPrefix(r, "/iaas/api/cloud-accounts/")
	}
	return m
}

// flattenAndNormalizeCLoudAccountGcpRegionIds will return region id's in the same order as regionOrder
func flattenAndNormalizeCLoudAccountGcpRegionIds(regionOrder []string, cloudAccount *models.CloudAccountGcp) ([]string, error) {
	returnOrder := cloudAccount.EnabledRegionIds
	refStrings := cloudAccount.Links["regions"].Hrefs
	m := make([]string, len(regionOrder))
	for i, r := range regionOrder {
		index, err := indexOf(r, returnOrder)
		if err != nil {
			return []string{}, err
		}
		m[i] = strings.TrimPrefix(refStrings[index], "/iaas/api/regions/")
	}
	return m, nil
}

// expandInputs will convert the interface  into a map of interface
func expandInputs(configInputs interface{}) map[string]interface{} {
	inputs := make(map[string]interface{})
	for key, value := range configInputs.(map[string]interface{}) {
		if value != nil {
			inputs[key] = fmt.Sprint(value)
		}
	}

	return inputs
}
