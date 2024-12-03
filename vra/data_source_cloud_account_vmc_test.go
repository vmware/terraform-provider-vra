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

func TestAccDataSourceVRACloudAccountVMC(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName1 := "vra_cloud_account_vmc.this"
	dataSourceName1 := "data.vra_cloud_account_vmc.this"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVMC(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRACloudAccountVMCNotFound(rInt),
				ExpectError: regexp.MustCompile("cloud account foobar not found"),
			},
			{
				Config: testAccDataSourceVRACloudAccountVMC(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "dc_id", dataSourceName1, "dc_id"),
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
					resource.TestCheckResourceAttrPair(resourceName1, "nsx_hostname", dataSourceName1, "nsx_hostname"),
					resource.TestCheckResourceAttrPair(resourceName1, "regions", dataSourceName1, "regions"),
					resource.TestCheckResourceAttrPair(resourceName1, "tags", dataSourceName1, "tags"),
					resource.TestCheckResourceAttrPair(resourceName1, "sddc_name", dataSourceName1, "sddc_name"),
					resource.TestCheckResourceAttrPair(resourceName1, "vcenter_hostname", dataSourceName1, "vcenter_hostname"),
					resource.TestCheckResourceAttrPair(resourceName1, "vcenter_username", dataSourceName1, "vcenter_username"),
				),
			},
		},
	})
}

func testAccDataSourceVRACloudAccountVMCBase(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	apiToken := os.Getenv("VRA_VMC_API_TOKEN")
	sddcName := os.Getenv("VRA_VMC_SDDC_NAME")
	vCenterHostName := os.Getenv("VRA_VMC_VCENTER_HOSTNAME")
	vCenterUserName := os.Getenv("VRA_VMC_VCENTER_USERNAME")
	vCenterPassword := os.Getenv("VRA_VMC_VCENTER_PASSWORD")
	nsxHostName := os.Getenv("VRA_VMC_NSX_HOSTNAME")
	dcName := os.Getenv("VRA_VMC_DATA_COLLECTOR_NAME")
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

	resource "vra_cloud_account_vmc" "this" {
	  name        = "test-vmc-account-%d"
	  description = "tf test vmc cloud account"

	  api_token = "%s"
	  sddc_name = "%s"

	  vcenter_username    = "%s"
	  vcenter_password    = "%s"
	  vcenter_hostname    = "%s"
	  nsx_hostname        = "%s"
	  dc_id               = data.vra_data_collector.dc.id

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
	}`, dcName, vCenterUserName, vCenterPassword, vCenterHostName, rInt, apiToken, sddcName, vCenterUserName, vCenterPassword, vCenterHostName, nsxHostName)
}

func testAccDataSourceVRACloudAccountVMCNotFound(rInt int) string {
	return testAccDataSourceVRACloudAccountVMCBase(rInt) + `
	data "vra_cloud_account_vmc" "this" {
		name = "foobar"
	}`
}

func testAccDataSourceVRACloudAccountVMC(rInt int) string {
	return testAccDataSourceVRACloudAccountVMCBase(rInt) + `
	data "vra_cloud_account_vmc" "this" {
		name = "${vra_cloud_account_vmc.this.name}"
	}`
}
