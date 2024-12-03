// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVRANetwork(t *testing.T) {
	rInt := acctest.RandInt()
	dataSourceName1 := "data.vra_network.test-network"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVra(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRANetworkNoneConfig(rInt),
				ExpectError: regexp.MustCompile("network invalid-name not found"),
			},
			{
				Config: testAccDataSourceVRANetworkOneConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName1, "id", "6d25dcb5d510875582822c89a1d4"),
				),
			},
		},
	})
}

func testAccDataSourceVRANetworkNoneConfig(_ int) string {
	return `
	    data "vra_network" "test-network" {
			name = "invalid-name"
		}`
}

func testAccDataSourceVRANetworkOneConfig(_ int) string {
	return `
		data "vra_network" "test-network" {
			name = "foo1-mcm653-56201379059"
		}`
}
