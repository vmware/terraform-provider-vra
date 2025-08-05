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

	id := d.Get("id").(string)
	name := d.Get("name").(string)

	if id == "" && name == "" {
		return fmt.Errorf("one of 'id' or 'name' must be set")
	}

	var cloudAccountNsxT *models.CloudAccountNsxT
	if id != "" {
		getResp, err := apiClient.CloudAccount.GetNsxTCloudAccount(cloud_account.NewGetNsxTCloudAccountParams().WithID(id))
		if err != nil {
			switch err.(type) {
			case *cloud_account.GetNsxTCloudAccountNotFound:
				return fmt.Errorf("nsxt cloud account with id '%s' not found", id)
			default:
				// nop
			}
			return err
		}

		cloudAccountNsxT = getResp.GetPayload()
	} else {
		getResp, err := apiClient.CloudAccount.GetNsxTCloudAccounts(cloud_account.NewGetNsxTCloudAccountsParams())
		if err != nil {
			return err
		}

		for _, account := range getResp.Payload.Content {
			if account.Name == name {
				cloudAccountNsxT = account
			}
		}

		if cloudAccountNsxT == nil {
			return fmt.Errorf("nsxt cloud account with name '%s' not found", name)
		}
	}

	d.SetId(*cloudAccountNsxT.ID)
	_ = d.Set("associated_cloud_account_ids", flattenAssociatedCloudAccountIDs(cloudAccountNsxT.Links))
	_ = d.Set("created_at", cloudAccountNsxT.CreatedAt)
	_ = d.Set("dc_id", cloudAccountNsxT.Dcid)
	_ = d.Set("description", cloudAccountNsxT.Description)
	_ = d.Set("hostname", cloudAccountNsxT.HostName)
	_ = d.Set("manager_mode", cloudAccountNsxT.ManagerMode)
	_ = d.Set("name", cloudAccountNsxT.Name)
	_ = d.Set("org_id", cloudAccountNsxT.OrgID)
	_ = d.Set("owner", cloudAccountNsxT.Owner)
	_ = d.Set("updated_at", cloudAccountNsxT.UpdatedAt)
	_ = d.Set("username", cloudAccountNsxT.Username)

	if err := d.Set("links", flattenLinks(cloudAccountNsxT.Links)); err != nil {
		return fmt.Errorf("error setting cloud_account_nsxt links - error: %#v", err)
	}

	if err := d.Set("tags", flattenTags(cloudAccountNsxT.Tags)); err != nil {
		return fmt.Errorf("error setting cloud_account_nsxt tags - error: %#v", err)
	}

	return nil
}
