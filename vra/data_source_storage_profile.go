// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/storage_profile"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"log"
)

func dataSourceStorageProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceStorageProfileRead,

		Schema: map[string]*schema.Schema{
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
				Description: "Indicates if this storage profile is a default profile.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A human-friendly description.",
			},
			"disk_properties": {
				Type:        schema.TypeMap,
				Computed:    true,
				Optional:    true,
				Description: "Map of storage properties that are to be applied on disk while provisioning.",
			},
			"external_region_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the region as seen in the cloud provider for which this profile is defined.",
			},
			"filter": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id"},
				Description:   "Search criteria to filter the list of storage profiles.",
			},
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"filter"},
				Description:   "The id of the storage profile.",
			},
			"links": linksSchema(),
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly name for storage profile.",
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
			"supports_encryption": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this storage profile supports encryption or not.",
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

func dataSourceStorageProfileRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading the vra_storage_profile data source with filter %s", d.Get("filter"))
	apiClient := meta.(*Client).apiClient

	var storageProfile *models.StorageProfile

	id := d.Get("id").(string)
	filter := d.Get("filter").(string)

	if id == "" && filter == "" {
		return fmt.Errorf("one of id or filter is required")
	}

	if id != "" {
		getResp, err := apiClient.StorageProfile.GetStorageProfile(storage_profile.NewGetStorageProfileParams().WithID(id))

		if err != nil {
			return err
		}
		storageProfile = getResp.GetPayload()
	} else {
		getResp, err := apiClient.StorageProfile.GetStorageProfiles(storage_profile.NewGetStorageProfilesParams().WithDollarFilter(withString(filter)))
		if err != nil {
			return err
		}

		storageProfiles := *getResp.Payload
		if len(storageProfiles.Content) > 1 {
			return fmt.Errorf("vra_storage_profile must filter to a storage profile")
		}
		if len(storageProfiles.Content) == 0 {
			return fmt.Errorf("vra_storage_profile filter did not match any storage profile")
		}

		storageProfile = storageProfiles.Content[0]
	}

	d.SetId(*storageProfile.ID)
	d.Set("created_at", storageProfile.CreatedAt)
	d.Set("default_item", storageProfile.DefaultItem)
	d.Set("description", storageProfile.Description)
	d.Set("disk_properties", storageProfile.DiskProperties)
	d.Set("external_region_id", storageProfile.ExternalRegionID)
	d.Set("name", storageProfile.Name)
	d.Set("org_id", storageProfile.OrgID)
	d.Set("owner", storageProfile.Owner)
	d.Set("supports_encryption", storageProfile.SupportsEncryption)
	d.Set("updated_at", storageProfile.UpdatedAt)

	if err := d.Set("tags", flattenTags(storageProfile.Tags)); err != nil {
		return fmt.Errorf("error setting storage profile tags - error: %v", err)
	}

	if err := d.Set("links", flattenLinks(storageProfile.Links)); err != nil {
		return fmt.Errorf("error setting storage profile links - error: %#v", err)
	}

	log.Printf("Finished reading the vra_storage_profile data source with filter %s", d.Get("filter"))
	return nil
}
