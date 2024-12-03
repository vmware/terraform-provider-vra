// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"

	"github.com/vmware/vra-sdk-go/pkg/client/catalog_sources"

	"testing"

	"github.com/go-openapi/strfmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVRACatalogSourceBlueprint_Valid(t *testing.T) {
	rInt := acctest.RandInt()
	resource1 := "vra_catalog_source_blueprint.this"
	project := "vra_project.this"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckBlueprint(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRACatalogSourceBlueprintDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRACatalogSourceBlueprintValidContentConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(resource1, "project_id", project, "id"),
					resource.TestCheckResourceAttrPair(resource1, "config.sourceProjectId", project, "id"),
					resource.TestCheckResourceAttr(resource1, "name", "tf-test-catalog-source"),
					resource.TestCheckResourceAttr(resource1, "type_id", "com.vmw.blueprint"),
				),
			},
		},
	})
}

func testAccCheckVRACatalogSourceBlueprintDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vra_catalog_source_blueprint" {
			continue
		}

		_, err := apiClient.CatalogSources.GetUsingGET2(catalog_sources.NewGetUsingGET2Params().WithSourceID(strfmt.UUID(rs.Primary.ID)))
		if err == nil {
			return fmt.Errorf("resource 'vra_catalog_source_blueprint' still exists with id %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckVRACatalogSourceBlueprintValidContentConfig(_ int) string {
	return `
	resource "vra_project" "this" {
	  name = "tf-test-project"
	}

	resource "vra_catalog_source_blueprint" "this" {
 	  name       = "tf-test-catalog-source"
  	  project_id = vra_project.this.id
	}`
}
