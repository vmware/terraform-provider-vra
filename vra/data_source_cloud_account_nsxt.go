// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceCloudAccountNSXT() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudAccountNSXTRead,

		Schema: map[string]*schema.Schema{
			// Optional arguments
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
				Description:   "The id of this resource instance.",
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"id"},
				Description:   "The name of this resource instance.",
			},

			// Computed attributes
			"associated_cloud_account_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was created. The date is in ISO 8601 and UTC.",
			},
			"dc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identifier of a data collector vm deployed in the on premise infrastructure",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly description.",
			},
			"hostname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Host name for the NSX-T endpoint.",
			},
			"links": linksSchema(),
			"manager_mode": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if NSX-T cloud account is created in Manager (legacy) mode. When set to true, NSX-T cloud account is created in Manager mode. Mode cannot be changed after cloud account is created. Default value is false.",
			},
			"org_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the organization this entity belongs to.",
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email of the user that owns the entity.",
			},
			"tags": tagsSchema(),
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC.",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Username to authenticate with the cloud account.",
			},
		},
	}
}

func dataSourceCloudAccountNSXTRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if !idOk && !nameOk {
		return fmt.Errorf("one of 'id' or 'name' must be assigned")
	}

	getResp, err := apiClient.CloudAccount.GetNsxTCloudAccounts(cloud_account.NewGetNsxTCloudAccountsParams())
	if err != nil {
		return err
	}

	setFields := func(account *models.CloudAccountNsxT) error {
		d.SetId(*account.ID)
		d.Set("associated_cloud_account_ids", flattenAssociatedCloudAccountIDs(account.Links))
		d.Set("created_at", account.CreatedAt)
		d.Set("dc_id", account.Dcid)
		d.Set("description", account.Description)
		d.Set("hostname", account.HostName)
		d.Set("manager_mode", account.ManagerMode)
		d.Set("name", account.Name)
		d.Set("org_id", account.OrgID)
		d.Set("owner", account.Owner)
		d.Set("updated_at", account.UpdatedAt)
		d.Set("username", account.Username)

		if err := d.Set("links", flattenLinks(account.Links)); err != nil {
			return fmt.Errorf("error setting cloud_account_nsxt links - error: %#v", err)
		}

		if err := d.Set("tags", flattenTags(account.Tags)); err != nil {
			return fmt.Errorf("error setting cloud_account_nsxt tags - error: %#v", err)
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
