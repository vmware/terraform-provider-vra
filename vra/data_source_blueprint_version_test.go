// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"testing"
)

func TestAccDataSourceVRABlueprintVersion_Valid(t *testing.T) {
	rInt := acctest.RandInt()
	resource1 := "vra_blueprint_version.this"
	dataSource := "data.vra_blueprint_version.this"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckBlueprint(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRABlueprintVersionFoundConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(resource1, "blueprint_id", dataSource, "blueprint_id"),
					resource.TestCheckResourceAttrPair(resource1, "id", dataSource, "id"),
					resource.TestCheckResourceAttrPair(resource1, "description", dataSource, "description"),
					resource.TestCheckResourceAttrPair(resource1, "blueprint_description", dataSource, "blueprint_description"),
					resource.TestCheckResourceAttrPair(resource1, "project_id", dataSource, "project_id"),
					resource.TestCheckResourceAttrPair(resource1, "project_name", dataSource, "project_name"),
				),
			},
			{
				Config:      testAccCheckVRABlueprintVersionNotFoundConfig(rInt),
				ExpectError: regexp.MustCompile("blueprint version '10' is not found"),
			},
		},
	})
}

func testAccCheckVRABlueprintVersionFoundConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "vra_project" "this" {
	  name = "terraform-test-project-%d"
	}

	resource "vra_blueprint" "this" {
	  name        = "test-blueprint-%d"
	  description = "terraform test blueprint"

	  project_id = vra_project.this.id

	  content = <<-EOT
				 formatVersion: 1
				 inputs: {}
				 resources:  {}
				EOT
	}

   resource "vra_blueprint_version" "this" {
	  blueprint_id = vra_blueprint.this.id
     description  = "Released from vRA terraform provider"
     version      = 1
     release      = true
     change_log   = "First version"

   }

	data "vra_blueprint_version" "this" {
 	  blueprint_id = vra_blueprint.this.id
 	  id           = vra_blueprint_version.this.id
	}`, rInt, rInt)
}

func testAccCheckVRABlueprintVersionNotFoundConfig(rInt int) string {
	return testAccCheckVRABlueprintVersionFoundConfig(rInt) + `
	data "vra_blueprint_version" "invalid" {
  	  blueprint_id = vra_blueprint.this.id
  	  id           = 10
	}`
}
