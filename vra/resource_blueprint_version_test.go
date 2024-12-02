// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"

	"github.com/vmware/vra-sdk-go/pkg/client/blueprint"

	"regexp"
	"strconv"
	"testing"

	"github.com/go-openapi/strfmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVRABlueprintVersion_Valid(t *testing.T) {
	rInt := acctest.RandInt()
	resource1 := "vra_blueprint_version.this"
	blueprint := "vra_blueprint.this"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckBlueprint(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRABlueprintVersionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRABlueprintVersionValidContentConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resource1, "name", regexp.MustCompile("^test-blueprint-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttrPair(resource1, "blueprint_id", blueprint, "id"),
					resource.TestCheckResourceAttrPair(resource1, "name", blueprint, "name"),
					resource.TestCheckResourceAttrPair(resource1, "blueprint_description", blueprint, "description"),
					resource.TestCheckResourceAttrPair(resource1, "id", resource1, "version"),
					resource.TestCheckResourceAttr(resource1, "description", "Released from vRA terraform provider"),
					resource.TestCheckResourceAttr(resource1, "change_log", "First version"),
				),
			},
		},
	})
}

func testAccCheckVRABlueprintVersionDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vra_blueprint" {
			continue
		}

		_, err := apiClient.Blueprint.GetBlueprintUsingGET1(blueprint.NewGetBlueprintUsingGET1Params().WithBlueprintID(strfmt.UUID(rs.Primary.ID)))
		if err == nil {
			return fmt.Errorf("resource 'vra_blueprint' still exists with id %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckVRABlueprintVersionValidContentConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "vra_project" "this" {
	  name = "terraform-test-project"
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

    }`, rInt)
}
