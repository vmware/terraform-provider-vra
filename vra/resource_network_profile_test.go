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
	"github.com/vmware/vra-sdk-go/pkg/client/network_profile"
)

func TestAccVRANetworkProfileBasic(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRANetworkProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRANetworkProfileConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRANetworkProfileExists("vra_network_profile.this"),
					resource.TestMatchResourceAttr(
						"vra_network_profile.this", "name", regexp.MustCompile("^my-vra-network-profile-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"vra_network_profile.this", "description", "my network profile"),
				),
			},
		},
	})
}

func testAccCheckVRANetworkProfileExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no network profile ID is set")
		}

		return nil
	}
}

func testAccCheckVRANetworkProfileDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "vra_cloud_account_aws" {
			_, err := apiClient.CloudAccount.GetAwsCloudAccount(cloud_account.NewGetAwsCloudAccountParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_cloud_account_aws' still exists with id %s", rs.Primary.ID)
			}
		}
		if rs.Type == "vra_network_profile" {
			_, err := apiClient.NetworkProfile.GetNetworkProfile(network_profile.NewGetNetworkProfileParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_network_profile' still exists with id %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckVRANetworkProfileConfig(rInt int) string {
	// Need valid credentials since this is creating a real cloud account and network profile
	id := os.Getenv("VRA_AWS_ACCESS_KEY_ID")
	secret := os.Getenv("VRA_AWS_SECRET_ACCESS_KEY")
	return fmt.Sprintf(`
resource "vra_cloud_account_aws" "this" {
	name = "my-cloud-account-%d"
	description = "test cloud account"
	access_key = "%s"
	secret_key = "%s"
	regions = ["us-east-1"]
 }

data "vra_region" "this" {
    cloud_account_id = "${vra_cloud_account_aws.this.id}"
    region = "us-east-1"
}

resource "vra_zone" "this" {
    name = "my-vra-zone-%d"
    description = "description my-vra-zone"
	region_id = "${data.vra_region.this.id}"
}

resource "vra_network_profile" "this" {
	name = "my-vra-network-profile-%d"
	description = "my network profile"
	region_id = "${data.vra_region.this.id}"
	isolation_type = "NONE"
}`, rInt, id, secret, rInt, rInt)
}
