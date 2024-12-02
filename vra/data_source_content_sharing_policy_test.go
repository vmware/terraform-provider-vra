// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"testing"
)

func TestAccDataSourceVRAContentSharingPolicy(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "vra_content_sharing_policy.test"
	dataSourceName := "data.vra_content_sharing_policy.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckContentSharingPolicy(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVRAContentSharingPolicyByID(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrPair(
						resourceName, "description", dataSourceName, "description"),
					resource.TestCheckResourceAttrPair(
						resourceName, "project_id", dataSourceName, "project_id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "catalog_item_ids.#", dataSourceName, "catalog_item_ids.#"),
					resource.TestCheckTypeSetElemAttrPair(
						resourceName, "catalog_item_ids.0", dataSourceName, "catalog_item_ids.0"),
					resource.TestCheckResourceAttrPair(
						resourceName, "catalog_source_ids.#", dataSourceName, "catalog_source_ids.#"),
					resource.TestCheckTypeSetElemAttrPair(
						resourceName, "catalog_source_ids.0", dataSourceName, "catalog_source_ids.0"),
				),
			},
			{
				Config: testAccDataSourceVRAContentSharingPolicyByName(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrPair(
						resourceName, "description", dataSourceName, "description"),
					resource.TestCheckResourceAttrPair(
						resourceName, "project_id", dataSourceName, "project_id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "catalog_item_ids.#", dataSourceName, "catalog_item_ids.#"),
					resource.TestCheckTypeSetElemAttrPair(
						resourceName, "catalog_item_ids.0", dataSourceName, "catalog_item_ids.0"),
					resource.TestCheckResourceAttrPair(
						resourceName, "catalog_source_ids.#", dataSourceName, "catalog_source_ids.#"),
					resource.TestCheckTypeSetElemAttrPair(
						resourceName, "catalog_source_ids.0", dataSourceName, "catalog_source_ids.0"),
				),
			},
		},
	})
}

func testAccDataSourceVRAContentSharingPolicy(rInt int) string {
	projectID := os.Getenv("VRA_PROJECT_ID")
	catalogItemID := os.Getenv("VRA_CATALOG_ITEM_ID")
	catalogSourceID := os.Getenv("VRA_CATALOG_SOURCE_ID")
	return fmt.Sprintf(`
	resource "vra_content_sharing_policy" "test" {
		name = "content-sharing-policy-%d"
		description = "Content Sharing Policy %d"
		project_id = "%s"
		catalog_item_ids = ["%s"]
		catalog_source_ids = ["%s"]
	}
	`, rInt, rInt, projectID, catalogItemID, catalogSourceID)
}

func testAccDataSourceVRAContentSharingPolicyByID(rInt int) string {
	return testAccDataSourceVRAContentSharingPolicy(rInt) + `
	data "vra_content_sharing_policy" "test" {
		id = vra_content_sharing_policy.test.id
	}`
}

func testAccDataSourceVRAContentSharingPolicyByName(rInt int) string {
	return testAccDataSourceVRAContentSharingPolicy(rInt) + `
	data "vra_content_sharing_policy" "test" {
		name = vra_content_sharing_policy.test.name
	}`
}
