package vra

import (
	"fmt"
	"log"

	"github.com/vmware/vra-sdk-go/pkg/client/storage_profile"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceStorageProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceStorageProfileCreate,
		Read:   resourceStorageProfileRead,
		Update: resourceStorageProfileUpdate,
		Delete: resourceStorageProfileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
				Required:    true,
				Description: "Indicates if this storage profile is a default profile.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A human-friendly description.",
			},
			"disk_properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Map of storage properties that are to be applied on disk while provisioning.",
			},
			"disk_target_properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Map of storage placements to know where the disk is provisioned.",
			},
			"external_region_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the region as seen in the cloud provider for which this profile is defined.",
			},
			"links": linksSchema(),
			"name": {
				Type:        schema.TypeString,
				Required:    true,
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
			"region_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the region that is associated with the storage profile.",
			},
			"supports_encryption": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
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

func resourceStorageProfileCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("Starting to create vra_storage_profile resource")
	apiClient := m.(*Client).apiClient

	name := d.Get("name").(string)
	regionID := d.Get("region_id").(string)
	defaultItem := d.Get("default_item").(bool)

	storageProfileSpecification := models.StorageProfileSpecification{
		DefaultItem:          &defaultItem,
		DiskProperties:       expandCustomProperties(d.Get("disk_properties").(map[string]interface{})),
		DiskTargetProperties: expandCustomProperties(d.Get("disk_target_properties").(map[string]interface{})),
		Name:                 &name,
		RegionID:             &regionID,
		SupportsEncryption:   d.Get("supports_encryption").(bool),
		Tags:                 expandTags(d.Get("tags").(*schema.Set).List()),
	}

	if v, ok := d.GetOk("description"); ok {
		storageProfileSpecification.Description = v.(string)
	}

	log.Printf("[DEBUG] create storage profile: %#v", storageProfileSpecification)
	createStorageProfileCreated, err := apiClient.StorageProfile.CreateStorageProfile(storage_profile.NewCreateStorageProfileParams().WithBody(&storageProfileSpecification))
	if err != nil {
		return err
	}

	d.SetId(*createStorageProfileCreated.Payload.ID)
	log.Printf("Finished to create vra_storage_profile resource with name %s", d.Get("name"))

	return resourceStorageProfileRead(d, m)
}

func resourceStorageProfileRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("Reading the vra_storage_profile resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	resp, err := apiClient.StorageProfile.GetStorageProfile(storage_profile.NewGetStorageProfileParams().WithID(id))
	if err != nil {
		return err
	}

	storageProfile := *resp.Payload
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

	log.Printf("Finished reading the vra_storage_profile resource with name %s", d.Get("name"))
	return nil
}

func resourceStorageProfileUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("Starting to update the vra_storage_profile resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	name := d.Get("name").(string)
	regionID := d.Get("region_id").(string)
	defaultItem := d.Get("default_item").(bool)

	storageProfileSpecification := models.StorageProfileSpecification{
		DefaultItem:          &defaultItem,
		DiskProperties:       expandCustomProperties(d.Get("disk_properties").(map[string]interface{})),
		DiskTargetProperties: expandCustomProperties(d.Get("disk_target_properties").(map[string]interface{})),
		Name:                 &name,
		RegionID:             &regionID,
		SupportsEncryption:   d.Get("supports_encryption").(bool),
		Tags:                 expandTags(d.Get("tags").(*schema.Set).List()),
	}

	if v, ok := d.GetOk("description"); ok {
		storageProfileSpecification.Description = v.(string)
	}
	_, err := apiClient.StorageProfile.ReplaceStorageProfile(storage_profile.NewReplaceStorageProfileParams().WithID(id).WithBody(&storageProfileSpecification))
	if err != nil {
		return err
	}

	log.Printf("Finished updating the vra_storage_profile resource with name %s", d.Get("name"))
	return resourceStorageProfileRead(d, m)
}

func resourceStorageProfileDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("Starting to delete the vra_storage_profile resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	_, err := apiClient.StorageProfile.DeleteStorageProfile(storage_profile.NewDeleteStorageProfileParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *storage_profile.GetStorageProfileNotFound:
			log.Printf("vra_storage_profile resource with id '%s' is not found", d.Get("name"))
		default:
			return err
		}
	}

	d.SetId("")
	log.Printf("Finished deleting the vra_storage_profile resource with name %s", d.Get("name"))
	return nil
}
