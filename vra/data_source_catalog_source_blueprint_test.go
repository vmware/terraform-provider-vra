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

func TestAccDataSourceVRACatalogSourceBlueprint_Valid(t *testing.T) {
	rInt := acctest.RandInt()
	resource1 := "vra_catalog_source_blueprint.this"
	dataSourceWithProjectID := "data.vra_catalog_source_blueprint.with_project_id"
	dataSourceWithID := "data.vra_catalog_source_blueprint.with_id"
	dataSourceWithName := "data.vra_catalog_source_blueprint.with_name"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckBlueprint(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRACatalogSourceBlueprintFoundConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(resource1, "project_id", dataSourceWithProjectID, "project_id"),
					resource.TestCheckResourceAttrPair(resource1, "project_id", dataSourceWithProjectID, "config.sourceProjectId"),
					resource.TestCheckResourceAttrPair(resource1, "id", dataSourceWithProjectID, "id"),
					resource.TestCheckResourceAttrPair(resource1, "name", dataSourceWithProjectID, "name"),
					resource.TestCheckResourceAttrPair(resource1, "type_id", dataSourceWithProjectID, "type_id"),
					resource.TestCheckResourceAttrPair(resource1, "project_id", dataSourceWithID, "project_id"),
					resource.TestCheckResourceAttrPair(resource1, "project_id", dataSourceWithID, "config.sourceProjectId"),
					resource.TestCheckResourceAttrPair(resource1, "id", dataSourceWithID, "id"),
					resource.TestCheckResourceAttrPair(resource1, "name", dataSourceWithID, "name"),
					resource.TestCheckResourceAttrPair(resource1, "type_id", dataSourceWithID, "type_id"),
					resource.TestCheckResourceAttrPair(resource1, "project_id", dataSourceWithName, "project_id"),
					resource.TestCheckResourceAttrPair(resource1, "project_id", dataSourceWithName, "config.sourceProjectId"),
					resource.TestCheckResourceAttrPair(resource1, "id", dataSourceWithName, "id"),
					resource.TestCheckResourceAttrPair(resource1, "name", dataSourceWithName, "name"),
					resource.TestCheckResourceAttrPair(resource1, "type_id", dataSourceWithName, "type_id"),
				),
			},
			{
				Config:      testAccCheckVRACatalogSourceBlueprintByIDNotFoundConfig(),
				ExpectError: regexp.MustCompile("blueprint catalog source with id '22b52754-09d5-4706-b42e-fe487da520f3' is not found"),
			},
			{
				Config:      testAccCheckVRACatalogSourceBlueprintByProjectIDNotFoundConfig(),
				ExpectError: regexp.MustCompile("blueprint catalog source with project_id 'invalid-id' is not found"),
			},
			{
				Config:      testAccCheckVRACatalogSourceBlueprintByNameNotFoundConfig(),
				ExpectError: regexp.MustCompile("blueprint catalog source with name 'invalid-name' is not found"),
			},
		},
	})
}

func testAccCheckVRACatalogSourceBlueprintFoundConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "vra_project" "this" {
	  name = "terraform-test-project-%d"
	}

	resource "vra_catalog_source_blueprint" "this" {
	  name       = "tf-test-catalog-source-%d"
  	  project_id = vra_project.this.id
	}

	data "vra_catalog_source_blueprint" "with_project_id" {
 	  project_id = vra_catalog_source_blueprint.this.project_id
	}

	data "vra_catalog_source_blueprint" "with_id" {
 	  id = vra_catalog_source_blueprint.this.id
	}

	data "vra_catalog_source_blueprint" "with_name" {
 	  name = vra_catalog_source_blueprint.this.name
	}`, rInt, rInt)
}

func testAccCheckVRACatalogSourceBlueprintByIDNotFoundConfig() string {
	return `
	data "vra_catalog_source_blueprint" "invalid_id" {
 	  id = "22b52754-09d5-4706-b42e-fe487da520f3"
	}`
}

func testAccCheckVRACatalogSourceBlueprintByProjectIDNotFoundConfig() string {
	return `
	data "vra_catalog_source_blueprint" "invalid_project_id" {
 	  project_id = "invalid-id"
	}`
}

func testAccCheckVRACatalogSourceBlueprintByNameNotFoundConfig() string {
	return `
	data "vra_catalog_source_blueprint" "invalid_name" {
 	  name = "invalid-name"
	}`
}
