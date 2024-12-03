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

func dataSourceCloudAccountVsphere() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudAccountVsphereRead,

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
				Description: "Identifier of a data collector vm deployed in the on premise infrastructure.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly description.",
			},
			"hostname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address or FQDN of the vCenter Server.",
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
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Username of the vCenter Server.",
			},
		},
	}
}

func dataSourceCloudAccountVsphereRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	id := d.Get("id").(string)
	name := d.Get("name").(string)

	if id == "" && name == "" {
		return fmt.Errorf("one of 'id' or 'name' must be set")
	}

	var cloudAccountVsphere *models.CloudAccountVsphere
	if id != "" {
		getResp, err := apiClient.CloudAccount.GetVSphereCloudAccount(cloud_account.NewGetVSphereCloudAccountParams().WithID(id))
		if err != nil {
			switch err.(type) {
			case *cloud_account.GetVSphereCloudAccountNotFound:
				return fmt.Errorf("vsphere cloud account with id '%s' not found", id)
			default:
				// nop
			}
			return err
		}

		cloudAccountVsphere = getResp.GetPayload()
	} else {
		getResp, err := apiClient.CloudAccount.GetVSphereCloudAccounts(cloud_account.NewGetVSphereCloudAccountsParams())
		if err != nil {
			return err
		}

		for _, account := range getResp.Payload.Content {
			if account.Name == name {
				cloudAccountVsphere = account
			}
		}

		if cloudAccountVsphere == nil {
			return fmt.Errorf("vsphere cloud account with name '%s' not found", name)
		}
	}

	d.SetId(*cloudAccountVsphere.ID)
	d.Set("associated_cloud_account_ids", flattenAssociatedCloudAccountIDs(cloudAccountVsphere.Links))
	d.Set("created_at", cloudAccountVsphere.CreatedAt)
	d.Set("dc_id", cloudAccountVsphere.Dcid)
	d.Set("description", cloudAccountVsphere.Description)
	d.Set("hostname", cloudAccountVsphere.HostName)
	d.Set("name", cloudAccountVsphere.Name)
	d.Set("org_id", cloudAccountVsphere.OrgID)
	d.Set("owner", cloudAccountVsphere.Owner)
	d.Set("updated_at", cloudAccountVsphere.UpdatedAt)
	d.Set("username", cloudAccountVsphere.Username)

	if err := d.Set("links", flattenLinks(cloudAccountVsphere.Links)); err != nil {
		return fmt.Errorf("error setting cloud_account_vsphere links - error: %#v", err)
	}

	if err := d.Set("regions", extractIDsFromRegion(cloudAccountVsphere.EnabledRegions)); err != nil {
		return fmt.Errorf("error setting cloud_account_vsphere regions - error: %#v", err)
	}

	if err := d.Set("tags", flattenTags(cloudAccountVsphere.Tags)); err != nil {
		return fmt.Errorf("error setting cloud_account_vsphere tags - error: %#v", err)
	}

	return nil
}
