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
	"github.com/vmware/vra-sdk-go/pkg/client/location"
)

func TestAccVRAZoneBasic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRAZoneConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAZoneExists("vra_zone.my-zone"),
					resource.TestMatchResourceAttr(
						"vra_zone.my-zone", "name", regexp.MustCompile("^my-vra-zone-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"vra_zone.my-zone", "description", "description my-vra-zone"),
					resource.TestCheckResourceAttr(
						"vra_zone.my-zone", "placement_policy", "DEFAULT"),
					resource.TestCheckResourceAttr(
						"vra_zone.my-zone", "tags.#", "3"),
				),
			},
			{
				Config: testAccCheckVRAZoneUpdateConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAZoneExists("vra_zone.my-zone"),
					resource.TestCheckResourceAttr(
						"vra_zone.my-zone", "description", "description my-vra-zone-update"),
					resource.TestCheckResourceAttr(
						"vra_zone.my-zone", "placement_policy", "BINPACK"),
					resource.TestCheckResourceAttr(
						"vra_zone.my-zone", "tags.#", "2"),
				),
			},
		},
	})
}

func testAccCheckVRAZoneExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no cloud account ID is set")
		}

		return nil
	}
}

func testAccCheckVRAZoneDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "vra_cloud_account_aws" {
			_, err := apiClient.CloudAccount.GetAwsCloudAccount(cloud_account.NewGetAwsCloudAccountParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_cloud_account_aws' still exists with id %s", rs.Primary.ID)
			}
		}
		if rs.Type == "vra_zone" {
			_, err := apiClient.Location.GetZone(location.NewGetZoneParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_zone' still exists with id %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckVRAZoneConfig(rInt int) string {
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
    name = "my-vra-zone-%d"
	description = "description my-vra-zone"
	region_id = "${data.vra_region.us-east-1-region.id}"

    tags {
        key = "mykey"
        value = "myvalue"
    }
    tags {
        key = "foo"
        value = "bar"
    }
    tags {
        key = "faz"
        value = "baz"
    }
}`, id, secret, rInt)
}

func testAccCheckVRAZoneUpdateConfig(rInt int) string {
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
		name = "my-vra-zone-update-%d"
		description = "description my-vra-zone-update"
		region_id = "${data.vra_region.us-east-1-region.id}"
		placement_policy = "BINPACK"
		tags {
			key = "mykey"
			value = "myvalue"
		}
		tags {
			key = "foo"
			value = "bar"
		}
	}`, id, secret, rInt)
}

func TestAccVRAZoneInvalidPlacementPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckVRAZoneConfigInvalidPlacementPolicy(),
				ExpectError: regexp.MustCompile("\"placement_policy\" must be one of 'DEFAULT', 'SPREAD', 'BINPACK'"),
			},
		},
	})
}

func testAccCheckVRAZoneConfigInvalidPlacementPolicy() string {
	return `
 resource "vra_zone" "my-zone" {
	name = "my-vra-zone-update"
	description = "description my-vra-zone-update"
	region_id = "fakeid"
	placement_policy = "INVALID"
}`
}
