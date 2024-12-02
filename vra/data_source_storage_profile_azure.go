// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"log"

	"github.com/vmware/vra-sdk-go/pkg/client/storage_profile"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceStorageProfileAzure() *schema.Resource {
	return &schema.Resource{
		Read: datasourceStorageProfileAzureRead,

		Schema: map[string]*schema.Schema{
			"filter": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"filter"},
			},
			"cloud_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Id of the cloud account this storage profile belongs to.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was created. The date is in ISO 8601 and UTC.",
			},
			"default_item": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this storage profile is default or not…",
			},
			"data_disk_caching": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the caching mechanism for additional disk.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly description.",
			},
			"disk_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the performance tier for the storage type. Premium disks are SSD backed and Standard disks are HDD backed.",
			},
			"external_region_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the region for which this profile is defined",
			},
			"links": linksSchema(),
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly name used as an identifier in APIs that support this option.",
			},
			"org_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the organization this entity belongs to.",
			},
			"os_disk_caching": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the caching mechanism for OS disk. Default policy for OS disks is Read/Write.",
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email of the user that owns the entity.",
			},
			"supports_encryption": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this storage profile should support encryption or not.",
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

func datasourceStorageProfileAzureRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("Reading the vra_storage_profile_azure data source with name %s", d.Get("filter"))
	apiClient := m.(*Client).apiClient

	var azureStorageProfile *models.AzureStorageProfile

	id := d.Get("id").(string)
	filter := d.Get("filter").(string)

	if id == "" && filter == "" {
		return fmt.Errorf("one of id or filter is required")
	}

	if id != "" {
		getResp, err := apiClient.StorageProfile.GetAzureStorageProfile(storage_profile.NewGetAzureStorageProfileParams().WithID(id))
		if err != nil {
			return err
		}
		azureStorageProfile = getResp.GetPayload()
	} else {
		getResp, err := apiClient.StorageProfile.GetAzureStorageProfiles(storage_profile.NewGetAzureStorageProfilesParams().WithDollarFilter(withString(filter)))
		if err != nil {
			return err
		}
		storageProfiles := *getResp.Payload
		if len(storageProfiles.Content) > 1 {
			return fmt.Errorf("vra_storage_profile_azure must filter to a storage profile")
		}
		if len(storageProfiles.Content) == 0 {
			return fmt.Errorf("vra_storage_profile_azure filter did not match any storage profile")
		}
		azureStorageProfile = storageProfiles.Content[0]
	}

	d.SetId(*azureStorageProfile.ID)
	d.Set("cloud_account_id", azureStorageProfile.CloudAccountID)
	d.Set("created_at", azureStorageProfile.CreatedAt)
	d.Set("data_disk_caching", azureStorageProfile.DataDiskCaching)
	d.Set("default_item", azureStorageProfile.DefaultItem)
	d.Set("description", azureStorageProfile.Description)
	d.Set("disk_type", azureStorageProfile.DiskType)
	d.Set("external_region_id", azureStorageProfile.ExternalRegionID)
	d.Set("name", azureStorageProfile.Name)
	d.Set("org_id", azureStorageProfile.OrgID)
	d.Set("os_disk_caching", azureStorageProfile.OsDiskCaching)
	d.Set("owner", azureStorageProfile.Owner)
	d.Set("supports_encryption", azureStorageProfile.SupportsEncryption)
	d.Set("updated_at", azureStorageProfile.UpdatedAt)

	if err := d.Set("tags", flattenTags(azureStorageProfile.Tags)); err != nil {
		return fmt.Errorf("error setting azure storage profile tags - error: %v", err)
	}

	if err := d.Set("links", flattenLinks(azureStorageProfile.Links)); err != nil {
		return fmt.Errorf("error setting azure storage profile links - error: %#v", err)
	}

	log.Printf("Finished reading the vra_storage_profile_azure data source with name %s", d.Get("name"))
	return nil
}
