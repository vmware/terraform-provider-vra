package vra

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceCloudAccountGCP() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudAccountGCPRead,

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
			"client_email": {
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
			"private_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
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
			"tags": tagsSchema(),
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCloudAccountGCPRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if !idOk && !nameOk {
		return fmt.Errorf("one of id or name must be assigned")
	}

	getResp, err := apiClient.CloudAccount.GetGcpCloudAccounts(cloud_account.NewGetGcpCloudAccountsParams())
	if err != nil {
		return err
	}

	setFields := func(account *models.CloudAccountGcp) error {
		d.SetId(*account.ID)
		d.Set("client_email", account.ClientEmail)
		d.Set("created_at", account.CreatedAt)
		d.Set("description", account.Description)
		d.Set("name", account.Name)
		d.Set("org_id", account.OrgID)
		d.Set("owner", account.Owner)
		d.Set("private_key_id", account.PrivateKeyID)
		d.Set("project_id", account.ProjectID)
		d.Set("regions", account.EnabledRegionIds)
		d.Set("updated_at", account.UpdatedAt)

		if err := d.Set("links", flattenLinks(account.Links)); err != nil {
			return fmt.Errorf("error setting cloud_account_gcp links - error: %#v", err)
		}

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
