package vra

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceRegionEnumerationGCP() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRegionEnumerationGCPRead,

		Schema: map[string]*schema.Schema{
			"client_email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"private_key_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceRegionEnumerationGCPRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	getResp, err := apiClient.CloudAccount.EnumerateGcpRegions(cloud_account.NewEnumerateGcpRegionsParams().WithBody(&models.CloudAccountGcpSpecification{
		ClientEmail:  withString(d.Get("client_email").(string)),
		PrivateKey:   withString(d.Get("private_key").(string)),
		PrivateKeyID: withString(d.Get("private_key_id").(string)),
		ProjectID:    withString(d.Get("project_id").(string)),
	}))

	if err != nil {
		return err
	}

	d.Set("regions", getResp.Payload.ExternalRegionIds)
	d.SetId(d.Get("private_key_id").(string))

	return nil
}
