package vra

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceRegionEnumerationAzure() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRegionEnumerationAzureRead,

		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"application_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"subscription_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceRegionEnumerationAzureRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	getResp, err := apiClient.CloudAccount.EnumerateAzureRegions(cloud_account.NewEnumerateAzureRegionsParams().WithBody(&models.CloudAccountAzureSpecification{
		ClientApplicationID:        withString(d.Get("application_id").(string)),
		ClientApplicationSecretKey: withString(d.Get("application_key").(string)),
		SubscriptionID:             withString(d.Get("subscription_id").(string)),
		TenantID:                   withString(d.Get("tenant_id").(string)),
	}))

	if err != nil {
		return err
	}

	d.Set("regions", getResp.Payload.ExternalRegionIds)
	d.SetId(d.Get("application_id").(string))

	return nil
}
