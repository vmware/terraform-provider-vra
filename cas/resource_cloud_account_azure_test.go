package cas

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/cas-sdk-go/pkg/client/cloud_account"
	tango "github.com/vmware/terraform-provider-cas/sdk"
)

func TestAccCASCloudAccountAzure_Basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAzure(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCASCloudAccountAzureDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCASCloudAccountAzureConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCASCloudAccountAzureExists("cas_cloud_account_azure.my-cloud-account"),
					resource.TestMatchResourceAttr(
						"cas_cloud_account_azure.my-cloud-account", "name", regexp.MustCompile("^my-cloud-account-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"cas_cloud_account_azure.my-cloud-account", "description", "test cloud account"),
					resource.TestCheckResourceAttr(
						"cas_cloud_account_azure.my-cloud-account", "application_id", os.Getenv("CAS_ARM_CLIENT_APP_ID")),
					resource.TestCheckResourceAttr(
						"cas_cloud_account_azure.my-cloud-account", "subscription_id", os.Getenv("CAS_ARM_SUBSCRIPTION_ID")),
					resource.TestCheckResourceAttr(
						"cas_cloud_account_azure.my-cloud-account", "tenant_id", os.Getenv("CAS_ARM_TENANT_ID")),
					resource.TestCheckResourceAttr(
						"cas_cloud_account_azure.my-cloud-account", "tags.#", "2"),
				),
			},
			{
				Config: testAccCheckCASCloudAccountAzureUpdateConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCASCloudAccountAzureExists("cas_cloud_account_azure.my-cloud-account"),
					resource.TestCheckResourceAttr(
						"cas_cloud_account_azure.my-cloud-account", "description", "your test cloud account"),
					resource.TestCheckResourceAttr(
						"cas_cloud_account_azure.my-cloud-account", "application_id", os.Getenv("CAS_ARM_CLIENT_APP_ID")),
					resource.TestCheckResourceAttr(
						"cas_cloud_account_azure.my-cloud-account", "subscription_id", os.Getenv("CAS_ARM_SUBSCRIPTION_ID")),
					resource.TestCheckResourceAttr(
						"cas_cloud_account_azure.my-cloud-account", "tenant_id", os.Getenv("CAS_ARM_TENANT_ID")),
				),
			},
		},
	})
}
func TestAccCASCloudAccountAzure_Duplicate(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAzure(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCASCloudAccountAzureDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckCASCloudAccountAzureConfigDuplicateRegion(rInt),
				ExpectError: regexp.MustCompile("Specified regions are not unique"),
			},
		},
	})
}

func testAccCheckCASCloudAccountAzureExists(n string) resource.TestCheckFunc {
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

func testAccCheckCASCloudAccountAzureDestroy(s *terraform.State) error {
	client := testAccProviderCAS.Meta().(*tango.Client)
	apiClient := client.GetAPIClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cas_cloud_account_azure" {
			continue
		}

		_, err := apiClient.CloudAccount.GetAzureCloudAccount(cloud_account.NewGetAzureCloudAccountParams().WithID(rs.Primary.ID))
		if err == nil {
			return err
		}
	}

	return nil
}

func testAccCheckCASCloudAccountAzureConfig(rInt int) string {
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

func testAccCheckCASCloudAccountAzureUpdateConfig(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	subscriptionID := os.Getenv("CAS_ARM_SUBSCRIPTION_ID")
	tenantID := os.Getenv("CAS_ARM_TENANT_ID")
	applicationID := os.Getenv("CAS_ARM_CLIENT_APP_ID")
	applicationKey := os.Getenv("CAS_ARM_CLIENT_APP_KEY")
	return fmt.Sprintf(`
resource "cas_cloud_account_azure" "my-cloud-account" {
	name = "my-cloud-account-%d"
	description = "your test cloud account"
	subscription_id = "%s"
	tenant_id = "%s"
	application_id = "%s"
	application_key = "%s"
	regions = ["centralus"]
 }`, rInt, subscriptionID, tenantID, applicationID, applicationKey)
}

func testAccCheckCASCloudAccountAzureConfigDuplicateRegion(rInt int) string {
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
	regions = ["centralus", "centralus"]
 }`, rInt, subscriptionID, tenantID, applicationID, applicationKey)
}
