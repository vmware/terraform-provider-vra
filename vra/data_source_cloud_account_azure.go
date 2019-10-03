package vra

import (
	"fmt"

	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceCloudAccountAzure() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudAccountAzureRead,

		Schema: map[string]*schema.Schema{

			"description": {
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
			"application_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subscription_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCloudAccountAzureRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if !idOk && !nameOk {
		return fmt.Errorf("One of id or name must be assigned")
	}

	getResp, err := apiClient.CloudAccount.GetAzureCloudAccounts(cloud_account.NewGetAzureCloudAccountsParams())
	if err != nil {
		return err
	}

	setFields := func(account *models.CloudAccountAzure) {
		d.SetId(*account.ID)
		d.Set("description", account.Description)
		d.Set("name", account.Name)
		d.Set("application_id", account.ClientApplicationID)
		d.Set("subscription_id", account.SubscriptionID)
		d.Set("tenant_id", account.TenantID)
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
