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
	"github.com/vmware/cas-sdk-go/pkg/client/location"
)

func TestAccCASZoneBasic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCASZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCASZoneConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCASZoneExists("cas_zone.my-zone"),
					resource.TestMatchResourceAttr(
						"cas_zone.my-zone", "name", regexp.MustCompile("^my-cas-zone-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"cas_zone.my-zone", "description", "description my-cas-zone"),
					resource.TestCheckResourceAttr(
						"cas_zone.my-zone", "placement_policy", "DEFAULT"),
					resource.TestCheckResourceAttr(
						"cas_zone.my-zone", "tags.#", "3"),
				),
			},
			{
				Config: testAccCheckCASZoneUpdateConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCASZoneExists("cas_zone.my-zone"),
					resource.TestCheckResourceAttr(
						"cas_zone.my-zone", "description", "description my-cas-zone-update"),
					resource.TestCheckResourceAttr(
						"cas_zone.my-zone", "placement_policy", "BINPACK"),
					resource.TestCheckResourceAttr(
						"cas_zone.my-zone", "tags.#", "2"),
				),
			},
		},
	})
}

func testAccCheckCASZoneExists(n string) resource.TestCheckFunc {
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

func testAccCheckCASZoneDestroy(s *terraform.State) error {
	apiClient := testAccProviderCAS.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "cas_cloud_account_aws" {
			_, err := apiClient.CloudAccount.GetAwsCloudAccount(cloud_account.NewGetAwsCloudAccountParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'cas_cloud_account_aws' still exists with id %s", rs.Primary.ID)
			}
		}
		if rs.Type == "cas_zone" {
			_, err := apiClient.Location.GetZone(location.NewGetZoneParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'cas_zone' still exists with id %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckCASZoneConfig(rInt int) string {
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

 resource "cas_zone" "my-zone" {
    name = "my-cas-zone-%d"
    description = "description my-cas-zone"
    region_id = "${element(cas_cloud_account_aws.my-cloud-account.0.region_ids, 0)}"
    tags {
        key = "mykey"
        value = "myvalue"
    }
    tags {
        key = "foo"
        value = "bar"
    }
    tags {
        key = "faz"
        value = "baz"
    }
}`, id, secret, rInt)
}

func testAccCheckCASZoneUpdateConfig(rInt int) string {
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

	 resource "cas_zone" "my-zone" {
		name = "my-cas-zone-update-%d"
		description = "description my-cas-zone-update"
		region_id = "${element(cas_cloud_account_aws.my-cloud-account.0.region_ids, 0)}"
		placement_policy = "BINPACK"
		tags {
			key = "mykey"
			value = "myvalue"
		}
		tags {
			key = "foo"
			value = "bar"
		}
	}`, id, secret, rInt)
}

func TestAccCASZoneInvalidPlacementPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCASZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckCASZoneConfigInvalidPlacementPolicy(),
				ExpectError: regexp.MustCompile("\"placement_policy\" must be one of 'DEFAULT', 'SPREAD', 'BINPACK'"),
			},
		},
	})
}

func testAccCheckCASZoneConfigInvalidPlacementPolicy() string {
	return `
 resource "cas_zone" "my-zone" {
	name = "my-cas-zone-update"
	description = "description my-cas-zone-update"
	region_id = "fakeid"
	placement_policy = "INVALID"
}`
}
