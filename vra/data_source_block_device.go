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

func dataSourceBlockDevice() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBlockDeviceRead,

		Schema: map[string]*schema.Schema{
			// Optional arguments
			"filter": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id"},
				Description:   "Search criteria to filter the list of block devices.",
			},
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"filter"},
				Description:   "The id of the block device.",
			},

			// Imported attributes
			"capacity_in_gb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Capacity of the block device in GB.",
			},
			"cloud_account_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Id of the cloud account this storage profile belongs to.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was created. The date is in ISO 8601 and UTC.",
			},
			"custom_properties": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "\"Additional custom properties that may be used to extend the block device.",
			},
			"deployment_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the deployment that is associated with this resource.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A human-friendly description.",
			},
			"expand_snapshots": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether the snapshots of the block-devices should be included in the state. Applicable only for first class block devices.",
			},
			"external_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "External entity id on the cloud provider side.",
			},
			"external_region_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The external regionId of the resource.",
			},
			"external_zone_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The external zoneId of the resource.",
			},
			"links": linksSchema(),
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly name for the block device.",
			},
			"org_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the organization this block device belongs to.",
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email of the user that owns this block device.",
			},
			"persistent": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the block device survives a delete action.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the project this resource belongs to.",
			},
			"snapshots": snapshotsSchema(),
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the block device.",
			},
			"tags": tagsSchema(),
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC.",
			},
		},
	}
}

func dataSourceBlockDeviceRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading the vra_block_device data source with filter %s", d.Get("filter"))
	apiClient := meta.(*Client).apiClient

	var blockDevice *models.BlockDevice

	id := d.Get("id").(string)
	filter := d.Get("filter").(string)

	if id == "" && filter == "" {
		return fmt.Errorf("one of id or filter is required")
	}

	if id != "" {
		getResp, err := apiClient.Disk.GetBlockDevice(disk.NewGetBlockDeviceParams().WithID(id))

		if err != nil {
			return err
		}
		blockDevice = getResp.GetPayload()
	} else {
		getResp, err := apiClient.Disk.GetBlockDevices(disk.NewGetBlockDevicesParams().WithDollarFilter(withString(filter)))
		if err != nil {
			return err
		}

		blockDevices := *getResp.Payload
		if len(blockDevices.Content) > 1 {
			return fmt.Errorf("vra_block_device must filter to a block device")
		}
		if len(blockDevices.Content) == 0 {
			return fmt.Errorf("vra_block_device filter did not match any block device")
		}

		blockDevice = blockDevices.Content[0]
	}

	d.SetId(*blockDevice.ID)
	d.Set("capacity_in_gb", blockDevice.CapacityInGB)
	d.Set("cloud_account_ids", blockDevice.CloudAccountIds)
	d.Set("created_at", blockDevice.CreatedAt)
	d.Set("custom_properties", blockDevice.CustomProperties)
	d.Set("deployment_id", blockDevice.DeploymentID)
	d.Set("description", blockDevice.Description)
	d.Set("external_id", blockDevice.ExternalID)
	d.Set("external_region_id", blockDevice.ExternalRegionID)
	d.Set("external_zone_id", blockDevice.ExternalZoneID)
	d.Set("name", blockDevice.Name)
	d.Set("org_id", blockDevice.OrgID)
	d.Set("owner", blockDevice.Owner)
	d.Set("persistent", blockDevice.Persistent)
	d.Set("project_id", blockDevice.ProjectID)
	d.Set("status", blockDevice.Status)
	d.Set("updated_at", blockDevice.UpdatedAt)

	if err := d.Set("tags", flattenTags(blockDevice.Tags)); err != nil {
		return fmt.Errorf("error setting block device tags - error: %v", err)
	}

	if err := d.Set("links", flattenLinks(blockDevice.Links)); err != nil {
		return fmt.Errorf("error setting block device links - error: %#v", err)
	}

	expandSnapshots := d.Get("expand_snapshots").(bool)
	if expandSnapshots {
		snapshots, err := apiClient.Disk.GetDiskSnapshots(disk.NewGetDiskSnapshotsParams().WithID(d.Id()))
		if err != nil {
			return fmt.Errorf("error getting block device snapshots - error: %#v", err)
		}

		if err := d.Set("snapshots", flattenSnapshots(snapshots.Payload)); err != nil {
			return fmt.Errorf("error setting block device snapshots - error: %#v", err)
		}
	}

	log.Printf("Finished reading the vra_block_device data source with filter %s", d.Get("filter"))
	return nil
}
