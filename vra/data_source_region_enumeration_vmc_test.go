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

func TestAccVRARegionEnumerationVMC_Basic(t *testing.T) {
	dcName := os.Getenv("VRA_VMC_DATA_COLLECTOR_NAME")
	datasourceName := "data.vra_region_enumeration_vmc.this"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVMC(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRARegionEnumerationVMCConfig(dcName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRARegionEnumerationVMCExists(datasourceName),
				),
			},
		},
	})
}

func testAccCheckVRARegionEnumerationVMCExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("vsphere_hostname is not set")
		}

		if rs.Primary.Attributes["regions.#"] == "0" {
			return fmt.Errorf("regions are not set")
		}

		return nil
	}
}

func testAccCheckVRARegionEnumerationVMCConfig(_ string) string {
	apiToken := os.Getenv("VRA_VMC_API_TOKEN")
	sddcName := os.Getenv("VRA_VMC_SDDC_NAME")
	vCenterHostName := os.Getenv("VRA_VMC_VCENTER_HOSTNAME")
	vCenterUserName := os.Getenv("VRA_VMC_VCENTER_USERNAME")
	vCenterPassword := os.Getenv("VRA_VMC_VCENTER_PASSWORD")
	//nsxHostName := os.Getenv("VRA_VMC_NSX_HOSTNAME")
	dcName := os.Getenv("VRA_VMC_DATA_COLLECTOR_NAME")
	return fmt.Sprintf(`
	data "vra_data_collector" "dc" {
		name = "%s"
	}

	data "vra_region_enumeration_vmc" "this" {
	  api_token = "%s"
	  sddc_name = "%s"

	  vcenter_username    = "%s"
	  vcenter_password    = "%s"
	  vcenter_hostname    = "%s"
	  dc_id               = data.vra_data_collector.dc.id

	  accept_self_signed_cert = true
	}`, dcName, apiToken, sddcName, vCenterUserName, vCenterPassword, vCenterHostName)
}
