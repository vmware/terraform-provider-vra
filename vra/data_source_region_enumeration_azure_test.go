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

func TestAccVRARegionEnumerationAzure_Basic(t *testing.T) {
	dataSourceName := "data.vra_region_enumeration_azure.this"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckAzure(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRARegionEnumerationAzureConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRARegionEnumerationAzureExists(dataSourceName),
				),
			},
		},
	})
}

func testAccCheckVRARegionEnumerationAzureExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("application_id is not set")
		}

		if rs.Primary.Attributes["regions.#"] == "0" {
			return fmt.Errorf("regions are not set")
		}

		return nil
	}
}

func testAccCheckVRARegionEnumerationAzureConfig() string {
	// Need valid credentials since this is creating a real cloud account
	subscriptionID := os.Getenv("VRA_ARM_SUBSCRIPTION_ID")
	tenantID := os.Getenv("VRA_ARM_TENANT_ID")
	applicationID := os.Getenv("VRA_ARM_CLIENT_APP_ID")
	applicationKey := os.Getenv("VRA_ARM_CLIENT_APP_KEY")
	return fmt.Sprintf(`
data "vra_region_enumeration_azure" "this" {
	application_id = "%s"
	application_key = "%s"
	subscription_id = "%s"
	tenant_id = "%s"
 }`, applicationID, applicationKey, subscriptionID, tenantID)
}
