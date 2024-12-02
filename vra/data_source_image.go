// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"errors"
	"fmt"

	"github.com/vmware/vra-sdk-go/pkg/client/fabric_images"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceImageRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Computed:      true,
				ConflictsWith: []string{"filter"},
				Description:   "The id of the image resource instance.",
				Optional:      true,
			},
			"filter": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"id"},
				Description:   "Search criteria to narrow down the image resource instance.",
				Optional:      true,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was created. The date is in ISO 8601 and UTC.",
			},
			"custom_properties": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "A list of key value pair of custom properties for the image resource.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly description.",
			},
			"external_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "External entity Id on the provider side.",
			},
			"links": linksSchema(),
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly name used as an identifier for the image resource instance.",
			},
			"org_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the organization this entity belongs to.",
			},
			"os_family": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Operating System family of the image.",
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email of the user that owns the entity.",
			},
			"private": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this image is private. For vSphere, private images are considered to be templates and snapshots and public are Content Library Items.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region of the image.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC.",
			},
		},
	}
}

func dataSourceImageRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	id := d.Get("id").(string)
	filter := d.Get("filter").(string)

	if id == "" && filter == "" {
		return errors.New("one of id or filter is required")
	}

	var image *models.FabricImage
	if id != "" {
		getResp, err := apiClient.FabricImages.GetFabricImage(fabric_images.NewGetFabricImageParams().WithID(id))
		if err != nil {
			switch err.(type) {
			case *fabric_images.GetFabricImageNotFound:
				return fmt.Errorf("image '%s' not found", id)
			default:
				// nop
			}
			return err
		}

		image = getResp.GetPayload()
	} else {
		getResp, err := apiClient.FabricImages.GetFabricImages(fabric_images.NewGetFabricImagesParams().WithDollarFilter(withString(filter)))
		if err != nil {
			return err
		}

		images := getResp.Payload
		if len(images.Content) > 1 {
			return fmt.Errorf("vra_image must filter to a single image")
		}
		if len(images.Content) == 0 {
			return fmt.Errorf("vra_image filter did not match any images")
		}

		image = images.Content[0]
	}

	d.SetId(*image.ID)
	d.Set("created_at", image.CreatedAt)
	d.Set("custom_properties", image.CustomProperties)
	d.Set("description", image.Description)
	d.Set("external_id", image.ExternalID)
	d.Set("id", image.ID)
	d.Set("name", image.Name)
	d.Set("org_id", image.OrgID)
	d.Set("os_family", image.OsFamily)
	d.Set("owner", image.Owner)
	d.Set("private", image.IsPrivate)
	d.Set("region", image.ExternalRegionID)
	d.Set("updated_at", image.UpdatedAt)

	if err := d.Set("links", flattenLinks(image.Links)); err != nil {
		return fmt.Errorf("error setting image links - error: %#v", err)
	}

	return nil
}
