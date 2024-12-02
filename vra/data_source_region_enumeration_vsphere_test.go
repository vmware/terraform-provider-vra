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

func TestAccVRARegionEnumerationVsphere_Basic(t *testing.T) {
	dcname := os.Getenv("VRA_VSPHERE_DATACOLLECTOR_NAME")
	datasourceName := "data.vra_region_enumeration_vsphere.this"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVsphere(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRARegionEnumerationVsphereConfig(dcname),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRARegionEnumerationVsphereExists(datasourceName),
				),
			},
		},
	})
}

func testAccCheckVRARegionEnumerationVsphereExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("hostname is not set")
		}

		if rs.Primary.Attributes["regions.#"] == "0" {
			return fmt.Errorf("regions are not set")
		}

		return nil
	}
}

func testAccCheckVRARegionEnumerationVsphereConfig(dcname string) string {
	username := os.Getenv("VRA_VSPHERE_USERNAME")
	password := os.Getenv("VRA_VSPHERE_PASSWORD")
	hostname := os.Getenv("VRA_VSPHERE_HOSTNAME")
	return fmt.Sprintf(`
	data "vra_data_collector" "dc" {
		name = "%s"
	}

	data "vra_region_enumeration_vsphere" "this" {
	  username    = "%s"
	  password    = "%s"
	  hostname    = "%s"
	  dc_id       = data.vra_data_collector.dc.id
	}`, dcname, username, password, hostname)
}
