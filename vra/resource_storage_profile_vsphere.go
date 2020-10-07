package vra

import (
	"fmt"
	"log"

	"github.com/vmware/vra-sdk-go/pkg/client/storage_profile"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceStorageProfileVsphere() *schema.Resource {
	return &schema.Resource{
		Create: resourceStorageProfileVsphereCreate,
		Read:   resourceStorageProfileVsphereRead,
		Update: resourceStorageProfileVsphereUpdate,
		Delete: resourceStorageProfileVsphereDelete,
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
			"datastore_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of the vSphere Datastore for placing disk and VM.",
			},
			"default_item": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if a storage profile acts as a default storage profile for a disk.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A human-friendly description.",
			},
			"disk_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Type of mode for the disk. Omitting this value will set it to dependent. " +
					"example: dependent / independent-persistent / independent-nonpersistent.",
			},
			"disk_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Disk types are specified as standard or first class, empty value is considered as standard.",
			},
			"external_region_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the region for which this profile is defined",
			},
			"limit_iops": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The upper bound for the I/O operations per second allocated for each virtual disk.",
			},
			"links": linksSchema(),
			"name": {
				Type:        schema.TypeString,
				Required:    true,
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
				Optional:    true,
				Description: "Type of provisioning policy for the disk.",
			},
			"region_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Id of the region that is associated with the storage profile.",
			},
			"shares": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A specific number of shares assigned to each virtual machine.",
			},
			"shares_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Indicates whether this storage profile supports encryption or not.",
			},
			"storage_policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of the vSphere Storage Policy to be applied.",
			},
			"supports_encryption": {
				Type:        schema.TypeBool,
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

func resourceStorageProfileVsphereCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("Starting to create vra_storage_profile_vsphere resource")
	apiClient := m.(*Client).apiClient

	name := d.Get("name").(string)
	regionID := d.Get("region_id").(string)
	defaultItem := d.Get("default_item").(bool)

	StorageProfileVsphereSpecification := models.StorageProfileVsphereSpecification{
		DatastoreID:        d.Get("datastore_id").(string),
		DefaultItem:        &defaultItem,
		DiskMode:           d.Get("disk_mode").(string),
		DiskType:           d.Get("disk_type").(string),
		LimitIops:          d.Get("limit_iops").(string),
		Name:               &name,
		ProvisioningType:   d.Get("provisioning_type").(string),
		RegionID:           &regionID,
		Shares:             d.Get("shares").(string),
		SharesLevel:        d.Get("shares_level").(string),
		StoragePolicyID:    d.Get("storage_policy_id").(string),
		SupportsEncryption: d.Get("supports_encryption").(bool),
		Tags:               expandTags(d.Get("tags").(*schema.Set).List()),
	}

	if v, ok := d.GetOk("description"); ok {
		StorageProfileVsphereSpecification.Description = v.(string)
	}

	log.Printf("[DEBUG] create vsphere storage profile: %#v", StorageProfileVsphereSpecification)
	createVsphereStorageProfileCreated, err := apiClient.StorageProfile.CreateVSphereStorageProfile(storage_profile.NewCreateVSphereStorageProfileParams().WithBody(&StorageProfileVsphereSpecification))
	if err != nil {
		return err
	}

	d.SetId(*createVsphereStorageProfileCreated.Payload.ID)
	log.Printf("Finished to create vra_storage_profile_vsphere resource with name %s", d.Get("name"))

	return resourceStorageProfileVsphereRead(d, m)
}

func resourceStorageProfileVsphereRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("Reading the vra_storage_profile_vsphere resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	resp, err := apiClient.StorageProfile.GetVSphereStorageProfile(storage_profile.NewGetVSphereStorageProfileParams().WithID(id))
	if err != nil {
		return err
	}

	VsphereStorageProfile := *resp.Payload
	d.Set("created_at", VsphereStorageProfile.CreatedAt)
	d.Set("default_item", VsphereStorageProfile.DefaultItem)
	d.Set("description", VsphereStorageProfile.Description)
	d.Set("disk_mode", VsphereStorageProfile.DiskMode)
	d.Set("disk_type", VsphereStorageProfile.DiskType)
	d.Set("external_region_id", VsphereStorageProfile.ExternalRegionID)
	d.Set("limit_iops", VsphereStorageProfile.LimitIops)
	d.Set("name", VsphereStorageProfile.Name)
	d.Set("organization_id", VsphereStorageProfile.OrganizationID)
	d.Set("provisioning_type", VsphereStorageProfile.ProvisioningType)
	d.Set("owner", VsphereStorageProfile.Owner)
	d.Set("shares", VsphereStorageProfile.Shares)
	d.Set("shares_level", VsphereStorageProfile.SharesLevel)
	d.Set("supports_encryption", VsphereStorageProfile.SupportsEncryption)
	d.Set("updated_at", VsphereStorageProfile.UpdatedAt)

	if err := d.Set("tags", flattenTags(VsphereStorageProfile.Tags)); err != nil {
		return fmt.Errorf("error setting vsphere storage profile vsphere tags - error: %v", err)
	}

	if err := d.Set("links", flattenLinks(VsphereStorageProfile.Links)); err != nil {
		return fmt.Errorf("error setting vsphere storage profile vsphere links - error: %#v", err)
	}

	log.Printf("Finished reading the vra_storage_profile_vsphere resource with name %s", d.Get("name"))
	return nil
}

func resourceStorageProfileVsphereUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	name := d.Get("name").(string)
	regionID := d.Get("region_id").(string)
	defaultItem := d.Get("default_item").(bool)

	StorageProfileVsphereSpecification := models.StorageProfileVsphereSpecification{
		DatastoreID:        d.Get("datastore_id").(string),
		DefaultItem:        &defaultItem,
		Description:        d.Get("description").(string),
		DiskMode:           d.Get("disk_mode").(string),
		LimitIops:          d.Get("limit_iops").(string),
		Name:               &name,
		ProvisioningType:   d.Get("provisioning_type").(string),
		RegionID:           &regionID,
		Shares:             d.Get("shares").(string),
		SharesLevel:        d.Get("shares_level").(string),
		StoragePolicyID:    d.Get("storage_policy_id").(string),
		SupportsEncryption: d.Get("supports_encryption").(bool),
		Tags:               expandTags(d.Get("tags").(*schema.Set).List()),
	}

	if v, ok := d.GetOk("description"); ok {
		StorageProfileVsphereSpecification.Description = v.(string)
	}
	_, err := apiClient.StorageProfile.UpdateVSphereStorageProfile(storage_profile.NewUpdateVSphereStorageProfileParams().WithID(id).WithBody(&StorageProfileVsphereSpecification))
	if err != nil {
		return err
	}

	return resourceStorageProfileVsphereRead(d, m)
}

func resourceStorageProfileVsphereDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("Starting to delete the vra_storage_profile_vsphere resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	_, err := apiClient.StorageProfile.DeleteVSphereStorageProfile(storage_profile.NewDeleteVSphereStorageProfileParams().WithID(id))
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("Finished deleting the vra_storage_profile_vsphere resource with name %s", d.Get("name"))
	return nil
}
