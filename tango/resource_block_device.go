package tango

import (
	"fmt"
	"log"
	"strings"

	"tango-terraform-provider/tango/client"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBlockDevice() *schema.Resource {
	return &schema.Resource{
		Create: resourceBlockDeviceCreate,
		Read:   resourceBlockDeviceRead,
		Update: resourceBlockDeviceUpdate,
		Delete: resourceBlockDeviceDelete,

		Schema: map[string]*schema.Schema{
			"capacity_in_gb": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"encrypted": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"source_reference": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk_content_base_64": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_properties": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Optional: true,
			},
			"constraints": constraintsSchema(),
			"tags":        tagsSchema(),
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

func resourceBlockDeviceCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)

	blockDeviceSpecification := tango.BlockDeviceSpecification{
		Name:             d.Get("name").(string),
		ProjectID:        client.GetProjectID(),
		CapacityInGB:     d.Get("capacity_in_gb").(int),
		Constraints:      expandConstraints(d.Get("constraints").([]interface{})),
		Tags:             expandTags(d.Get("tags").([]interface{})),
		CustomProperties: expandCustomProperties(d.Get("custom_properties").(map[string]interface{})),
	}

	blockDeviceSpecification.CustomProperties["__composition_context_id"] = client.GetDeploymentID()

	if v, ok := d.GetOk("encrypted"); ok {
		blockDeviceSpecification.Encrypted = v.(bool)
	}

	if v, ok := d.GetOk("source_reference"); ok {
		blockDeviceSpecification.SourceReference = v.(string)
	}

	if v, ok := d.GetOk("disk_content_base_64"); ok {
		blockDeviceSpecification.DiskContentBase64 = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		blockDeviceSpecification.Description = v.(string)
	}

	log.Printf("[DEBUG] record create block device: %#v", blockDeviceSpecification)
	resourceObject, err := client.CreateResource(blockDeviceSpecification)
	if err != nil {
		return err
	}

	blockDeviceObject := resourceObject.(*tango.BlockDevice)

	d.SetId(blockDeviceObject.ID)
	d.Set("name", blockDeviceObject.Name)
	d.Set("capacity_in_gb", blockDeviceObject.CapacityInGB)
	d.Set("status", blockDeviceObject.Status)
	d.Set("external_zone_id", blockDeviceObject.ExternalZoneID)
	d.Set("external_region_id", blockDeviceObject.ExternalRegionID)
	d.Set("external_id", blockDeviceObject.ExternalID)
	d.Set("self_link", blockDeviceObject.SelfLink)
	d.Set("created_at", blockDeviceObject.CreatedAt)
	d.Set("updated_at", blockDeviceObject.UpdatedAt)
	d.Set("owner", blockDeviceObject.Owner)
	d.Set("organization_id", blockDeviceObject.OrganizationID)
	d.Set("custom_properties", blockDeviceObject.CustomProperties)

	if err := d.Set("tags", flattenTags(blockDeviceObject.Tags)); err != nil {
		return fmt.Errorf("Error setting Block Device tags - error: %#v", err)
	}

	if err := d.Set("links", flattenLinks(blockDeviceObject.Links)); err != nil {
		return fmt.Errorf("Error setting Block Device links - error: %#v", err)
	}

	return nil
}

func resourceBlockDeviceRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)

	resourceObject, err := client.ReadResource(getSelfLink(d.Get("links").([]interface{})))
	if err != nil {
		d.SetId("")
		return nil
	}

	blockDeviceObject := resourceObject.(*tango.BlockDevice)

	d.Set("capacity_in_gb", blockDeviceObject.CapacityInGB)
	d.Set("status", blockDeviceObject.Status)
	d.Set("external_zone_id", blockDeviceObject.ExternalZoneID)
	d.Set("external_region_id", blockDeviceObject.ExternalRegionID)
	d.Set("external_id", blockDeviceObject.ExternalID)
	d.Set("name", blockDeviceObject.Name)
	d.Set("description", blockDeviceObject.Description)
	d.Set("self_link", blockDeviceObject.SelfLink)
	d.Set("created_at", blockDeviceObject.CreatedAt)
	d.Set("updated_at", blockDeviceObject.UpdatedAt)
	d.Set("owner", blockDeviceObject.Owner)
	d.Set("organization_id", blockDeviceObject.OrganizationID)
	d.Set("custom_properties", blockDeviceObject.CustomProperties)

	if err := d.Set("tags", flattenTags(blockDeviceObject.Tags)); err != nil {
		return fmt.Errorf("Error setting Block Device tags - error: %#v", err)
	}

	if err := d.Set("links", flattenLinks(blockDeviceObject.Links)); err != nil {
		return fmt.Errorf("Error setting Block Device links - error: %#v", err)
	}

	return nil
}

func resourceBlockDeviceUpdate(d *schema.ResourceData, m interface{}) error {

	return nil
}

func resourceBlockDeviceDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)
	err := client.DeleteResource(getSelfLink(d.Get("links").([]interface{})))

	if err != nil && strings.Contains(err.Error(), "404") { // already deleted
		return nil
	}

	return err
}
