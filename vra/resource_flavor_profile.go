package vra

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/flavor_profile"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func resourceFlavorProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceFlavorProfileCreate,
		Read:   resourceFlavorProfileRead,
		Update: resourceFlavorProfileUpdate,
		Delete: resourceFlavorProfileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cloud_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the cloud account this entity belongs to.",
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"external_region_id": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"flavor_mapping": {
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
			"links": linksSchema(),
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceFlavorProfileCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

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
	apiClient := m.(*Client).apiClient

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

	d.Set("cloud_account_id", flavor.CloudAccountID)
	d.Set("created_at", flavor.CreatedAt)
	d.Set("description", flavor.Description)
	d.Set("external_region_id", flavor.ExternalRegionID)

	if err := d.Set("flavor_mapping", flattenFlavors(flavor.FlavorMappings.Mapping)); err != nil {
		return fmt.Errorf("error setting flavor mapping - error: %#v", err)
	}

	if err := d.Set("links", flattenLinks(flavor.Links)); err != nil {
		return fmt.Errorf("error setting flavor_profile links - error: %#v", err)
	}

	d.Set("name", flavor.Name)
	d.Set("org_id", flavor.OrgID)
	d.Set("owner", flavor.Owner)

	if regionLink, ok := flavor.Links["region"]; ok {
		if regionLink.Href != "" {
			d.Set("region_id", strings.TrimPrefix(regionLink.Href, "/iaas/api/regions/"))
		}
	}

	d.Set("updated_at", flavor.UpdatedAt)

	return nil
}

func resourceFlavorProfileUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	description := d.Get("description").(string)
	name := d.Get("name").(string)
	flavorMapping := expandFlavors(d.Get("flavor_mapping").(*schema.Set).List())

	_, err := apiClient.FlavorProfile.UpdateFlavorProfile(flavor_profile.NewUpdateFlavorProfileParams().WithID(id).WithBody(&models.UpdateFlavorProfileSpecification{
		Description:   description,
		Name:          &name,
		FlavorMapping: flavorMapping,
	}))
	if err != nil {
		return err
	}

	return resourceFlavorProfileRead(d, m)
}

func resourceFlavorProfileDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

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
	for flavorName, flavor := range list {
		l := map[string]interface{}{
			"cpu_count":     flavor.CPUCount,
			"instance_type": flavor.ID,
			"memory":        flavor.MemoryInMB,
			"name":          flavorName,
		}

		result = append(result, l)
	}
	return result
}
