// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVRADataCollector_Basic(t *testing.T) {
	dcname := os.Getenv("VRA_VSPHERE_DATACOLLECTOR_NAME")
	// resourceName := "vra_data_collector.dc"
	datasourceName := "data.vra_data_collector.dc"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVsphere(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckVRADataCollectorConfig("foobar123"),
				ExpectError: regexp.MustCompile("vra_data_collector \"foobar123\" not found"),
			},
			{
				Config: testAccCheckVRADataCollectorConfig(dcname),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRADataCollectorExists(datasourceName),
					resource.TestCheckResourceAttr(datasourceName, "name", dcname),
				),
			},
		},
	})
}

func testAccCheckVRADataCollectorExists(n string) resource.TestCheckFunc {
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

func testAccCheckVRADataCollectorConfig(dcname string) string {
	return fmt.Sprintf(`
	data "vra_data_collector" "dc" {
		name = "%s"
	}`, dcname)
}
