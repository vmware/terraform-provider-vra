// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vra-sdk-go/pkg/client/disk"
)

func TestAccVRABlockDeviceSnapshotBasic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckBlockDeviceSnapshotResource(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRABlockDeviceSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRABlockDeviceSnapshotConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRABlockDeviceSnapshotExists("vra_block_device_snapshot.this"),
					resource.TestCheckResourceAttr(
						"vra_block_device_snapshot.this", "description", "terraform block device snapshot"),
					resource.TestCheckResourceAttrSet(
						"vra_block_device_snapshot.this", "is_current"),
				),
			},
		},
	})
}

func testAccCheckVRABlockDeviceSnapshotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no block_device_snapshot ID is set")
		}

		return nil
	}
}

func testAccCheckVRABlockDeviceSnapshotDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "vra_block_device" {
			_, err := apiClient.Disk.GetBlockDevice(disk.NewGetBlockDeviceParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_block_device' still exists with id %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckVRABlockDeviceSnapshotConfig(rInt int) string {
	return testAccCheckVRABlockDeviceBasicConfig(rInt) + `
	resource "vra_block_device_snapshot" "this" {
	  block_device_id = vra_block_device.disk1.id
	  description = "terraform block device snapshot"
	}`
}

func testAccCheckVRABlockDeviceBasicConfig(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	name := os.Getenv("VRA_VSPHERE_CLOUD_ACCOUNT_NAME")
	region := os.Getenv("VRA_REGION")
	projectName := os.Getenv("VRA_PROJECT")
	datastoreName := os.Getenv("VRA_VSPHERE_DATASTORE_NAME")
	storagePolicyName := os.Getenv("VRA_VSPHERE_STORAGE_POLICY_NAME")

	return fmt.Sprintf(`
	data "vra_cloud_account_vsphere" "this" {
		name = "%s"
	  }

	data "vra_region" "this" {
		cloud_account_id = data.vra_cloud_account_vsphere.this.id
		region = "%s"
	}

	data "vra_project" "this" {
		name = "%s"
	 }

	# Lookup vSphere fabric datastore using its name
	data "vra_fabric_datastore_vsphere" "this" {
	  filter = "name eq '%s' and cloudAccountId eq '${data.vra_cloud_account_vsphere.this.id}'"
	}

	# Lookup vSphere fabric storage policy using its name
	data "vra_fabric_storage_policy_vsphere" "this" {
	  filter = "name eq '%s' and cloudAccountId eq '${data.vra_cloud_account_vsphere.this.id}'"
	}

	resource "vra_storage_profile" "this" {
	  name         = "vSphere-first-class-disk"
	  description  = "vSphere Storage Profile with first class disk."
	  region_id    = data.vra_region.this.id
	  default_item = true

	  disk_properties = {
		diskType         = "firstClass"
		provisioningType = "thin" // Supported values: "thin", "thick", "eagerZeroedThick"
	  }

	  disk_target_properties = {
		datastoreId     = data.vra_fabric_datastore_vsphere.this.id
		storagePolicyId = data.vra_fabric_storage_policy_vsphere.this.id
	  }

	  tags {
		key   = "foo"
		value = "bar"
	  }
	}

	resource "vra_block_device" "disk1" {
	  capacity_in_gb = 10
      name = "block-device-%d"
      description = "terraform block device for snapshot"
	  project_id = data.vra_project.this.id
      depends_on = [vra_storage_profile.this]
      persistent = true
      purge = true

	  constraints {
		mandatory  = true
		expression = "foo:bar"
	  }
	}
`, name, region, projectName, datastoreName, storagePolicyName, rInt)
}
