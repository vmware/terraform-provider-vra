// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVRACloudAccountAWS(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName1 := "vra_cloud_account_aws.my-cloud-account"
	dataSourceName1 := "data.vra_cloud_account_aws.test-cloud-account"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckAWS(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVRACloudAccountAWS(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "access_key", dataSourceName1, "access_key"),
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
				),
			},
		},
	})
}

func testAccDataSourceVRACloudAccountAWS(rInt int) string {
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

	 data "vra_cloud_account_aws" "test-cloud-account" {
     name = "${vra_cloud_account_aws.my-cloud-account.name}"
	 }`, rInt, id, secret)
}
