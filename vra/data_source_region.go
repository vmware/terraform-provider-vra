// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"log"

	"github.com/vmware/vra-sdk-go/pkg/client/location"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRegion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRegionRead,

		Schema: map[string]*schema.Schema{
			"cloud_account_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"id", "filter"},
				Description:   "The id of the cloud account the region belongs to.",
				RequiredWith:  []string{"region"},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was created. The date is in ISO 8601 and UTC.",
			},
			"external_region_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier of region on the provider side.",
			},
			"filter": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"cloud_account_id", "id", "region"},
				Optional:      true,
				Description:   "Search criteria to narrow down Regions.",
			},
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"cloud_account_id", "filter", "region"},
				Description:   "The id of the region instance.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of region on the provider side. In vSphere, the name of the region is different from its id.",
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
			"region": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "filter"},
				Description:   "The specific region associated with the cloud account. On vSphere, this is the external ID.",
				RequiredWith:  []string{"cloud_account_id"},
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC.",
			},
		},
	}
}

func dataSourceRegionRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	cloudAccountID, cloudAccountIDOk := d.GetOk("cloud_account_id")
	region, regionOk := d.GetOk("region")
	filter, filterOk := d.GetOk("filter")
	id, idOk := d.GetOk("id")

	if !idOk && !cloudAccountIDOk && !regionOk && !filterOk {
		return fmt.Errorf("one of the following are required: (`id`, `filter`, or `region` and `cloud_account_id`)")
	}

	setFields := func(region *models.Region) {
		d.SetId(*region.ID)
		d.Set("cloud_account_id", region.CloudAccountID)
		d.Set("created_at", region.CreatedAt)
		d.Set("external_region_id", region.ExternalRegionID)
		d.Set("name", region.Name)
		d.Set("org_id", region.OrgID)
		d.Set("owner", region.Owner)
		d.Set("updated_at", region.UpdatedAt)
	}

	if idOk {
		// config includes id, using id to get region details
		getResp, err := apiClient.Location.GetRegion(location.NewGetRegionParams().WithID(id.(string)))
		if err != nil {
			switch err.(type) {
			case *location.GetRegionNotFound:
				return fmt.Errorf("region %s not found", id.(string))
			default:
				return err
			}
		}

		setFields(getResp.Payload)
		return nil
	}

	if filterOk {
		// config includes filter.
		getResp, err := apiClient.Location.GetRegions(location.NewGetRegionsParams().WithDollarFilter(withString(filter.(string))))
		if err != nil {
			return err
		}

		var region *models.Region
		regions := getResp.Payload
		if len(regions.Content) > 1 {
			log.Printf("received more than one result with the filter provided")
			name, nameOk := d.GetOk("name")

			if !nameOk {
				return fmt.Errorf("vra_region must filter to a single region. Provide the 'name' argument to filter more")
			}

			for _, reg := range regions.Content {
				if reg.Name == name.(string) {
					setFields(reg)
					return nil
				}
			}

			return fmt.Errorf("more than one region found with the filter criteria, but the name provided did not match any regions")
		}

		if len(regions.Content) == 0 {
			return fmt.Errorf("vra_region filter did not match any regions")
		}

		if len(regions.Content) == 1 {
			region = regions.Content[0]

			name, nameOk := d.GetOk("name")

			if nameOk && region.Name != name.(string) {
				return fmt.Errorf("one region found with the filter criteria, but the name provided did not match")
			}
		}

		setFields(region)
		return nil
	}

	if cloudAccountIDOk && regionOk {
		getResp, err := apiClient.CloudAccount.GetCloudAccount(cloud_account.NewGetCloudAccountParams().WithID(cloudAccountID.(string)))
		if err != nil {
			switch err.(type) {
			case *cloud_account.GetCloudAccountNotFound:
				return fmt.Errorf("cloud account %s not found", cloudAccountID.(string))
			default:
				return err
			}
		}

		var id string
		cloudAccount := getResp.Payload
		for _, enabledRegion := range cloudAccount.EnabledRegions {
			// Look for the external region ID instead of the region name to be backwards compatible.
			if *enabledRegion.ExternalRegionID == region {
				id = *enabledRegion.ID
				break
			}
		}

		if id == "" {
			return fmt.Errorf("region %s not found", region)
		}

		resp, err := apiClient.Location.GetRegion(location.NewGetRegionParams().WithID(id))
		if err != nil {
			switch err.(type) {
			case *location.GetRegionNotFound:
				return fmt.Errorf("region %s not found", id)
			default:
				return err
			}
		}

		setFields(resp.Payload)
		return nil
	}

	return fmt.Errorf("region %s not found in cloud account", region)
}
