// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"
	"log"

	"github.com/vmware/vra-sdk-go/pkg/client/storage_profile"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceStorageProfileAzure() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStorageProfileAzureCreate,
		ReadContext:   resourceStorageProfileAzureRead,
		UpdateContext: resourceStorageProfileAzureUpdate,
		DeleteContext: resourceStorageProfileAzureDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_item": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_disk_caching": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"disk_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"os_disk_caching": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"storage_account_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"supports_encryption": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"tags": tagsSchema(),
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": linksSchema(),
			"organization_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceStorageProfileAzureCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to create vra_azure_storage_profile resource")
	apiClient := m.(*Client).apiClient

	name := d.Get("name").(string)
	regionID := d.Get("region_id").(string)

	StorageProfileAzureSpecification := models.StorageProfileAzureSpecification{
		DefaultItem:        d.Get("default_item").(bool),
		DiskType:           d.Get("disk_type").(string),
		DataDiskCaching:    d.Get("data_disk_caching").(string),
		Name:               &name,
		OsDiskCaching:      d.Get("os_disk_caching").(string),
		RegionID:           &regionID,
		StorageAccountID:   d.Get("storage_account_id").(string),
		SupportsEncryption: d.Get("supports_encryption").(bool),
		Tags:               expandTags(d.Get("tags").(*schema.Set).List()),
	}

	if v, ok := d.GetOk("description"); ok {
		StorageProfileAzureSpecification.Description = v.(string)
	}

	log.Printf("[DEBUG] create azure storage profile: %#v", StorageProfileAzureSpecification)
	createAzureStorageProfileCreated, err := apiClient.StorageProfile.CreateAzureStorageProfile(storage_profile.NewCreateAzureStorageProfileParams().WithBody(&StorageProfileAzureSpecification))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*createAzureStorageProfileCreated.Payload.ID)
	log.Printf("Finished to create vra_azure_storage_profile resource with name %s", d.Get("name"))

	return resourceStorageProfileAzureRead(ctx, d, m)
}

func resourceStorageProfileAzureRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Reading the vra_azure_storage_profile resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	resp, err := apiClient.StorageProfile.GetAzureStorageProfile(storage_profile.NewGetAzureStorageProfileParams().WithID(id))
	if err != nil {
		return diag.FromErr(err)
	}

	AzureStorageProfile := *resp.Payload
	d.Set("created_at", AzureStorageProfile.CreatedAt)
	d.Set("default_item", AzureStorageProfile.DefaultItem)
	d.Set("description", AzureStorageProfile.Description)
	d.Set("disk_type", AzureStorageProfile.DiskType)
	d.Set("data_disk_caching", AzureStorageProfile.DataDiskCaching)
	d.Set("external_region_id", AzureStorageProfile.ExternalRegionID)
	d.Set("name", AzureStorageProfile.Name)
	d.Set("organization_id", AzureStorageProfile.OrgID)
	d.Set("os_disk_caching", AzureStorageProfile.OsDiskCaching)
	d.Set("owner", AzureStorageProfile.Owner)
	d.Set("supports_encryption", AzureStorageProfile.SupportsEncryption)
	d.Set("updated_at", AzureStorageProfile.UpdatedAt)

	if err := d.Set("tags", flattenTags(AzureStorageProfile.Tags)); err != nil {
		return diag.Errorf("error setting azure storage profile tags - error: %v", err)
	}

	if err := d.Set("links", flattenLinks(AzureStorageProfile.Links)); err != nil {
		return diag.Errorf("error setting azure storage profile links - error: %#v", err)
	}

	log.Printf("Finished reading the vra_azure_storage_profile resource with name %s", d.Get("name"))
	return nil
}

func resourceStorageProfileAzureUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	name := d.Get("name").(string)
	regionID := d.Get("region_id").(string)

	StorageProfileAzureSpecification := models.StorageProfileAzureSpecification{
		DefaultItem:        d.Get("default_item").(bool),
		DiskType:           d.Get("disk_type").(string),
		DataDiskCaching:    d.Get("data_disk_caching").(string),
		Name:               &name,
		OsDiskCaching:      d.Get("os_disk_caching").(string),
		RegionID:           &regionID,
		StorageAccountID:   d.Get("storage_account_id").(string),
		SupportsEncryption: d.Get("supports_encryption").(bool),
		Tags:               expandTags(d.Get("tags").(*schema.Set).List()),
	}

	if v, ok := d.GetOk("description"); ok {
		StorageProfileAzureSpecification.Description = v.(string)
	}
	_, err := apiClient.StorageProfile.UpdateAzureStorageProfile(storage_profile.NewUpdateAzureStorageProfileParams().WithID(id).WithBody(&StorageProfileAzureSpecification))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceStorageProfileAzureRead(ctx, d, m)
}

func resourceStorageProfileAzureDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to delete the vra_azure_storage_profile resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	_, err := apiClient.StorageProfile.DeleteAzureStorageProfile(storage_profile.NewDeleteAzureStorageProfileParams().WithID(id))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("Finished deleting the vra_azure_storage_profile resource with name %s", d.Get("name"))
	return nil
}
