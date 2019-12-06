package vra

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/storage_profile"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"log"
)

func dataSourceStorageProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceStorageProfileRead,

		Schema: map[string]*schema.Schema{
			"default_item": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_properties": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"external_region_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": linksSchema(),
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"org_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"supports_encryption": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"tags": tagsSchema(),
			"updated_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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
