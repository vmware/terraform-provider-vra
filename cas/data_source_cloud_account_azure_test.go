package cas

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceCASCloudAccountAzure(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName1 := "cas_cloud_account_azure.my-cloud-account"
	dataSourceName1 := "data.cas_cloud_account_azure.test-cloud-account"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckAzure(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceCASCloudAccountAzureNotFound(rInt),
				ExpectError: regexp.MustCompile("cloud account foobar not found"),
			},
			{
				Config: testAccDataSourceCASCloudAccountAzure(rInt),
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

func testAccDataSourceCASCloudAccountAzureBase(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	subscriptionID := os.Getenv("CAS_ARM_SUBSCRIPTION_ID")
	tenantID := os.Getenv("CAS_ARM_TENANT_ID")
	applicationID := os.Getenv("CAS_ARM_CLIENT_APP_ID")
	applicationKey := os.Getenv("CAS_ARM_CLIENT_APP_KEY")
	return fmt.Sprintf(`
	resource "cas_cloud_account_azure" "my-cloud-account" {
		name = "my-cloud-account-%d"
		description = "test cloud account"
		subscription_id = "%s"
		tenant_id = "%s"
		application_id = "%s"
		application_key = "%s"
		regions = ["centralus"]
	 }`, rInt, subscriptionID, tenantID, applicationID, applicationKey)
}

func testAccDataSourceCASCloudAccountAzureNotFound(rInt int) string {
	return testAccDataSourceCASCloudAccountAzureBase(rInt) + fmt.Sprintf(`
	data "cas_cloud_account_azure" "test-cloud-account" {
		name = "foobar"
	}`)
}

func testAccDataSourceCASCloudAccountAzure(rInt int) string {
	return testAccDataSourceCASCloudAccountAzureBase(rInt) + fmt.Sprintf(`
	data "cas_cloud_account_azure" "test-cloud-account" {
		name = "${cas_cloud_account_azure.my-cloud-account.name}"
	}`)
}
