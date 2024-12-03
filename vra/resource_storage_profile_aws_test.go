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
	"github.com/vmware/vra-sdk-go/pkg/client/storage_profile"
)

func TestAccVRAStorageProfileAwsBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAStorageProfileAwsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRAStorageProfileAwsConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAStorageProfileAwsExists("vra_storage_profile_aws.my-storage-profile-aws"),
					resource.TestCheckResourceAttr(
						"vra_storage_profile_aws.my-storage-profile-aws", "name", "my-vra-storage-profile-aws"),
					resource.TestCheckResourceAttr(
						"vra_storage_profile_aws.my-storage-profile-aws", "description", "my storage profile aws"),
					resource.TestCheckResourceAttr(
						"vra_storage_profile_aws.my-storage-profile-aws", "default_item", "true"),
					resource.TestCheckResourceAttr(
						"vra_storage_profile_aws.my-storage-profile-aws", "device_type", "EBS"),
				),
			},
		},
	})
}

func testAccCheckVRAStorageProfileAwsExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no storage profile aws ID is set")
		}

		return nil
	}
}

func testAccCheckVRAStorageProfileAwsDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "vra_cloud_account_aws" {
			_, err := apiClient.CloudAccount.GetAwsCloudAccount(cloud_account.NewGetAwsCloudAccountParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_cloud_account_aws' still exists with id %s", rs.Primary.ID)
			}
		}
		if rs.Type == "vra_storage_profile_aws" {
			_, err := apiClient.StorageProfile.GetStorageProfile(storage_profile.NewGetStorageProfileParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_storage_profile_aws' still exists with id %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckVRAStorageProfileAwsConfig() string {
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

resource "vra_storage_profile_aws" "my-storage-profile-aws" {
	name = "my-vra-storage-profile-aws"
	description = "my storage profile aws"
	region_id = "${data.vra_region.us-east-1-region.id}"
	default_item = true
	device_type = "EBS"
}`, id, secret)
}
