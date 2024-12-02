// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"os"
	"regexp"
	"testing"
)

func TestAccDataSourceVRACatalogItem(t *testing.T) {
	dataSource := "data.vra_catalog_item.this"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCatalogItem(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRACatalogItemNotFound(),
				ExpectError: regexp.MustCompile("catalog item foobar not found"),
			},
			{
				Config: testAccDataSourceVRACatalogItemFound(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSource, "name", os.Getenv("VRA_CATALOG_ITEM_NAME")),
				),
			},
		},
	})
}

func testAccDataSourceVRACatalogItemBase(catalogItemName string) string {
	return fmt.Sprintf(`
	data "vra_catalog_item" "this" {
      name = "%s"
	  expand_versions = true
	}`, catalogItemName)
}

func testAccDataSourceVRACatalogItemNotFound() string {
	return testAccDataSourceVRACatalogItemBase("foobar")
}

func testAccDataSourceVRACatalogItemFound() string {
	// Need valid catalog item name since this is looking for real catalog item
	catalogItemName := os.Getenv("VRA_CATALOG_ITEM_NAME")
	return testAccDataSourceVRACatalogItemBase(catalogItemName)
}
