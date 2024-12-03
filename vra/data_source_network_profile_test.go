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

func TestAccDataSourceVRANetworkProfile(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName1 := "vra_network_profile.this"
	dataSourceName1 := "data.vra_network_profile.this"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRANetworkProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRANetworkProfileNotFound(rInt),
				ExpectError: regexp.MustCompile("vra_network_profile filter did not match any network profile"),
			},
			{
				Config: testAccDataSourceVRANetworkProfileNameFilter(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
					resource.TestCheckResourceAttrPair(resourceName1, "isolation_type", dataSourceName1, "isolation_type"),
					resource.TestCheckResourceAttrPair(resourceName1, "region_id", dataSourceName1, "region_id"),
				),
			},
			{
				Config: testAccDataSourceVRANetworkProfileRegionIDFilter(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
					resource.TestCheckResourceAttrPair(resourceName1, "isolation_type", dataSourceName1, "isolation_type"),
					resource.TestCheckResourceAttrPair(resourceName1, "region_id", dataSourceName1, "region_id"),
				),
			},
			{
				Config: testAccDataSourceVRANetworkProfileByID(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
					resource.TestCheckResourceAttrPair(resourceName1, "isolation_type", dataSourceName1, "isolation_type"),
					resource.TestCheckResourceAttrPair(resourceName1, "region_id", dataSourceName1, "region_id"),
				),
			},
		},
	})
}

func testAccDataSourceVRANetworkProfileNotFound(rInt int) string {
	return testAccCheckVRANetworkProfileConfig(rInt) + `
	data "vra_network_profile" "this" {
		filter = "name eq 'foobar'"
	}`
}

func testAccDataSourceVRANetworkProfileNameFilter(rInt int) string {
	return testAccCheckVRANetworkProfileConfig(rInt) + `
	data "vra_network_profile" "this" {
		filter = "name eq '${vra_network_profile.this.name}'"
	}`
}

func testAccDataSourceVRANetworkProfileRegionIDFilter(rInt int) string {
	return testAccCheckVRANetworkProfileConfig(rInt) + `
	data "vra_network_profile" "this" {
		filter = "regionId eq '${data.vra_region.this.id}'"
	}`
}

func testAccDataSourceVRANetworkProfileByID(rInt int) string {
	return testAccCheckVRANetworkProfileConfig(rInt) + `
	data "vra_network_profile" "this" {
		id = vra_network_profile.this.id
	}`
}
