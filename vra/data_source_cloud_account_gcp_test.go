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
)

func TestAccDataSourceVRACloudAccountGCP(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName1 := "vra_cloud_account_gcp.my-cloud-account"
	dataSourceName1 := "data.vra_cloud_account_gcp.test-cloud-account"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckGCP(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRACloudAccountGCPNotFound(rInt),
				ExpectError: regexp.MustCompile("cloud account foobar not found"),
			},
			{
				Config: testAccDataSourceVRACloudAccountGCP(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
					resource.TestCheckResourceAttrPair(resourceName1, "client_email", dataSourceName1, "client_email"),
					resource.TestCheckResourceAttrPair(resourceName1, "private_key_id", dataSourceName1, "private_key_id"),
					resource.TestCheckResourceAttrPair(resourceName1, "project_id", dataSourceName1, "project_id"),
					resource.TestCheckResourceAttrPair(resourceName1, "regions", dataSourceName1, "regions"),
					resource.TestCheckResourceAttrPair(resourceName1, "tags", dataSourceName1, "tags"),
				),
			},
		},
	})
}

func testAccDataSourceVRACloudAccountGCPBase(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	clientEmail := os.Getenv("VRA_GCP_CLIENT_EMAIL")
	privateKeyID := os.Getenv("VRA_GCP_PRIVATE_KEY_ID")
	privateKey := os.Getenv("VRA_GCP_PRIVATE_KEY")
	projectID := os.Getenv("VRA_GCP_PROJECT_ID")
	return fmt.Sprintf(`
resource "vra_cloud_account_gcp" "my-cloud-account" {
	name = "my-cloud-account-%d"
	description = "test cloud account"
	client_email = "%s"
	private_key_id = "%s"
	private_key = "%s"
	project_id = "%s"
	regions = ["us-west2"]
	tags {
		key = "foo"
		value = "bar"
	}
	tags {
		key = "where"
		value = "waldo"
	}
 }`, rInt, clientEmail, privateKeyID, privateKey, projectID)
}

func testAccDataSourceVRACloudAccountGCPNotFound(rInt int) string {
	return testAccDataSourceVRACloudAccountGCPBase(rInt) + `
	data "vra_cloud_account_gcp" "test-cloud-account" {
		name = "foobar"
	}`
}

func testAccDataSourceVRACloudAccountGCP(rInt int) string {
	return testAccDataSourceVRACloudAccountGCPBase(rInt) + `
	data "vra_cloud_account_gcp" "test-cloud-account" {
		name = "${vra_cloud_account_gcp.my-cloud-account.name}"
	}`
}
