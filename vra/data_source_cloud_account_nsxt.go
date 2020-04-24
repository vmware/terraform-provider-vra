package vra

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceCloudAccountNSXT() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudAccountNSXTRead,

		Schema: map[string]*schema.Schema{
			"associated_cloud_account_ids": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"dc_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"hostname": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tagsSchema(),
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCloudAccountNSXTRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if !idOk && !nameOk {
		return fmt.Errorf("one of id or name must be assigned")
	}

	getResp, err := apiClient.CloudAccount.GetNsxTCloudAccounts(cloud_account.NewGetNsxTCloudAccountsParams())
	if err != nil {
		return err
	}

	setFields := func(account *models.CloudAccountNsxT) error {
		d.SetId(*account.ID)
		d.Set("associated_cloud_account_ids", flattenAssociatedCloudAccountIds(account.Links))
		d.Set("dc_id", account.Dcid)
		d.Set("description", account.Description)
		d.Set("hostname", account.HostName)
		d.Set("name", account.Name)
		d.Set("username", account.Username)

		if err := d.Set("tags", flattenTags(account.Tags)); err != nil {
			return fmt.Errorf("error setting cloud account tags - error: %#v", err)
		}
		return nil
	}
	for _, account := range getResp.Payload.Content {
		if idOk && account.ID == id {
			return setFields(account)
		}
		if nameOk && account.Name == name {
			return setFields(account)
		}
	}

	return fmt.Errorf("cloud account %s not found", name)
}
