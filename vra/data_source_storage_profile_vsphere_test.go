// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"testing"
)

func TestAccDataSourceStorageProfileVsphere(t *testing.T) {
	resourceName1 := "vra_storage_profile_vsphere.this"
	dataSourceName1 := "data.vra_storage_profile_vsphere.this"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckStorageProfileVsphere(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAStorageProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRAStorageProfileVsphereNotFound(),
				ExpectError: regexp.MustCompile("vra_storage_profile_vsphere filter did not match any storage profile"),
			},
			//{
			// TODO: Enable filter by name once this is fixed https://jira.eng.vmware.com/browse/VCOM-13947
			//Config: testAccDataSourceVRAStorageProfileVsphereExternalRegionIDFilter(),
			//Check: resource.ComposeTestCheckFunc(
			//	resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
			//	resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
			//	resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
			//	resource.TestCheckResourceAttrPair(resourceName1, "default_item", dataSourceName1, "default_item"),
			//	resource.TestCheckResourceAttrPair(resourceName1, "external_region_id", dataSourceName1, "external_region_id"),
			//),
			//},
			{
				Config: testAccDataSourceVRAStorageProfileVsphereByID(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
					resource.TestCheckResourceAttrPair(resourceName1, "default_item", dataSourceName1, "default_item"),
					resource.TestCheckResourceAttrPair(resourceName1, "external_region_id", dataSourceName1, "external_region_id"),
				),
			},
		},
	})
}

func testAccDataSourceVRAStorageProfileVsphereNotFound() string {
	return testAccCheckVRAStorageProfileVsphereConfig() + `
	data "vra_storage_profile_vsphere" "this" {
		filter = "externalRegionId eq 'foobar'"
	}`
}

func testAccDataSourceVRAStorageProfileVsphereByID() string {
	return testAccCheckVRAStorageProfileVsphereConfig() + `
	data "vra_storage_profile_vsphere" "this" {
		id = vra_storage_profile_vsphere.this.id
	}`
}
