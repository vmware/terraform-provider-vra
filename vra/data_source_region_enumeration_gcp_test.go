// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVRARegionEnumerationGCP_Basic(t *testing.T) {
	dataSourceName := "data.vra_region_enumeration_gcp.this"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckGCP(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRARegionEnumerationGCPConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRARegionEnumerationGCPExists(dataSourceName),
				),
			},
		},
	})
}

func testAccCheckVRARegionEnumerationGCPExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("private_key_id is not set")
		}

		if rs.Primary.Attributes["regions.#"] == "0" {
			return fmt.Errorf("regions are not set")
		}

		return nil
	}
}

func testAccCheckVRARegionEnumerationGCPConfig() string {
	clientEmail := os.Getenv("VRA_GCP_CLIENT_EMAIL")
	projectID := os.Getenv("VRA_GCP_PROJECT_ID")
	privateKeyID := os.Getenv("VRA_GCP_PRIVATE_KEY_ID")
	privateKey := os.Getenv("VRA_GCP_PRIVATE_KEY")
	return fmt.Sprintf(`
data "vra_region_enumeration_gcp" "this" {
	client_email = "%s"
	project_id = "%s"
	private_key_id = "%s"
	private_key = "%s"
 }`, clientEmail, projectID, privateKeyID, privateKey)
}
