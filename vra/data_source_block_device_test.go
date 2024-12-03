// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"regexp"
	"testing"
)

func TestAccDataSourceVRABlockDevice(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "vra_block_device.this"
	dataSourceName := "data.vra_block_device.this"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckBlockDevice(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRABlockDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRABlockDeviceNotFound(rInt),
				ExpectError: regexp.MustCompile("vra_block_device filter did not match any block device"),
			},
			{
				Config: testAccDataSourceVRABlockDeviceNameFilter(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "id", dataSourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceName, "description"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "capacity_in_gb", dataSourceName, "capacity_in_gb"),
					resource.TestCheckResourceAttrPair(resourceName, "external_region_id", dataSourceName, "external_region_id"),
				),
			},
			{
				Config: testAccDataSourceVRABlockDeviceByID(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "id", dataSourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceName, "description"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "capacity_in_gb", dataSourceName, "capacity_in_gb"),
					resource.TestCheckResourceAttrPair(resourceName, "external_region_id", dataSourceName, "external_region_id"),
				),
			},
		},
	})
}

func testAccDataSourceVRABlockDeviceNotFound(rInt int) string {
	return testAccCheckVRABlockDeviceConfig(rInt) + `
	data "vra_block_device" "this" {
		filter = "name eq 'foobar'"
	}`
}

func testAccDataSourceVRABlockDeviceNameFilter(rInt int) string {
	return testAccCheckVRABlockDeviceConfig(rInt) + `
	data "vra_block_device" "this" {
		filter = "name eq '${vra_block_device.this.name}'"
	}`
}

func testAccDataSourceVRABlockDeviceByID(rInt int) string {
	return testAccCheckVRABlockDeviceConfig(rInt) + `
	data "vra_block_device" "this" {
		id = vra_block_device.this.id
	}`
}
