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
			"client_email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GCP Client email.",
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
			"private_key_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GCP Private key ID.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GCP Project ID.",
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

func dataSourceCloudAccountGCPRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	id := d.Get("id").(string)
	name := d.Get("name").(string)

	if id == "" && name == "" {
		return fmt.Errorf("one of 'id' or 'name' must be set")
	}

	var cloudAccountGcp *models.CloudAccountGcp
	if id != "" {
		getResp, err := apiClient.CloudAccount.GetGcpCloudAccount(cloud_account.NewGetGcpCloudAccountParams().WithID(id))
		if err != nil {
			switch err.(type) {
			case *cloud_account.GetGcpCloudAccountNotFound:
				return fmt.Errorf("gcp cloud account with id '%s' not found", id)
			default:
				// nop
			}
			return err
		}

		cloudAccountGcp = getResp.GetPayload()
	} else {
		getResp, err := apiClient.CloudAccount.GetGcpCloudAccounts(cloud_account.NewGetGcpCloudAccountsParams())
		if err != nil {
			return err
		}

		for _, account := range getResp.Payload.Content {
			if account.Name == name {
				cloudAccountGcp = account
			}
		}

		if cloudAccountGcp == nil {
			return fmt.Errorf("gcp cloud account with name '%s' not found", name)
		}
	}

	d.SetId(*cloudAccountGcp.ID)
	d.Set("client_email", cloudAccountGcp.ClientEmail)
	d.Set("created_at", cloudAccountGcp.CreatedAt)
	d.Set("description", cloudAccountGcp.Description)
	d.Set("name", cloudAccountGcp.Name)
	d.Set("org_id", cloudAccountGcp.OrgID)
	d.Set("owner", cloudAccountGcp.Owner)
	d.Set("private_key_id", cloudAccountGcp.PrivateKeyID)
	d.Set("project_id", cloudAccountGcp.ProjectID)
	d.Set("updated_at", cloudAccountGcp.UpdatedAt)

	if err := d.Set("links", flattenLinks(cloudAccountGcp.Links)); err != nil {
		return fmt.Errorf("error setting cloud_account_gcp links - error: %#v", err)
	}

	if err := d.Set("regions", extractIDsFromRegion(cloudAccountGcp.EnabledRegions)); err != nil {
		return fmt.Errorf("error setting cloud_account_gcp regions - error: %#v", err)
	}

	if err := d.Set("tags", flattenTags(cloudAccountGcp.Tags)); err != nil {
		return fmt.Errorf("error setting cloud_account_gcp tags - error: %#v", err)
	}

	return nil
}
