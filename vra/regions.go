// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func extractIDsFromRegion(regions []*models.Region) []string {
	var regionIDs []string

	for _, region := range regions {
		regionIDs = append(regionIDs, *region.ExternalRegionID)
	}

	return regionIDs
}

func extractIDsFromRegionSpecification(regions []*models.RegionSpecification) []string {
	var regionIDs []string

	for _, region := range regions {
		regionIDs = append(regionIDs, *region.ExternalRegionID)
	}

	return regionIDs
}

func expandRegionSpecificationList(rList []any) []*models.RegionSpecification {
	rsList := make([]*models.RegionSpecification, 0, len(rList))

	for _, r := range rList {
		rs := models.RegionSpecification{
			ExternalRegionID: withString(r.(string)),
			Name:             withString(r.(string)),
		}
		rsList = append(rsList, &rs)
	}

	return rsList
}

func expandEnabledRegions(enabledRegionsMap []any) []*models.RegionSpecification {
	rsList := make([]*models.RegionSpecification, 0, len(enabledRegionsMap))

	for _, region := range enabledRegionsMap {
		regionMap := region.(map[string]interface{})
		helper := models.RegionSpecification{
			ExternalRegionID: withString(regionMap["external_region_id"].(string)),
			Name:             withString(regionMap["name"].(string)),
		}
		rsList = append(rsList, &helper)
	}

	return rsList
}

func flattenEnabledRegions(regions []*models.Region) []any {
	if len(regions) == 0 {
		return make([]any, 0)
	}

	regionsMap := make([]any, 0, len(regions))

	for _, region := range regions {
		helper := make(map[string]any)
		helper["external_region_id"] = region.ExternalRegionID
		helper["id"] = region.ID
		helper["name"] = region.Name

		regionsMap = append(regionsMap, helper)
	}

	return regionsMap
}

func flattenExternalRegions(regions []*models.RegionSpecification) []any {
	if len(regions) == 0 {
		return make([]any, 0)
	}

	regionsMap := make([]any, 0, len(regions))

	for _, region := range regions {
		helper := make(map[string]any)
		helper["external_region_id"] = region.ExternalRegionID
		helper["name"] = region.Name

		regionsMap = append(regionsMap, helper)
	}

	return regionsMap
}
