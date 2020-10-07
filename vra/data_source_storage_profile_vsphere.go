package vra

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/storage_profile"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"log"
)

func dataSourceStorageProfileVsphere() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceStorageProfileVsphereRead,

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
				Optional:    true,
				Description: "Indicates if a storage profile acts as a default storage profile for a disk.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly description.",
			},
			"disk_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of mode for the disk.",
			},
			"disk_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Disk types are specified as standard or first class, empty value is considered as standard.",
			},
			"external_region_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the region for which this profile is defined",
			},
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
			"limit_iops": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The upper bound for the I/O operations per second allocated for each virtual disk.",
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
			"provisioning_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of provisioning policy for the disk.",
			},
			"shares": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A specific number of shares assigned to each virtual machine.",
			},
			"shares_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Indicates whether this storage profile supports encryption or not.",
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

func dataSourceStorageProfileVsphereRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading the vra_storage_profile_vsphere data source with filter %s", d.Get("filter"))
	apiClient := meta.(*Client).apiClient

	var storageProfile *models.VsphereStorageProfile

	id := d.Get("id").(string)
	filter := d.Get("filter").(string)

	if id == "" && filter == "" {
		return fmt.Errorf("one of id or filter is required")
	}

	if id != "" {
		getResp, err := apiClient.StorageProfile.GetVSphereStorageProfile(storage_profile.NewGetVSphereStorageProfileParams().WithID(id))
		if err != nil {
			return err
		}
		storageProfile = getResp.GetPayload()
	} else {
		getResp, err := apiClient.StorageProfile.GetVSphereStorageProfiles(storage_profile.NewGetVSphereStorageProfilesParams().WithDollarFilter(withString(filter)))
		if err != nil {
			return err
		}

		storageProfiles := *getResp.Payload
		if len(storageProfiles.Content) > 1 {
			return fmt.Errorf("vra_storage_profile_vsphere must filter to a storage profile")
		}
		if len(storageProfiles.Content) == 0 {
			return fmt.Errorf("vra_storage_profile_vsphere filter did not match any storage profile")
		}

		storageProfile = storageProfiles.Content[0]
	}

	d.SetId(*storageProfile.ID)
	d.Set("created_at", storageProfile.CreatedAt)
	d.Set("default_item", storageProfile.DefaultItem)
	d.Set("description", storageProfile.Description)
	d.Set("disk_mode", storageProfile.DiskMode)
	d.Set("disk_type", storageProfile.DiskType)
	d.Set("external_region_id", storageProfile.ExternalRegionID)
	d.Set("limit_iops", storageProfile.LimitIops)
	d.Set("name", storageProfile.Name)
	d.Set("org_id", storageProfile.OrgID)
	d.Set("owner", storageProfile.Owner)
	d.Set("provisioning_type", storageProfile.ProvisioningType)
	d.Set("shares", storageProfile.Shares)
	d.Set("shares_level", storageProfile.SharesLevel)
	d.Set("supports_encryption", storageProfile.SupportsEncryption)
	d.Set("updated_at", storageProfile.UpdatedAt)

	if err := d.Set("tags", flattenTags(storageProfile.Tags)); err != nil {
		return fmt.Errorf("error setting storage profile tags - error: %v", err)
	}

	if err := d.Set("links", flattenLinks(storageProfile.Links)); err != nil {
		return fmt.Errorf("error setting storage profile links - error: %#v", err)
	}

	log.Printf("Finished reading the vra_storage_profile_vsphere data source with filter %s", d.Get("filter"))
	return nil
}
