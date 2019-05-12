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
)

func TestAccCASCloudAccountAWS_Basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCASCloudAccountAWSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCASCloudAccountAWSConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCASCloudAccountAWSExists("cas_cloud_account_aws.my-cloud-account"),
					resource.TestMatchResourceAttr(
						"cas_cloud_account_aws.my-cloud-account", "name", regexp.MustCompile("^my-cloud-account-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"cas_cloud_account_aws.my-cloud-account", "description", "test cloud account"),
					resource.TestCheckResourceAttr(
						"cas_cloud_account_aws.my-cloud-account", "access_key", os.Getenv("CAS_AWS_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr(
						"cas_cloud_account_aws.my-cloud-account", "secret_key", os.Getenv("CAS_AWS_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr(
						"cas_cloud_account_aws.my-cloud-account", "tags.#", "2"),
				),
			},
			{
				Config: testAccCheckCASCloudAccountAWSUpdateConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCASCloudAccountAWSExists("cas_cloud_account_aws.my-cloud-account"),
					resource.TestCheckResourceAttr(
						"cas_cloud_account_aws.my-cloud-account", "description", "your test cloud account"),
					resource.TestCheckResourceAttr(
						"cas_cloud_account_aws.my-cloud-account", "access_key", os.Getenv("CAS_AWS_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr(
						"cas_cloud_account_aws.my-cloud-account", "secret_key", os.Getenv("CAS_AWS_SECRET_ACCESS_KEY")),
				),
			},
		},
	})
}
func TestAccCASCloudAccountAWS_Duplicate(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCASCloudAccountAWSDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckCASCloudAccountAWSConfigDuplicateRegion(rInt),
				ExpectError: regexp.MustCompile("Specified regions are not unique"),
			},
		},
	})
}

func testAccCheckCASCloudAccountAWSExists(n string) resource.TestCheckFunc {
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

func testAccCheckCASCloudAccountAWSDestroy(s *terraform.State) error {
	apiClient := testAccProviderCAS.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cas_cloud_account_aws" {
			continue
		}

		_, err := apiClient.CloudAccount.GetAwsCloudAccount(cloud_account.NewGetAwsCloudAccountParams().WithID(rs.Primary.ID))
		if err == nil {
			return err
		}
	}

	return nil
}

func testAccCheckCASCloudAccountAWSConfig(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	id := os.Getenv("CAS_AWS_ACCESS_KEY_ID")
	secret := os.Getenv("CAS_AWS_SECRET_ACCESS_KEY")
	return fmt.Sprintf(`
resource "cas_cloud_account_aws" "my-cloud-account" {
	name = "my-cloud-account-%d"
	description = "test cloud account"
	access_key = "%s"
	secret_key = "%s"
	regions = ["us-east-1"]
	tags {
		key = "foo"
		value = "bar"
	}
	tags {
		key = "where"
		value = "waldo"
	}
 }`, rInt, id, secret)
}

func testAccCheckCASCloudAccountAWSUpdateConfig(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	id := os.Getenv("CAS_AWS_ACCESS_KEY_ID")
	secret := os.Getenv("CAS_AWS_SECRET_ACCESS_KEY")
	return fmt.Sprintf(`
resource "cas_cloud_account_aws" "my-cloud-account" {
	name = "my-cloud-account-%d"
	description = "your test cloud account"
	access_key = "%s"
	secret_key = "%s"
	regions = ["us-east-1"]
 }`, rInt, id, secret)
}

func testAccCheckCASCloudAccountAWSConfigDuplicateRegion(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	id := os.Getenv("CAS_AWS_ACCESS_KEY_ID")
	secret := os.Getenv("CAS_AWS_SECRET_ACCESS_KEY")
	return fmt.Sprintf(`
resource "cas_cloud_account_aws" "my-cloud-account" {
	name = "my-cloud-account-%d"
	description = "test cloud account"
	access_key = "%s"
	secret_key = "%s"
	regions = ["us-east-1", "us-west-1", "us-east-1"]
 }`, rInt, id, secret)
}
