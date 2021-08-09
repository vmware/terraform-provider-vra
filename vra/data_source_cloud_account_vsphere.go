package vra

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceCloudAccountVsphere() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudAccountVsphereRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"associated_cloud_account_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_properties": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"dcid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled_region_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
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
			"tags": tagsSchema(),
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceCloudAccountVsphereRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if !idOk && !nameOk {
		return fmt.Errorf("One of id or name must be assigned")
	}
	getResp, err := apiClient.CloudAccount.GetVSphereCloudAccounts(cloud_account.NewGetVSphereCloudAccountsParams())
	if err != nil {
		return err
	}

	setFields := func(account *models.CloudAccountVsphere) {
		d.SetId(*account.ID)
		d.Set("associated_cloud_account_ids", flattenAssociatedCloudAccountIds(account.Links))
		d.Set("created_at", account.CreatedAt)
		d.Set("custom_properties", account.CustomProperties)
		d.Set("dcid", account.Dcid)
		d.Set("description", account.Description)
		d.Set("enabled_region_ids", account.EnabledRegionIds)
		d.Set("hostname", account.HostName)
		d.Set("name", account.Name)
		d.Set("org_id", account.OrgID)
		d.Set("owner", account.Owner)
		d.Set("tags", account.Tags)
		d.Set("updated_at", account.UpdatedAt)
	}
	for _, account := range getResp.Payload.Content {
		if idOk && account.ID == id {
			setFields(account)
			return nil
		}
		if nameOk && account.Name == name {
			setFields(account)
			return nil
		}
	}

	return fmt.Errorf("cloud account %s not found", name)
}
