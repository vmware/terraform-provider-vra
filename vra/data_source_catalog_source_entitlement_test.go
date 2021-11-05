package vra

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVRACatalogSourceEntitlement_Valid(t *testing.T) {
	rInt := acctest.RandInt()
	resource1 := "vra_catalog_source_entitlement.this"
	dataSource1 := "data.vra_catalog_source_entitlement.with_id"
	dataSource2 := "data.vra_catalog_source_entitlement.with_catalog_source_id"
	catalogSource := "vra_catalog_source_blueprint.this"
	project := "vra_project.this"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckBlueprint(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRACatalogSourceEntitlementFoundConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(resource1, "id", dataSource1, "id"),
					resource.TestCheckResourceAttrPair(resource1, "definition.0.id", dataSource1, "definition.0.id"),
					resource.TestCheckResourceAttrPair(catalogSource, "id", dataSource1, "definition.0.id"),
					resource.TestCheckResourceAttrPair(catalogSource, "name", dataSource1, "definition.0.name"),
					resource.TestCheckResourceAttrPair(project, "id", dataSource1, "project_id"),
					resource.TestCheckResourceAttrPair(resource1, "id", dataSource2, "id"),
					resource.TestCheckResourceAttrPair(resource1, "definition.0.id", dataSource2, "definition.0.id"),
					resource.TestCheckResourceAttrPair(catalogSource, "id", dataSource2, "definition.0.id"),
					resource.TestCheckResourceAttrPair(catalogSource, "name", dataSource2, "definition.0.name"),
					resource.TestCheckResourceAttrPair(project, "id", dataSource2, "project_id"),
				),
			},
		},
	})
}

func testAccCheckVRACatalogSourceEntitlementFoundConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "vra_project" "this" {
	  name = "tf-test-project-%d"
	}

	resource "vra_catalog_source_blueprint" "this" {
 	  name       = "tf-test-catalog-source-%d"
  	  project_id = vra_project.this.id
	}

	resource "vra_catalog_source_entitlement" "this" {
  	  catalog_source_id = vra_catalog_source_blueprint.this.id
	  project_id        = vra_project.this.id
	}

	data "vra_catalog_source_entitlement" "with_id" {
 	  id = vra_catalog_source_entitlement.this.id
	  project_id = vra_catalog_source_entitlement.this.project_id
	}

	data "vra_catalog_source_entitlement" "with_catalog_source_id" {
 	  catalog_source_id = vra_catalog_source_entitlement.this.catalog_source_id
	  project_id = vra_catalog_source_entitlement.this.project_id
	}`, rInt, rInt)
}
