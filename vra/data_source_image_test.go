// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
)

func TestAccDataSourceVRAImageBasic(t *testing.T) {
	rInt := acctest.RandInt()
	dataSourceName := "data.vra_image.ubuntu"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAImageDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRAImageNoneConfig(rInt),
				ExpectError: regexp.MustCompile("vra_image filter did not match any images"),
			},
			{
				Config:      testAccDataSourceVRAImageManyConfig(rInt),
				ExpectError: regexp.MustCompile("vra_image must filter to a single image"),
			},
			{
				Config: testAccDataSourceVRAImageOneConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "description", "Canonical, Ubuntu, 16.04 LTS, amd64 xenial image build on 2017-10-26"),
					resource.TestCheckResourceAttr(dataSourceName, "external_id", "ami-da05a4a0"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "ubuntu/images/hvm-ssd/ubuntu-xenial-16.04-amd64-server-20171026.1"),
					resource.TestCheckResourceAttr(dataSourceName, "private", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "region", "us-east-1"),
				),
			},
		},
	})
}

func testAccDataSourceVRAImageBaseConfig(rInt int) string {
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
	}
	`, rInt, id, secret)
}

func testAccDataSourceVRAImageNoneConfig(rInt int) string {
	return testAccDataSourceVRAImageBaseConfig(rInt) + `
		data "vra_image" "ubuntu" {
			filter = "name eq 'foobar ubuntu/images/hvm-ssd/ubuntu-xenial-16.04-amd64-server-20171026.1'"
		}`
}

func testAccDataSourceVRAImageManyConfig(rInt int) string {
	return testAccDataSourceVRAImageBaseConfig(rInt) + `
		data "vra_image" "ubuntu" {
			filter = "name eq 'ubuntu/images/hvm-ssd/ubuntu-xenial-16.04-amd64-server-20171026.1'"
		}`
}

func testAccDataSourceVRAImageOneConfig(rInt int) string {
	return testAccDataSourceVRAImageBaseConfig(rInt) + `
		data "vra_image" "ubuntu" {
			filter = "name eq 'ubuntu/images/hvm-ssd/ubuntu-xenial-16.04-amd64-server-20171026.1' and externalRegionId eq 'us-east-1'"
		}`
}

func testAccCheckVRAImageDestroy(s *terraform.State) error {
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
