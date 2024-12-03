// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVRACloudAccountVsphere(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName1 := "vra_cloud_account_vsphere.my-cloud-account"
	dataSourceName1 := "data.vra_cloud_account_vsphere.test-cloud-account"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVsphere(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRACloudAccountVsphereNotFound(rInt),
				ExpectError: regexp.MustCompile("cloud account foobar not found"),
			},
			{
				Config: testAccDataSourceVRACloudAccountVsphere(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
					resource.TestCheckResourceAttrPair(resourceName1, "hostname", dataSourceName1, "hostname"),
					resource.TestCheckResourceAttrPair(resourceName1, "dc_id", dataSourceName1, "dc_id"),
				),
			},
		},
	})
}

func testAccDataSourceVRACloudAccountVsphereBase(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	username := os.Getenv("VRA_VSPHERE_USERNAME")
	password := os.Getenv("VRA_VSPHERE_PASSWORD")
	hostname := os.Getenv("VRA_VSPHERE_HOSTNAME")
	dcname := os.Getenv("VRA_VSPHERE_DATACOLLECTOR_NAME")
	return fmt.Sprintf(`
	data "vra_data_collector" "dc" {
		name = "%s"
	}

	data "vra_region_enumeration" "dc_regions" {
	  username    = "%s"
	  password    = "%s"
	  hostname    = "%s"
	  dc_id       = data.vra_data_collector.dc.id
	}

	resource "vra_cloud_account_vsphere" "my_vsphere_account" {
	  name        = "my_vsphere_account_%d"
	  description = "test cloud account"
	  username    = "%s"
	  password    = "%s"
	  hostname    = "%s"
	  dc_id       = data.vra_data_collector.dc.id

	  regions                 = data.vra_region_enumeration.dc_regions.regions
	  accept_self_signed_cert = true
	  tags {
		key   = "foo"
		value = "bar"
	  }
	  tags {
		key = "where"
		value = "waldo"
	  }
	}`, dcname, username, password, hostname, rInt, username, password, hostname)
}

func testAccDataSourceVRACloudAccountVsphereNotFound(rInt int) string {
	return testAccDataSourceVRACloudAccountVsphereBase(rInt) + `
	data "vra_cloud_account_vsphere" "test-cloud-account" {
		name = "foobar"
	}`
}

func testAccDataSourceVRACloudAccountVsphere(rInt int) string {
	return testAccDataSourceVRACloudAccountVsphereBase(rInt) + `
	data "vra_cloud_account_vsphere" "test-cloud-account" {
		name = vra_cloud_account_vsphere.my-cloud-account.name
	}`
}
