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

func TestAccDataSourceVRAMachine(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "vra_machine.my_machine"
	dataSourceName := "data.vra_machine.this"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckMachine(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRAMachineNotFound(rInt),
				ExpectError: regexp.MustCompile("vra_machine filter did not match any machine"),
			},
			{
				Config: testAccDataSourceVRAMachineNameFilter(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "id", dataSourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceName, "description"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "image", dataSourceName, "image"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor", dataSourceName, "flavor"),
				),
			},
			{
				Config: testAccDataSourceVRAMachineByID(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "id", dataSourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceName, "description"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "image", dataSourceName, "image"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor", dataSourceName, "flavor"),
				),
			},
		},
	})
}

func testAccDataSourceVRAMachineNotFound(rInt int) string {
	return testAccCheckVRAMachineConfig(rInt) + `
	data "vra_machine" "this" {
		filter = "name eq 'foobar'"
	}`
}

func testAccDataSourceVRAMachineNameFilter(rInt int) string {
	return testAccCheckVRAMachineConfig(rInt) + `
	data "vra_machine" "this" {
		filter = "name eq '${vra_machine.my_machine.name}'"
	}`
}

func testAccDataSourceVRAMachineByID(rInt int) string {
	return testAccCheckVRAMachineConfig(rInt) + `
	data "vra_machine" "this" {
		id = vra_machine.my_machine.id
	}`
}
