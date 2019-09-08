package vra

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceCloudAccountVMC() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudAccountVMCRead,

		Schema: map[string]*schema.Schema{
			"dc_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
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
			"nsx_hostname": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"regions": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sddc_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"vcenter_hostname": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vcenter_username": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCloudAccountVMCRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if !idOk && !nameOk {
		return fmt.Errorf("one of id or name must be assigned")
	}

	getResp, err := apiClient.CloudAccount.GetCloudAccounts(cloud_account.NewGetCloudAccountsParams())
	if err != nil {
		return err
	}

	setFields := func(account *models.CloudAccount) error {
		cloudAccountProperties := account.CloudAccountProperties

		d.SetId(*account.ID)
		d.Set("dc_id", cloudAccountProperties["dcId"])
		d.Set("description", account.Description)
		d.Set("name", account.Name)
		d.Set("nsx_hostname", cloudAccountProperties["nsxHostName"])
		d.Set("regions", account.EnabledRegionIds)
		d.Set("sddc_name", cloudAccountProperties["sddcId"])
		d.Set("vcenter_hostname", cloudAccountProperties["hostName"])
		d.Set("vcenter_username", cloudAccountProperties["privateKeyId"])

		if err := d.Set("tags", flattenTags(account.Tags)); err != nil {
			return fmt.Errorf("error setting cloud account tags - error: %#v", err)
		}
		return nil
	}
	for _, account := range getResp.Payload.Content {
		if idOk && account.ID == id && *account.CloudAccountType == "vmc" {
			return setFields(account)
		}
		if nameOk && account.Name == name && *account.CloudAccountType == "vmc" {
			return setFields(account)
		}
	}

	return fmt.Errorf("cloud account %s not found", name)
}
