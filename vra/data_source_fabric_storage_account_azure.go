// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"github.com/vmware/vra-sdk-go/pkg/client/fabric_azure_storage_account"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFabricStorageAccountAzure() *schema.Resource {
	return &schema.Resource{
		Read: resourceFabricStorageAccountAzureRead,

		Schema: map[string]*schema.Schema{
			"cloud_account_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Set of ids of the cloud accounts this entity belongs to.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was created. The date is in ISO 8601 and UTC.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly description of the fabric Azure storage account.",
			},
			"external_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "External entity Id on the provider side.",
			},
			"external_region_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the ID of region.",
			},
			"filter": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id"},
				Description:   "Search criteria to narrow down the fabric Azure storage accounts.",
			},
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"filter"},
				Description:   "The id of this fabric Azure storage account.",
			},
			"links": linksSchema(),
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly name of the fabric Azure storage account.",
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
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the performance tier for the storage type. Premium disks are SSD backed and Standard disks are HDD backed. example: Standard_LRS / Premium_LRS",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC.",
			},
		},
	}
}

func resourceFabricStorageAccountAzureRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading the vra_fabric_storage_account_azure data source with filter [%s]", d.Get("filter"))
	apiClient := meta.(*Client).apiClient

	var fabricAzureStorageAccount *models.FabricAzureStorageAccount

	id := d.Get("id").(string)
	filter := d.Get("filter").(string)

	if id == "" && filter == "" {
		return fmt.Errorf("one of id or filter is required")
	}

	if id != "" {
		getResp, err := apiClient.FabricAzureStorageAccount.GetFabricAzureStorageAccount(fabric_azure_storage_account.NewGetFabricAzureStorageAccountParams().WithID(id))

		if err != nil {
			return err
		}
		fabricAzureStorageAccount = getResp.GetPayload()
	} else {
		getResp, err := apiClient.FabricAzureStorageAccount.GetFabricAzureStorageAccounts(fabric_azure_storage_account.NewGetFabricAzureStorageAccountsParams().WithDollarFilter(withString(filter)))
		if err != nil {
			return err
		}

		fabricAzureStorageAccounts := *getResp.Payload
		if len(fabricAzureStorageAccounts.Content) > 1 {
			return fmt.Errorf("vra_fabric_storage_account_azure must filter to a fabric Azure storage account")
		}
		if len(fabricAzureStorageAccounts.Content) == 0 {
			return fmt.Errorf("vra_fabric_storage_account_azure filter did not match any fabric Azure storage accounts")
		}

		fabricAzureStorageAccount = fabricAzureStorageAccounts.Content[0]
	}

	d.SetId(*fabricAzureStorageAccount.ID)
	d.Set("cloud_account_ids", fabricAzureStorageAccount.CloudAccountIds)
	d.Set("created_at", fabricAzureStorageAccount.CreatedAt)
	d.Set("description", fabricAzureStorageAccount.Description)
	d.Set("external_id", fabricAzureStorageAccount.ExternalID)
	d.Set("external_region_id", fabricAzureStorageAccount.ExternalRegionID)
	d.Set("name", fabricAzureStorageAccount.Name)
	d.Set("org_id", fabricAzureStorageAccount.OrgID)
	d.Set("owner", fabricAzureStorageAccount.Owner)
	d.Set("type", fabricAzureStorageAccount.Type)
	d.Set("updated_at", fabricAzureStorageAccount.UpdatedAt)

	if err := d.Set("links", flattenLinks(fabricAzureStorageAccount.Links)); err != nil {
		return fmt.Errorf("error getting fabric Azure storage account links - error: %#v", err)
	}

	log.Printf("Finished reading the vra_fabric_storage_account_azure data source with filter [%s]", d.Get("filter"))
	return nil
}
