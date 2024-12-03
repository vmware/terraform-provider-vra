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

func TestAccDataSourceVRAImageProfile(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName1 := "vra_image_profile.this"
	dataSourceName1 := "data.vra_image_profile.this"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckImageProfile(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAImageProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRAImageProfileNotFound(rInt),
				ExpectError: regexp.MustCompile("vra_image_profile filter did not match any image profile"),
			},
			{
				Config: testAccDataSourceVRAImageProfileByName(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
					resource.TestCheckResourceAttrPair(resourceName1, "image_mapping", dataSourceName1, "image_mapping"),
					resource.TestCheckResourceAttrPair(resourceName1, "region_id", dataSourceName1, "region_id"),
					resource.TestCheckResourceAttrPair(resourceName1, "external_region_id", dataSourceName1, "external_region_id"),
				),
			},
			{
				Config: testAccDataSourceVRAImageProfileByRegionID(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
					resource.TestCheckResourceAttrPair(resourceName1, "image_mapping", dataSourceName1, "image_mapping"),
					resource.TestCheckResourceAttrPair(resourceName1, "region_id", dataSourceName1, "region_id"),
					resource.TestCheckResourceAttrPair(resourceName1, "external_region_id", dataSourceName1, "external_region_id"),
				),
			},
			{
				Config: testAccDataSourceVRAImageProfileByFilter(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
					resource.TestCheckResourceAttrPair(resourceName1, "image_mapping", dataSourceName1, "image_mapping"),
					resource.TestCheckResourceAttrPair(resourceName1, "region_id", dataSourceName1, "region_id"),
					resource.TestCheckResourceAttrPair(resourceName1, "external_region_id", dataSourceName1, "external_region_id"),
				),
			},
			{
				Config: testAccDataSourceVRAImageProfileByID(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
					resource.TestCheckResourceAttrPair(resourceName1, "image_mapping", dataSourceName1, "image_mapping"),
					resource.TestCheckResourceAttrPair(resourceName1, "region_id", dataSourceName1, "region_id"),
					resource.TestCheckResourceAttrPair(resourceName1, "external_region_id", dataSourceName1, "external_region_id"),
				),
			},
		},
	})
}

func testAccDataSourceVRAImageProfileNotFound(rInt int) string {
	return testAccCheckVRAImageProfileConfig(rInt) + `
	data "vra_image_profile" "this" {
		filter = "name eq 'foobar'"
	}`
}

func testAccDataSourceVRAImageProfileByName(rInt int) string {
	return testAccCheckVRAImageProfileConfig(rInt) + `
	data "vra_image_profile" "this" {
		name = vra_image_profile.this.name
	}`
}

func testAccDataSourceVRAImageProfileByRegionID(rInt int) string {
	return testAccCheckVRAImageProfileConfig(rInt) + `
	data "vra_image_profile" "this" {
		region_id = vra_image_profile.this.region_id
	}`
}

func testAccDataSourceVRAImageProfileByFilter(rInt int) string {
	return testAccCheckVRAImageProfileConfig(rInt) + `
	data "vra_image_profile" "this" {
		filter = "regionId eq '${vra_image_profile.this.region_id}'"
	}`
}

func testAccDataSourceVRAImageProfileByID(rInt int) string {
	return testAccCheckVRAImageProfileConfig(rInt) + `
	data "vra_image_profile" "this" {
		id = vra_image_profile.this.id
	}`
}
