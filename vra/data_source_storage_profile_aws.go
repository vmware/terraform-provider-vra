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

func datasourceStorageProfileAws() *schema.Resource {
	return &schema.Resource{
		Read: datasourceStorageProfileAwsRead,

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
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly description.",
			},
			"device_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the type of storage device.",
			},
			"external_region_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the region for which this profile is defined",
			},
			"iops": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates maximum I/O operations per second in range(1-20,000).",
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
			"volume_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the type of volume associated with type of storage device.",
			},
		},
	}
}

func datasourceStorageProfileAwsRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("Reading the vra_storage_profile_aws data source with name %s", d.Get("filter"))
	apiClient := m.(*Client).apiClient

	var awsStorageProfile *models.AwsStorageProfile

	id := d.Get("id").(string)
	filter := d.Get("filter").(string)

	if id == "" && filter == "" {
		return fmt.Errorf("one of id or filter is required")
	}

	if id != "" {
		getResp, err := apiClient.StorageProfile.GetAwsStorageProfile(storage_profile.NewGetAwsStorageProfileParams().WithID(id))
		if err != nil {
			return err
		}
		awsStorageProfile = getResp.GetPayload()
	} else {
		getResp, err := apiClient.StorageProfile.GetAwsStorageProfiles(storage_profile.NewGetAwsStorageProfilesParams().WithDollarFilter(withString(filter)))
		if err != nil {
			return err
		}
		storageProfiles := *getResp.Payload
		if len(storageProfiles.Content) > 1 {
			return fmt.Errorf("vra_storage_profile_aws must filter to a storage profile")
		}
		if len(storageProfiles.Content) == 0 {
			return fmt.Errorf("vra_storage_profile_aws filter did not match any storage profile")
		}
		awsStorageProfile = storageProfiles.Content[0]
	}

	d.SetId(*awsStorageProfile.ID)
	d.Set("cloud_account_id", awsStorageProfile.CloudAccountID)
	d.Set("created_at", awsStorageProfile.CreatedAt)
	d.Set("default_item", awsStorageProfile.DefaultItem)
	d.Set("description", awsStorageProfile.Description)
	d.Set("device_type", awsStorageProfile.DeviceType)
	d.Set("external_region_id", awsStorageProfile.ExternalRegionID)
	d.Set("iops", awsStorageProfile.Iops)
	d.Set("name", awsStorageProfile.Name)
	d.Set("org_id", awsStorageProfile.OrgID)
	d.Set("owner", awsStorageProfile.Owner)
	d.Set("supports_encryption", awsStorageProfile.SupportsEncryption)
	d.Set("updated_at", awsStorageProfile.UpdatedAt)
	d.Set("volume_type", awsStorageProfile.VolumeType)

	if err := d.Set("tags", flattenTags(awsStorageProfile.Tags)); err != nil {
		return fmt.Errorf("error setting aws storage profile tags - error: %v", err)
	}

	if err := d.Set("links", flattenLinks(awsStorageProfile.Links)); err != nil {
		return fmt.Errorf("error setting aws storage profile links - error: %#v", err)
	}

	log.Printf("Finished reading the vra_storage_profile_aws data source with name %s", d.Get("name"))
	return nil
}
