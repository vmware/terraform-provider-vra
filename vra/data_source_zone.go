package vra

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/location"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceZoneRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_properties": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"folder": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"org_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"placement_policy": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags":          tagsSchema(),
			"tags_to_match": tagsSchema(),
			"updated_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceZoneRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if !idOk && !nameOk {
		return fmt.Errorf("One of id or name must be assigned")
	}
	getResp, err := apiClient.Location.GetZones(location.NewGetZonesParams())

	if err != nil {
		return err
	}

	setFields := func(zone *models.Zone) {
		d.SetId(*zone.ID)
		d.Set("description", zone.Description)
		d.Set("name", zone.Name)
		d.Set("created_at", zone.CreatedAt)
		d.Set("custom_properties", zone.CustomProperties)
		d.Set("folder", zone.Folder)
		d.Set("org_id", zone.OrgID)
		d.Set("owner", zone.Owner)
		d.Set("placement_policy", zone.PlacementPolicy)
		d.Set("tags", zone.Tags)
		d.Set("tags_to_match", zone.TagsToMatch)
		d.Set("updated_at", zone.UpdatedAt)
	}
	for _, zone := range getResp.Payload.Content {
		if idOk && zone.ID == id {
			setFields(zone)
			return nil
		}
		if nameOk && zone.Name == name {
			setFields(zone)
			return nil
		}
	}

	return fmt.Errorf("zone %s not found", name)
}
