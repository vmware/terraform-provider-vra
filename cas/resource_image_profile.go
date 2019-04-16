package cas

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/cas-sdk-go/pkg/client/image_profile"
	"github.com/vmware/cas-sdk-go/pkg/models"

	tango "github.com/vmware/terraform-provider-cas/sdk"
)

func resourceImageProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceImageProfileCreate,
		Read:   resourceImageProfileRead,
		Update: resourceImageProfileUpdate,
		Delete: resourceImageProfileDelete,

		Schema: map[string]*schema.Schema{
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"image_mapping": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"image_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cloud_config": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"external_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"external_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_family": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"region_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceImageProfileCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)
	apiClient := client.GetAPIClient()

	imageMapping := expandImageMapping(d.Get("image_mapping").(*schema.Set).List())

	createResp, err := apiClient.ImageProfile.CreateImageProfile(image_profile.NewCreateImageProfileParams().WithBody(&models.ImageProfileSpecification{
		Description:  d.Get("description").(string),
		Name:         withString(d.Get("name").(string)),
		RegionID:     withString(d.Get("region_id").(string)),
		ImageMapping: imageMapping,
	}))
	if err != nil {
		return err
	}

	d.SetId(*createResp.Payload.ID)

	return resourceImageProfileRead(d, m)
}

func resourceImageProfileRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)
	apiClient := client.GetAPIClient()

	id := d.Id()
	ret, err := apiClient.ImageProfile.GetImageProfile(image_profile.NewGetImageProfileParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *image_profile.GetImageProfileNotFound:
			d.SetId("")
			return nil
		}
		return err
	}
	image := *ret.Payload
	d.Set("description", image.Description)
	d.Set("name", image.Name)
	d.Set("image_mapping", flattenImageMapping(image.ImageMappings.Mapping))

	return nil
}

func resourceImageProfileUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)
	apiClient := client.GetAPIClient()

	id := d.Id()
	imageMapping := expandImageMapping(d.Get("image_mapping").(*schema.Set).List())

	_, err := apiClient.ImageProfile.UpdateImageProfile(image_profile.NewUpdateImageProfileParams().WithID(id).WithBody(&models.ImageProfileSpecification{
		Description:  d.Get("description").(string),
		Name:         withString(d.Get("name").(string)),
		RegionID:     withString(d.Get("region_id").(string)),
		ImageMapping: imageMapping,
	}))
	if err != nil {
		return err
	}

	return resourceImageProfileRead(d, m)
}

func resourceImageProfileDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)
	apiClient := client.GetAPIClient()

	id := d.Id()
	_, err := apiClient.ImageProfile.DeleteImageProfile(image_profile.NewDeleteImageProfileParams().WithID(id))
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func expandImageMapping(configImageMappings []interface{}) map[string]models.FabricImageDescription {
	images := make(map[string]models.FabricImageDescription)

	for _, configImageMapping := range configImageMappings {
		image := configImageMapping.(map[string]interface{})

		i := models.FabricImageDescription{
			CloudConfig: image["cloud_config"].(string),
			ID:          image["image_id"].(string),
			Name:        image["image_name"].(string),
		}
		images[image["name"].(string)] = i
	}

	return images
}

func flattenImageMapping(list map[string]models.ImageMappingDescription) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, image := range list {
		l := map[string]interface{}{
			"description":        image.Description,
			"id":                 image.ID,
			"external_id":        image.ExternalID,
			"external_region_id": image.ExternalRegionID,
			"organization":       image.OrganizationID,
			"os_family":          image.OsFamily,
			"owner":              image.Owner,
			"private":            image.IsPrivate,
		}
		result = append(result, l)
	}
	return result
}
