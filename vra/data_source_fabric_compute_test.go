// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFabricCompute_Basic(t *testing.T) {
	regionName := os.Getenv("VRA_FABRIC_COMPUTE_NAME")
	dataSourceName := "data.vra_fabric_compute.my-fabric-compute"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckFabricCompute(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRAFabricComputeNoConfig(),
				ExpectError: regexp.MustCompile("one of id or filter is required"),
			},
			{
				Config:      testAccDataSourceVRAFabricComputeNoneConfig(),
				ExpectError: regexp.MustCompile("filter doesn't match to any fabric compute"),
			},
			{
				Config: testAccDataSourceVRAFabricComputeOneConfig(regionName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "external_id", regionName),
					resource.TestCheckResourceAttr(dataSourceName, "name", regionName),
				),
			},
		},
	})
}

func testAccDataSourceVRAFabricComputeNoConfig() string {
	return `
		data "vra_fabric_compute" "my-fabric-compute" {
		}`
}

func testAccDataSourceVRAFabricComputeNoneConfig() string {
	return `
		data "vra_fabric_compute" "my-fabric-compute" {
			filter = "name eq 'foobar'"
		}`
}

func testAccDataSourceVRAFabricComputeOneConfig(regionName string) string {
	return fmt.Sprintf(`
		data "vra_fabric_compute" "my-fabric-compute" {
			filter = "name eq '%s'"
		}`, regionName)
}
