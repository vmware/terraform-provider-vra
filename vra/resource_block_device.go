package vra

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/vmware/vra-sdk-go/pkg/client"
	"github.com/vmware/vra-sdk-go/pkg/client/disk"
	"github.com/vmware/vra-sdk-go/pkg/client/request"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBlockDevice() *schema.Resource {
	return &schema.Resource{
		Create: resourceBlockDeviceCreate,
		Read:   resourceBlockDeviceRead,
		Update: resourceBlockDeviceUpdate,
		Delete: resourceBlockDeviceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"capacity_in_gb": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !strings.HasPrefix(new, old)
				},
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"constraints": constraintsSchema(),
			"custom_properties": {
				Type:     schema.TypeMap,
				Computed: true,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"deployment_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"disk_content_base_64": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"persistent": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"source_reference": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
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
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": linksSchema(),
			"status": {
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

func resourceBlockDeviceCreate(d *schema.ResourceData, m interface{}) error {
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
		return err
	}

	stateChangeFunc := resource.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    blockDeviceStateRefreshFunc(*apiClient, *createBlockDeviceCreated.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
	}

	resourceIDs, err := stateChangeFunc.WaitForState()
	if err != nil {
		return err
	}

	blockDeviceIDs := resourceIDs.([]string)
	i := strings.LastIndex(blockDeviceIDs[0], "/")
	blockDeviceID := blockDeviceIDs[0][i+1 : len(blockDeviceIDs[0])]
	d.SetId(blockDeviceID)
	log.Printf("Finished to create vra_block_device resource with name %s", d.Get("name"))

	return resourceBlockDeviceRead(d, m)
}

func blockDeviceStateRefreshFunc(apiClient client.MulticloudIaaS, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ret, err := apiClient.Request.GetRequestTracker(request.NewGetRequestTrackerParams().WithID(id))
		if err != nil {
			return "", models.RequestTrackerStatusFAILED, err
		}

		status := ret.Payload.Status
		switch *status {
		case models.RequestTrackerStatusFAILED:
			return []string{""}, *status, fmt.Errorf(ret.Payload.Message)
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

func resourceBlockDeviceRead(d *schema.ResourceData, m interface{}) error {
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
		return err
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
		return fmt.Errorf("error setting block device tags - error: %v", err)
	}

	if err := d.Set("links", flattenLinks(blockDevice.Links)); err != nil {
		return fmt.Errorf("error setting block device links - error: %#v", err)
	}

	log.Printf("Finished reading the vra_block_device resource with name %s", d.Get("name"))
	return nil
}

func resourceBlockDeviceUpdate(d *schema.ResourceData, m interface{}) error {

	log.Printf("Starting to update the vra_block_device resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	if d.HasChange("capacity_in_gb") {
		err := resizeDisk(d, apiClient, id)
		if err != nil {
			return err
		}
	}

	log.Printf("Finished updating vra_block_device resource with name %s", d.Get("name"))
	return resourceBlockDeviceRead(d, m)
}

func resizeDisk(d *schema.ResourceData, apiClient *client.MulticloudIaaS, id string) error {

	log.Printf("Starting resize of vra_block_device resource with name %s", d.Get("name"))

	capacityInGB := int32(d.Get("capacity_in_gb").(int))
	resizeBlockDeviceAccepted, resizeBlockDeviceNoContent, err := apiClient.Disk.ResizeBlockDevice(disk.NewResizeBlockDeviceParams().WithID(id).WithCapacityInGB(capacityInGB))
	if err != nil {
		return err
	}
	if resizeBlockDeviceNoContent != nil {
		return nil
	}

	stateChangeFunc := resource.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    blockDeviceStateRefreshFunc(*apiClient, *resizeBlockDeviceAccepted.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 5 * time.Second,
	}

	_, err = stateChangeFunc.WaitForState()
	if err != nil {
		return err
	}

	return nil
}

func resourceBlockDeviceDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("Starting to delete the vra_block_device resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	deleteBlockDeviceAccepted, deleteBlockDeviceCompleted, err := apiClient.Disk.DeleteBlockDevice(disk.NewDeleteBlockDeviceParams().WithID(id))
	if err != nil {
		return err
	}

	// Handle non-request tracker case
	if deleteBlockDeviceCompleted != nil {
		return nil
	}

	stateChangeFunc := resource.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    blockDeviceStateRefreshFunc(*apiClient, *deleteBlockDeviceAccepted.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 5 * time.Second,
	}

	_, err = stateChangeFunc.WaitForState()
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("Finished deleting the vra_block_device resource with name %s", d.Get("name"))
	return nil
}
