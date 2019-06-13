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

func TestAccCASStorageProfileBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCASStorageProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCASStorageProfileConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCASStorageProfileExists("cas_storage_profile.my-storage-profile"),
					resource.TestCheckResourceAttr(
						"cas_storage_profile.my-storage-profile", "name", "my-cas-storage-profile"),
					resource.TestCheckResourceAttr(
						"cas_storage_profile.my-storage-profile", "description", "my storage profile"),
					resource.TestCheckResourceAttr(
						"cas_storage_profile.my-storage-profile", "default_item", true),
				),
			},
		},
	})
}

func testAccCheckCASStorageProfileExists(n string) resource.TestCheckFunc {
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

func testAccCheckCASStorageProfileDestroy(s *terraform.State) error {
	client := testAccProviderCAS.Meta().(*tango.Client)
	apiClient := client.GetAPIClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "cas_cloud_account_aws" {
			_, err := apiClient.CloudAccount.GetAwsCloudAccount(cloud_account.NewGetAwsCloudAccountParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'cas_cloud_account_aws' still exists with id %s", rs.Primary.ID)
			}
		}
		if rs.Type == "cas_storage_profile" {
			_, err := apiClient.StorageProfile.GetStorageProfile(storage_profile.NewGetStorageProfileParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'cas_storage_profile' still exists with id %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckCASStorageProfileConfig() string {
	// Need valid credentials since this is creating a real cloud account
	id := os.Getenv("CAS_AWS_ACCESS_KEY_ID")
	secret := os.Getenv("CAS_AWS_SECRET_ACCESS_KEY")
	return fmt.Sprintf(`
resource "cas_cloud_account_aws" "my-cloud-account" {
	name = "my-cloud-account"
	description = "test cloud account"
	access_key = "%s"
	secret_key = "%s"
	regions = ["us-east-1"]
 }

data "cas_region" "us-east-1-region" {
    cloud_account_id = "${cas_cloud_account_aws.my-cloud-account.id}"
    region = "us-east-1"
}

resource "cas_zone" "my-zone" {
    name = "my-cas-zone"
    description = "description my-cas-zone"
	region_id = "${data.cas_region.us-east-1-region.id}"
}

resource "cas_storage_profile" "my-storage-profile" {
	name = "my-cas-storage-profile"
	description = "my storage profile"
	region_id = "${data.cas_region.us-east-1-region.id}"
	default_item = true
}`, id, secret)
}
