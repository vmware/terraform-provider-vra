// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vra-sdk-go/pkg/client/integration"
)

func TestAccVRAIntegration_Basic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "vra_integration.this"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckIntegration(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRAIntegrationConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAIntegrationExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "name", "my-integration-"+strconv.Itoa(rInt)),
					resource.TestCheckResourceAttr(
						resourceName, "description", "my integration"),
					resource.TestCheckResourceAttr(
						resourceName, "integration_type", "com.github.saas"),
					resource.TestCheckResourceAttr(
						resourceName, "tags.#", "1"),
				),
			},
			{
				Config: testAccCheckVRAIntegrationUpdateConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAIntegrationExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "name", "my-integration-"+strconv.Itoa(rInt)),
					resource.TestCheckResourceAttr(
						resourceName, "description", "my integration update"),
					resource.TestCheckResourceAttr(
						resourceName, "integration_type", "com.github.saas"),
					resource.TestCheckResourceAttr(
						resourceName, "tags.#", "1"),
				),
			},
		},
	})
}

func testAccCheckVRAIntegrationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no integration ID is set")
		}

		return nil
	}
}

func testAccCheckVRAIntegrationDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "vra_integration" {
			if _, err := apiClient.Integration.GetIntegration(integration.NewGetIntegrationParams().WithID(rs.Primary.ID)); err == nil {
				return fmt.Errorf("Resource 'vra_integration' still exists with id %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckVRAIntegrationConfig(rInt int) string {
	githubToken := os.Getenv("VRA_GITHUB_TOKEN")
	return fmt.Sprintf(`
	resource "vra_integration" "this" {
		name = "my-integration-%d"
		description = "my integration"
		integration_properties = {
			url: "https://api.github.com"
		}
		integration_type = "com.github.saas"
		private_key = "%s"
		tags {
			key = "created_by"
			value = "vra-terraform-provider"
		}
    }`, rInt, githubToken)
}

func testAccCheckVRAIntegrationUpdateConfig(rInt int) string {
	githubToken := os.Getenv("VRA_GITHUB_TOKEN")
	return fmt.Sprintf(`
	resource "vra_integration" "this" {
		name = "my-integration-%d"
		description = "my integration update"
		integration_properties = {
			url: "https://api.github.com"
		}
		integration_type = "com.github.saas"
		private_key = "%s"
		tags {
			key = "created_by"
			value = "vra-terraform-provider"
		}
	}`, rInt, githubToken)
}
