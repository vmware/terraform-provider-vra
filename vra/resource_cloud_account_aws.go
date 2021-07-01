package vra

import (
	"fmt"

	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudAccountAWS() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudAccountAWSCreate,
		Read:   resourceCloudAccountAWSRead,
		Update: resourceCloudAccountAWSUpdate,
		Delete: resourceCloudAccountAWSDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// Required arguments
			"access_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"secret_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			// Optional arguments
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"regions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": tagsSchema(),
			// Computed attributes
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": linksSchema(),
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudAccountAWSCreate(d *schema.ResourceData, m interface{}) error {
	var regions []string

	apiClient := m.(*Client).apiClient

	accessKey := d.Get("access_key").(string)
	description := d.Get("description").(string)
	name := d.Get("name").(string)
	secretAccessKey := d.Get("secret_key").(string)
	tags := expandTags(d.Get("tags").(*schema.Set).List())

	if v, ok := d.GetOk("regions"); ok {
		if !compareUnique(v.([]interface{})) {
			return fmt.Errorf("Specified regions are not unique")
		}
		regions = expandStringList(v.([]interface{}))
	}

	createResp, err := apiClient.CloudAccount.CreateAwsCloudAccount(cloud_account.NewCreateAwsCloudAccountParams().WithBody(&models.CloudAccountAwsSpecification{
		AccessKeyID:        &accessKey,
		CreateDefaultZones: false,
		Description:        description,
		Name:               &name,
		SecretAccessKey:    &secretAccessKey,
		RegionIds:          regions,
		Tags:               tags,
	}))

	if err != nil {
		return err
	}

	// The returned EnabledRegionIds and Hrefs containing the region ids can be in a different order than the request order.
	// Call a routine to normalize the order to correspond with the users region order.
	regionsIds, err := flattenAndNormalizeCLoudAccountAWSRegionIds(regions, createResp.Payload)
	if err != nil {
		return err
	}
	d.Set("region_ids", regionsIds)

	if err := d.Set("tags", flattenTags(tags)); err != nil {
		return fmt.Errorf("Error setting cloud account tags - error: %#v", err)
	}
	d.SetId(*createResp.Payload.ID)

	return resourceCloudAccountAWSRead(d, m)
}

func resourceCloudAccountAWSRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	ret, err := apiClient.CloudAccount.GetAwsCloudAccount(cloud_account.NewGetAwsCloudAccountParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *cloud_account.GetAwsCloudAccountNotFound:
			d.SetId("")
			return nil
		}
		return err
	}
	awsAccount := *ret.Payload
	regions := awsAccount.EnabledRegionIds

	d.Set("access_key", awsAccount.AccessKeyID)
	d.Set("created_at", awsAccount.CreatedAt)
	d.Set("description", awsAccount.Description)
	d.Set("name", awsAccount.Name)
	d.Set("org_id", awsAccount.OrgID)
	d.Set("owner", awsAccount.Owner)
	d.Set("regions", regions)
	d.Set("updated_at", awsAccount.UpdatedAt)

	if err := d.Set("links", flattenLinks(awsAccount.Links)); err != nil {
		return fmt.Errorf("error setting cloud_account_aws links - error: %#v", err)
	}

	// The returned EnabledRegionIds and Hrefs containing the region ids can be in a different order than the request order.
	// Call a routine to normalize the order to correspond with the users region order.
	regionsIds, err := flattenAndNormalizeCLoudAccountAWSRegionIds(regions, &awsAccount)
	if err != nil {
		return err
	}
	d.Set("region_ids", regionsIds)

	if err := d.Set("tags", flattenTags(awsAccount.Tags)); err != nil {
		return fmt.Errorf("Error setting cloud account tags - error: %#v", err)
	}

	return nil
}

func resourceCloudAccountAWSUpdate(d *schema.ResourceData, m interface{}) error {
	var regions []string

	apiClient := m.(*Client).apiClient

	id := d.Id()
	description := d.Get("description").(string)
	tags := expandTags(d.Get("tags").(*schema.Set).List())

	if v, ok := d.GetOk("regions"); ok {
		if !compareUnique(v.([]interface{})) {
			return fmt.Errorf("Specified regions are not unique")
		}
		regions = expandStringList(v.([]interface{}))
	}
	_, err := apiClient.CloudAccount.UpdateAwsCloudAccount(cloud_account.NewUpdateAwsCloudAccountParams().WithID(id).WithBody(&models.UpdateCloudAccountAwsSpecification{
		CreateDefaultZones: false,
		Description:        description,
		RegionIds:          regions,
		Tags:               tags,
	}))
	if err != nil {
		return err
	}

	return resourceCloudAccountAWSRead(d, m)
}

func resourceCloudAccountAWSDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	_, err := apiClient.CloudAccount.DeleteAwsCloudAccount(cloud_account.NewDeleteAwsCloudAccountParams().WithID(id))
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
