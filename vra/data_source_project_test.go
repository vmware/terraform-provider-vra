// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVRAProject(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName1 := "vra_project.my-project"
	dataSourceName1 := "data.vra_project.test-project"
	dataSourceName2 := "data.vra_project.this"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVra(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRAProjectNoneConfig(rInt),
				ExpectError: regexp.MustCompile("project invalid-name not found"),
			},
			{
				Config: testAccDataSourceVRAProjectOneConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
					resource.TestCheckResourceAttr(resourceName1, "constraints.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName1, "constraints.#", "1"),
				),
			},
			{
				Config: testAccDataSourceVRAProjectIDConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName2, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName2, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName2, "name"),
					resource.TestCheckResourceAttr(resourceName1, "constraints.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName2, "constraints.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceVRAProject(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	return fmt.Sprintf(`
	resource "vra_project" "my-project" {
		name = "my-project-%d"
		description = "test project"

		constraints {
    		extensibility {
      			expression = "foo:bar"
      			mandatory  = false
    		}
    		extensibility {
      			expression = "environment:Test"
      			mandatory  = true
			}

    		network {
      			expression = "foo:bar"
      			mandatory  = false
    		}
    		network {
      			expression = "environment:Test"
      			mandatory  = true
    		}

    		storage {
      			expression = "foo:bar"
      			mandatory  = false
			}
    		storage {
      			expression = "environment:Test"
      			mandatory  = true
    		}
  		}
	 }`, rInt)
}

func testAccDataSourceVRAProjectNoneConfig(rInt int) string {
	return testAccDataSourceVRAProject(rInt) + `
		data "vra_project" "test-project" {
			name = "invalid-name"
		}`
}

func testAccDataSourceVRAProjectOneConfig(rInt int) string {
	return testAccDataSourceVRAProject(rInt) + `
		data "vra_project" "test-project" {
			name = "${vra_project.my-project.name}"
		}`
}

func testAccDataSourceVRAProjectIDConfig(rInt int) string {
	return testAccDataSourceVRAProject(rInt) + `
		data "vra_project" "this" {
			id = "${vra_project.my-project.id}"
		}`
}
