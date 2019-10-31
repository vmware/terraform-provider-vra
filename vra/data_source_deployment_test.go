package vra

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"os"
	"regexp"
	"testing"
)

func TestAccDataSourceVRADeployment(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName1 := "vra_deployment.this"
	dataSourceName1 := "data.vra_deployment.this"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDeploymentDataSource(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRADeploymentNoneConfig(rInt),
				ExpectError: regexp.MustCompile("deployment invalid-name not found"),
			},
			{
				Config: testAccDataSourceVRADeploymentOneConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
					resource.TestCheckResourceAttrPair(resourceName1, "project_id", dataSourceName1, "project_id"),
					resource.TestCheckResourceAttrPair(resourceName1, "catalog_item_id", dataSourceName1, "catalog_item_id"),
				),
			},
		},
	})
}

func testAccDataSourceVRADeployment(rInt int) string {
	// Need valid details since this is creating a real deployment
	catalogItemName := os.Getenv("VRA_CATALOG_ITEM_NAME")
	projectName := os.Getenv("VRA_PROJECT_NAME")

	return fmt.Sprintf(`
	data "vra_project" "this" {
	  name = "%s"
	}

	data "vra_catalog_item" "this" {
	  name = "%s"
	}

	resource "vra_deployment" "this" {
	  name        = "test-deployment-%d"
	  description = "terraform test deployment"

	  catalog_item_id = data.vra_catalog_item.this.id
	  project_id = data.vra_project.this.id
	}`, projectName, catalogItemName, rInt)
}

func testAccDataSourceVRADeploymentNoneConfig(rInt int) string {
	return testAccDataSourceVRADeployment(rInt) + `
		data "vra_deployment" "this" {
			name = "invalid-name"
		}`
}

func testAccDataSourceVRADeploymentOneConfig(rInt int) string {
	return testAccDataSourceVRADeployment(rInt) + `
		data "vra_deployment" "this" {
			name = "${vra_deployment.this.name}"
		}`
}
