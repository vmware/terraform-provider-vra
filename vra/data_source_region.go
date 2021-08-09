package vra

import (
	"fmt"
	"log"
	"strings"

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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
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
		return fmt.Errorf("one of the following are required: (id, filter, region and cloudAccountId)")
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
			return err
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
			return err
		}

		cloudAccount := getResp.Payload
		for i, enabledRegion := range cloudAccount.EnabledRegionIds {
			if enabledRegion == region {
				d.SetId(strings.TrimPrefix(cloudAccount.Links["regions"].Hrefs[i], "/iaas/api/regions/"))
				return nil
			}
		}

		resp, err := apiClient.Location.GetRegion(location.NewGetRegionParams().WithID(id.(string)))
		if err != nil {
			return err
		}

		setFields(resp.Payload)
		return nil
	}

	return fmt.Errorf("region %s not found in cloud account", region)
}
