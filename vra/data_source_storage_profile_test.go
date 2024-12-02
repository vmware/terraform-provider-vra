// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"testing"
)

func TestAccDataSourceVRAStorageProfile(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "vra_storage_profile.this"
	dataSourceName := "data.vra_storage_profile.this"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckStorageProfile(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAStorageProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRAStorageProfileNotFound(rInt),
				ExpectError: regexp.MustCompile("vra_storage_profile filter did not match any storage profile"),
			},
			// TODO: Enable filter by name once this is fixed https://jira.eng.vmware.com/browse/VCOM-13947
			// {
			// 	Config: testAccDataSourceVRAStorageProfileNameFilter(rInt),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestCheckResourceAttrPair(resourceName, "id", dataSourceName, "id"),
			// 		resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceName, "description"),
			// 		resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
			// 		resource.TestCheckResourceAttrPair(resourceName, "default_item", dataSourceName, "default_item"),
			// 		resource.TestCheckResourceAttrPair(resourceName, "external_region_id", dataSourceName, "region_id"),
			// 	),
			// },
			// TODO: Enable once https://jira.eng.vmware.com/browse/VCOM-13947 is fixed and include name in the filter to narrow the results to 1.
			// This works only if there is one storage profile with the externalRegionId.
			//{
			//	Config: testAccDataSourceVRAStorageProfileExternalRegionIDFilter(rInt),
			//	Check: resource.ComposeTestCheckFunc(
			//		resource.TestCheckResourceAttrPair(resourceName, "id", dataSourceName, "id"),
			//		resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceName, "description"),
			//		resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
			//		resource.TestCheckResourceAttrPair(resourceName, "default_item", dataSourceName, "default_item"),
			//		resource.TestCheckResourceAttrPair(resourceName, "external_region_id", dataSourceName, "external_region_id"),
			//	),
			//},
			{
				Config: testAccDataSourceVRAStorageProfileByID(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "id", dataSourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceName, "description"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "default_item", dataSourceName, "default_item"),
					resource.TestCheckResourceAttrPair(resourceName, "external_region_id", dataSourceName, "external_region_id"),
				),
			},
		},
	})
}

func testAccDataSourceVRAStorageProfileNotFound(rInt int) string {
	return testAccCheckVRAStorageProfileAWSConfig(rInt) + `
	data "vra_storage_profile" "this" {
		filter = "externalRegionId eq 'foobar'"
	}`
}

// TBD: Enable filter by name once this is fixed https://jira.eng.vmware.com/browse/VCOM-13947
// func testAccDataSourceVRAStorageProfileNameFilter(rInt int) string {
// 	return testAccCheckVRAStorageProfileAWSConfig(rInt) + fmt.Sprintf(`
// 	data "vra_storage_profile" "this" {
// 		filter = "name eq '${vra_storage_profile.my-storage-profile.name}'"
// 	}`)
// }

//func testAccDataSourceVRAStorageProfileExternalRegionIDFilter(rInt int) string {
//	return testAccCheckVRAStorageProfileAWSConfig(rInt) + `
//	data "vra_storage_profile" "this" {
//		filter = "externalRegionId eq '${data.vra_region.this.region}' and cloudAccountId eq '${data.vra_cloud_account_aws.this.id}'"
//	}`
//}

func testAccDataSourceVRAStorageProfileByID(rInt int) string {
	return testAccCheckVRAStorageProfileAWSConfig(rInt) + `
	data "vra_storage_profile" "this" {
		id = vra_storage_profile.this.id
	}`
}
