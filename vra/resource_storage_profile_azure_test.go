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

func TestAccVRAStorageProfileAzureBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckStorageProfileAzure(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAStorageProfileAzureDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRAStorageProfileAzureConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAStorageProfileAzureExists("vra_storage_profile_azure.my-storage-profile-azure"),
					resource.TestCheckResourceAttr(
						"vra_storage_profile_azure.my-storage-profile-azure", "name", "my-vra-storage-profile-azure"),
					resource.TestCheckResourceAttr(
						"vra_storage_profile_azure.my-storage-profile-azure", "description", "my storage profile azure"),
					resource.TestCheckResourceAttr(
						"vra_storage_profile_azure.my-storage-profile-azure", "default_item", "true"),
					resource.TestCheckResourceAttr(
						"vra_storage_profile_azure.my-storage-profile-azure", "disk_type", "Standard HDD"),
					resource.TestCheckResourceAttr(
						"vra_storage_profile_azure.my-storage-profile-azure", "os_disk_caching", "Read Only"),
					resource.TestCheckResourceAttr(
						"vra_storage_profile_azure.my-storage-profile-azure", "data_disk_caching", "Read Only"),
				),
			},
		},
	})
}

func testAccCheckVRAStorageProfileAzureExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no storage profile azure ID is set")
		}

		return nil
	}
}

func testAccCheckVRAStorageProfileAzureDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "vra_cloud_account_azure" {
			_, err := apiClient.CloudAccount.GetAzureCloudAccount(cloud_account.NewGetAzureCloudAccountParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_cloud_account_azure' still exists with id %s", rs.Primary.ID)
			}
		}
		if rs.Type == "vra_storage_profile_azure" {
			_, err := apiClient.StorageProfile.GetStorageProfile(storage_profile.NewGetStorageProfileParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_storage_profile_azure' still exists with id %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckVRAStorageProfileAzureConfig() string {
	// Need valid credentials since this is creating a real cloud account
	subscriptionID := os.Getenv("VRA_ARM_SUBSCRIPTION_ID")
	tenantID := os.Getenv("VRA_ARM_TENANT_ID")
	applicationID := os.Getenv("VRA_ARM_CLIENT_APP_ID")
	applicationKey := os.Getenv("VRA_ARM_CLIENT_APP_KEY")
	region := os.Getenv("VRA_ARM_REGION")
	return fmt.Sprintf(`
resource "vra_cloud_account_azure" "my-cloud-account" {
	name = "my-cloud-account"
	description = "test cloud account"
	subscription_id = "%s"
	tenant_id = "%s"
	application_id = "%s"
	application_key = "%s"
	regions = ["%s"]
 }

data "vra_region" "my-region" {
    cloud_account_id = "${vra_cloud_account_azure.my-cloud-account.id}"
    region = "%s"
}

resource "vra_storage_profile_azure" "my-storage-profile-azure" {
	name = "my-vra-storage-profile-azure"
	description = "my storage profile azure"
	region_id = "${data.vra_region.my-region.id}"
	default_item = true
	disk_type = "Standard HDD"
	os_disk_caching = "Read Only"
    data_disk_caching = "Read Only"
}`, subscriptionID, tenantID, applicationID, applicationKey, region, region)
}
