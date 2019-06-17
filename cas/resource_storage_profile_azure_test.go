package cas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/cas-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/cas-sdk-go/pkg/client/storage_profile"

	tango "github.com/vmware/terraform-provider-cas/sdk"
)

func TestAccCASStorageProfileAzureBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCASStorageProfileAzureDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCASStorageProfileAzureConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCASStorageProfileAzureExists("cas_storage_profile_azure.my-storage-profile-azure"),
					resource.TestCheckResourceAttr(
						"cas_storage_profile_azure.my-storage-profile-azure", "name", "my-cas-storage-profile-azure"),
					resource.TestCheckResourceAttr(
						"cas_storage_profile_azure.my-storage-profile-azure", "description", "my storage profile azure"),
					resource.TestCheckResourceAttr(
						"cas_storage_profile_azure.my-storage-profile-azure", "default_item", true),
					resource.TestCheckResourceAttr(
						"cas_storage_profile_azure.my-storage-profile-azure", "disk_type", "Standard HDD"),
					resource.TestCheckResourceAttr(
						"cas_storage_profile_azure.my-storage-profile-azure", "os_disk_caching", "Read Only"),
					resource.TestCheckResourceAttr(
						"cas_storage_profile_azure.my-storage-profile-azure", "data_disk_caching", "Read Only"),
				),
			},
		},
	})
}

func testAccCheckCASStorageProfileAzureExists(n string) resource.TestCheckFunc {
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

func testAccCheckCASStorageProfileAzureDestroy(s *terraform.State) error {
	client := testAccProviderCAS.Meta().(*tango.Client)
	apiClient := client.GetAPIClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "cas_cloud_account_azure" {
			_, err := apiClient.CloudAccount.GetAzureCloudAccount(cloud_account.NewGetAzureCloudAccountParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'cas_cloud_account_azure' still exists with id %s", rs.Primary.ID)
			}
		}
		if rs.Type == "cas_storage_profile_azure" {
			_, err := apiClient.StorageProfile.GetStorageProfile(storage_profile.NewGetStorageProfileParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'cas_storage_profile_azure' still exists with id %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckCASStorageProfileAzureConfig() string {
	// Need valid credentials since this is creating a real cloud account
	subscriptionID := os.Getenv("CAS_ARM_SUBSCRIPTION_ID")
	tenantID := os.Getenv("CAS_ARM_TENANT_ID")
	applicationID := os.Getenv("CAS_ARM_CLIENT_APP_ID")
	applicationKey := os.Getenv("CAS_ARM_CLIENT_APP_KEY")
	return fmt.Sprintf(`
resource "cas_cloud_account_azure" "my-cloud-account" {
	name = "my-cloud-account"
	description = "test cloud account"
	subscription_id = "%s"
	tenant_id = "%s"
	application_id = "%s"
	application_key = "%s"
	regions = ["eastus"]
 }

data "cas_region" "us-east-azure-region" {
    cloud_account_id = "${cas_cloud_account_azure.my-cloud-account.id}"
    region = "eastus"
}

resource "cas_storage_profile_azure" "my-storage-profile-azure" {
	name = "my-cas-storage-profile-azure"
	description = "my storage profile azure"
	region_id = "${data.cas_region.us-east-azure-region.id}"
	default_item = true
	disk_type = "Standard HDD"
	os_disk_caching = "Read Only"
    data_disk_caching = "Read Only"
}`, subscriptionID, tenantID, applicationID, applicationKey)
}
