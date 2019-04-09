package cas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/cas-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/cas-sdk-go/pkg/client/flavor_profile"

	tango "github.com/vmware/terraform-provider-cas/sdk"
)

func TestAccCASFlavorBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCASFlavorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCASFlavorConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCASFlavorExists("cas_flavor.my-flavor"),
					resource.TestCheckResourceAttr(
						"cas_flavor.my-flavor", "name", "AWS"),
					resource.TestCheckResourceAttr(
						"cas_flavor.my-flavor", "description", "my flavor"),
					resource.TestCheckResourceAttr(
						"cas_flavor.my-flavor", "flavor_mapping.#", "2"),
					resource.TestCheckResourceAttr(
						"cas_flavor.my-flavor", "flavor_mapping.2163174927.name", "small"),
					resource.TestCheckResourceAttr(
						"cas_flavor.my-flavor", "flavor_mapping.2163174927.instance_type", "t2.small"),
					resource.TestCheckResourceAttr(
						"cas_flavor.my-flavor", "flavor_mapping.310071531.name", "medium"),
					resource.TestCheckResourceAttr(
						"cas_flavor.my-flavor", "flavor_mapping.310071531.instance_type", "t2.medium"),
				),
			},
		},
	})
}

func testAccCheckCASFlavorExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no flavor ID is set")
		}

		return nil
	}
}

func testAccCheckCASFlavorDestroy(s *terraform.State) error {
	client := testAccProviderCAS.Meta().(*tango.Client)
	apiClient := client.GetAPIClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "cas_cloud_account_aws" {
			_, err := apiClient.CloudAccount.GetAwsCloudAccount(cloud_account.NewGetAwsCloudAccountParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'cas_cloud_account_aws' still exists with id %s", rs.Primary.ID)
			}
		}
		if rs.Type == "cas_flavor" {
			_, err := apiClient.FlavorProfile.GetFlavorProfile(flavor_profile.NewGetFlavorProfileParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'cas_Flavor' still exists with id %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckCASFlavorConfig() string {
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

resource "cas_flavor" "my-flavor" {
	name = "AWS"
	description = "my flavor"
	region_id = "${data.cas_region.us-east-1-region.id}"
	flavor_mapping {
		name = "small"
		instance_type = "t2.small"
	}
	flavor_mapping {
		name = "medium"
		instance_type = "t2.medium"
	}
}`, id, secret)
}
