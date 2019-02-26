package cas

import (
	"fmt"
	"log"
	"strings"

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
			"image": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			//"image_ref": &schema.Schema{
			//	Type:     schema.TypeString,
			//	Required: true,
			//},
			"power_state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"machine_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"constraints": constraintsSchema(),
			"tags":        tagsSchema(),
			"custom_properties": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Optional: true,
			},
			"nics": nicsSchema(false),
			"disks": &schema.Schema{
				Type:     schema.TypeList,
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
				Type:     schema.TypeList,
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
			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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
	client := m.(*tango.Client)

	machineSpecification := tango.MachineSpecification{
		Name:             d.Get("name").(string),
		Image:            d.Get("image").(string),
		Flavor:           d.Get("flavor").(string),
		ProjectID:        client.GetProjectID(),
		MachineCount:     d.Get("machine_count").(int),
		Constraints:      expandConstraints(d.Get("constraints").([]interface{})),
		Tags:             expandTags(d.Get("tags").([]interface{})),
		CustomProperties: expandCustomProperties(d.Get("custom_properties").(map[string]interface{})),
		Nics:             expandNics(d.Get("nics").([]interface{})),
		Disks:            expandDisks(d.Get("disks").([]interface{})),
	}

	machineSpecification.CustomProperties["__composition_context_id"] = client.GetDeploymentID()

	if v, ok := d.GetOk("description"); ok {
		machineSpecification.Description = v.(string)
	}

	if v, ok := d.GetOk("boot_config"); ok {
		configBootConfig := v.([]interface{})[0].(map[string]interface{})

		bootConfig := make(map[string]string)
		bootConfig["content"] = configBootConfig["content"].(string)

		machineSpecification.BootConfig = bootConfig
	}

	log.Printf("[DEBUG] record create machine: %#v", machineSpecification)
	resourceObject, err := client.CreateResource(machineSpecification)
	if err != nil {
		return err
	}

	machineObject := resourceObject.(*tango.Machine)

	d.SetId(machineObject.ID)
	d.Set("name", machineObject.Name)
	d.Set("power_state", machineObject.PowerState)
	d.Set("address", machineObject.Address)
	d.Set("external_zone_id", machineObject.ExternalZoneID)
	d.Set("external_region_id", machineObject.ExternalRegionID)
	d.Set("external_id", machineObject.ExternalID)
	d.Set("self_link", machineObject.SelfLink)
	d.Set("created_at", machineObject.CreatedAt)
	d.Set("updated_at", machineObject.UpdatedAt)
	d.Set("owner", machineObject.Owner)
	d.Set("organization_id", machineObject.OrganizationID)
	d.Set("custom_properties", machineObject.CustomProperties)

	if err := d.Set("tags", flattenTags(machineObject.Tags)); err != nil {
		return fmt.Errorf("Error setting Machine tags - error: %#v", err)
	}

	if err := d.Set("links", flattenLinks(machineObject.Links)); err != nil {
		return fmt.Errorf("Error setting Machine links - error: %#v", err)
	}

	return nil
}

func expandDisks(configDisks []interface{}) []tango.Disk {
	disks := make([]tango.Disk, 0, len(configDisks))

	for _, configDisk := range configDisks {
		diskMap := configDisk.(map[string]interface{})

		disk := tango.Disk{
			BlockDeviceID: diskMap["block_device_id"].(string),
		}

		if v, ok := diskMap["name"].(string); ok && v != "" {
			disk.Name = v
		}

		if v, ok := diskMap["description"].(string); ok && v != "" {
			disk.Description = v
		}

		disks = append(disks, disk)
	}

	return disks
}

func resourceMachineRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)

	resourceObject, err := client.ReadResource(getSelfLink(d.Get("links").([]interface{})))
	if err != nil {
		d.SetId("")
		return nil
	}

	machineObject := resourceObject.(*tango.Machine)

	d.Set("power_state", machineObject.PowerState)
	d.Set("address", machineObject.Address)
	d.Set("external_zone_id", machineObject.ExternalZoneID)
	d.Set("external_region_id", machineObject.ExternalRegionID)
	d.Set("external_id", machineObject.ExternalID)
	d.Set("name", machineObject.Name)
	d.Set("description", machineObject.Description)
	d.Set("self_link", machineObject.SelfLink)
	d.Set("created_at", machineObject.CreatedAt)
	d.Set("updated_at", machineObject.UpdatedAt)
	d.Set("owner", machineObject.Owner)
	d.Set("organization_id", machineObject.OrganizationID)
	d.Set("custom_properties", machineObject.CustomProperties)

	if err := d.Set("tags", flattenTags(machineObject.Tags)); err != nil {
		return fmt.Errorf("Error setting Machine tags - error: %#v", err)
	}

	if err := d.Set("links", flattenLinks(machineObject.Links)); err != nil {
		return fmt.Errorf("Error setting Machine links - error: %#v", err)
	}

	return nil
}

func resourceMachineUpdate(d *schema.ResourceData, m interface{}) error {

	return nil
}

func resourceMachineDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)

	// TODO: Confirm that deleting the machine also deletes all the attached disks automatically.
	// Commented the code as the machine deletion is failing while trying to detach the boot-disk.
	//resourceObject, err := client.ReadResource(getSelfLink(d.Get("links").([]interface{})) + "/disks")
	//if err != nil {
	//	return err
	//}
	//
	//machineAttachedDisksObject := resourceObject.(*tango.MachineAttachedDisks)
	//
	//peripheralDevicesRegex := regexp.MustCompile("^(CD/DVD|Floppy) drive")
	//for _, blockDevice := range machineAttachedDisksObject.Content {
	//	if blockDevice.Name != "boot-disk" && !peripheralDevicesRegex.MatchString(blockDevice.Name) {
	//		err := client.DeleteResource(blockDevice.Links["self"].Href) // detach disks first
	//		if err != nil {
	//			return err
	//		}
	//	}
	//}

	err := client.DeleteResource(getSelfLink(d.Get("links").([]interface{})))

	if err != nil && strings.Contains(err.Error(), "404") { // already deleted
		return nil
	}

	return err
}
