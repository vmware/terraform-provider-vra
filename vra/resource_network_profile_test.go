package vra

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/client/network_profile"
)

func TestAccVRANetworkProfileBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRANetworkProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRANetworkProfileConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRANetworkProfileExists("vra_network_profile.my-network-profile"),
					resource.TestCheckResourceAttr(
						"vra_network_profile.my-network-profile", "name", "my-vra-network-profile"),
					resource.TestCheckResourceAttr(
						"vra_network_profile.my-network-profile", "description", "my network profile"),
				),
			},
		},
	})
}

func testAccCheckVRANetworkProfileExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no network profile ID is set")
		}

		return nil
	}
}

func testAccCheckVRANetworkProfileDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "vra_cloud_account_aws" {
			_, err := apiClient.CloudAccount.GetAwsCloudAccount(cloud_account.NewGetAwsCloudAccountParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_cloud_account_aws' still exists with id %s", rs.Primary.ID)
			}
		}
		if rs.Type == "vra_network_profile" {
			_, err := apiClient.NetworkProfile.GetNetworkProfile(network_profile.NewGetNetworkProfileParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_network_profile' still exists with id %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckVRANetworkProfileConfig() string {
	// Need valid credentials since this is creating a real cloud account
	id := os.Getenv("VRA_AWS_ACCESS_KEY_ID")
	secret := os.Getenv("VRA_AWS_SECRET_ACCESS_KEY")
	return fmt.Sprintf(`
resource "vra_cloud_account_aws" "my-cloud-account" {
	name = "my-cloud-account"
	description = "test cloud account"
	access_key = "%s"
	secret_key = "%s"
	regions = ["us-east-1"]
 }

data "vra_region" "us-east-1-region" {
    cloud_account_id = "${vra_cloud_account_aws.my-cloud-account.id}"
    region = "us-east-1"
}

resource "vra_zone" "my-zone" {
    name = "my-vra-zone"
    description = "description my-vra-zone"
	region_id = "${data.vra_region.us-east-1-region.id}"
}

resource "vra_network_profile" "my-network-profile" {
	name = "my-vra-network-profile"
	description = "my network profile"
	region_id = "${data.vra_region.us-east-1-region.id}"
	isolation_type = "NONE"
}`, id, secret)
}
