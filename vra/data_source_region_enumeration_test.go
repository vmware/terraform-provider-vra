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

func TestAccVRARegionEnumeration_Basic(t *testing.T) {
	dcname := os.Getenv("VRA_VSPHERE_DATACOLLECTOR_NAME")
	datasourceName := "data.vra_region_enumeration.dc_regions"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVsphere(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRARegionEnumerationConfig(dcname),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRARegionEnumerationExists(datasourceName),
				),
			},
		},
	})
}

func testAccCheckVRARegionEnumerationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no datacollector ID is set")
		}

		return nil
	}
}

func testAccCheckVRARegionEnumerationConfig(dcname string) string {
	username := os.Getenv("VRA_VSPHERE_USERNAME")
	password := os.Getenv("VRA_VSPHERE_PASSWORD")
	hostname := os.Getenv("VRA_VSPHERE_HOSTNAME")
	return fmt.Sprintf(`
	data "vra_data_collector" "dc" {
		name = "%s"
	}

	data "vra_region_enumeration" "dc_regions" {
	  username    = "%s"
	  password    = "%s"
	  hostname    = "%s"
	  dc_id       = data.vra_data_collector.dc.id
	}`, dcname, username, password, hostname)
}
