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

func TestAccDataSourceVRARegionBasic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "vra_cloud_account_aws.my-cloud-account"
	dataSourceName := "data.vra_region.east-region"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRARegionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVRARegionConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "id", dataSourceName, "cloud_account_id"),
					resource.TestCheckResourceAttrPair(resourceName, "region_ids.0", dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "region", "us-east-1"),
				),
			},
		},
	})
}

func testAccDataSourceVRARegionConfig(rInt int) string {
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

	data "vra_region" "east-region" {
		cloud_account_id = "${vra_cloud_account_aws.my-cloud-account.id}"
		region = "us-east-1"
	}`, rInt, id, secret)
}

func TestAccDataSourceVRARegionInvalidRegion(t *testing.T) {
	rInt := acctest.RandInt()
	dataSourceName := "data.vra_region.east-region"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRARegionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRARegionInvalidConfig(rInt),
				ExpectError: regexp.MustCompile("region us-west-1 not found in cloud account"),
			},
			{
				Config: testAccDataSourceVRARegionConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "region", "us-east-1"),
				),
			},
		},
	})
}

func testAccDataSourceVRARegionInvalidConfig(rInt int) string {
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

	data "vra_region" "east-region" {
		cloud_account_id = "${vra_cloud_account_aws.my-cloud-account.id}"
		region = "us-west-1"
	}`, rInt, id, secret)
}

func testAccCheckVRARegionDestroy(s *terraform.State) error {
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
