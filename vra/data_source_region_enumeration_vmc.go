package vra

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceRegionEnumerationVMC() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRegionEnumerationVMCRead,

		Schema: map[string]*schema.Schema{
			"accept_self_signed_cert": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"api_token": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dc_id": {
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
			"sddc_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vcenter_hostname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vcenter_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"vcenter_username": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceRegionEnumerationVMCRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	getResp, err := apiClient.CloudAccount.EnumerateVmcRegions(
		cloud_account.NewEnumerateVmcRegionsParams().
			WithTimeout(IncreasedTimeOut).
			WithBody(&models.CloudAccountVmcSpecification{
				AcceptSelfSignedCertificate: d.Get("accept_self_signed_cert").(bool),
				APIKey:                      d.Get("api_token").(string),
				DcID:                        d.Get("dc_id").(string),
				HostName:                    withString(d.Get("vcenter_hostname").(string)),
				Password:                    withString(d.Get("vcenter_password").(string)),
				SddcID:                      d.Get("sddc_name").(string),
				Username:                    withString(d.Get("vcenter_username").(string)),
			}))

	if err != nil {
		return err
	}

	d.Set("regions", getResp.Payload.ExternalRegionIds)
	d.SetId(d.Get("vcenter_hostname").(string))

	return nil
}
