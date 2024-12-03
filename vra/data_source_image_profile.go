// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/image_profile"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"fmt"
	"log"
	"strings"
)

func dataSourceImageProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceImageProfileRead,

		Schema: map[string]*schema.Schema{
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was created. The date is in ISO 8601 and UTC.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A human-friendly description.",
			},
			"external_region_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the region for which this profile is defined as in the cloud provider.",
			},
			"filter": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Filter query string that is supported by vRA multi-cloud IaaS API. Example: regionId eq '<regionId>' and cloudAccountId eq '<cloudAccountId>'. Only one of 'filter', 'id', 'name' or 'region_id' must be specified.",
				ConflictsWith: []string{"id", "name", "region_id"},
			},
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Description:   "The id of the image profile instance.  Only one of 'filter', 'id', 'name' or 'region_id' must be specified.",
				ConflictsWith: []string{"filter", "name", "region_id"},
			},
			"image_mapping": imageMappingSchema(),
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Description:   "A human-friendly name used as an identifier in APIs that support this option.  Only one of 'filter', 'id', 'name' or 'region_id' must be specified.",
				ConflictsWith: []string{"filter", "id", "region_id"},
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email of the user that owns the entity.",
			},
			"region_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Description:   "The id of the region for which this profile is defined as in vRealize Automation(vRA).  Only one of 'filter', 'id', 'name' or 'region_id' must be specified.",
				ConflictsWith: []string{"filter", "id", "name"},
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC.",
			},
		},
	}
}

func dataSourceImageProfileRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("Reading the vra_image_profile data source with filter %s", d.Get("filter"))

	apiClient := m.(*Client).apiClient
	var imageProfile *models.ImageProfile
	var filter string

	id := d.Get("id").(string)
	name := d.Get("name").(string)
	configFilter := d.Get("filter").(string)
	regionID := d.Get("region_id").(string)

	if id == "" && name == "" && configFilter == "" && regionID == "" {
		return fmt.Errorf("one of (id, name, region_id, filter) is required")
	}

	setFields := func(account *models.ImageProfile) error {
		d.SetId(*account.ID)
		d.Set("created_at", imageProfile.CreatedAt)
		d.Set("description", imageProfile.Description)
		d.Set("external_region_id", imageProfile.ExternalRegionID)
		d.Set("name", imageProfile.Name)
		d.Set("owner", imageProfile.Owner)
		d.Set("updated_at", imageProfile.UpdatedAt)

		if regionLink, ok := imageProfile.Links["region"]; ok {
			if regionLink.Href != "" {
				d.Set("region_id", strings.TrimPrefix(regionLink.Href, "/iaas/api/regions/"))
			}
		}

		if err := d.Set("image_mapping", flattenImageMappings(imageProfile.ImageMappings.Mapping)); err != nil {
			return fmt.Errorf("error setting image mappings - error: %#v", err)
		}
		return nil
	}

	if id != "" {
		getResp, err := apiClient.ImageProfile.GetImageProfile(image_profile.NewGetImageProfileParams().WithID(id))

		if err != nil {
			return err
		}

		imageProfile = getResp.GetPayload()
		return setFields(imageProfile)

	} else if regionID != "" {
		filter = fmt.Sprintf("regionId eq '%v'", regionID)
	} else if name != "" {
		filter = fmt.Sprintf("name eq '%v'", name)
	} else if configFilter != "" {
		filter = configFilter
	}

	getResp, err := apiClient.ImageProfile.GetImageProfiles(image_profile.NewGetImageProfilesParams().WithDollarFilter(withString(filter)))
	if err != nil {
		return err
	}

	imageProfiles := *getResp.Payload
	if len(imageProfiles.Content) > 1 {
		return fmt.Errorf("vra_image_profile must filter to an image profile")
	}
	if len(imageProfiles.Content) == 0 {
		return fmt.Errorf("vra_image_profile filter did not match any image profile")
	}

	imageProfile = imageProfiles.Content[0]

	return setFields(imageProfile)
}
