// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vra-sdk-go/pkg/client/fabric_network"
)

func TestAccVRAFabricNetworkVsphere_importBasic(t *testing.T) {
	resourceName := "vra_fabric_network_vsphere.this"
	fabricNetworkID := os.Getenv("VRA_FABRIC_NETWORK_VSPHERE_ID")
	createErrorRegex := regexp.MustCompile("vra_fabric_network_vsphere resources are only importable")

	checkFn := func(s []*terraform.InstanceState) error {
		if len(s) != 1 {
			return fmt.Errorf("expected 1 state %#v", s)
		}

		fabricNetworkVsphereState := s[0]
		if fabricNetworkID != fabricNetworkVsphereState.ID {
			return fmt.Errorf("expected fabric network ID of %s,%s received instead", fabricNetworkID, fabricNetworkVsphereState.ID)
		}
		return nil
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckVsphere(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAFabricNetworkVsphereDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckVRAFabricNetworkVsphereConfig(),
				ExpectError: createErrorRegex,
			},
			{
				Config:           testAccCheckVRAFabricNetworkVsphereConfig(),
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateId:    fabricNetworkID,
				ImportStateCheck: checkFn,
			},
		},
	})
}

func testAccCheckVRAFabricNetworkVsphereDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "vra_fabric_network_vsphere" {
			_, err := apiClient.FabricNetwork.GetVsphereFabricNetwork(fabric_network.NewGetVsphereFabricNetworkParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_fabric_network_vsphere' with id %s does not exist", rs.Primary.ID)
			}
		}

	}

	return nil
}

func testAccCheckVRAFabricNetworkVsphereConfig() string {
	CIDR := os.Getenv("VRA_FABRIC_NETWORK_CIDR")
	gateway := os.Getenv("VRA_FABRIC_NETWORK_GW")
	domain := os.Getenv("VRA_FABRIC_NETWORK_DOMAIN")
	return fmt.Sprintf(`
	resource "vra_fabric_network_vsphere" "this" {
		cidr = "%s"
		default_gateway = "%s"
		domain = "%s"
	}

`, CIDR, gateway, domain)
}
