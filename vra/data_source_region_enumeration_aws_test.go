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

func TestAccVRARegionEnumerationAWS_Basic(t *testing.T) {
	dataSourceName := "data.vra_region_enumeration_aws.this"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckAWS(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRARegionEnumerationAWSConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRARegionEnumerationAWSExists(dataSourceName),
				),
			},
		},
	})
}

func testAccCheckVRARegionEnumerationAWSExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("access_key is not set")
		}

		if rs.Primary.Attributes["regions.#"] == "0" {
			return fmt.Errorf("regions are not set")
		}

		return nil
	}
}

func testAccCheckVRARegionEnumerationAWSConfig() string {
	id := os.Getenv("VRA_AWS_ACCESS_KEY_ID")
	secret := os.Getenv("VRA_AWS_SECRET_ACCESS_KEY")
	return fmt.Sprintf(`
data "vra_region_enumeration_aws" "this" {
	access_key = "%s"
	secret_key = "%s"
 }`, id, secret)
}
