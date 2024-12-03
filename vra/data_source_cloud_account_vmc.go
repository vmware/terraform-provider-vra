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
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was created. The date is in ISO 8601 and UTC.",
			},
			"dc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identifier of a data collector vm deployed in the on premise infrastructure.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly description.",
			},
			"links": linksSchema(),
			"nsx_hostname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address or FQDN of the NSX-T Server in the specified SDDC.",
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
			"regions": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "The set of region ids that are enabled for this cloud account.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sddc_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identifier of the on-premise SDDC.",
			},
			"tags": tagsSchema(),
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC.",
			},
			"vcenter_hostname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address or FQDN of the vCenter Server in the specified SDDC.",
			},
			"vcenter_username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Username of the vCenter Server in the specified SDDC.",
			},
		},
	}
}

func dataSourceCloudAccountVMCRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	id := d.Get("id").(string)
	name := d.Get("name").(string)

	if id == "" && name == "" {
		return fmt.Errorf("one of 'id' or 'name' must be set")
	}

	var cloudAccount *models.CloudAccount
	if id != "" {
		getResp, err := apiClient.CloudAccount.GetCloudAccount(cloud_account.NewGetCloudAccountParams().WithID(id))
		if err != nil {
			switch err.(type) {
			case *cloud_account.GetCloudAccountNotFound:
				return fmt.Errorf("vmc cloud account with id '%s' not found", id)
			default:
				// nop
			}
			return err
		}

		cloudAccount = getResp.GetPayload()
	} else {
		getResp, err := apiClient.CloudAccount.GetCloudAccounts(cloud_account.NewGetCloudAccountsParams())
		if err != nil {
			return err
		}

		for _, account := range getResp.Payload.Content {
			if account.Name == name {
				cloudAccount = account
			}
		}

		if cloudAccount == nil {
			return fmt.Errorf("vmc cloud account with name '%s' not found", name)
		}
	}

	cloudAccountProperties := cloudAccount.CloudAccountProperties

	d.SetId(*cloudAccount.ID)
	d.Set("created_at", cloudAccount.CreatedAt)
	d.Set("dc_id", cloudAccountProperties["dcId"])
	d.Set("description", cloudAccount.Description)
	d.Set("name", cloudAccount.Name)
	d.Set("nsx_hostname", cloudAccountProperties["nsxHostName"])
	d.Set("org_id", cloudAccount.OrgID)
	d.Set("owner", cloudAccount.Owner)
	d.Set("sddc_name", cloudAccountProperties["sddcId"])
	d.Set("updated_at", cloudAccount.UpdatedAt)
	d.Set("vcenter_hostname", cloudAccountProperties["hostName"])
	d.Set("vcenter_username", cloudAccountProperties["privateKeyId"])

	if err := d.Set("links", flattenLinks(cloudAccount.Links)); err != nil {
		return fmt.Errorf("error setting cloud_account_vmc links - error: %#v", err)
	}

	if err := d.Set("regions", extractIDsFromRegion(cloudAccount.EnabledRegions)); err != nil {
		return fmt.Errorf("error setting cloud_account_vmc regions - error: %#v", err)
	}

	if err := d.Set("tags", flattenTags(cloudAccount.Tags)); err != nil {
		return fmt.Errorf("error setting cloud_account_vmc tags - error: %#v", err)
	}

	return nil
}
