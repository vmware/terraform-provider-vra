// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vra-sdk-go/pkg/client/policies"
)

func TestAccVRAContentSharingPolicyBasic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "vra_content_sharing_policy.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckContentSharingPolicy(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAContentSharingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRAContentSharingPolicyConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAContentSharingPolicyExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "name", "content-sharing-policy-"+strconv.Itoa(rInt)),
					resource.TestCheckResourceAttr(
						resourceName, "description", "Content Sharing Policy "+strconv.Itoa(rInt)),
					resource.TestCheckResourceAttr(
						resourceName, "project_id", os.Getenv("VRA_PROJECT_ID")),
					resource.TestCheckResourceAttr(
						resourceName, "catalog_item_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "catalog_item_ids.*", os.Getenv("VRA_CATALOG_ITEM_ID")),
					resource.TestCheckResourceAttr(
						resourceName, "catalog_source_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "catalog_source_ids.*", os.Getenv("VRA_CATALOG_SOURCE_ID")),
				),
			},
			{
				Config: testAccCheckVRAContentSharingPolicyUpdateConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAContentSharingPolicyExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "name", "content-sharing-policy-"+strconv.Itoa(rInt)),
					resource.TestCheckResourceAttr(
						resourceName, "description", "Updated Content Sharing Policy "+strconv.Itoa(rInt)),
					resource.TestCheckResourceAttr(
						resourceName, "project_id", os.Getenv("VRA_PROJECT_ID")),
					resource.TestCheckResourceAttr(
						resourceName, "catalog_item_ids.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "catalog_source_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "catalog_source_ids.*", os.Getenv("VRA_CATALOG_SOURCE_ID")),
				),
			},
		},
	})
}

func testAccCheckVRAContentSharingPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no content sharing policy id is set")
		}

		return nil
	}
}

func testAccCheckVRAContentSharingPolicyDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "vra_content_sharing_policy" {
			if _, err := apiClient.Policies.GetPolicyUsingGET5(policies.NewGetPolicyUsingGET5Params().WithID(strfmt.UUID(rs.Primary.ID))); err == nil {
				return fmt.Errorf("Resource 'ra_content_sharing_policy' still exists with id %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckVRAContentSharingPolicyConfig(rInt int) string {
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
    }`, rInt, rInt, projectID, catalogItemID, catalogSourceID)
}

func testAccCheckVRAContentSharingPolicyUpdateConfig(rInt int) string {
	projectID := os.Getenv("VRA_PROJECT_ID")
	catalogSourceID := os.Getenv("VRA_CATALOG_SOURCE_ID")
	return fmt.Sprintf(`
	resource "vra_content_sharing_policy" "test" {
		name = "content-sharing-policy-%d"
		description = "Updated Content Sharing Policy %d"
		project_id = "%s"
		catalog_item_ids = []
		catalog_source_ids = ["%s"]
	}`, rInt, rInt, projectID, catalogSourceID)
}
