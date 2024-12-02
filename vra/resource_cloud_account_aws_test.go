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

func TestAccVRACloudAccountAWS_Basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRACloudAccountAWSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRACloudAccountAWSConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRACloudAccountAWSExists("vra_cloud_account_aws.my-cloud-account"),
					resource.TestMatchResourceAttr(
						"vra_cloud_account_aws.my-cloud-account", "name", regexp.MustCompile("^my-cloud-account-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_aws.my-cloud-account", "description", "test cloud account"),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_aws.my-cloud-account", "access_key", os.Getenv("VRA_AWS_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_aws.my-cloud-account", "secret_key", os.Getenv("VRA_AWS_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_aws.my-cloud-account", "tags.#", "2"),
				),
			},
			{
				Config: testAccCheckVRACloudAccountAWSUpdateConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRACloudAccountAWSExists("vra_cloud_account_aws.my-cloud-account"),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_aws.my-cloud-account", "description", "your test cloud account"),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_aws.my-cloud-account", "access_key", os.Getenv("VRA_AWS_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_aws.my-cloud-account", "secret_key", os.Getenv("VRA_AWS_SECRET_ACCESS_KEY")),
				),
			},
		},
	})
}
func TestAccVRACloudAccountAWS_Duplicate(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRACloudAccountAWSDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckVRACloudAccountAWSConfigDuplicateRegion(rInt),
				ExpectError: regexp.MustCompile("Specified regions are not unique"),
			},
		},
	})
}

func testAccCheckVRACloudAccountAWSExists(n string) resource.TestCheckFunc {
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

func testAccCheckVRACloudAccountAWSDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vra_cloud_account_aws" {
			continue
		}

		_, err := apiClient.CloudAccount.GetAwsCloudAccount(cloud_account.NewGetAwsCloudAccountParams().WithID(rs.Primary.ID))
		if err == nil {
			return err
		}
	}

	return nil
}

func testAccCheckVRACloudAccountAWSConfig(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	id := os.Getenv("VRA_AWS_ACCESS_KEY_ID")
	secret := os.Getenv("VRA_AWS_SECRET_ACCESS_KEY")
	return fmt.Sprintf(`
resource "vra_cloud_account_aws" "my-cloud-account" {
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

func testAccCheckVRACloudAccountAWSUpdateConfig(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	id := os.Getenv("VRA_AWS_ACCESS_KEY_ID")
	secret := os.Getenv("VRA_AWS_SECRET_ACCESS_KEY")
	return fmt.Sprintf(`
resource "vra_cloud_account_aws" "my-cloud-account" {
	name = "my-cloud-account-%d"
	description = "your test cloud account"
	access_key = "%s"
	secret_key = "%s"
	regions = ["us-east-1"]
 }`, rInt, id, secret)
}

func testAccCheckVRACloudAccountAWSConfigDuplicateRegion(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	id := os.Getenv("VRA_AWS_ACCESS_KEY_ID")
	secret := os.Getenv("VRA_AWS_SECRET_ACCESS_KEY")
	return fmt.Sprintf(`
resource "vra_cloud_account_aws" "my-cloud-account" {
	name = "my-cloud-account-%d"
	description = "test cloud account"
	access_key = "%s"
	secret_key = "%s"
	regions = ["us-east-1", "us-west-1", "us-east-1"]
 }`, rInt, id, secret)
}
