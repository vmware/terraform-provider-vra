package vra

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/disk"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"log"
)

func dataSourceBlockDevice() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBlockDeviceRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"capacity_in_gb": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cloud_account_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_properties": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"deployment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"external_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_zone_id": {
				Type:     schema.TypeString,
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
			"persistent": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
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

	log.Printf("Finished reading the vra_block_device data source with filter %s", d.Get("filter"))
	return nil
}
