package vra

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceCloudAccountVMC() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudAccountVMCRead,

		Schema: map[string]*schema.Schema{
			// Optional arguments
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"id"},
			},
			// Computed attributes
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": linksSchema(),
			"nsx_hostname": {
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
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sddc_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vcenter_hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vcenter_username": {
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
		d.Set("created_at", account.CreatedAt)
		d.Set("dc_id", cloudAccountProperties["dcId"])
		d.Set("description", account.Description)
		d.Set("name", account.Name)
		d.Set("nsx_hostname", cloudAccountProperties["nsxHostName"])
		d.Set("owner", account.Owner)
		d.Set("org_id", account.OrgID)
		d.Set("regions", account.EnabledRegionIds)
		d.Set("sddc_name", cloudAccountProperties["sddcId"])
		d.Set("updated_at", account.UpdatedAt)
		d.Set("vcenter_hostname", cloudAccountProperties["hostName"])
		d.Set("vcenter_username", cloudAccountProperties["privateKeyId"])

		if err := d.Set("links", flattenLinks(account.Links)); err != nil {
			return fmt.Errorf("error setting cloud_account_vmc links - error: %#v", err)
		}

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
