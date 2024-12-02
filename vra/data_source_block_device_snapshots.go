// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/disk"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"log"
)

func dataSourceBlockDeviceSnapshots() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBlockDeviceSnapshotsRead,
		Schema: map[string]*schema.Schema{
			"block_device_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"snapshots": snapshotsSchema(),
		},
	}
}

func dataSourceBlockDeviceSnapshotsRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading the vra_block_device_snapshots data source with block_device_id %s", d.Get("block_device_id"))
	apiClient := meta.(*Client).apiClient

	var diskSnapshot []*models.DiskSnapshot

	id := d.Get("block_device_id").(string)

	getResp, err := apiClient.Disk.GetDiskSnapshots(disk.NewGetDiskSnapshotsParams().WithID(id))

	if err != nil {
		return err
	}
	diskSnapshot = getResp.GetPayload()

	if len(diskSnapshot) == 0 {
		log.Printf("No snapshots were found with block_device_id %s", id)
		return nil
	}

	d.SetId(id)
	if err := d.Set("snapshots", flattenSnapshots(diskSnapshot)); err != nil {
		return fmt.Errorf("error setting vra_block_device_snapshots - error: %#v", err)
	}

	log.Printf("Finished reading the vra_block_device_snapshots data source with block_device_id %s", id)
	return nil
}
