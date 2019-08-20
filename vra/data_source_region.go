package vra

import (
	"fmt"
	"strings"

	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceRegion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRegionRead,

		Schema: map[string]*schema.Schema{
			"cloud_account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceRegionRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	cloudAccountID := d.Get("cloud_account_id").(string)
	region := d.Get("region").(string)

	getResp, err := apiClient.CloudAccount.GetCloudAccount(cloud_account.NewGetCloudAccountParams().WithID(cloudAccountID))
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

	return fmt.Errorf("region %s not found in cloud account", region)
}
