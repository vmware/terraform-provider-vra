package vra

import (
	"fmt"

	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudAccountAzure() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudAccountAzureCreate,
		Read:   resourceCloudAccountAzureRead,
		Update: resourceCloudAccountAzureUpdate,
		Delete: resourceCloudAccountAzureDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required arguments
			"application_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"application_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subscription_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
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
			//Computed attributes
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

func resourceCloudAccountAzureCreate(d *schema.ResourceData, m interface{}) error {
	var regions []string

	apiClient := m.(*Client).apiClient

	if v, ok := d.GetOk("regions"); ok {
		if !compareUnique(v.([]interface{})) {
			return fmt.Errorf("Specified regions are not unique")
		}
		regions = expandStringList(v.([]interface{}))
	}

	applicationKey := d.Get("application_key").(string)

	createResp, err := apiClient.CloudAccount.CreateAzureCloudAccount(cloud_account.NewCreateAzureCloudAccountParams().WithBody(&models.CloudAccountAzureSpecification{
		Description:                d.Get("description").(string),
		Name:                       withString(d.Get("name").(string)),
		ClientApplicationID:        withString(d.Get("application_id").(string)),
		ClientApplicationSecretKey: &applicationKey,
		SubscriptionID:             withString(d.Get("subscription_id").(string)),
		TenantID:                   withString(d.Get("tenant_id").(string)),
		CreateDefaultZones:         false,
		RegionIds:                  regions,
		Tags:                       expandTags(d.Get("tags").(*schema.Set).List()),
	}))

	if err != nil {
		return err
	}

	d.Set("application_key", applicationKey)
	d.SetId(*createResp.Payload.ID)

	return resourceCloudAccountAzureRead(d, m)
}

func resourceCloudAccountAzureRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	ret, err := apiClient.CloudAccount.GetAzureCloudAccount(cloud_account.NewGetAzureCloudAccountParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *cloud_account.GetAzureCloudAccountNotFound:
			d.SetId("")
			return nil
		}
		return err
	}
	azureAccount := *ret.Payload
	regions := azureAccount.EnabledRegionIds

	d.Set("application_id", azureAccount.ClientApplicationID)
	d.Set("created_at", azureAccount.CreatedAt)
	d.Set("description", azureAccount.Description)
	d.Set("name", azureAccount.Name)
	d.Set("org_id", azureAccount.OrgID)
	d.Set("owner", azureAccount.Owner)
	d.Set("regions", regions)
	d.Set("subscription_id", azureAccount.SubscriptionID)
	d.Set("tenant_id", azureAccount.TenantID)
	d.Set("updated_at", azureAccount.UpdatedAt)

	if err := d.Set("links", flattenLinks(azureAccount.Links)); err != nil {
		return fmt.Errorf("error setting cloud_account_azure links - error: %#v", err)
	}

	// The returned EnabledRegionIds and Hrefs containing the region ids can be in a different order than the request order.
	// Call a routine to normalize the order to correspond with the users region order.
	regionsIds, err := flattenAndNormalizeCLoudAccountAzureRegionIds(regions, &azureAccount)
	if err != nil {
		return err
	}
	d.Set("region_ids", regionsIds)

	if err := d.Set("tags", flattenTags(azureAccount.Tags)); err != nil {
		return fmt.Errorf("Error setting cloud account tags - error: %#v", err)
	}

	return nil
}

func resourceCloudAccountAzureUpdate(d *schema.ResourceData, m interface{}) error {
	var regions []string

	apiClient := m.(*Client).apiClient

	id := d.Id()

	if v, ok := d.GetOk("regions"); ok {
		if !compareUnique(v.([]interface{})) {
			return fmt.Errorf("Specified regions are not unique")
		}
		regions = expandStringList(v.([]interface{}))
	}
	tags := expandTags(d.Get("tags").(*schema.Set).List())

	_, err := apiClient.CloudAccount.UpdateAzureCloudAccount(cloud_account.NewUpdateAzureCloudAccountParams().WithID(id).WithBody(&models.UpdateCloudAccountAzureSpecification{
		Description:        d.Get("description").(string),
		CreateDefaultZones: false,
		RegionIds:          regions,
		Tags:               tags,
	}))
	if err != nil {
		return err
	}

	return resourceCloudAccountAzureRead(d, m)
}

func resourceCloudAccountAzureDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	_, err := apiClient.CloudAccount.DeleteAzureCloudAccount(cloud_account.NewDeleteAzureCloudAccountParams().WithID(id))
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
