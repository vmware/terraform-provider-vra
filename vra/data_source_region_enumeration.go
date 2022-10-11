package vra

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/client/request"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceRegionEnumeration() *schema.Resource {
	return &schema.Resource{
		ReadContext:        dataSourceRegionEnumerationRead,
		DeprecationMessage: "'region_enumeration' is deprecated. Use 'region_enumeration_vsphere' instead.",

		Schema: map[string]*schema.Schema{
			"accept_self_signed_cert": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to accept self signed certificate when connecting to the vCenter Server.",
			},
			"dcid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identifier of a data collector vm deployed in the on premise infrastructure.",
				Deprecated:  "Please use `dc_id` instead.",
			},
			"dc_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"dcid"},
				Description:   "Identifier of a data collector vm deployed in the on premise infrastructure.",
			},
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP address or FQDN of the vCenter Server.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password of the vCenter Server.",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username of the vCenter Server.",
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

func dataSourceRegionEnumerationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*Client).apiClient

	dcid := ""
	if v, ok := d.GetOk("dc_id"); ok {
		dcid = v.(string)
	} else if v, ok := d.GetOk("dcid"); ok {
		dcid = v.(string)
	}
	enumResp, err := apiClient.CloudAccount.EnumerateVSphereRegionsAsync(
		cloud_account.NewEnumerateVSphereRegionsAsyncParams().
			WithAPIVersion(withString(IaaSAPIVersion)).
			WithTimeout(IncreasedTimeOut).
			WithBody(&models.CloudAccountVsphereRegionEnumerationSpecification{
				AcceptSelfSignedCertificate: d.Get("accept_self_signed_cert").(bool),
				Dcid:                        dcid,
				HostName:                    d.Get("hostname").(string),
				Password:                    d.Get("password").(string),
				Username:                    d.Get("username").(string),
			}))
	if err != nil {
		return diag.FromErr(err)
	}

	stateChangeFunc := resource.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    dataSourceRegionEnumerationReadRefreshFunc(*apiClient, *enumResp.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutRead),
		MinTimeout: 5 * time.Second,
	}

	resourceIds, err := stateChangeFunc.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	enumID := (resourceIds.([]string))[0]

	getResp, err := apiClient.CloudAccount.GetRegionEnumerationResult(
		cloud_account.NewGetRegionEnumerationResultParams().
			WithAPIVersion(withString(IaaSAPIVersion)).
			WithTimeout(IncreasedTimeOut).
			WithID(enumID))
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("regions", extractIdsFromRegionSpecification(getResp.Payload.ExternalRegions))
	d.SetId(d.Get("hostname").(string))

	return nil
}

func dataSourceRegionEnumerationReadRefreshFunc(apiClient client.API, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		reqResp, err := apiClient.Request.GetRequestTracker(request.NewGetRequestTrackerParams().WithID(id))
		if err != nil {
			return "", models.RequestTrackerStatusFAILED, err
		}

		status := reqResp.Payload.Status
		switch *status {
		case models.RequestTrackerStatusFAILED:
			return []string{""}, *status, fmt.Errorf(reqResp.Payload.Message)
		case models.RequestTrackerStatusINPROGRESS:
			return [...]string{id}, *status, nil
		case models.RequestTrackerStatusFINISHED:
			regionEnumerationIds := make([]string, len(reqResp.Payload.Resources))
			for i, r := range reqResp.Payload.Resources {
				regionEnumerationIds[i] = strings.TrimPrefix(r, "/iaas/api/cloud-accounts/region-enumeration/")
			}
			return regionEnumerationIds, *status, nil
		default:
			return [...]string{id}, reqResp.Payload.Message, fmt.Errorf("dataSourceRegionEnumerationReadRefreshFunc: unknown status %v", *status)
		}
	}
}
