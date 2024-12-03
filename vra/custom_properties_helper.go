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
