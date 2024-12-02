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

func resourceStorageProfileAws() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStorageProfileAwsCreate,
		ReadContext:   resourceStorageProfileAwsRead,
		UpdateContext: resourceStorageProfileAwsUpdate,
		DeleteContext: resourceStorageProfileAwsDelete,
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"device_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"iops": {
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
			"volume_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
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

func resourceStorageProfileAwsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to create vra_aws_storage_profile resource")
	apiClient := m.(*Client).apiClient

	name := d.Get("name").(string)
	regionID := d.Get("region_id").(string)
	deviceType := d.Get("device_type").(string)

	StorageProfileAwsSpecification := models.StorageProfileAwsSpecification{
		DefaultItem:        d.Get("default_item").(bool),
		DeviceType:         &deviceType,
		Iops:               d.Get("iops").(string),
		Name:               &name,
		RegionID:           &regionID,
		SupportsEncryption: d.Get("supports_encryption").(bool),
		Tags:               expandTags(d.Get("tags").(*schema.Set).List()),
		VolumeType:         d.Get("volume_type").(string),
	}

	if v, ok := d.GetOk("description"); ok {
		StorageProfileAwsSpecification.Description = v.(string)
	}

	log.Printf("[DEBUG] create aws storage profile: %#v", StorageProfileAwsSpecification)
	createAwsStorageProfileCreated, err := apiClient.StorageProfile.CreateAwsStorageProfile(storage_profile.NewCreateAwsStorageProfileParams().WithBody(&StorageProfileAwsSpecification))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*createAwsStorageProfileCreated.Payload.ID)
	log.Printf("Finished to create vra_Aws_storage_profile resource with name %s", d.Get("name"))

	return resourceStorageProfileAwsRead(ctx, d, m)
}

func resourceStorageProfileAwsRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Reading the vra_Aws_storage_profile resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	resp, err := apiClient.StorageProfile.GetAwsStorageProfile(storage_profile.NewGetAwsStorageProfileParams().WithID(id))
	if err != nil {
		return diag.FromErr(err)
	}

	awsStorageProfile := *resp.Payload
	d.Set("created_at", awsStorageProfile.CreatedAt)
	d.Set("default_item", awsStorageProfile.DefaultItem)
	d.Set("description", awsStorageProfile.Description)
	d.Set("device_type", awsStorageProfile.DeviceType)
	d.Set("external_region_id", awsStorageProfile.ExternalRegionID)
	d.Set("iops", awsStorageProfile.Iops)
	d.Set("name", awsStorageProfile.Name)
	d.Set("organization_id", awsStorageProfile.OrgID)
	d.Set("owner", awsStorageProfile.Owner)
	d.Set("supports_encryption", awsStorageProfile.SupportsEncryption)
	d.Set("updated_at", awsStorageProfile.UpdatedAt)
	d.Set("volume_type", awsStorageProfile.VolumeType)

	if err := d.Set("tags", flattenTags(awsStorageProfile.Tags)); err != nil {
		return diag.Errorf("error setting swa storage profile tags - error: %v", err)
	}

	if err := d.Set("links", flattenLinks(awsStorageProfile.Links)); err != nil {
		return diag.Errorf("error setting aws storage profile links - error: %#v", err)
	}

	log.Printf("Finished reading the vra_aws_storage_profile resource with name %s", d.Get("name"))
	return nil
}

func resourceStorageProfileAwsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	name := d.Get("name").(string)
	regionID := d.Get("region_id").(string)
	deviceType := d.Get("device_type").(string)

	StorageProfileAwsSpecification := models.StorageProfileAwsSpecification{
		DefaultItem:        d.Get("default_item").(bool),
		DeviceType:         &deviceType,
		Iops:               d.Get("iops").(string),
		Name:               &name,
		RegionID:           &regionID,
		SupportsEncryption: d.Get("supports_encryption").(bool),
		Tags:               expandTags(d.Get("tags").(*schema.Set).List()),
		VolumeType:         d.Get("volume_type").(string),
	}

	if v, ok := d.GetOk("description"); ok {
		StorageProfileAwsSpecification.Description = v.(string)
	}
	_, err := apiClient.StorageProfile.UpdateAwsStorageProfile(storage_profile.NewUpdateAwsStorageProfileParams().WithID(id).WithBody(&StorageProfileAwsSpecification))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceStorageProfileAwsRead(ctx, d, m)
}

func resourceStorageProfileAwsDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to delete the vra_aws_storage_profile resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	_, err := apiClient.StorageProfile.DeleteAwsStorageProfile(storage_profile.NewDeleteAwsStorageProfileParams().WithID(id))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("Finished deleting the vra_aws_storage_profile resource with name %s", d.Get("name"))
	return nil
}
