// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVRANetwork_Basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRANetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRANetworkConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRANetworkExists("vra_network.my_network"),
					resource.TestMatchResourceAttr(
						"vra_network.my_network", "name", regexp.MustCompile("^terraform_vra_network-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"vra_network.my_network", "outbound_access", "false"),
					resource.TestCheckResourceAttr(
						"vra_network.my_network", "constraints.#", "1"),
					resource.TestCheckResourceAttr(
						"vra_network.my_network", "constraints.0.mandatory", "true"),
					resource.TestCheckResourceAttr(
						"vra_network.my_network", "constraints.0.expression", "pci"),

					// Currently not working as expected - possible bug on Provisioning service's side
					resource.TestCheckResourceAttr(
						"vra_network.my_network", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"vra_network.my_network", "tags.0.key", "stoyan"),
					resource.TestCheckResourceAttr(
						"vra_network.my_network", "tags.0.value", "genchev"),
				),
			},
		},
	})
}

func testAccCheckVRANetworkExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Network ID is set")
		}

		return nil
	}
}

func testAccCheckVRANetworkDestroy(_ *terraform.State) error {
	/*
		apiClient := testAccProviderVRA.Meta().(*Client).apiClient

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "vra_network" {
				continue
			}

			_, err := apiClient.ReadResource("/iaas/networks/" + rs.Primary.ID)

			if err != nil && !strings.Contains(err.Error(), "404") {
				return fmt.Errorf(
					"Error waiting for network (%s) to be destroyed: %s",
					rs.Primary.ID, err)
			}
		}
	*/

	return nil
}

func testAccCheckVRANetworkConfig(rInt int) string {
	return fmt.Sprintf(`
resource "vra_network" "my_network" {
  name = "terraform_vra_network-%d"
  outbound_access = false

  tags {
	key = "stoyan"
    value = "genchev"
  }

  constraints {
	  mandatory = true
	  expression = "pci"
  }
}`, rInt)
}
