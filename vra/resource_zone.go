// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/location"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func resourceZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceZoneCreate,
		ReadContext:   resourceZoneRead,
		UpdateContext: resourceZoneUpdate,
		DeleteContext: resourceZoneDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required arguments
			"name": {
				Type:        schema.TypeString,
				Description: "A human-friendly name used as an identifier for the zone resource instance.",
				Required:    true,
			},
			"region_id": {
				Type:        schema.TypeString,
				Description: "The id of the region for which this zone is created.",
				Required:    true,
			},

			// Optional arguments
			"compute_ids": {
				Type:        schema.TypeSet,
				Computed:    true, // it needs to be computed because vRA will add compute ids besides the ones specified in the terraform plan
				Description: "The ids of the compute resources that will be explicitly assigned to this zone.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"custom_properties": {
				Type:        schema.TypeMap,
				Computed:    true, // it needs to be computed because vRA will add its own custom properties besides the ones specified in the terraform plan
				Description: "A list of key value pair of properties that will be used.",
				Optional:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "A human-friendly description.",
				Optional:    true,
			},
			"folder": {
				Type:        schema.TypeString,
				Description: "The folder relative path to the datacenter where resources are deployed to (only applicable for vSphere cloud zones).",
				Optional:    true,
			},
			"placement_policy": {
				Type:        schema.TypeString,
				Default:     "DEFAULT",
				Description: "The placement policy for the zone. One of DEFAULT, SPREAD or BINPACK.",
				Optional:    true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "DEFAULT" && value != "SPREAD" && value != "BINPACK" {
						errors = append(errors, fmt.Errorf(
							"%q must be one of 'DEFAULT', 'SPREAD', 'BINPACK'", k))
					}
					return
				},
			},
			"tags":          tagsSchema(),
			"tags_to_match": tagsSchema(),

			// Computed attributes
			"cloud_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the cloud account this zone belongs to.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was created. The date is in ISO 8601 and UTC.",
			},
			"external_region_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the region for which this zone is defined.",
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
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC.",
			},
		},
	}
}

func resourceZoneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	name := d.Get("name").(string)
	regionID := d.Get("region_id").(string)

	var computeIDs []string
	if v, ok := d.GetOk("compute_ids"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.FromErr(errors.New("specified compute_ids are not unique"))
		}
		computeIDs = expandStringList(v.(*schema.Set).List())
	}

	createResp, err := apiClient.Location.CreateZone(location.NewCreateZoneParams().WithBody(&models.ZoneSpecification{
		ComputeIds:       computeIDs,
		CustomProperties: expandCustomProperties(d.Get("custom_properties").(map[string]interface{})),
		Description:      d.Get("description").(string),
		Folder:           d.Get("folder").(string),
		Name:             &name,
		PlacementPolicy:  d.Get("placement_policy").(string),
		RegionID:         &regionID,
		Tags:             expandTags(d.Get("tags").(*schema.Set).List()),
		TagsToMatch:      expandTags(d.Get("tags_to_match").(*schema.Set).List()),
	}))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*createResp.Payload.ID)

	return resourceZoneRead(ctx, d, m)
}

func resourceZoneRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	getResp, err := apiClient.Location.GetZone(location.NewGetZoneParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *location.GetZoneNotFound:
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	zone := *getResp.Payload
	d.Set("cloud_account_id", zone.CloudAccountID)
	d.Set("created_at", zone.CreatedAt)
	d.Set("custom_properties", zone.CustomProperties)
	d.Set("description", zone.Description)
	d.Set("external_region_id", zone.ExternalRegionID)
	d.Set("folder", zone.Folder)
	d.Set("name", zone.Name)
	d.Set("org_id", zone.OrgID)
	d.Set("owner", zone.Owner)
	d.Set("placement_policy", zone.PlacementPolicy)
	d.Set("updated_at", zone.UpdatedAt)

	if err := d.Set("links", flattenLinks(zone.Links)); err != nil {
		return diag.Errorf("error setting zone links - error: %#v", err)
	}

	if err := d.Set("tags", flattenTags(zone.Tags)); err != nil {
		return diag.Errorf("error setting zone tags - error: %v", err)
	}

	if err := d.Set("tags_to_match", flattenTags(zone.TagsToMatch)); err != nil {
		return diag.Errorf("error setting zone tags to match - error: %v", err)
	}

	getComputesResp, err := apiClient.Location.GetComputes(location.NewGetComputesParams().WithID(id))
	if err != nil {
		return diag.Errorf("error getting zone computes - error: %v", err)
	}

	var computeIDs []string
	for _, compute := range getComputesResp.Payload.Content {
		computeIDs = append(computeIDs, *compute.ID)
	}
	d.Set("compute_ids", computeIDs)

	return nil
}

func resourceZoneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	name := d.Get("name").(string)
	regionID := d.Get("region_id").(string)

	var computeIDs []string
	if v, ok := d.GetOk("compute_ids"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.FromErr(errors.New("specified compute_ids are not unique"))
		}
		computeIDs = expandStringList(v.(*schema.Set).List())
	}

	if _, err := apiClient.Location.UpdateZone(location.NewUpdateZoneParams().WithID(id).WithBody(&models.ZoneSpecification{
		ComputeIds:       computeIDs,
		CustomProperties: expandCustomProperties(d.Get("custom_properties").(map[string]interface{})),
		Description:      d.Get("description").(string),
		Folder:           d.Get("folder").(string),
		Name:             &name,
		PlacementPolicy:  d.Get("placement_policy").(string),
		RegionID:         &regionID,
		Tags:             expandTags(d.Get("tags").(*schema.Set).List()),
		TagsToMatch:      expandTags(d.Get("tags_to_match").(*schema.Set).List()),
	})); err != nil {
		return diag.FromErr(err)
	}

	return resourceZoneRead(ctx, d, m)
}

func resourceZoneDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	if _, err := apiClient.Location.DeleteZone(location.NewDeleteZoneParams().WithID(id)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
