// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/client/location"
)

func TestAccDataSourceVRAZoneBasic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName1 := "vra_zone.my-zone"
	dataSourceName1 := "data.vra_zone.test-zone"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRADataSourceZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRAZoneNoneConfig(rInt),
				ExpectError: regexp.MustCompile("zone with id `` or name `invalid-name` not found"),
			},
			{
				Config: testAccDataSourceVRAZoneOneConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
				),
			},
		},
	})
}

func testAccDataSourceVRAZone(rInt int) string {
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

	 resource "vra_zone" "my-zone" {
		name = "my-vra-zone-%d"
		description = "description my-vra-zone"
		region_id = "${element(vra_cloud_account_aws.my-cloud-account.regions, 0)}"
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

func testAccDataSourceVRAZoneNoneConfig(rInt int) string {
	return testAccDataSourceVRAZone(rInt) + `
		data "vra_zone" "test-zone" {
			name = "invalid-name"
		}`
}

func testAccDataSourceVRAZoneOneConfig(rInt int) string {
	return testAccDataSourceVRAZone(rInt) + `
		data "vra_zone" "test-zone" {
			name = "${vra_zone.my-zone.name}"
		}`
}

func testAccCheckVRADataSourceZoneDestroy(s *terraform.State) error {
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
