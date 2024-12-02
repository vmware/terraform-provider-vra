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

func TestAccVRACloudAccountvSphere_Basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckVsphere(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRACloudAccountvSphereDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRACloudAccountvSphereConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRACloudAccountvSphereExists("vra_cloud_account_vsphere.my_vsphere_account"),
					resource.TestMatchResourceAttr(
						"vra_cloud_account_vsphere.my_vsphere_account", "name", regexp.MustCompile("^my_vsphere_account_"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_vsphere.my_vsphere_account", "description", "test cloud account"),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_vsphere.my_vsphere_account", "tags.#", "2"),
				),
			},
		},
	})
}

func testAccCheckVRACloudAccountvSphereExists(n string) resource.TestCheckFunc {
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

func testAccCheckVRACloudAccountvSphereDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vra_cloud_account_vsphere" {
			continue
		}

		_, err := apiClient.CloudAccount.GetVSphereCloudAccount(cloud_account.NewGetVSphereCloudAccountParams().WithID(rs.Primary.ID))
		if err == nil {
			return err
		}
	}

	return nil
}

func testAccCheckVRACloudAccountvSphereConfig(rInt int) string {
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
