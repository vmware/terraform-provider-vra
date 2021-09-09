package vra

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/location"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceZoneRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Computed:      true,
				ConflictsWith: []string{"name"},
				Description:   "The id of the zone resource instance.",
				Optional:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Computed:      true,
				ConflictsWith: []string{"id"},
				Description:   "A human-friendly name used as an identifier for the zone resource instance.",
				Optional:      true,
			},
			"cloud_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the cloud account this zone belongs to.",
			},
			"compute_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The ids of the compute resources that has been explicitly assigned to this zone.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was created. The date is in ISO 8601 and UTC.",
			},
			"custom_properties": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "A list of key value pair of properties that will be used.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly description.",
			},
			"external_region_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the region for which this zone is defined.",
			},
			"folder": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The folder relative path to the datacenter where resources are deployed to (only applicable for vSphere cloud zones).",
			},
			"links": linksSchema(),
			"org_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the organization this entity belongs to.",
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email of the user that owns the entity.",
			},
			"placement_policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The placement policy for the zone.",
			},
			"tags":          tagsSchema(),
			"tags_to_match": tagsSchema(),
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC.",
			},
		},
	}
}

func dataSourceZoneRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if !idOk && !nameOk {
		return errors.New("one of id or name must be assigned")
	}

	getResp, err := apiClient.Location.GetZones(location.NewGetZonesParams())
	if err != nil {
		return err
	}

	setFields := func(zone *models.Zone) error {
		d.SetId(*zone.ID)
		d.Set("name", zone.Name)
		d.Set("cloud_account_id", zone.CloudAccountID)
		d.Set("created_at", zone.CreatedAt)
		d.Set("custom_properties", zone.CustomProperties)
		d.Set("description", zone.Description)
		d.Set("external_region_id", zone.ExternalRegionID)
		d.Set("folder", zone.Folder)
		d.Set("org_id", zone.OrgID)
		d.Set("owner", zone.Owner)
		d.Set("placement_policy", zone.PlacementPolicy)
		d.Set("updated_at", zone.UpdatedAt)

		if err := d.Set("links", flattenLinks(zone.Links)); err != nil {
			return fmt.Errorf("error setting zone links - error: %#v", err)
		}

		if err := d.Set("tags", flattenTags(zone.Tags)); err != nil {
			return fmt.Errorf("error setting zone tags - error: %v", err)
		}

		if err := d.Set("tags_to_match", flattenTags(zone.TagsToMatch)); err != nil {
			return fmt.Errorf("error setting zone tags to match - error: %v", err)
		}

		getComputesResp, err := apiClient.Location.GetComputes(location.NewGetComputesParams().WithID(*zone.ID))
		if err != nil {
			return fmt.Errorf("error getting zone computes - error: %v", err)
		}

		var computeIds []string
		for _, compute := range getComputesResp.Payload.Content {
			computeIds = append(computeIds, *compute.ID)
		}
		d.Set("compute_ids", computeIds)

		return nil
	}

	for _, zone := range getResp.Payload.Content {
		if idOk && *zone.ID == id {
			return setFields(zone)
		}
		if nameOk && zone.Name == name {
			return setFields(zone)
		}
	}

	return fmt.Errorf("zone with id `%s` or name `%s` not found", id, name)
}
