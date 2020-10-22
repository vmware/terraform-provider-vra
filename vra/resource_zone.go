package vra

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/location"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func resourceZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceZoneCreate,
		Read:   resourceZoneRead,
		Update: resourceZoneUpdate,
		Delete: resourceZoneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A human-friendly description for the zone",
			},
			"folder": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The folder relative path to the datacenter where resources are deployed to. (only applicable for vSphere cloud zones)",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A human-friendly name for the zone",
			},
			"placement_policy": {
				Type:     schema.TypeString,
				Required: true,
				Default:  "DEFAULT",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "DEFAULT" && value != "SPREAD" && value != "BINPACK" {
						errors = append(errors, fmt.Errorf(
							"%q must be one of 'DEFAULT', 'SPREAD', 'BINPACK'", k))
					}
					return
				},
				Description: "Placement policy for the zone. One of DEFAULT, SPREAD or BINPACK.",
			},
			"region_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the region for which this profile is created",
			},
			"tags":          tagsSchema(),
			"tags_to_match": tagsSchema(),
		},
	}
}

func resourceZoneCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	description := d.Get("description").(string)
	folder := d.Get("folder").(string)
	name := d.Get("name").(string)
	placementPolicy := d.Get("placement_policy").(string)
	regionID := d.Get("region_id").(string)
	tags := expandTags(d.Get("tags").(*schema.Set).List())
	tagsToMatch := expandTags(d.Get("tags_to_match").(*schema.Set).List())

	createResp, err := apiClient.Location.CreateZone(location.NewCreateZoneParams().WithBody(&models.ZoneSpecification{
		Description:     description,
		Folder:          folder,
		Name:            &name,
		PlacementPolicy: placementPolicy,
		RegionID:        &regionID,
		Tags:            tags,
		TagsToMatch:     tagsToMatch,
	}))
	if err != nil {
		return err
	}

	if err := d.Set("tags", flattenTags(tags)); err != nil {
		return fmt.Errorf("Error setting zone tags - error: %#v", err)
	}
	if err := d.Set("tags_to_match", flattenTags(tagsToMatch)); err != nil {
		return fmt.Errorf("Error setting zone tags_to_match - error: %#v", err)
	}
	d.SetId(*createResp.Payload.ID)

	return resourceZoneRead(d, m)
}

func resourceZoneRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	ret, err := apiClient.Location.GetZone(location.NewGetZoneParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *location.GetZoneNotFound:
			d.SetId("")
			return nil
		}
		return err
	}
	zone := *ret.Payload
	d.Set("description", zone.Description)
	d.Set("name", zone.Name)
	d.Set("placement_policy", zone.PlacementPolicy)
	if err := d.Set("tags", flattenTags(zone.Tags)); err != nil {
		return fmt.Errorf("Error setting zone tags - error: %#v", err)
	}
	if err := d.Set("tags_to_match", flattenTags(zone.TagsToMatch)); err != nil {
		return fmt.Errorf("Error setting zone tags_to_match - error: %#v", err)
	}
	return nil
}

func resourceZoneUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	description := d.Get("description").(string)
	folder := d.Get("folder").(string)
	name := d.Get("name").(string)
	placementPolicy := d.Get("placement_policy").(string)
	regionID := d.Get("region_id").(string)
	tags := expandTags(d.Get("tags").(*schema.Set).List())
	tagsToMatch := expandTags(d.Get("tags_to_match").(*schema.Set).List())

	_, err := apiClient.Location.UpdateZone(location.NewUpdateZoneParams().WithID(id).WithBody(&models.ZoneSpecification{
		Description:     description,
		Folder:          folder,
		Name:            &name,
		PlacementPolicy: placementPolicy,
		RegionID:        &regionID,
		Tags:            tags,
		TagsToMatch:     tagsToMatch,
	}))
	if err != nil {
		return err
	}

	return resourceZoneRead(d, m)
}

func resourceZoneDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	_, err := apiClient.Location.DeleteZone(location.NewDeleteZoneParams().WithID(id))
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
