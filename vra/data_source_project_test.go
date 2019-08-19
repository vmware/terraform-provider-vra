package vra

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceVRAProject(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName1 := "vra_project.my-project"
	dataSourceName1 := "data.vra_project.test-project"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCas(t) },
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
