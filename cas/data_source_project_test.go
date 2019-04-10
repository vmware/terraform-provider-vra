package cas

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceCASProject(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName1 := "cas_project.my-project"
	dataSourceName1 := "data.cas_project.test-project"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCasProject(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceCASProjectNoneConfig(rInt),
				ExpectError: regexp.MustCompile("project invalid-name not found"),
			},
			{
				Config: testAccDataSourceCASProjectOneConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
				),
			},
		},
	})
}

func testAccDataSourceCASProject(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	return fmt.Sprintf(`
	resource "cas_project" "my-project" {
		name = "my-project-%d"
		description = "test project"
	 }`, rInt)
}

func testAccDataSourceCASProjectNoneConfig(rInt int) string {
	return testAccDataSourceCASProject(rInt) + `
		data "cas_project" "test-project" {
			name = "invalid-name"
		}`
}

func testAccDataSourceCASProjectOneConfig(rInt int) string {
	return testAccDataSourceCASProject(rInt) + `
		data "cas_project" "test-project" {
			name = "${cas_project.my-project.name}"
		}`
}
