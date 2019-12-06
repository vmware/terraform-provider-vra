package vra

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/client/storage_profile"
)

func TestAccVRAStorageProfileBasic(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckStorageProfile(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAStorageProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRAStorageProfileConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAStorageProfileExists("vra_storage_profile.this"),
					resource.TestMatchResourceAttr(
						"vra_storage_profile.this", "name", regexp.MustCompile("^my-profile-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"vra_storage_profile.this", "description", "my storage profile"),
					resource.TestCheckResourceAttr(
						"vra_storage_profile.this", "default_item", "true"),
					resource.TestCheckResourceAttr(
						"vra_storage_profile.this", "external_region_id", os.Getenv("VRA_REGION")),
				),
			},
		},
	})
}

func testAccCheckVRAStorageProfileExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no storage profile ID is set")
		}

		return nil
	}
}

func testAccCheckVRAStorageProfileDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "vra_cloud_account_aws" {
			_, err := apiClient.CloudAccount.GetAwsCloudAccount(cloud_account.NewGetAwsCloudAccountParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_cloud_account_aws' still exists with id %s", rs.Primary.ID)
			}
		}
		if rs.Type == "vra_storage_profile" {
			_, err := apiClient.StorageProfile.GetStorageProfile(storage_profile.NewGetStorageProfileParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_storage_profile' still exists with id %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckVRAStorageProfileConfig(rInt int) string {
	// Need valid credentials since this is creating a real cloud account and network profile
	id := os.Getenv("VRA_AWS_ACCESS_KEY_ID")
	secret := os.Getenv("VRA_AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("VRA_REGION")
	return fmt.Sprintf(`
	resource "vra_cloud_account_aws" "this" {
		name = "my-cloud-account-%d"
		description = "test cloud account"
		access_key = "%s"
		secret_key = "%s"
		regions = ["%s"]
	 }

	data "vra_region" "this" {
		cloud_account_id = "${vra_cloud_account_aws.this.id}"
		region = "%s"
	}

	resource "vra_zone" "this" {
		name = "my-vra-zone-%d"
		description = "description my-vra-zone"
		region_id = "${data.vra_region.this.id}"
	}

resource "vra_storage_profile" "this" {
	name = "my-profile-%d"
	description = "my storage profile"
	region_id = "${data.vra_region.this.id}"
	default_item = true
	disk_properties = {
		deviceType = "instance-store"
	}
}`, rInt, id, secret, region, region, rInt, rInt)
}
