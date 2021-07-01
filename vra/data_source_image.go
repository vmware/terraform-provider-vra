package vra

import (
	"fmt"

	"github.com/vmware/vra-sdk-go/pkg/client/fabric_images"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceImageRead,

		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceImageRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	filter := d.Get("filter").(string)

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

	image := images.Content[0]
	d.Set("description", image.Description)
	d.Set("external_id", image.ExternalID)
	d.Set("id", image.ID)
	d.Set("name", image.Name)
	d.Set("private", image.IsPrivate)
	d.Set("region", image.ExternalRegionID)

	d.SetId(*image.ID)

	return nil
}
