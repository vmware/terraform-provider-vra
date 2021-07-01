package vra

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceRegionEnumerationAWS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRegionEnumerationAWSRead,

		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"secret_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceRegionEnumerationAWSRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	getResp, err := apiClient.CloudAccount.EnumerateAwsRegions(cloud_account.NewEnumerateAwsRegionsParams().WithBody(&models.CloudAccountAwsSpecification{
		AccessKeyID:     withString(d.Get("access_key").(string)),
		SecretAccessKey: withString(d.Get("secret_key").(string)),
	}))

	if err != nil {
		return err
	}

	d.Set("regions", getResp.Payload.ExternalRegionIds)
	d.SetId(d.Get("access_key").(string))

	return nil
}
