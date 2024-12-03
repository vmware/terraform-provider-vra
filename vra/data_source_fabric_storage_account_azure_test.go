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
)

func TestAccDataSourceVRAFabricStorageAccountAzure(t *testing.T) {
	dataSourceName := "data.vra_fabric_storage_account_azure.this"
	regionName := os.Getenv("VRA_ARM_REGION_NAME")
	fabricStorageAccountName := os.Getenv("VRA_ARM_FABRIC_STORAGE_ACCOUNT_NAME")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckFabricStorageAccountAzure(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRAFabricStorageAccountAzureNoConfig(),
				ExpectError: regexp.MustCompile("one of id or filter is required"),
			},
			{
				Config:      testAccDataSourceVRAFabricStorageAccountAzureNoneConfig(),
				ExpectError: regexp.MustCompile("vra_fabric_storage_account_azure filter did not match any fabric Azure storage accounts"),
			},
			{
				Config: testAccDataSourceVRAFabricStorageAccountOneConfig(regionName, fabricStorageAccountName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "external_region_id", regionName),
					resource.TestCheckResourceAttr(dataSourceName, "name", fabricStorageAccountName),
					resource.TestCheckResourceAttr(dataSourceName, "cloud_account_ids.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "external_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "links.#", "3"), // One each for self, cloud account and region
					resource.TestCheckResourceAttrSet(dataSourceName, "org_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "type"),
				),
			},
		},
	})
}

func testAccDataSourceVRAFabricStorageAccountAzureBaseConfig() string {
	// Need valid credentials since this is creating a real cloud account
	cloudAccountName := os.Getenv("VRA_ARM_CLOUD_ACCOUNT_NAME")
	return fmt.Sprintf(`
	data "vra_cloud_account_azure" "this" {
		name = "%s"
	}`, cloudAccountName)
}

func testAccDataSourceVRAFabricStorageAccountAzureNoConfig() string {
	return testAccDataSourceVRAFabricStorageAccountAzureBaseConfig() + `
	data "vra_fabric_storage_account_azure" "this" {
	}`
}

func testAccDataSourceVRAFabricStorageAccountAzureNoneConfig() string {
	return testAccDataSourceVRAFabricStorageAccountAzureBaseConfig() + `
	data "vra_fabric_storage_account_azure" "this" {
		filter = "name eq 'foobar'"
	}`
}

func testAccDataSourceVRAFabricStorageAccountOneConfig(regionName, fabricStorageAccountName string) string {
	return testAccDataSourceVRAFabricStorageAccountAzureBaseConfig() + fmt.Sprintf(`
	data "vra_fabric_storage_account_azure" "this" {
		filter = "name eq '%s' and externalRegionId eq '%s' and cloudAccountId eq '${data.vra_cloud_account_azure.this.id}'"
	}`, fabricStorageAccountName, regionName)
}
