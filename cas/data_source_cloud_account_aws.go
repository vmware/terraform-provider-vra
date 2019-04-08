package cas

import (
	"fmt"

	"github.com/vmware/cas-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/cas-sdk-go/pkg/models"
	tango "github.com/vmware/terraform-provider-cas/sdk"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceCloudAccountAWS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudAccountAWSRead,

		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
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
		},
	}
}

func dataSourceCloudAccountAWSRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*tango.Client)
	apiClient := client.GetAPIClient()

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk == false && nameOk == false {
		return fmt.Errorf("One of id or name must be assigned")
	}

	getResp, err := apiClient.CloudAccount.GetAwsCloudAccounts(cloud_account.NewGetAwsCloudAccountsParams())
	if err != nil {
		return err
	}

	setFields := func(account *models.CloudAccountAws) {
		d.SetId(*account.ID)
		d.Set("access_key", account.AccessKeyID)
		d.Set("description", account.Description)
		d.Set("name", account.Name)
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
