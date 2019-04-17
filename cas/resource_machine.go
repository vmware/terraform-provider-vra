package cas

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/vmware/cas-sdk-go/pkg/client"
	"github.com/vmware/cas-sdk-go/pkg/client/compute"
	"github.com/vmware/cas-sdk-go/pkg/client/request"
	"github.com/vmware/cas-sdk-go/pkg/models"
	"log"
	"strings"
	"time"

	"github.com/vmware/terraform-provider-cas/sdk"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceMachineCreate,
		Read:   resourceMachineRead,
		Update: resourceMachineUpdate,
		Delete: resourceMachineDelete,

		Schema: map[string]*schema.Schema{
			"flavor": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"image": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_ref": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"power_state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"constraints": constraintsSDKSchema(),
			"tags":        tagsSDKSchema(),
			"custom_properties": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Optional: true,
			},
			"nics": nicsSDKSchema(false),
			"disks": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"block_device_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"boot_config": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"external_zone_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_region_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !strings.HasPrefix(new, old)
				},
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": linksSchema(),
		},
	}
}

func resourceMachineCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("Starting to create cas_machine resource")
	client := m.(*tango.Client)
	apiClient := client.GetAPIClient()

	name := d.Get("name").(string)
	flavor := d.Get("flavor").(string)
	projectId := d.Get("project_id").(string)
	constraints := expandSDKConstraints(d.Get("constraints").(*schema.Set).List())
	tags := expandSDKTags(d.Get("tags").(*schema.Set).List())
	customProperties := expandCustomProperties(d.Get("custom_properties").(map[string]interface{}))
	nics := expandSDKNics(d.Get("nics").(*schema.Set).List())
	disks := expandSDKDisks(d.Get("disks").(*schema.Set).List())

	machineSpecification := models.MachineSpecification{
		Name:             &name,
		Flavor:           &flavor,
		ProjectID:        &projectId,
		Constraints:      constraints,
		Tags:             tags,
		CustomProperties: customProperties,
		Nics:             nics,
		Disks:            disks,
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
		return errors.New("image or image_ref required")
	}

	if v, ok := d.GetOk("description"); ok {
		machineSpecification.Description = v.(string)
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
		return err
	}

	stateChangeFunc := resource.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    machineStateRefreshFunc(*apiClient, *createMachineCreated.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    5 * time.Minute,
		MinTimeout: 5 * time.Second,
	}

	resourceIds, err := stateChangeFunc.WaitForState()
	if err != nil {
		return err
	}

	machineIds := resourceIds.([]string)
	d.SetId(machineIds[0])
	log.Printf("Finished to create cas_machine resource with name %s", d.Get("name"))

	return resourceMachineRead(d, m)
}

func machineStateRefreshFunc(apiClient client.MulticloudIaaS, id string) resource.StateRefreshFunc {
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
			machineIds := make([]string, len(ret.Payload.Resources))
			for i, r := range ret.Payload.Resources {
				machineIds[i] = strings.TrimPrefix(r, "/iaas/api/machines/")
			}
			return machineIds, *status, nil
		default:
			return [...]string{id}, ret.Payload.Message, fmt.Errorf("machineStateRefreshFunc: unknown status %v", *status)
		}
	}
}

func resourceMachineRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("Reading the cas_machine resource with name %s", d.Get("name"))
	client := m.(*tango.Client)
	apiClient := client.GetAPIClient()

	id := d.Id()
	resp, err := apiClient.Compute.GetMachine(compute.NewGetMachineParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *compute.GetMachineNotFound:
			d.SetId("")
			return nil
		}
		return err
	}

	machine := *resp.Payload
	d.Set("name", machine.Name)
	d.Set("description", machine.Description)
	d.Set("power_state", machine.PowerState)
	d.Set("address", machine.Address)
	d.Set("project_id", machine.ProjectID)
	d.Set("external_zone_id", machine.ExternalZoneID)
	d.Set("external_region_id", machine.ExternalRegionID)
	d.Set("external_id", machine.ExternalID)
	d.Set("created_at", machine.CreatedAt)
	d.Set("updated_at", machine.UpdatedAt)
	d.Set("owner", machine.Owner)
	d.Set("organization_id", machine.OrganizationID)
	d.Set("custom_properties", machine.CustomProperties)

	if err := d.Set("tags", flattenSDKTags(machine.Tags)); err != nil {
		return fmt.Errorf("error setting machine tags - error: %v", err)
	}

	if err := d.Set("links", flattenSDKLinks(machine.Links)); err != nil {
		return fmt.Errorf("error setting machine links - error: %#v", err)
	}

	log.Printf("Finished reading the cas_machine resource with name %s", d.Get("name"))
	return nil
}

func resourceMachineUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("Starting to update the cas_machine resource with name %s", d.Get("name"))
	client := m.(*tango.Client)
	apiClient := client.GetAPIClient()

	image, imageRef := "", ""

	if v, ok := d.GetOk("image"); ok {
		image = v.(string)
	}

	if v, ok := d.GetOk("image_ref"); ok {
		imageRef = v.(string)
	}

	id := d.Id()
	name := d.Get("name").(string)
	flavor := d.Get("flavor").(string)
	projectId := d.Get("project_id").(string)
	description := d.Get("description").(string)
	constraints := expandSDKConstraints(d.Get("constraints").(*schema.Set).List())
	tags := expandSDKTags(d.Get("tags").(*schema.Set).List())
	customProperties := expandCustomProperties(d.Get("custom_properties").(map[string]interface{}))
	nics := expandSDKNics(d.Get("nics").(*schema.Set).List())
	disks := expandSDKDisks(d.Get("disks").(*schema.Set).List())

	machineSpecification := models.MachineSpecification{
		Name:             &name,
		Image:            &image,
		ImageRef:         &imageRef,
		Description:      description,
		Flavor:           &flavor,
		ProjectID:        &projectId,
		Constraints:      constraints,
		Tags:             tags,
		CustomProperties: customProperties,
		Nics:             nics,
		Disks:            disks,
	}

	log.Printf("[DEBUG] update machine: %#v", machineSpecification)
	_, err := apiClient.Compute.UpdateMachine(compute.NewUpdateMachineParams().WithID(id).WithBody(&machineSpecification))
	if err != nil {
		return err
	}

	log.Printf("Finished updating the cas_machine resource with name %s", d.Get("name"))
	return resourceMachineRead(d, m)
}

func resourceMachineDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("Starting to delete the cas_machine resource with name %s", d.Get("name"))
	client := m.(*tango.Client)
	apiClient := client.GetAPIClient()

	id := d.Id()
	_, err := apiClient.Compute.DeleteMachine(compute.NewDeleteMachineParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *compute.DeleteMachineNotFound:
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId("")
	log.Printf("Finished deleting the cas_machine resource with name %s", d.Get("name"))
	return nil
}

func expandSDKDisks(configDisks []interface{}) []*models.DiskAttachmentSpecification {
	disks := make([]*models.DiskAttachmentSpecification, 0, len(configDisks))

	for _, configDisk := range configDisks {
		diskMap := configDisk.(map[string]interface{})

		disk := models.DiskAttachmentSpecification{
			BlockDeviceID: withString(diskMap["block_device_id"].(string)),
		}

		if v, ok := diskMap["name"].(string); ok && v != "" {
			disk.Name = withString(v)
		}

		if v, ok := diskMap["description"].(string); ok && v != "" {
			disk.Description = v
		}

		disks = append(disks, &disk)
	}

	return disks
}
