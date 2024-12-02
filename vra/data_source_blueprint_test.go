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

func TestAccDataSourceVRABlueprint_NotFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckBlueprint(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRABlueprintNotFound(),
				ExpectError: regexp.MustCompile("blueprint foobar not found"),
			},
		},
	})
}

func TestAccDataSourceVRABlueprint_Found(t *testing.T) {
	resource1 := "vra_blueprint.this"
	dataSource := "data.vra_blueprint.this"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckBlueprint(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVRABlueprintFound(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSource, "name", resource1, "name"),
					resource.TestCheckResourceAttrPair(dataSource, "id", resource1, "id"),
					resource.TestCheckResourceAttrPair(dataSource, "description", resource1, "description"),
					resource.TestCheckResourceAttrPair(dataSource, "content", resource1, "content"),
					resource.TestCheckResourceAttrPair(dataSource, "project_id", resource1, "project_id"),
				),
			},
		},
	})
}

func testAccDataSourceVRABlueprintBase(rInt int) string {
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
	}`, rInt, rInt)
}

func testAccDataSourceVRABlueprintNotFound() string {
	return `data "vra_blueprint" "this" {
			name = "foobar"
		}`
}

func testAccDataSourceVRABlueprintFound() string {
	rInt := acctest.RandInt()
	return testAccDataSourceVRABlueprintBase(rInt) + `
		data "vra_blueprint" "this" {
			name = vra_blueprint.this.name
		}`
}
