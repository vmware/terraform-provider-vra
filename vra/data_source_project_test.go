package vra

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceVRAProject(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName1 := "vra_project.my-project"
	dataSourceName1 := "data.vra_project.test-project"
	dataSourceName2 := "data.vra_project.this"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVra(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRAProjectNoneConfig(rInt),
				ExpectError: regexp.MustCompile("project invalid-name not found"),
			},
			{
				Config: testAccDataSourceVRAProjectOneConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
				),
			},
			{
				Config: testAccDataSourceVRAProjectIdConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName2, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName2, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName2, "name"),
				),
			},
		},
	})
}

func testAccDataSourceVRAProject(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	return fmt.Sprintf(`
	resource "vra_project" "my-project" {
		name = "my-project-%d"
		description = "test project"
	 }`, rInt)
}

func testAccDataSourceVRAProjectNoneConfig(rInt int) string {
	return testAccDataSourceVRAProject(rInt) + `
		data "vra_project" "test-project" {
			name = "invalid-name"
		}`
}

func testAccDataSourceVRAProjectOneConfig(rInt int) string {
	return testAccDataSourceVRAProject(rInt) + `
		data "vra_project" "test-project" {
			name = "${vra_project.my-project.name}"
		}`
}

func testAccDataSourceVRAProjectIdConfig(rInt int) string {
	return testAccDataSourceVRAProject(rInt) + `
		data "vra_project" "this" {
			id = "${vra_project.my-project.id}"
		}`
}
