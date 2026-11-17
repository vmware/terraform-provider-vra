// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

func expandCustomProperties(configCustomProperties map[string]interface{}) map[string]string {
	customProperties := make(map[string]string)

	for key, value := range configCustomProperties {
		customProperties[key] = value.(string)
	}

	return customProperties
}

// removedCustomPropertyKeys returns the keys present in the old custom properties but absent from the new ones.
func removedCustomPropertyKeys(oldCustomProperties, newCustomProperties map[string]interface{}) []string {
	removed := make([]string, 0)

	for key := range oldCustomProperties {
		if _, ok := newCustomProperties[key]; !ok {
			removed = append(removed, key)
		}
	}

	return removed
}
