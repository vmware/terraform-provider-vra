// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"os"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"testing"
)

func TestAccDataSourceStorageProfileAzure(t *testing.T) {
	resourceName1 := "vra_storage_profile_azure.this"
	dataSourceName1 := "data.vra_storage_profile_azure.this"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckStorageProfileAzure(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAStorageProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRAStorageProfileAzureNotFound(),
				ExpectError: regexp.MustCompile("vra_storage_profile_azure filter did not match any storage profile"),
			},
			{
				Config: testAccDataSourceVRAStorageProfileAzureByID(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
					resource.TestCheckResourceAttrPair(resourceName1, "default_item", dataSourceName1, "default_item"),
					resource.TestCheckResourceAttrPair(resourceName1, "external_region_id", dataSourceName1, "external_region_id"),
				),
			},
		},
	})
}

func testAccDataSourceVRAStorageProfileAzureNotFound() string {
	return testAccCheckVRAStorageProfileAzureConfig() + `
	data "vra_storage_profile_azure" "this" {
		filter = "externalRegionId eq 'foobar'"
	}`
}

func testAccDataSourceVRAStorageProfileAzureByID() string {
	return testAccDataSourceVRAStorageProfileAzureConfig() + `
	data "vra_storage_profile_azure" "this" {
		id = vra_storage_profile_azure.this.id
	}`
}

func testAccDataSourceVRAStorageProfileAzureConfig() string {
	azureCloudAccountName := os.Getenv("VRA_AZURE_CLOUD_ACCOUNT_NAME")
	azureRegionName := os.Getenv("VRA_AZURE_REGION_NAME")
	return fmt.Sprintf(`
	data "vra_cloud_account_azure" "this" {
  		name = "%s"
	}

	data "vra_region" "this" {
  		cloud_account_id = data.vra_cloud_account_azure.this.id
  		region           = "%s"
	}

	resource "vra_storage_profile_azure" "this" {
  		name                = "azure-with-managed-disks-1"
  		description         = "Azure Storage Profile with managed disks."
  		region_id           = data.vra_region.this.id
  		default_item        = false
  		supports_encryption = false

  		data_disk_caching   = "None"         // Supported Values: None, ReadOnly, ReadWrite
  		disk_type           = "Standard_LRS" // Supported Values: Standard_LRS, StandardSSD_LRS, Premium_LRS
  		os_disk_caching     = "None"         // Supported Values: None, ReadOnly, ReadWrite

  		tags {
    		key   = "foo"
    		value = "bar"
  		}
	}`, azureCloudAccountName, azureRegionName)
}
