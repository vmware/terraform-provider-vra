package vra

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"testing"
)

func TestAccDataSourceVRAStorageProfile(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName1 := "vra_storage_profile.this"
	dataSourceName1 := "data.vra_storage_profile.this"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckStorageProfile(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAStorageProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRAStorageProfileNotFound(rInt),
				ExpectError: regexp.MustCompile("vra_storage_profile filter did not match any storage profile"),
			},
			// TBD: Enable filter by name once this is fixed https://jira.eng.vmware.com/browse/VCOM-13947
			// {
			// 	Config: testAccDataSourceVRAStorageProfileNameFilter(rInt),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
			// 		resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
			// 		resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
			// 		resource.TestCheckResourceAttrPair(resourceName1, "default_item", dataSourceName1, "default_item"),
			// 		resource.TestCheckResourceAttrPair(resourceName1, "external_region_id", dataSourceName1, "region_id"),
			// 	),
			// },
			{
				Config: testAccDataSourceVRAStorageProfileExternalRegionIDFilter(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
					resource.TestCheckResourceAttrPair(resourceName1, "default_item", dataSourceName1, "default_item"),
					resource.TestCheckResourceAttrPair(resourceName1, "external_region_id", dataSourceName1, "external_region_id"),
				),
			},
			{
				Config: testAccDataSourceVRAStorageProfileByID(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
					resource.TestCheckResourceAttrPair(resourceName1, "default_item", dataSourceName1, "default_item"),
					resource.TestCheckResourceAttrPair(resourceName1, "external_region_id", dataSourceName1, "external_region_id"),
				),
			},
		},
	})
}

func testAccDataSourceVRAStorageProfileNotFound(rInt int) string {
	return testAccCheckVRAStorageProfileConfig(rInt) + fmt.Sprintf(`
	data "vra_storage_profile" "this" {
		filter = "externalRegionId eq 'foobar'"
	}`)
}

// func testAccDataSourceVRAStorageProfileNameFilter(rInt int) string {
// 	return testAccCheckVRAStorageProfileConfig(rInt) + fmt.Sprintf(`
// 	data "vra_storage_profile" "this" {
// 		filter = "name eq '${vra_storage_profile.my-storage-profile.name}'"
// 	}`)
// }

func testAccDataSourceVRAStorageProfileExternalRegionIDFilter(rInt int) string {
	return testAccCheckVRAStorageProfileConfig(rInt) + fmt.Sprintf(`
	data "vra_storage_profile" "this" {
		filter = "externalRegionId eq '${data.vra_region.this.id}'"
	}`)
}

func testAccDataSourceVRAStorageProfileByID(rInt int) string {
	return testAccCheckVRAStorageProfileConfig(rInt) + fmt.Sprintf(`
	data "vra_storage_profile" "this" {
		id = vra_storage_profile.this.id
	}`)
}
