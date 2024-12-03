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

func dataSourceCloudAccountAWS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudAccountAWSRead,

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
			"access_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Aws Access key ID.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was created. The date is in ISO 8601 and UTC.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly description.",
			},
			"links": linksSchema(),
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
			"tags": tagsSchema(),
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC.",
			},
		},
	}
}

func dataSourceCloudAccountAWSRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	id := d.Get("id").(string)
	name := d.Get("name").(string)

	if id == "" && name == "" {
		return fmt.Errorf("one of 'id' or 'name' must be set")
	}

	var cloudAccountAws *models.CloudAccountAws
	if id != "" {
		getResp, err := apiClient.CloudAccount.GetAwsCloudAccount(cloud_account.NewGetAwsCloudAccountParams().WithID(id))
		if err != nil {
			switch err.(type) {
			case *cloud_account.GetAwsCloudAccountNotFound:
				return fmt.Errorf("aws cloud account with id '%s' not found", id)
			default:
				// nop
			}
			return err
		}

		cloudAccountAws = getResp.GetPayload()
	} else {
		getResp, err := apiClient.CloudAccount.GetAwsCloudAccounts(cloud_account.NewGetAwsCloudAccountsParams())
		if err != nil {
			return err
		}

		for _, account := range getResp.Payload.Content {
			if account.Name == name {
				cloudAccountAws = account
			}
		}

		if cloudAccountAws == nil {
			return fmt.Errorf("aws cloud account with name '%s' not found", name)
		}
	}

	d.SetId(*cloudAccountAws.ID)
	d.Set("access_key", cloudAccountAws.AccessKeyID)
	d.Set("created_at", cloudAccountAws.CreatedAt)
	d.Set("description", cloudAccountAws.Description)
	d.Set("name", cloudAccountAws.Name)
	d.Set("org_id", cloudAccountAws.OrgID)
	d.Set("owner", cloudAccountAws.Owner)
	d.Set("updated_at", cloudAccountAws.UpdatedAt)

	if err := d.Set("links", flattenLinks(cloudAccountAws.Links)); err != nil {
		return fmt.Errorf("error setting cloud_account_aws links - error: %#v", err)
	}

	if err := d.Set("regions", extractIDsFromRegion(cloudAccountAws.EnabledRegions)); err != nil {
		return fmt.Errorf("error setting cloud_account_aws regions - error: %#v", err)
	}

	if err := d.Set("tags", flattenTags(cloudAccountAws.Tags)); err != nil {
		return fmt.Errorf("error setting cloud_account_aws tags - error: %v", err)
	}

	return nil
}
