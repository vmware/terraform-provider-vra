// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/client/flavor_profile"
)

func TestAccVRAFlavorProfileBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAFlavorProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRAFlavorProfileConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAFlavorProfileExists("vra_flavor_profile.my-flavor-profile"),
					resource.TestCheckResourceAttr(
						"vra_flavor_profile.my-flavor-profile", "name", "AWS"),
					resource.TestCheckResourceAttr(
						"vra_flavor_profile.my-flavor-profile", "description", "my flavor"),
					resource.TestCheckResourceAttr(
						"vra_flavor_profile.my-flavor-profile", "flavor_mapping.#", "2"),
					resource.TestCheckResourceAttr(
						"vra_flavor_profile.my-flavor-profile", "flavor_mapping.2163174927.name", "small"),
					resource.TestCheckResourceAttr(
						"vra_flavor_profile.my-flavor-profile", "flavor_mapping.2163174927.instance_type", "t2.small"),
					resource.TestCheckResourceAttr(
						"vra_flavor_profile.my-flavor-profile", "flavor_mapping.310071531.name", "medium"),
					resource.TestCheckResourceAttr(
						"vra_flavor_profile.my-flavor-profile", "flavor_mapping.310071531.instance_type", "t2.medium"),
				),
			},
		},
	})
}

func testAccCheckVRAFlavorProfileExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no flavor profile ID is set")
		}

		return nil
	}
}

func testAccCheckVRAFlavorProfileDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "vra_cloud_account_aws" {
			_, err := apiClient.CloudAccount.GetAwsCloudAccount(cloud_account.NewGetAwsCloudAccountParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_cloud_account_aws' still exists with id %s", rs.Primary.ID)
			}
		}
		if rs.Type == "vra_flavor_profile" {
			_, err := apiClient.FlavorProfile.GetFlavorProfile(flavor_profile.NewGetFlavorProfileParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_flavor_profile' still exists with id %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckVRAFlavorProfileConfig() string {
	// Need valid credentials since this is creating a real cloud account
	id := os.Getenv("VRA_AWS_ACCESS_KEY_ID")
	secret := os.Getenv("VRA_AWS_SECRET_ACCESS_KEY")
	return fmt.Sprintf(`
resource "vra_cloud_account_aws" "my-cloud-account" {
	name = "my-cloud-account"
	description = "test cloud account"
	access_key = "%s"
	secret_key = "%s"
	regions = ["us-east-1"]
 }

data "vra_region" "us-east-1-region" {
    cloud_account_id = "${vra_cloud_account_aws.my-cloud-account.id}"
    region = "us-east-1"
}

resource "vra_zone" "my-zone" {
    name = "my-vra-zone"
    description = "description my-vra-zone"
	region_id = "${data.vra_region.us-east-1-region.id}"
}

resource "vra_flavor_profile" "my-flavor-profile" {
	name = "AWS"
	description = "my flavor"
	region_id = "${data.vra_region.us-east-1-region.id}"
	flavor_mapping {
		name = "small"
		instance_type = "t2.small"
	}
	flavor_mapping {
		name = "medium"
		instance_type = "t2.medium"
	}
}`, id, secret)
}
