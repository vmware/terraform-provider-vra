// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/image_profile"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"strings"
)

func resourceImageProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImageProfileCreate,
		ReadContext:   resourceImageProfileRead,
		UpdateContext: resourceImageProfileUpdate,
		DeleteContext: resourceImageProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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
			"image_mapping": imageMappingSchema(),
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email of the user that owns the entity.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A human-friendly name used as an identifier in APIs that support this option.",
			},
			"region_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The if of the region for which this profile is defined as in vRealize Automation(vRA).",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC.",
			},
		},
	}
}

func resourceImageProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	imageMapping := expandImageMapping(d.Get("image_mapping").(*schema.Set).List())

	createResp, err := apiClient.ImageProfile.CreateImageProfile(image_profile.NewCreateImageProfileParams().WithBody(&models.ImageProfileSpecification{
		Description:  d.Get("description").(string),
		Name:         withString(d.Get("name").(string)),
		RegionID:     withString(d.Get("region_id").(string)),
		ImageMapping: imageMapping,
	}))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*createResp.Payload.ID)

	return resourceImageProfileRead(ctx, d, m)
}

func resourceImageProfileRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	ret, err := apiClient.ImageProfile.GetImageProfile(image_profile.NewGetImageProfileParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *image_profile.GetImageProfileNotFound:
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	imageProfile := *ret.Payload
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
		return diag.Errorf("error setting image mappings - error: %#v", err)
	}

	return nil
}

func resourceImageProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	if d.HasChange("region_id") {
		return diag.FromErr(errors.New("change detected in region_id, but not supported on an image profile"))
	}

	id := d.Id()
	imageMapping := expandImageMapping(d.Get("image_mapping").(*schema.Set).List())

	_, err := apiClient.ImageProfile.UpdateImageProfile(image_profile.NewUpdateImageProfileParams().WithID(id).WithBody(&models.UpdateImageProfileSpecification{
		Description:  d.Get("description").(string),
		Name:         d.Get("name").(string),
		ImageMapping: imageMapping,
	}))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceImageProfileRead(ctx, d, m)
}

func resourceImageProfileDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	_, err := apiClient.ImageProfile.DeleteImageProfile(image_profile.NewDeleteImageProfileParams().WithID(id))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
