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
)

func TestAccVRACloudAccountAzure_Basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAzure(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRACloudAccountAzureDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRACloudAccountAzureConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRACloudAccountAzureExists("vra_cloud_account_azure.my-cloud-account"),
					resource.TestMatchResourceAttr(
						"vra_cloud_account_azure.my-cloud-account", "name", regexp.MustCompile("^my-cloud-account-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_azure.my-cloud-account", "description", "test cloud account"),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_azure.my-cloud-account", "application_id", os.Getenv("VRA_ARM_CLIENT_APP_ID")),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_azure.my-cloud-account", "subscription_id", os.Getenv("VRA_ARM_SUBSCRIPTION_ID")),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_azure.my-cloud-account", "tenant_id", os.Getenv("VRA_ARM_TENANT_ID")),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_azure.my-cloud-account", "tags.#", "2"),
				),
			},
			{
				Config: testAccCheckVRACloudAccountAzureUpdateConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRACloudAccountAzureExists("vra_cloud_account_azure.my-cloud-account"),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_azure.my-cloud-account", "description", "your test cloud account"),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_azure.my-cloud-account", "application_id", os.Getenv("VRA_ARM_CLIENT_APP_ID")),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_azure.my-cloud-account", "subscription_id", os.Getenv("VRA_ARM_SUBSCRIPTION_ID")),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_azure.my-cloud-account", "tenant_id", os.Getenv("VRA_ARM_TENANT_ID")),
				),
			},
		},
	})
}
func TestAccVRACloudAccountAzure_Duplicate(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAzure(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRACloudAccountAzureDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckVRACloudAccountAzureConfigDuplicateRegion(rInt),
				ExpectError: regexp.MustCompile("Specified regions are not unique"),
			},
		},
	})
}

func testAccCheckVRACloudAccountAzureExists(n string) resource.TestCheckFunc {
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

func testAccCheckVRACloudAccountAzureDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vra_cloud_account_azure" {
			continue
		}

		_, err := apiClient.CloudAccount.GetAzureCloudAccount(cloud_account.NewGetAzureCloudAccountParams().WithID(rs.Primary.ID))
		if err == nil {
			return err
		}
	}

	return nil
}

func testAccCheckVRACloudAccountAzureConfig(rInt int) string {
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
	tags {
		key = "foo"
		value = "bar"
	}
	tags {
		key = "where"
		value = "waldo"
	}
 }`, rInt, subscriptionID, tenantID, applicationID, applicationKey)
}

func testAccCheckVRACloudAccountAzureUpdateConfig(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	subscriptionID := os.Getenv("VRA_ARM_SUBSCRIPTION_ID")
	tenantID := os.Getenv("VRA_ARM_TENANT_ID")
	applicationID := os.Getenv("VRA_ARM_CLIENT_APP_ID")
	applicationKey := os.Getenv("VRA_ARM_CLIENT_APP_KEY")
	return fmt.Sprintf(`
resource "vra_cloud_account_azure" "my-cloud-account" {
	name = "my-cloud-account-%d"
	description = "your test cloud account"
	subscription_id = "%s"
	tenant_id = "%s"
	application_id = "%s"
	application_key = "%s"
	regions = ["centralus"]
 }`, rInt, subscriptionID, tenantID, applicationID, applicationKey)
}

func testAccCheckVRACloudAccountAzureConfigDuplicateRegion(rInt int) string {
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
	regions = ["centralus", "centralus"]
 }`, rInt, subscriptionID, tenantID, applicationID, applicationKey)
}
