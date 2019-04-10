package cas

import (
	"fmt"

	"github.com/vmware/cas-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/cas-sdk-go/pkg/models"
	tango "github.com/vmware/terraform-provider-cas/sdk"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudAccountAWS() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudAccountAWSCreate,
		Read:   resourceCloudAccountAWSRead,
		Update: resourceCloudAccountAWSUpdate,
		Delete: resourceCloudAccountAWSDelete,

		Schema: map[string]*schema.Schema{
			"access_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"regions": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"region_ids": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"secret_key": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"tags": tagsSDKSchema(),
		},
	}
}

func resourceCloudAccountAWSCreate(d *schema.ResourceData, m interface{}) error {
	var regions []string

	client := m.(*tango.Client)
	apiClient := client.GetAPIClient()

	accessKey := d.Get("access_key").(string)
	description := d.Get("description").(string)
	name := d.Get("name").(string)
	secretAccessKey := d.Get("secret_key").(string)
	tags := expandSDKTags(d.Get("tags").(*schema.Set).List())

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
	regionsIds, err := flattenAndNormalizeCLoudAccountRegionIds(regions, createResp.Payload)
	if err != nil {
		return err
	}
	d.Set("region_ids", regionsIds)

	if err := d.Set("tags", flattenSDKTags(tags)); err != nil {
		return fmt.Errorf("Error setting cloud account tags - error: %#v", err)
	}
	d.SetId(*createResp.Payload.ID)

	return resourceCloudAccountAWSRead(d, m)
}

func resourceCloudAccountAWSRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)
	apiClient := client.GetAPIClient()

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

	d.Set("access_key", awsAccount.AccessKeyID)
	d.Set("description", awsAccount.Description)
	d.Set("name", awsAccount.Name)
	regions := awsAccount.EnabledRegionIds
	d.Set("regions", regions)

	// The returned EnabledRegionIds and Hrefs containing the region ids can be in a different order than the request order.
	// Call a routine to normalize the order to correspond with the users region order.
	regionsIds, err := flattenAndNormalizeCLoudAccountAWSRegionIds(regions, &awsAccount)
	if err != nil {
		return err
	}
	d.Set("region_ids", regionsIds)

	if err := d.Set("tags", flattenSDKTags(awsAccount.Tags)); err != nil {
		return fmt.Errorf("Error setting cloud account tags - error: %#v", err)
	}

	return nil
}

func resourceCloudAccountAWSUpdate(d *schema.ResourceData, m interface{}) error {
	var regions []string

	client := m.(*tango.Client)
	apiClient := client.GetAPIClient()

	id := d.Id()
	accessKey := d.Get("access_key").(string)
	description := d.Get("description").(string)
	name := d.Get("name").(string)
	secretAccessKey := d.Get("secret_key").(string)
	tags := expandSDKTags(d.Get("tags").(*schema.Set).List())

	if v, ok := d.GetOk("regions"); ok {
		if !compareUnique(v.([]interface{})) {
			return fmt.Errorf("Specified regions are not unique")
		}
		regions = expandStringList(v.([]interface{}))
	}
	_, err := apiClient.CloudAccount.UpdateCloudAccount(cloud_account.NewUpdateCloudAccountParams().WithID(id).WithBody(&models.CloudAccountSpecification{
		PrivateKeyID:       &accessKey,
		CreateDefaultZones: false,
		Description:        description,
		Name:               &name,
		PrivateKey:         &secretAccessKey,
		RegionIds:          regions,
		Tags:               tags,
	}))
	if err != nil {
		return err
	}

	return resourceCloudAccountAWSRead(d, m)
}

func resourceCloudAccountAWSDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)
	apiClient := client.GetAPIClient()

	id := d.Id()
	_, err := apiClient.CloudAccount.DeleteAwsCloudAccount(cloud_account.NewDeleteAwsCloudAccountParams().WithID(id))
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
