// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceRegionEnumerationVMC() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRegionEnumerationVMCRead,

		Schema: map[string]*schema.Schema{
			"accept_self_signed_cert": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to accept self signed certificate when connecting to the vCenter Server.",
			},
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VMC API access key.",
			},
			"dc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identifier of a data collector vm deployed in the on premise infrastructure.",
			},
			"nsx_hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The IP address of the NSX Manager server in the specified SDDC / FQDN.",
			},
			"sddc_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifier of the on-premise SDDC.",
			},
			"vcenter_hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP address or FQDN of the vCenter Server in the specified SDDC.",
			},
			"vcenter_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password of the vCenter Server in the specified SDDC.",
			},
			"vcenter_username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username of the vCenter Server in the specified SDDC.",
			},
			"regions": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "A set of region ids that can be enabled for this cloud account.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceRegionEnumerationVMCRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*Client).apiClient

	enumResp, err := apiClient.CloudAccount.EnumerateVmcRegionsAsync(
		cloud_account.NewEnumerateVmcRegionsAsyncParams().
			WithAPIVersion(IaaSAPIVersion).
			WithBody(&models.CloudAccountVmcRegionEnumerationSpecification{
				AcceptSelfSignedCertificate: d.Get("accept_self_signed_cert").(bool),
				APIKey:                      d.Get("api_token").(string),
				DcID:                        d.Get("dc_id").(string),
				HostName:                    d.Get("vcenter_hostname").(string),
				NsxHostName:                 d.Get("nsx_hostname").(string),
				Password:                    d.Get("vcenter_password").(string),
				SddcID:                      d.Get("sddc_name").(string),
				Username:                    d.Get("vcenter_username").(string),
			}))
	if err != nil {
		return diag.FromErr(err)
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    dataSourceRegionEnumerationReadRefreshFunc(*apiClient, *enumResp.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutRead),
		MinTimeout: 5 * time.Second,
	}

	resourceIDs, err := stateChangeFunc.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	enumID := (resourceIDs.([]string))[0]

	getResp, err := apiClient.CloudAccount.GetRegionEnumerationResult(
		cloud_account.NewGetRegionEnumerationResultParams().
			WithAPIVersion(IaaSAPIVersion).
			WithID(enumID))
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("regions", extractIDsFromRegionSpecification(getResp.Payload.ExternalRegions))
	d.SetId(d.Get("vcenter_hostname").(string))

	return nil
}
