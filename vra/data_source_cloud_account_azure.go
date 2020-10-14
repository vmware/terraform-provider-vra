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

			"application_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": linksSchema(),
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
			"subscription_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
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

	setFields := func(account *models.CloudAccountAzure) error {
		d.SetId(*account.ID)
		d.Set("application_id", account.ClientApplicationID)
		d.Set("created_at", account.CreatedAt)
		d.Set("description", account.Description)

		if err := d.Set("links", flattenLinks(account.Links)); err != nil {
			return fmt.Errorf("error setting cloud_account_azure links - error: %#v", err)
		}

		d.Set("name", account.Name)
		d.Set("org_id", account.OrgID)
		d.Set("regions", account.EnabledRegionIds)
		d.Set("subscription_id", account.SubscriptionID)

		if err := d.Set("tags", flattenTags(account.Tags)); err != nil {
			return fmt.Errorf("error setting cloud_account_azure tags - error: %v", err)
		}

		d.Set("tenant_id", account.TenantID)
		d.Set("updated_at", account.UpdatedAt)

		return nil
	}

	for _, account := range getResp.Payload.Content {
		if (idOk && account.ID == id) || (nameOk && account.Name == name) {
			return setFields(account)
		}
	}

	return fmt.Errorf("cloud account %s not found", name)
}
