package cas

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/cas-sdk-go/pkg/client/location"
	"github.com/vmware/cas-sdk-go/pkg/models"

	tango "github.com/vmware/terraform-provider-cas/sdk"
)

func resourceZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceZoneCreate,
		Read:   resourceZoneRead,
		Update: resourceZoneUpdate,
		Delete: resourceZoneDelete,

		Schema: map[string]*schema.Schema{
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"placement_policy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "DEFAULT",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "DEFAULT" && value != "SPREAD" && value != "BINPACK" {
						errors = append(errors, fmt.Errorf(
							"%q must be one of 'DEFAULT', 'SPREAD', 'BINPACK'", k))
					}
					return
				},
			},
			"region_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"tags":        tagsSDKSchema(),
			"tagstomatch": tagsSDKSchema(),
		},
	}
}

func resourceZoneCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)
	apiClient := client.GetAPIClient()

	description := d.Get("description").(string)
	name := d.Get("name").(string)
	placementPolicy := d.Get("placement_policy").(string)
	regionID := d.Get("region_id").(string)
	tags := expandSDKTags(d.Get("tags").(*schema.Set).List())
	tagsToMatch := expandSDKTags(d.Get("tagstomatch").(*schema.Set).List())

	createResp, err := apiClient.Location.CreateZone(location.NewCreateZoneParams().WithBody(&models.ZoneSpecification{
		Description:     description,
		Name:            &name,
		PlacementPolicy: placementPolicy,
		RegionID:        &regionID,
		Tags:            tags,
		TagsToMatch:     tagsToMatch,
	}))
	if err != nil {
		return err
	}

	if err := d.Set("tags", flattenSDKTags(tags)); err != nil {
		return fmt.Errorf("Error setting zone tags - error: %#v", err)
	}
	if err := d.Set("tagstomatch", flattenSDKTags(tagsToMatch)); err != nil {
		return fmt.Errorf("Error setting zone tagstomatch - error: %#v", err)
	}
	d.SetId(*createResp.Payload.ID)

	return resourceZoneRead(d, m)
}

func resourceZoneRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)
	apiClient := client.GetAPIClient()

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
	if err := d.Set("tags", flattenSDKTags(zone.Tags)); err != nil {
		return fmt.Errorf("Error setting zone tags - error: %#v", err)
	}
	if err := d.Set("tagstomatch", flattenSDKTags(zone.TagsToMatch)); err != nil {
		return fmt.Errorf("Error setting zone tagstomatch - error: %#v", err)
	}
	return nil
}

func resourceZoneUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)
	apiClient := client.GetAPIClient()

	id := d.Id()
	description := d.Get("description").(string)
	name := d.Get("name").(string)
	placementPolicy := d.Get("placement_policy").(string)
	regionID := d.Get("region_id").(string)
	tags := expandSDKTags(d.Get("tags").(*schema.Set).List())
	tagsToMatch := expandSDKTags(d.Get("tagstomatch").(*schema.Set).List())

	_, err := apiClient.Location.UpdateZone(location.NewUpdateZoneParams().WithID(id).WithBody(&models.ZoneSpecification{
		Description:     description,
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
	client := m.(*tango.Client)
	apiClient := client.GetAPIClient()

	id := d.Id()
	_, err := apiClient.Location.DeleteZone(location.NewDeleteZoneParams().WithID(id))
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
