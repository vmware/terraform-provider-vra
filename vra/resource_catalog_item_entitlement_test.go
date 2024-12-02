// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"

	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/vmware/vra-sdk-go/pkg/client/catalog_items"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVRACatalogItemEntitlement_Valid(t *testing.T) {
	rInt := acctest.RandInt()
	resource1 := "vra_catalog_item_entitlement.this"
	blueprintCatalogSource := "vra_catalog_item_blueprint.this"
	project := "vra_project.this"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckBlueprint(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRACatalogItemEntitlementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRACatalogItemEntitlementValidContentConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(resource1, "definition.0.id", blueprintCatalogSource, "id"),
					resource.TestCheckResourceAttrPair(resource1, "definition.0.name", blueprintCatalogSource, "name"),
					resource.TestCheckResourceAttrPair(resource1, "project_id", project, "id"),
				),
			},
		},
	})
}

func testAccCheckVRACatalogItemEntitlementDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vra_catalog_item_entitlement" {
			continue
		}

		_, err := apiClient.CatalogItems.GetCatalogItemUsingGET5(catalog_items.NewGetCatalogItemUsingGET5Params().WithID(strfmt.UUID(rs.Primary.ID)))
		if err == nil {
			return fmt.Errorf("resource 'vra_catalog_item_entitlement' still exists with id %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckVRACatalogItemEntitlementValidContentConfig(rInt int) string {
	return fmt.Sprintf(`
        resource "vra_project" "this" {
          name = "tf-test-project-%d"
        }

        resource "vra_catalog_item_blueprint" "this" {
          name       = "tf-test-catalog-item-%d"
          project_id = vra_project.this.id
        }

        resource "vra_catalog_item_entitlement" "this" {
          catalog_source_id = vra_catalog_item_blueprint.this.id
          project_id        = vra_project.this.id
        }`, rInt, rInt)
}
