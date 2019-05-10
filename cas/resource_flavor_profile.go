package cas

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/cas-sdk-go/pkg/client/flavor_profile"
	"github.com/vmware/cas-sdk-go/pkg/models"

	tango "github.com/vmware/terraform-provider-cas/sdk"
)

func resourceFlavorProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceFlavorProfileCreate,
		Read:   resourceFlavorProfileRead,
		Update: resourceFlavorProfileUpdate,
		Delete: resourceFlavorProfileDelete,

		Schema: map[string]*schema.Schema{
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor_mapping": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cpu_count": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"memory": {
							Type:     schema.TypeInt,
							Optional: true,
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

func resourceFlavorProfileCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)
	apiClient := client.GetAPIClient()

	description := d.Get("description").(string)
	name := d.Get("name").(string)
	regionID := d.Get("region_id").(string)
	flavorMapping := expandFlavors(d.Get("flavor_mapping").(*schema.Set).List())

	createResp, err := apiClient.FlavorProfile.CreateFlavorProfile(flavor_profile.NewCreateFlavorProfileParams().WithBody(&models.FlavorProfileSpecification{
		Description:   description,
		Name:          &name,
		RegionID:      &regionID,
		FlavorMapping: flavorMapping,
	}))
	if err != nil {
		return err
	}

	d.SetId(*createResp.Payload.ID)

	return resourceFlavorProfileRead(d, m)
}

func resourceFlavorProfileRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)
	apiClient := client.GetAPIClient()

	id := d.Id()
	ret, err := apiClient.FlavorProfile.GetFlavorProfile(flavor_profile.NewGetFlavorProfileParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *flavor_profile.GetFlavorProfileNotFound:
			d.SetId("")
			return nil
		}
		return err
	}
	flavor := *ret.Payload
	d.Set("description", flavor.Description)
	d.Set("name", flavor.Name)
	d.Set("flavor_mappings", flattenFlavors(flavor.FlavorMappings.Mapping))

	return nil
}

func resourceFlavorProfileUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)
	apiClient := client.GetAPIClient()

	id := d.Id()
	description := d.Get("description").(string)
	name := d.Get("name").(string)
	regionID := d.Get("region_id").(string)
	flavorMapping := expandFlavors(d.Get("flavor_mapping").(*schema.Set).List())

	_, err := apiClient.FlavorProfile.UpdateFlavorProfile(flavor_profile.NewUpdateFlavorProfileParams().WithID(id).WithBody(&models.FlavorProfileSpecification{
		Description:   description,
		Name:          &name,
		RegionID:      &regionID,
		FlavorMapping: flavorMapping,
	}))
	if err != nil {
		return err
	}

	return resourceFlavorProfileRead(d, m)
}

func resourceFlavorProfileDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)
	apiClient := client.GetAPIClient()

	id := d.Id()
	_, err := apiClient.FlavorProfile.DeleteFlavorProfile(flavor_profile.NewDeleteFlavorProfileParams().WithID(id))
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func expandFlavors(configFlavors []interface{}) map[string]models.FabricFlavorDescription {
	flavors := make(map[string]models.FabricFlavorDescription)

	for _, configFlavor := range configFlavors {
		flavor := configFlavor.(map[string]interface{})

		f := models.FabricFlavorDescription{
			CPUCount:   int32(flavor["cpu_count"].(int)),
			MemoryInMB: int64(flavor["memory"].(int)),
			Name:       flavor["instance_type"].(string),
		}
		flavors[flavor["name"].(string)] = f
	}

	return flavors
}

func flattenFlavors(list map[string]models.FabricFlavor) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, flavor := range list {
		l := map[string]interface{}{
			"cpu_count":     flavor.CPUCount,
			"instance_type": flavor.ID,
			"memory":        flavor.MemoryInMB,
			"name":          flavor.Name,
		}

		result = append(result, l)
	}
	return result
}
