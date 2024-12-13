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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vmware/vra-sdk-go/pkg/client"
	"github.com/vmware/vra-sdk-go/pkg/client/compute"
	"github.com/vmware/vra-sdk-go/pkg/client/disk"
	"github.com/vmware/vra-sdk-go/pkg/client/request"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMachine() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMachineCreate,
		ReadContext:   resourceMachineRead,
		UpdateContext: resourceMachineUpdate,
		DeleteContext: resourceMachineDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"boot_config": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Description: "Machine boot config that will be passed to the instance that can be used to perform common automated configuration tasks and even run scripts after the instance starts.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A valid cloud config data in json-escaped yaml syntax.",
						},
					},
				},
			},
			"constraints": constraintsSchema(),
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_properties": {
				Type:        schema.TypeMap,
				Computed:    true,
				Optional:    true,
				Description: "Additional custom properties that may be used to extend the machine.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Describes machine within the scope of your organization and is not propagated to the cloud.",
			},
			"deployment_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The id of the deployment that is associated with this resource.",
			},
			"attach_disks_before_boot": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "By default, disks are attached after the machine has been built. FCDs cannot be attached to machine as a day 0 task.",
			},
			"disks": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Specification for attaching/detaching disks to a machine.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A human-friendly block-device name used as an identifier in APIs that support this option.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A human-friendly description.",
						},
						"block_device_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The id of the existing block device.",
						},
						"scsi_controller": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The id of the SCSI controller. Example: SCSI_Controller_0",
						},
						"unit_number": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The unit number of the SCSI controller. Example: 2",
						},
					},
				},
			},
			"disks_list": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of all disks attached to a machine including boot disk, and additional block devices attached using the disks attribute.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A human-friendly block-device name used as an identifier in APIs that support this option.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A human-friendly description.",
						},
						"block_device_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The id of the existing block device.",
						},
						"scsi_controller": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The id of the SCSI controller. Example: SCSI_Controller_0",
						},
						"unit_number": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The unit number of the SCSI controller. Example: 2",
						},
					},
				},
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
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
			},
			"image": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Type of image used for this machine.",
			},
			"image_disk_constraints": constraintsSchema(),
			"image_ref": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"links": linksSchema(),
			"name": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					return !strings.HasPrefix(new, old)
				},
			},
			"nics": nicsSchema(false),
			"organization_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"power_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": tagsSchema(),
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

func resourceMachineCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to create vra_machine resource")
	apiClient := m.(*Client).apiClient

	name := d.Get("name").(string)
	flavor := d.Get("flavor").(string)
	projectID := d.Get("project_id").(string)
	constraints := expandConstraints(d.Get("constraints").(*schema.Set).List())
	tags := expandTags(d.Get("tags").(*schema.Set).List())
	customProperties := expandCustomProperties(d.Get("custom_properties").(map[string]interface{}))
	imageDiskConstraints := expandConstraints(d.Get("image_disk_constraints").(*schema.Set).List())
	nics := expandNics(d.Get("nics").(*schema.Set).List())
	disks := expandDisks(d.Get("disks").(*schema.Set).List())

	machineSpecification := models.MachineSpecification{
		Name:                 &name,
		Flavor:               &flavor,
		ProjectID:            &projectID,
		Constraints:          constraints,
		Tags:                 tags,
		CustomProperties:     customProperties,
		Nics:                 nics,
		ImageDiskConstraints: imageDiskConstraints,
	}

	if v, ok := d.GetOk("attach_disks_before_boot"); ok && v == true {
		machineSpecification.Disks = disks
	}

	image, imageRef := "", ""
	if v, ok := d.GetOk("image"); ok {
		image = v.(string)
		machineSpecification.Image = withString(image)
	}

	if v, ok := d.GetOk("image_ref"); ok {
		imageRef = v.(string)
		machineSpecification.ImageRef = withString(imageRef)
	}

	if image == "" && imageRef == "" {
		return diag.FromErr(errors.New("image or image_ref required"))
	}

	if v, ok := d.GetOk("description"); ok {
		machineSpecification.Description = v.(string)
	}

	if v, ok := d.GetOk("deployment_id"); ok {
		machineSpecification.DeploymentID = v.(string)
	}

	if v, ok := d.GetOk("boot_config"); ok {
		configBootConfig := v.(*schema.Set).List()[0].(map[string]interface{})

		bootConfig := models.MachineBootConfig{
			Content: configBootConfig["content"].(string),
		}

		machineSpecification.BootConfig = &bootConfig
	}

	log.Printf("[DEBUG] create machine: %#v", machineSpecification)
	createMachineCreated, err := apiClient.Compute.CreateMachine(compute.NewCreateMachineParams().WithBody(&machineSpecification))
	if err != nil {
		return diag.FromErr(err)
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    machineStateRefreshFunc(*apiClient, *createMachineCreated.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
	}

	resourceIDs, err := stateChangeFunc.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	machineID := (resourceIDs.([]string))[0]
	d.SetId(machineID)
	log.Printf("Finished to create vra_machine resource with name %s", d.Get("name"))

	// As FCDs cannot be attached to machine as day 0, the machine is first provisioned without requested disks attached.
	// Once the machine provisioning is complete, disks are attached one by one as day-2 action.
	if v, ok := d.GetOk("attach_disks_before_boot"); !ok || v != true {
		for i, diskAttachmentSpecification := range disks {
			log.Printf("Attaching the disk %v of %v (disk id: %v) to vra_machine resource %v", i+1, len(disks), diskAttachmentSpecification.BlockDeviceID, d.Get("name"))

			attachMachineDiskOk, err := apiClient.Disk.AttachMachineDisk(disk.NewAttachMachineDiskParams().WithID(machineID).WithBody(diskAttachmentSpecification))

			if err != nil {
				return diag.FromErr(err)
			}

			stateChangeFunc := retry.StateChangeConf{
				Delay:      5 * time.Second,
				Pending:    []string{models.RequestTrackerStatusINPROGRESS},
				Refresh:    machineStateRefreshFunc(*apiClient, *attachMachineDiskOk.Payload.ID),
				Target:     []string{models.RequestTrackerStatusFINISHED},
				Timeout:    d.Timeout(schema.TimeoutCreate),
				MinTimeout: 5 * time.Second,
			}

			if _, err := stateChangeFunc.WaitForStateContext(ctx); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceMachineRead(ctx, d, m)
}

func machineStateRefreshFunc(apiClient client.API, id string) retry.StateRefreshFunc {
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
			machineIDs := make([]string, len(ret.Payload.Resources))
			for i, r := range ret.Payload.Resources {
				machineIDs[i] = strings.TrimPrefix(r, "/iaas/api/machines/")
			}
			return machineIDs, *status, nil
		default:
			return [...]string{id}, ret.Payload.Message, fmt.Errorf("machineStateRefreshFunc: unknown status %v", *status)
		}
	}
}

func resourceMachineRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Reading the vra_machine resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	resp, err := apiClient.Compute.GetMachine(compute.NewGetMachineParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *compute.GetMachineNotFound:
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	machine := *resp.Payload
	d.Set("name", machine.Name)
	d.Set("description", machine.Description)
	d.Set("deployment_id", machine.DeploymentID)
	d.Set("power_state", machine.PowerState)
	d.Set("address", machine.Address)
	d.Set("project_id", machine.ProjectID)
	d.Set("external_zone_id", machine.ExternalZoneID)
	d.Set("external_region_id", machine.ExternalRegionID)
	d.Set("external_id", machine.ExternalID)
	d.Set("created_at", machine.CreatedAt)
	d.Set("updated_at", machine.UpdatedAt)
	d.Set("owner", machine.Owner)
	d.Set("organization_id", machine.OrgID)
	d.Set("custom_properties", machine.CustomProperties)

	if image, found := machine.CustomProperties["image"]; found {
		d.Set("image", image)
	}

	if imageRef, found := machine.CustomProperties["imageRef"]; found {
		d.Set("imageRef", imageRef)
	}

	if err := d.Set("tags", flattenTags(machine.Tags)); err != nil {
		return diag.Errorf("error setting machine tags - error: %v", err)
	}

	if err := d.Set("links", flattenLinks(machine.Links)); err != nil {
		return diag.Errorf("error setting machine links - error: %#v", err)
	}

	// get all the disks currently attached to the machine
	getMachineDisksOk, err := apiClient.Disk.GetMachineDisks(disk.NewGetMachineDisksParams().WithID(id))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("disks_list", flattenDisks(getMachineDisksOk.Payload.Content)); err != nil {
		return diag.Errorf("error setting machine disks list - error: %#v", err)
	}

	disksConfig := d.Get("disks").(*schema.Set).List()
	if err := d.Set("disks", filterDisks(disksConfig, getMachineDisksOk.Payload.Content)); err != nil {
		return diag.Errorf("error setting machine disks - error: %#v", err)
	}

	log.Printf("Finished reading the vra_machine resource with name %s", d.Get("name"))
	return nil
}

func resourceMachineUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to update the vra_machine resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	if d.HasChange("description") || d.HasChange("tags") {
		err := updateMachine(d, apiClient, id)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// machine resize operation
	if d.HasChange("flavor") {
		err := resizeMachine(ctx, d, apiClient, id)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// attach and/or detach disks if disks configuration is changed
	if d.HasChange("disks") {
		err := attachAndDetachDisks(ctx, d, apiClient, id)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("finished updating the vra_machine resource with name %s", d.Get("name"))
	return resourceMachineRead(ctx, d, m)
}

// attaches and detaches disks
func attachAndDetachDisks(ctx context.Context, d *schema.ResourceData, apiClient *client.API, id string) error {
	log.Printf("identified change in the disks configuration for the machine %s", d.Get("name"))

	oldValue, newValue := d.GetChange("disks")
	oldDisks := oldValue.(*schema.Set).List()
	newDisks := newValue.(*schema.Set).List()

	disksToDetach := disksDifference(oldDisks, newDisks)
	disksToAttach := disksDifference(newDisks, oldDisks)

	log.Printf("number of disks to detach:%v, %+v", len(disksToDetach), disksToDetach)
	log.Printf("number of disks to attach:%v, %+v", len(disksToAttach), disksToAttach)

	// detach the disks one by one
	for i, diskToDetach := range disksToDetach {
		diskID := diskToDetach["block_device_id"].(string)
		log.Printf("Detaching the disk %v of %v (disk id: %v) from vra_machine resource %v", i+1, len(disksToDetach), diskID, d.Get("name"))
		deleteMachineDiskAccepted, _, err := apiClient.Disk.DeleteMachineDisk(disk.NewDeleteMachineDiskParams().WithID(id).WithDiskID(diskID))

		if err != nil {
			return err
		}

		stateChangeFunc := retry.StateChangeConf{
			Delay:      5 * time.Second,
			Pending:    []string{models.RequestTrackerStatusINPROGRESS},
			Refresh:    machineStateRefreshFunc(*apiClient, *deleteMachineDiskAccepted.Payload.ID),
			Target:     []string{models.RequestTrackerStatusFINISHED},
			Timeout:    d.Timeout(schema.TimeoutCreate),
			MinTimeout: 5 * time.Second,
		}

		if _, err := stateChangeFunc.WaitForStateContext(ctx); err != nil {
			return err
		}
	}

	// get all the disks currently attached to the machine
	getMachineDisksOk, err := apiClient.Disk.GetMachineDisks(disk.NewGetMachineDisksParams().WithID(id))
	if err != nil {
		return err
	}

	diskIDs := make([]string, len(getMachineDisksOk.GetPayload().Content))

	for i, blockDevice := range getMachineDisksOk.GetPayload().Content {
		diskIDs[i] = *blockDevice.ID
	}

	log.Printf("disks currently attached to machine %v: %v", id, diskIDs)

	// attach the disks one by one
	for i, diskToAttach := range disksToAttach {
		diskID := diskToAttach["block_device_id"].(string)
		log.Printf("Attaching the disk %v of %v (disk id: %v) to vra_machine resource %v", i+1, len(diskToAttach), diskID, d.Get("name"))

		// attach the disk if it's not already attached to machine
		if index, _ := indexOf(diskID, diskIDs); index == -1 {
			diskAttachmentSpecification := models.DiskAttachmentSpecification{
				BlockDeviceID: withString(diskID),
				Description:   diskToAttach["description"].(string),
				Name:          diskToAttach["name"].(string),
			}

			if vScsiController, okScsiController := diskToAttach["scsi_controller"].(string); okScsiController && vScsiController != "" {
				if vUnitNumber, okUnitNumber := diskToAttach["unit_number"].(string); okUnitNumber && vUnitNumber != "" {
					diskAttachmentSpecification.DiskAttachmentProperties = map[string]string{"scsiController": diskToAttach["scsi_controller"].(string), "unitNumber": diskToAttach["unit_number"].(string)}
				}
			}

			attachMachineDiskOk, err := apiClient.Disk.AttachMachineDisk(disk.NewAttachMachineDiskParams().WithID(id).WithBody(&diskAttachmentSpecification))

			if err != nil {
				return err
			}

			stateChangeFunc := retry.StateChangeConf{
				Delay:      5 * time.Second,
				Pending:    []string{models.RequestTrackerStatusINPROGRESS},
				Refresh:    machineStateRefreshFunc(*apiClient, *attachMachineDiskOk.Payload.ID),
				Target:     []string{models.RequestTrackerStatusFINISHED},
				Timeout:    d.Timeout(schema.TimeoutCreate),
				MinTimeout: 5 * time.Second,
			}

			if _, err := stateChangeFunc.WaitForStateContext(ctx); err != nil {
				return err
			}
		} else {
			log.Printf("disk %v is already attached to machine %v, moving on to the next disk to attach", diskID, id)
		}

	}

	log.Printf("finished to attach/detach disks to vra_machine resource with name %s", d.Get("name"))
	return nil
}

// updates machine description and tags
func updateMachine(d *schema.ResourceData, apiClient *client.API, id string) error {
	log.Printf("identified change in the description and/or tags")
	description := d.Get("description").(string)
	tags := expandTags(d.Get("tags").(*schema.Set).List())

	updateMachineSpecification := models.UpdateMachineSpecification{
		Description: description,
		Tags:        tags,
	}

	log.Printf("[DEBUG] update machine: %#v", updateMachineSpecification)
	_, err := apiClient.Compute.UpdateMachine(compute.NewUpdateMachineParams().WithID(id).WithBody(&updateMachineSpecification))
	if err != nil {
		return err
	}

	log.Printf("finished updating description/tags in vra_machine resource with name %s", d.Get("name"))
	return nil
}

// returns the disks from a that are not in b i.e. a - b
func disksDifference(a, b []interface{}) (diff []map[string]interface{}) {
	m := make(map[string]bool)

	for _, item := range b {
		diskConfig := item.(map[string]interface{})
		blockDeviceID := diskConfig["block_device_id"].(string)
		m[blockDeviceID] = true
	}

	for _, item := range a {
		diskConfig := item.(map[string]interface{})
		blockDeviceID := diskConfig["block_device_id"].(string)
		if _, ok := m[blockDeviceID]; !ok {
			diff = append(diff, diskConfig)
		}
	}
	return
}

// resize machine when there is a change in the flavor
func resizeMachine(ctx context.Context, d *schema.ResourceData, apiClient *client.API, id string) error {
	log.Printf("identified change in the flavor, machine resize will be performed")
	flavor := d.Get("flavor").(string)
	resizeMachine, err := apiClient.Compute.ResizeMachine(compute.NewResizeMachineParams().WithID(id).WithName(&flavor))
	if err != nil {
		return err
	}
	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    machineStateRefreshFunc(*apiClient, *resizeMachine.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 5 * time.Second,
	}
	resourceIDs, err := stateChangeFunc.WaitForStateContext(ctx)
	if err != nil {
		return err
	}
	machineIDs := resourceIDs.([]string)
	d.SetId(machineIDs[0])
	log.Printf("Finished to resize vra_machine resource with name %s", d.Get("name"))
	return nil
}

func resourceMachineDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to delete the vra_machine resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	deleteMachine, err := apiClient.Compute.DeleteMachine(compute.NewDeleteMachineParams().WithID(id))
	if err != nil {
		return diag.FromErr(err)
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    machineStateRefreshFunc(*apiClient, *deleteMachine.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 5 * time.Second,
	}

	if _, err = stateChangeFunc.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("Finished deleting the vra_machine resource with name %s", d.Get("name"))
	return nil
}

func expandDisks(configDisks []interface{}) []*models.DiskAttachmentSpecification {
	disks := make([]*models.DiskAttachmentSpecification, 0, len(configDisks))

	for _, configDisk := range configDisks {
		diskMap := configDisk.(map[string]interface{})

		disk := models.DiskAttachmentSpecification{
			BlockDeviceID: withString(diskMap["block_device_id"].(string)),
		}

		if vScsiController, okScsiController := diskMap["scsi_controller"].(string); okScsiController && vScsiController != "" {
			if vUnitNumber, okUnitNumber := diskMap["unit_number"].(string); okUnitNumber && vUnitNumber != "" {
				disk.DiskAttachmentProperties = map[string]string{"scsiController": diskMap["scsi_controller"].(string), "unitNumber": diskMap["unit_number"].(string)}
			}
		}

		if v, ok := diskMap["name"].(string); ok && v != "" {
			disk.Name = v
		}

		if v, ok := diskMap["description"].(string); ok && v != "" {
			disk.Description = v
		}

		disks = append(disks, &disk)
	}

	return disks
}

func flattenDisks(blockDevices []*models.BlockDevice) []interface{} {
	if len(blockDevices) == 0 {
		return make([]interface{}, 0)
	}

	configDisks := make([]interface{}, 0, len(blockDevices))

	for _, blockDevice := range blockDevices {
		helper := make(map[string]interface{})
		helper["name"] = blockDevice.Name
		helper["description"] = blockDevice.Description
		helper["block_device_id"] = blockDevice.ID

		configDisks = append(configDisks, helper)
	}

	return configDisks
}

func filterDisks(disksConfig []interface{}, blockDevices []*models.BlockDevice) []interface{} {
	if len(disksConfig) == 0 {
		return make([]interface{}, 0)
	}

	disks := make([]interface{}, 0, len(disksConfig))

	// Look for existing disks configuration in the block devices received and map only those.
	// This filters the default boot disk, CD and Floppy drives that are attached by default to machine resource and avoid incorrect plan even when no changes are made to config file.
	for _, diskConfig := range disksConfig {
		diskConfigMap := diskConfig.(map[string]interface{})
		for _, blockDevice := range blockDevices {
			if diskConfigMap["block_device_id"].(string) == *blockDevice.ID {
				helper := make(map[string]interface{})
				helper["block_device_id"] = blockDevice.ID

				if diskConfigMap["name"].(string) != "" {
					helper["name"] = blockDevice.Name
				}

				if diskConfigMap["description"].(string) != "" {
					helper["description"] = blockDevice.Description
				}

				if diskConfigMap["scsi_controller"].(string) != "" {
					helper["scsi_controller"] = diskConfigMap["scsi_controller"]
				}

				if diskConfigMap["unit_number"].(string) != "" {
					helper["unit_number"] = diskConfigMap["unit_number"]
				}

				disks = append(disks, helper)
				break
			}
		}
	}

	return disks
}
