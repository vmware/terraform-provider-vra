// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func extractIDsFromRegion(regions []*models.Region) []string {
	regionIDs := []string{}

	for _, region := range regions {
		regionIDs = append(regionIDs, *region.ExternalRegionID)
	}

	return regionIDs
}

func extractIDsFromRegionSpecification(regions []*models.RegionSpecification) []string {
	regionIDs := []string{}

	for _, region := range regions {
		regionIDs = append(regionIDs, *region.ExternalRegionID)
	}

	return regionIDs
}

func expandRegionSpecificationList(rList []interface{}) []*models.RegionSpecification {
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
