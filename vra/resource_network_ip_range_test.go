// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/client/network_ip_range"
)

func TestAccVRANetworkIPRangeBasic(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRANetworkIPRangeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRANetworkIPRangeConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRANetworkIPRangeExists("vra_network_ip_range.this"),
					resource.TestMatchResourceAttr(
						"vra_network_ip_range.this", "name", regexp.MustCompile("^my-vra-network-ip-range-"+strconv.Itoa(rInt))),
				),
			},
		},
	})
}

func testAccCheckVRANetworkIPRangeExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no network ip range ID is set")
		}
		return nil
	}
}

func testAccCheckVRANetworkIPRangeDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "vra_cloud_account_aws" {
			_, err := apiClient.CloudAccount.GetAwsCloudAccount(cloud_account.NewGetAwsCloudAccountParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_cloud_account_aws' still exists with id %s", rs.Primary.ID)
			}
		}
		if rs.Type == "vra_network_ip_range" {
			_, err := apiClient.NetworkIPRange.GetInternalNetworkIPRange(network_ip_range.NewGetInternalNetworkIPRangeParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_network_ip_range' still exists with id %s", rs.Primary.ID)
			}
		}

	}

	return nil
}

func testAccCheckVRANetworkIPRangeConfig(rInt int) string {

	id := os.Getenv("VRA_AWS_ACCESS_KEY_ID")
	secret := os.Getenv("VRA_AWS_SECRET_ACCESS_KEY")

	fabricNetworkName := os.Getenv("VRA_FABRICNETWORK_NAME")

	startIP := os.Getenv("VRA_NWIPRANGE_STARTIP")
	endIP := os.Getenv("VRA_NWIPRANGE_ENDIP")
	ipVersion := os.Getenv("VRA_NWIPRANGE_VERSION")

	return fmt.Sprintf(`

	resource "vra_cloud_account_aws" "this" {
		name = "my-cloud-account-%d"
		description = "aws test cloud account"
		access_key = "%s"
		secret_key = "%s"
		regions = ["us-east-1"]
	 }

	data "vra_fabric_network" "this" {
		filter = "name eq '%s'"
	}

	resource "vra_network_ip_range" "this" {
		name               = "my-vra-network-range-%d"
		start_ip_address   = "%s"
		end_ip_address     = "%s"
		ip_version         = "%s"
		fabric_network_ids = [data.vra_fabric_network.this.id]
	}

`, rInt, id, secret, fabricNetworkName, rInt, startIP, endIP, ipVersion)
}
