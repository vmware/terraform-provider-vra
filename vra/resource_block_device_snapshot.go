// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/vmware/vra-sdk-go/pkg/client"
	"github.com/vmware/vra-sdk-go/pkg/client/disk"
	"github.com/vmware/vra-sdk-go/pkg/client/request"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBlockDeviceSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBlockDeviceSnapshotCreate,
		ReadContext:   resourceBlockDeviceSnapshotRead,
		UpdateContext: resourceBlockDeviceSnapshotUpdate,
		DeleteContext: resourceBlockDeviceSnapshotDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"block_device_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_current": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"links": linksSchema(),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceBlockDeviceSnapshotCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to create vra_block_device_snapshot resource")
	apiClient := m.(*Client).apiClient

	description := d.Get("description").(string)
	blockDeviceID := d.Get("block_device_id").(string)

	DiskSnapshotSpecification := models.DiskSnapshotSpecification{
		Description: description,
	}

	log.Printf("[DEBUG] create vra_block_device_snapshot: %#v", DiskSnapshotSpecification)
	createDiskSnapshotCreated, _, err := apiClient.Disk.CreateBlockDeviceSnapshot(disk.NewCreateBlockDeviceSnapshotParams().WithID(blockDeviceID).WithBody(&DiskSnapshotSpecification))
	if err != nil {
		return diag.FromErr(err)
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    BlockDeviceSnapshotStateRefreshFunc(*apiClient, *createDiskSnapshotCreated.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
	}

	if _, err = stateChangeFunc.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("Finished to create vra_block_device_snapshot resource with for vra_block_device: %s", blockDeviceID)

	snapshotID, err := findCreatedBlockDeviceSnapshot(blockDeviceID, m)
	d.SetId(snapshotID)

	if err != nil {
		return nil
	}

	return resourceBlockDeviceSnapshotRead(ctx, d, m)
}

func BlockDeviceSnapshotStateRefreshFunc(apiClient client.API, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ret, err := apiClient.Request.GetRequestTracker(request.NewGetRequestTrackerParams().WithID(id))
		if err != nil {
			return "", models.RequestTrackerStatusFAILED, err
		}

		status := ret.Payload.Status
		switch *status {
		case models.RequestTrackerStatusFAILED:
			return []string{""}, *status, errors.New(ret.Payload.Message)
		case models.RequestTrackerStatusINPROGRESS:
			return [...]string{id}, *status, nil
		case models.RequestTrackerStatusFINISHED:
			snapshotID := *ret.Payload.ID
			return snapshotID, *status, nil
		default:
			return [...]string{id}, ret.Payload.Message, fmt.Errorf("BlockDeviceSnapshotStateRefreshFunc: unknown status %v", *status)
		}
	}
}

func findCreatedBlockDeviceSnapshot(blockDeviceID string, m interface{}) (string, error) {

	log.Printf("Reading the vra_block_device_snapshot resource for vra_block_device %s ", blockDeviceID)
	apiClient := m.(*Client).apiClient

	errMsg := "failed to find the created snapshot for the vra_block_device_snapshot resource with id %s"

	resp, err := apiClient.Disk.GetDiskSnapshots(disk.NewGetDiskSnapshotsParams().WithID(blockDeviceID))
	if err != nil {
		return "", fmt.Errorf(errMsg, blockDeviceID)
	}

	diskSnapshots := resp.Payload
	if len(diskSnapshots) < 1 {
		return "", fmt.Errorf(errMsg, blockDeviceID)
	}

	for _, diskSnapshot := range diskSnapshots {
		if isCurrent, ok := diskSnapshot.SnapshotProperties["isCurrent"]; ok && isCurrent == "true" {
			return *diskSnapshot.ID, nil
		}
	}

	return "", fmt.Errorf(errMsg, blockDeviceID)
}

func resourceBlockDeviceSnapshotRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	blockDeviceID := d.Get("block_device_id").(string)
	log.Printf("Reading the vra_block_device_snapshot resource for vra_block_device %s ", blockDeviceID)
	apiClient := m.(*Client).apiClient

	resp, err := apiClient.Disk.GetDiskSnapshot(disk.NewGetDiskSnapshotParams().WithID(blockDeviceID).WithId1(d.Id()))
	if err != nil {
		switch err.(type) {
		case *disk.GetDiskSnapshotNotFound:
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	diskSnapshot := resp.Payload

	d.SetId(*diskSnapshot.ID)
	d.Set("created_at", diskSnapshot.CreatedAt)
	d.Set("description", diskSnapshot.Desc)
	d.Set("name", diskSnapshot.Name)
	d.Set("org_id", diskSnapshot.OrgID)
	d.Set("owner", diskSnapshot.Owner)
	d.Set("updated_at", diskSnapshot.UpdatedAt)

	d.Set("is_current", false)
	if isCurrent, ok := diskSnapshot.SnapshotProperties["isCurrent"]; ok && isCurrent == "true" {
		d.Set("is_current", true)
	}
	if err := d.Set("links", flattenLinks(diskSnapshot.Links)); err != nil {
		return diag.Errorf("error setting vra_block_device_snapshot links - error: %#v", err)
	}

	log.Printf("Finished reading the vra_block_device_snapshot resource with id %s", *diskSnapshot.ID)
	return nil
}

func resourceBlockDeviceSnapshotUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return diag.Errorf("update vra_block_device_snapshot is not supported")
}

func resourceBlockDeviceSnapshotDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	blockDeviceID := d.Get("block_device_id").(string)
	snapshotID := d.Id()
	log.Printf("Starting to delete the vra_block_device_snapshot of vra_block_device %s with ID: %s", blockDeviceID, snapshotID)
	apiClient := m.(*Client).apiClient

	deleteDiskSnapshotAccepted, deleteDiskSnapshotCompleted, err := apiClient.Disk.
		DeleteBlockDeviceSnapshot(
			disk.NewDeleteBlockDeviceSnapshotParams().WithID(blockDeviceID).WithId1(snapshotID))
	if err != nil {
		return diag.FromErr(err)
	}

	// Handle non-request tracker case
	if deleteDiskSnapshotCompleted != nil {
		d.SetId("")
		return nil
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    BlockDeviceSnapshotStateRefreshFunc(*apiClient, *deleteDiskSnapshotAccepted.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 5 * time.Second,
	}

	if _, err = stateChangeFunc.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("Finished deleting the vra_block_device_snapshot resource with name %s", d.Get("name"))
	return nil
}
