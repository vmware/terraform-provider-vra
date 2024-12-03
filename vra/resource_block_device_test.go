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

	"github.com/vmware/vra-sdk-go/pkg/client/disk"
	"github.com/vmware/vra-sdk-go/pkg/client/location"
	"github.com/vmware/vra-sdk-go/pkg/client/project"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVRABlockDeviceBasic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckBlockDevice(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRABlockDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRABlockDeviceConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRABlockDeviceExists("vra_block_device.this"),
					resource.TestMatchResourceAttr(
						"vra_block_device.this", "name", regexp.MustCompile("^block-device-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"vra_block_device.this", "description", "terraform block device"),
					resource.TestCheckResourceAttr(
						"vra_block_device.this", "capacity_in_gb", "4"),
				),
			},
			{
				Config: testAccCheckVRABlockDeviceUpdateConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRABlockDeviceExists("vra_block_device.this"),
					resource.TestMatchResourceAttr(
						"vra_block_device.this", "name", regexp.MustCompile("^block-device-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"vra_block_device.this", "description", "terraform block device"),
					resource.TestCheckResourceAttr(
						"vra_block_device.this", "capacity_in_gb", "8"),
				),
			},
		},
	})
}

func testAccCheckVRABlockDeviceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no block_device ID is set")
		}

		return nil
	}
}

func testAccCheckVRABlockDeviceDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "vra_project" {
			_, err := apiClient.Project.GetProject(project.NewGetProjectParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_project' still exists with id %s", rs.Primary.ID)
			}
		}

		if rs.Type == "vra_zone" {
			_, err := apiClient.Location.GetZone(location.NewGetZoneParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_zone' still exists with id %s", rs.Primary.ID)
			}
		}

		if rs.Type == "vra_block_device" {
			_, err := apiClient.Disk.GetBlockDevice(disk.NewGetBlockDeviceParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_block_device' still exists with id %s", rs.Primary.ID)

			}
		}
	}

	return nil
}

func testAccCheckVRABlockDeviceConfig(rInt int) string {

	return testAccCheckVRABlockDevice(rInt) + fmt.Sprintf(`
	resource "vra_block_device" "this" {
		name = "block-device-%d"
		description = "terraform block device"
		project_id = vra_project.this.id
		capacity_in_gb = 4
	  }`, rInt)
}

func testAccCheckVRABlockDeviceUpdateConfig(rInt int) string {

	return testAccCheckVRABlockDevice(rInt) + fmt.Sprintf(`
	resource "vra_block_device" "this" {
		name = "block-device-%d"
		description = "terraform block device"
		project_id = vra_project.this.id
		capacity_in_gb = 8
	  }`, rInt)
}

func testAccCheckVRABlockDevice(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	name := os.Getenv("VRA_AWS_CLOUD_ACCOUNT_NAME")
	region := os.Getenv("VRA_REGION")
	return fmt.Sprintf(`

	data "vra_cloud_account_aws" "this" {
		name = "%s"
	  }

data "vra_region" "this" {
    cloud_account_id = data.vra_cloud_account_aws.this.id
    region = "%s"
}

resource "vra_zone" "this" {
    name = "my-zone-%d"
    description = "description my-vra-zone"
	region_id = data.vra_region.this.id
}

resource "vra_project" "this" {
	name = "my-project-%d"
	description = "test project"
	zone_assignments {
		zone_id       = vra_zone.this.id
		priority      = 1
		max_instances = 2
	  }
 }`, name, region, rInt, rInt)
}
