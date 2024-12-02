// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/vmware/vra-sdk-go/pkg/client"
	"github.com/vmware/vra-sdk-go/pkg/client/disk"
	"github.com/vmware/vra-sdk-go/pkg/client/request"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBlockDevice() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBlockDeviceCreate,
		ReadContext:   resourceBlockDeviceRead,
		UpdateContext: resourceBlockDeviceUpdate,
		DeleteContext: resourceBlockDeviceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required arguments
			"capacity_in_gb": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Capacity of the block device in GB.",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					return !strings.HasPrefix(new, old)
				},
				Description: "A human-friendly name for the block device.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the project this resource belongs to.",
			},

			// Optional arguments
			"constraints": constraintsSchema(),
			"custom_properties": {
				Type:        schema.TypeMap,
				Computed:    true,
				Optional:    true,
				Description: "Additional custom properties that may be used to extend the block device. Following disk custom properties can be passed while creating a block device: \n\n1. dataStore: Defines name of the datastore in which the disk has to be provisioned.\n2. storagePolicy: Defines name of the storage policy in which the disk has to be provisioned. If name of the datastore is specified in the custom properties then, datastore takes precedence.\n3. provisioningType: Defines the type of provisioning. For eg. thick/thin.",
			},
			"deployment_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The id of the deployment that is associated with this resource.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A human-friendly description.",
			},
			"disk_content_base_64": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Content of a disk, base64 encoded.",
			},
			"encrypted": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether the block device should be encrypted or not.",
			},
			"expand_snapshots": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether the snapshots of the block-devices should be included in the resource state. Applicable only for first class block devices.",
			},
			"persistent": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether the block device survives a delete action.",
			},
			"purge": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the disk has to be completely destroyed or should be kept in the system. Valid only for block devices with ‘persistent’ set to true, only used for destroy the resource",
			},
			"source_reference": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Reference to URI using which the block device has to be created. Example: ami-0d4cfd66",
			},

			// Imported attributes
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was created. The date is in ISO 8601 and UTC.",
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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceBlockDeviceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to create vra_block_device resource")
	apiClient := m.(*Client).apiClient

	capacityInGB := int32(d.Get("capacity_in_gb").(int))
	name := d.Get("name").(string)
	projectID := d.Get("project_id").(string)
	constraints := expandConstraints(d.Get("constraints").(*schema.Set).List())
	customProperties := expandCustomProperties(d.Get("custom_properties").(map[string]interface{}))
	tags := expandTags(d.Get("tags").(*schema.Set).List())

	blockDeviceSpecification := models.BlockDeviceSpecification{
		CapacityInGB:     &capacityInGB,
		Name:             &name,
		ProjectID:        &projectID,
		Constraints:      constraints,
		CustomProperties: customProperties,
		Tags:             tags,
	}

	if v, ok := d.GetOk("description"); ok {
		blockDeviceSpecification.Description = v.(string)
	}

	if v, ok := d.GetOk("deployment_id"); ok {
		blockDeviceSpecification.DeploymentID = v.(string)
	}

	if v, ok := d.GetOk("encrypted"); ok {
		blockDeviceSpecification.Encrypted = v.(bool)
	}

	if v, ok := d.GetOk("persistent"); ok {
		blockDeviceSpecification.Persistent = v.(bool)
	}

	if v, ok := d.GetOk("source_reference"); ok {
		blockDeviceSpecification.SourceReference = v.(string)
	}

	if v, ok := d.GetOk("disk_content_base_64"); ok {
		blockDeviceSpecification.DiskContentBase64 = v.(string)
	}

	log.Printf("[DEBUG] create block device: %#v", blockDeviceSpecification)
	createBlockDeviceCreated, err := apiClient.Disk.CreateBlockDevice(disk.NewCreateBlockDeviceParams().WithBody(&blockDeviceSpecification))
	if err != nil {
		return diag.FromErr(err)
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    blockDeviceStateRefreshFunc(*apiClient, *createBlockDeviceCreated.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
	}

	resourceIDs, err := stateChangeFunc.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	blockDeviceIDs := resourceIDs.([]string)
	i := strings.LastIndex(blockDeviceIDs[0], "/")
	blockDeviceID := blockDeviceIDs[0][i+1 : len(blockDeviceIDs[0])]
	d.SetId(blockDeviceID)
	log.Printf("Finished to create vra_block_device resource with name %s", d.Get("name"))

	return resourceBlockDeviceRead(ctx, d, m)
}

func blockDeviceStateRefreshFunc(apiClient client.API, id string) retry.StateRefreshFunc {
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
			loadBalancerIDs := make([]string, len(ret.Payload.Resources))
			for i, r := range ret.Payload.Resources {
				loadBalancerIDs[i] = strings.TrimPrefix(r, "/iaas/api/block-device/")
			}
			return loadBalancerIDs, *status, nil
		default:
			return [...]string{id}, ret.Payload.Message, fmt.Errorf("blockDeviceStateRefreshFunc: unknown status %v", *status)
		}
	}
}

func resourceBlockDeviceRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Reading the vra_block_device resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	resp, err := apiClient.Disk.GetBlockDevice(disk.NewGetBlockDeviceParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *disk.GetBlockDeviceNotFound:
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	blockDevice := *resp.Payload
	d.Set("capacity_in_gb", blockDevice.CapacityInGB)
	d.Set("created_at", blockDevice.CreatedAt)
	d.Set("custom_properties", blockDevice.CustomProperties)
	d.Set("description", blockDevice.Description)
	d.Set("deployment_id", blockDevice.DeploymentID)
	d.Set("external_id", blockDevice.ExternalID)
	d.Set("external_region_id", blockDevice.ExternalRegionID)
	d.Set("external_zone_id", blockDevice.ExternalZoneID)
	d.Set("name", blockDevice.Name)
	d.Set("org_id", blockDevice.OrgID)
	d.Set("owner", blockDevice.Owner)
	d.Set("persistent", blockDevice.Persistent)
	d.Set("status", blockDevice.Status)
	d.Set("updated_at", blockDevice.UpdatedAt)

	if err := d.Set("tags", flattenTags(blockDevice.Tags)); err != nil {
		return diag.Errorf("error setting block device tags - error: %v", err)
	}

	if err := d.Set("links", flattenLinks(blockDevice.Links)); err != nil {
		return diag.Errorf("error setting block device links - error: %#v", err)
	}

	expandSnapshots := d.Get("expand_snapshots").(bool)
	if expandSnapshots {
		snapshots, err := apiClient.Disk.GetDiskSnapshots(disk.NewGetDiskSnapshotsParams().WithID(d.Id()))
		if err != nil {
			return diag.Errorf("error getting block device snapshots - error: %#v", err)
		}

		if err := d.Set("snapshots", flattenSnapshots(snapshots.Payload)); err != nil {
			return diag.Errorf("error setting block device snapshots - error: %#v", err)
		}
	} else {
		d.Set("snapshots", make([]map[string]interface{}, 0))
	}

	log.Printf("Finished reading the vra_block_device resource with name %s", d.Get("name"))
	return nil
}

func resourceBlockDeviceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	log.Printf("Starting to update the vra_block_device resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	if d.HasChange("capacity_in_gb") {
		err := resizeDisk(ctx, d, apiClient, id)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("Finished updating vra_block_device resource with name %s", d.Get("name"))
	return resourceBlockDeviceRead(ctx, d, m)
}

func resizeDisk(ctx context.Context, d *schema.ResourceData, apiClient *client.API, id string) error {

	log.Printf("Starting resize of vra_block_device resource with name %s", d.Get("name"))

	capacityInGB := int32(d.Get("capacity_in_gb").(int))
	resizeBlockDeviceAccepted, resizeBlockDeviceNoContent, err := apiClient.Disk.ResizeBlockDevice(disk.NewResizeBlockDeviceParams().WithID(id).WithCapacityInGB(capacityInGB))
	if err != nil {
		return err
	}
	if resizeBlockDeviceNoContent != nil {
		return nil
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    blockDeviceStateRefreshFunc(*apiClient, *resizeBlockDeviceAccepted.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 5 * time.Second,
	}

	if _, err = stateChangeFunc.WaitForStateContext(ctx); err != nil {
		return err
	}

	return nil
}

func resourceBlockDeviceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to delete the vra_block_device resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	deleteBlockDeviceParams := disk.NewDeleteBlockDeviceParams().WithID(id)
	purge := false
	persistent := false

	if v, ok := d.GetOk("persistent"); ok {
		persistent = v.(bool)
	}

	// If the disk is persistent type, have to pass purge query parameter
	if v, ok := d.GetOk("purge"); ok {
		purge = v.(bool)
	}

	if purge && persistent {
		log.Printf("The vra_block_device %s is persistent type and purge set to true, it will be purged", d.Get("name"))
		deleteBlockDeviceParams.WithPurge(&purge)
	}

	deleteBlockDeviceAccepted, deleteBlockDeviceCompleted, err := apiClient.Disk.DeleteBlockDevice(deleteBlockDeviceParams)
	if err != nil {
		return diag.FromErr(err)
	}

	// Handle non-request tracker case
	if deleteBlockDeviceCompleted != nil {
		return nil
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    blockDeviceStateRefreshFunc(*apiClient, *deleteBlockDeviceAccepted.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 5 * time.Second,
	}

	if _, err = stateChangeFunc.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("Finished deleting the vra_block_device resource with name %s", d.Get("name"))
	return nil
}
