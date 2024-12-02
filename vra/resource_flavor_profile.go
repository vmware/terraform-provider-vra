// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/flavor_profile"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func resourceFlavorProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFlavorProfileCreate,
		ReadContext:   resourceFlavorProfileRead,
		UpdateContext: resourceFlavorProfileUpdate,
		DeleteContext: resourceFlavorProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"cloud_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Id of the cloud account this flavor profile belongs to.",
			},
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
				Description: "The id of the region for which this profile is defined.",
			},
			"flavor_mapping": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A list of the flavor mappings defined for the corresponding cloud end-point region.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the flavor mapping.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The value of the instance type in the corresponding cloud. Mandatory for public clouds. Only `instance_type` or `cpu_count`/`memory` must be specified.",
						},
						"cpu_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of CPU cores. Mandatory for private clouds such as vSphere. Only `instance_type` or `cpu_count`/`memory` must be specified.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Total amount of memory (in megabytes). Mandatory for private clouds such as vSphere. Only `instance_type` or `cpu_count`/`memory` must be specified.",
						},
					},
				},
			},
			"links": linksSchema(),
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A human-friendly name used as an identifier in APIs that support this option.",
			},
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
			"region_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the region for which this profile is defined ",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC.",
			},
		},
	}
}

func resourceFlavorProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		return diag.FromErr(err)
	}

	d.SetId(*createResp.Payload.ID)

	return resourceFlavorProfileRead(ctx, d, m)
}

func resourceFlavorProfileRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	ret, err := apiClient.FlavorProfile.GetFlavorProfile(flavor_profile.NewGetFlavorProfileParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *flavor_profile.GetFlavorProfileNotFound:
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	flavor := *ret.Payload

	d.Set("cloud_account_id", flavor.CloudAccountID)
	d.Set("created_at", flavor.CreatedAt)
	d.Set("description", flavor.Description)
	d.Set("external_region_id", flavor.ExternalRegionID)

	if err := d.Set("flavor_mapping", flattenFlavors(flavor.FlavorMappings.Mapping)); err != nil {
		return diag.Errorf("error setting flavor mapping - error: %#v", err)
	}

	if err := d.Set("links", flattenLinks(flavor.Links)); err != nil {
		return diag.Errorf("error setting flavor_profile links - error: %#v", err)
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

func resourceFlavorProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	description := d.Get("description").(string)
	name := d.Get("name").(string)
	flavorMapping := expandFlavors(d.Get("flavor_mapping").(*schema.Set).List())

	if _, err := apiClient.FlavorProfile.UpdateFlavorProfile(flavor_profile.NewUpdateFlavorProfileParams().WithID(id).WithBody(&models.UpdateFlavorProfileSpecification{
		Description:   description,
		Name:          name,
		FlavorMapping: flavorMapping,
	})); err != nil {
		return diag.FromErr(err)
	}

	return resourceFlavorProfileRead(ctx, d, m)
}

func resourceFlavorProfileDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	if _, err := apiClient.FlavorProfile.DeleteFlavorProfile(flavor_profile.NewDeleteFlavorProfileParams().WithID(id)); err != nil {
		return diag.FromErr(err)
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
