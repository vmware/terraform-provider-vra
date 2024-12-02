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
)

func TestAccBlockDeviceSnapshotRead(t *testing.T) {
	blockDeviceID := os.Getenv("VRA_BLOCK_DEVICE_ID")
	blockDevice := "data.vra_block_device_snapshot.snapshot"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckBlockDeviceSnapshot(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccBlockDeviceSnapshotConfig(blockDeviceID + "foo"),
				ExpectError: regexp.MustCompile(fmt.Sprintf("Disk with id '%sfoo' does not exist", blockDeviceID)),
			},
			{
				Config: testAccBlockDeviceSnapshotConfig(blockDeviceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(blockDevice, "snapshots"),
					resource.TestCheckResourceAttrSet(blockDevice, "id"),
				),
			},
		},
	})
}

func testAccBlockDeviceSnapshotConfig(blockDeviceID string) string {
	return fmt.Sprintf(`
		data "vra_block_device_snapshot" "snapshot" {
		  block_device_id = "%s"
		}`, blockDeviceID)
}
