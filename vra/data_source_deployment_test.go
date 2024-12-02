// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"regexp"
	"testing"
)

func TestAccDataSourceVRADeployment(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName1 := "vra_deployment.this"
	dataSourceName1 := "data.vra_deployment.this"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDeploymentDataSource(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRADeploymentNoneConfig(),
				ExpectError: regexp.MustCompile("deployment invalid-name not found"),
			},
			{
				Config: testAccDataSourceVRADeploymentOneConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
					resource.TestCheckResourceAttrPair(resourceName1, "project_id", dataSourceName1, "project_id"),
					resource.TestCheckResourceAttrPair(resourceName1, "blueprint_id", dataSourceName1, "blueprint_id"),
				),
			},
		},
	})
}

func testAccDataSourceVRADeployment(rInt int) string {
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

	resource "vra_deployment" "this" {
	  name        = "test-deployment-%d"
	  description = "terraform test deployment"

	  blueprint_id = vra_blueprint.this.id
	  project_id = vra_project.this.id
	}`, rInt, rInt, rInt)
}

func testAccDataSourceVRADeploymentNoneConfig() string {
	return `data "vra_deployment" "invalid" {
			name = "invalid-name"
		}`
}

func testAccDataSourceVRADeploymentOneConfig(rInt int) string {
	return testAccDataSourceVRADeployment(rInt) + `
		data "vra_deployment" "this" {
			name = "${vra_deployment.this.name}"
		}`
}
