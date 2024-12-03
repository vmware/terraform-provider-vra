// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"

	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"application_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Azure Client Application ID.",
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
			"subscription_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Azure Subscription ID.",
			},
			"tags": tagsSchema(),
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Azure Tenant ID.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC.",
			},
		},
	}
}

func dataSourceCloudAccountAzureRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	id := d.Get("id").(string)
	name := d.Get("name").(string)

	if id == "" && name == "" {
		return fmt.Errorf("one of 'id' or 'name' must be set")
	}

	var cloudAccountAzure *models.CloudAccountAzure
	if id != "" {
		getResp, err := apiClient.CloudAccount.GetAzureCloudAccount(cloud_account.NewGetAzureCloudAccountParams().WithID(id))
		if err != nil {
			switch err.(type) {
			case *cloud_account.GetAzureCloudAccountNotFound:
				return fmt.Errorf("azure cloud account with id '%s' not found", id)
			default:
				// nop
			}
			return err
		}

		cloudAccountAzure = getResp.GetPayload()
	} else {
		getResp, err := apiClient.CloudAccount.GetAzureCloudAccounts(cloud_account.NewGetAzureCloudAccountsParams())
		if err != nil {
			return err
		}

		for _, account := range getResp.Payload.Content {
			if account.Name == name {
				cloudAccountAzure = account
			}
		}

		if cloudAccountAzure == nil {
			return fmt.Errorf("azure cloud account with name '%s' not found", name)
		}
	}

	d.SetId(*cloudAccountAzure.ID)
	d.Set("application_id", cloudAccountAzure.ClientApplicationID)
	d.Set("created_at", cloudAccountAzure.CreatedAt)
	d.Set("description", cloudAccountAzure.Description)
	d.Set("name", cloudAccountAzure.Name)
	d.Set("org_id", cloudAccountAzure.OrgID)
	d.Set("owner", cloudAccountAzure.Owner)
	d.Set("subscription_id", cloudAccountAzure.SubscriptionID)
	d.Set("tenant_id", cloudAccountAzure.TenantID)
	d.Set("updated_at", cloudAccountAzure.UpdatedAt)

	if err := d.Set("links", flattenLinks(cloudAccountAzure.Links)); err != nil {
		return fmt.Errorf("error setting cloud_account_azure links - error: %#v", err)
	}

	if err := d.Set("regions", extractIDsFromRegion(cloudAccountAzure.EnabledRegions)); err != nil {
		return fmt.Errorf("error setting cloud_account_azure regions - error: %#v", err)
	}

	if err := d.Set("tags", flattenTags(cloudAccountAzure.Tags)); err != nil {
		return fmt.Errorf("error setting cloud_account_azure tags - error: %v", err)
	}

	return nil
}
