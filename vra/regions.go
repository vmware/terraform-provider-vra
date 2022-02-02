package vra

import (
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func extractIdsFromRegion(regions []*models.Region) []string {
	regionIds := []string{}

	for _, region := range regions {
		regionIds = append(regionIds, *region.ExternalRegionID)
	}

	return regionIds
}

func extractIdsFromRegionSpecification(regions []*models.RegionSpecification) []string {
	regionIds := []string{}

	for _, region := range regions {
		regionIds = append(regionIds, *region.ExternalRegionID)
	}

	return regionIds
}

func expandRegionSpecificationList(rList []interface{}) []*models.RegionSpecification {
	rsList := make([]*models.RegionSpecification, 0, len(rList))

	for _, r := range rList {
		rs := models.RegionSpecification{
			ExternalRegionID: withString(r.(string)),
			Name:             withString(""),
		}
		rsList = append(rsList, &rs)
	}

	return rsList
}
