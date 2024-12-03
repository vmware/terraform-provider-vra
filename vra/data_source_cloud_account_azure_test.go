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
)

func TestAccDataSourceVRACloudAccountAzure(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName1 := "vra_cloud_account_azure.my-cloud-account"
	dataSourceName1 := "data.vra_cloud_account_azure.test-cloud-account"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckAzure(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRACloudAccountAzureNotFound(rInt),
				ExpectError: regexp.MustCompile("cloud account foobar not found"),
			},
			{
				Config: testAccDataSourceVRACloudAccountAzure(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
					resource.TestCheckResourceAttrPair(resourceName1, "application_id", dataSourceName1, "application_id"),
					resource.TestCheckResourceAttrPair(resourceName1, "subscription_id", dataSourceName1, "subscription_id"),
					resource.TestCheckResourceAttrPair(resourceName1, "tenant_id", dataSourceName1, "tenant_id"),
				),
			},
		},
	})
}

func testAccDataSourceVRACloudAccountAzureBase(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	subscriptionID := os.Getenv("VRA_ARM_SUBSCRIPTION_ID")
	tenantID := os.Getenv("VRA_ARM_TENANT_ID")
	applicationID := os.Getenv("VRA_ARM_CLIENT_APP_ID")
	applicationKey := os.Getenv("VRA_ARM_CLIENT_APP_KEY")
	return fmt.Sprintf(`
	resource "vra_cloud_account_azure" "my-cloud-account" {
		name = "my-cloud-account-%d"
		description = "test cloud account"
		subscription_id = "%s"
		tenant_id = "%s"
		application_id = "%s"
		application_key = "%s"
		regions = ["centralus"]
	 }`, rInt, subscriptionID, tenantID, applicationID, applicationKey)
}

func testAccDataSourceVRACloudAccountAzureNotFound(rInt int) string {
	return testAccDataSourceVRACloudAccountAzureBase(rInt) + `
	data "vra_cloud_account_azure" "test-cloud-account" {
		name = "foobar"
	}`
}

func testAccDataSourceVRACloudAccountAzure(rInt int) string {
	return testAccDataSourceVRACloudAccountAzureBase(rInt) + `
	data "vra_cloud_account_azure" "test-cloud-account" {
		name = "${vra_cloud_account_azure.my-cloud-account.name}"
	}`
}
