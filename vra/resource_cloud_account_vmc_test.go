// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
)

func TestAccVRACloudAccountVMC_Basic(t *testing.T) {
	rInt := acctest.RandInt()
	resource1 := "vra_cloud_account_vmc.this"
	dataSource1 := "data.vra_data_collector.dc"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckVMC(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRACloudAccountVMCDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRACloudAccountVMCConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRACloudAccountVMCExists(resource1),
					resource.TestMatchResourceAttr(
						resource1, "name", regexp.MustCompile("^test-vmc-account-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						resource1, "description", "my vmc cloud account"),
					resource.TestCheckResourceAttr(
						resource1, "sddc_name", os.Getenv("VRA_VMC_SDDC_NAME")),
					resource.TestCheckResourceAttr(
						resource1, "vcenter_hostname", os.Getenv("VRA_VMC_VCENTER_HOSTNAME")),
					resource.TestCheckResourceAttr(
						resource1, "vcenter_username", os.Getenv("VRA_VMC_VCENTER_USERNAME")),
					resource.TestCheckResourceAttr(
						resource1, "nsx_hostname", os.Getenv("VRA_VMC_NSX_HOSTNAME")),
					resource.TestCheckResourceAttrPair(
						resource1, "dc_id", dataSource1, "id"),
					resource.TestCheckResourceAttr(
						resource1, "tags.#", "2"),
				),
			},
			{
				Config: testAccCheckVRACloudAccountVMCUpdateConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRACloudAccountVMCExists(resource1),
					resource.TestMatchResourceAttr(
						resource1, "name", regexp.MustCompile("^test-vmc-account-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						resource1, "description", "your vmc cloud account"),
					resource.TestCheckResourceAttr(
						resource1, "sddc_name", os.Getenv("VRA_VMC_SDDC_NAME")),
					resource.TestCheckResourceAttr(
						resource1, "vcenter_hostname", os.Getenv("VRA_VMC_VCENTER_HOSTNAME")),
					resource.TestCheckResourceAttr(
						resource1, "vcenter_username", os.Getenv("VRA_VMC_VCENTER_USERNAME")),
					resource.TestCheckResourceAttr(
						resource1, "nsx_hostname", os.Getenv("VRA_VMC_NSX_HOSTNAME")),
					resource.TestCheckResourceAttrPair(
						resource1, "dc_id", dataSource1, "id"),
					resource.TestCheckResourceAttr(
						resource1, "tags.#", "2"),
				),
			},
		},
	})
}

func TestAccVRACloudAccountVMC_Duplicate(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckVMC(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRACloudAccountVMCDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckVRACloudAccountVMCConfigDuplicateRegion(rInt),
				ExpectError: regexp.MustCompile("specified regions are not unique"),
			},
		},
	})
}

func testAccCheckVRACloudAccountVMCExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no cloud account ID is set")
		}

		return nil
	}
}

func testAccCheckVRACloudAccountVMCDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vra_cloud_account_vmc" {
			continue
		}

		_, err := apiClient.CloudAccount.GetCloudAccount(cloud_account.NewGetCloudAccountParams().WithID(rs.Primary.ID))
		if err == nil {
			return err
		}
	}

	return nil
}

func testAccCheckVRACloudAccountVMCConfig(rInt int) string {
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
	  description = "my vmc cloud account"

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

func testAccCheckVRACloudAccountVMCUpdateConfig(rInt int) string {
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
	  description = "your vmc cloud account"

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

func testAccCheckVRACloudAccountVMCConfigDuplicateRegion(rInt int) string {
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
	  description = "my vmc cloud account"

	  api_token = "%s"
	  sddc_name = "%s"

	  vcenter_username    = "%s"
	  vcenter_password    = "%s"
	  vcenter_hostname    = "%s"
	  nsx_hostname        = "%s"
	  dc_id               = data.vra_data_collector.dc.id

	  regions                 = ["Datacenter:datacenter-2", "Datacenter:datacenter-2"]
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
