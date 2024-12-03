// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/vmware/vra-sdk-go/pkg/client"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/client/request"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudAccountVMC() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudAccountVMCCreate,
		ReadContext:   resourceCloudAccountVMCRead,
		UpdateContext: resourceCloudAccountVMCUpdate,
		DeleteContext: resourceCloudAccountVMCDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required arguments
			"api_token": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of this resource instance.",
			},
			"nsx_hostname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"regions": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The set of region ids that will be enabled for this cloud account.",
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

			// Optional arguments
			"accept_self_signed_cert": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to accept self signed certificate when connecting to the vCenter Server.",
			},
			"dc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identifier of a data collector vm deployed in the on premise infrastructure.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A human-friendly description.",
			},
			"tags": tagsSchema(),

			// Computed attributes
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was created. The date is in ISO 8601 and UTC.",
			},
			"links": linksSchema(),
			"org_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the organization this entity belongs to.",
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email of the user that owns the entity.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC.",
			},
		},
	}
}

func resourceCloudAccountVMCCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var regions []*models.RegionSpecification

	apiClient := m.(*Client).apiClient

	if v, ok := d.GetOk("regions"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.FromErr(errors.New("specified regions are not unique"))
		}
		regions = expandRegionSpecificationList(v.(*schema.Set).List())
	}

	createResp, err := apiClient.CloudAccount.CreateVmcCloudAccountAsync(
		cloud_account.NewCreateVmcCloudAccountAsyncParams().
			WithAPIVersion(IaaSAPIVersion).
			WithBody(&models.CloudAccountVmcSpecification{
				AcceptSelfSignedCertificate: d.Get("accept_self_signed_cert").(bool),
				APIKey:                      withString(d.Get("api_token").(string)),
				CreateDefaultZones:          false,
				DcID:                        withString(d.Get("dc_id").(string)),
				Description:                 d.Get("description").(string),
				HostName:                    withString(d.Get("vcenter_hostname").(string)),
				Name:                        withString(d.Get("name").(string)),
				NsxHostName:                 withString(d.Get("nsx_hostname").(string)),
				Password:                    withString(d.Get("vcenter_password").(string)),
				Regions:                     regions,
				SddcID:                      withString(d.Get("sddc_name").(string)),
				Tags:                        expandTags(d.Get("tags").(*schema.Set).List()),
				Username:                    withString(d.Get("vcenter_username").(string)),
			}))
	if err != nil {
		return diag.FromErr(err)
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    resourceCloudAccountVMCStateRefreshFunc(*apiClient, *createResp.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
	}

	resourceIDs, err := stateChangeFunc.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	cloudAccountVMC := (resourceIDs.([]string))[0]

	d.SetId(cloudAccountVMC)

	return resourceCloudAccountVMCRead(ctx, d, m)
}

func resourceCloudAccountVMCRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	ret, err := apiClient.CloudAccount.GetCloudAccount(cloud_account.NewGetCloudAccountParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *cloud_account.GetVmcCloudAccountNotFound:
			d.SetId("")
			return diag.Errorf("vmc cloud account '%s' not found", id)
		}
		return diag.FromErr(err)
	}

	vmcAccount := *ret.Payload
	d.Set("created_at", vmcAccount.CreatedAt)
	d.Set("dc_id", vmcAccount.CloudAccountProperties["dcId"])
	d.Set("description", vmcAccount.Description)
	d.Set("name", vmcAccount.Name)
	d.Set("nsx_hostname", vmcAccount.CloudAccountProperties["nsxHostName"])
	d.Set("org_id", vmcAccount.OrgID)
	d.Set("owner", vmcAccount.Owner)
	d.Set("sddc_name", vmcAccount.CloudAccountProperties["sddcId"])
	d.Set("updated_at", vmcAccount.UpdatedAt)
	d.Set("vcenter_hostname", vmcAccount.CloudAccountProperties["hostName"])
	d.Set("vcenter_username", vmcAccount.CloudAccountProperties["privateKeyId"])

	if err := d.Set("links", flattenLinks(vmcAccount.Links)); err != nil {
		return diag.Errorf("error setting cloud_account_vmc links - error: %#v", err)
	}

	if err := d.Set("regions", extractIDsFromRegion(vmcAccount.EnabledRegions)); err != nil {
		return diag.Errorf("error setting cloud_account_vmc regions - error: %#v", err)
	}

	if err := d.Set("tags", flattenTags(vmcAccount.Tags)); err != nil {
		return diag.Errorf("error setting cloud_account_vmc tags - error: %#v", err)
	}

	return nil
}

func resourceCloudAccountVMCUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var regions []*models.RegionSpecification

	apiClient := m.(*Client).apiClient

	if v, ok := d.GetOk("regions"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.FromErr(errors.New("specified regions are not unique"))
		}
		regions = expandRegionSpecificationList(v.(*schema.Set).List())
	}

	id := d.Id()
	updateResp, err := apiClient.CloudAccount.UpdateCloudAccountAsync(
		cloud_account.NewUpdateCloudAccountAsyncParams().
			WithAPIVersion(IaaSAPIVersion).
			WithID(id).
			WithBody(&models.UpdateCloudAccountSpecification{
				CreateDefaultZones: false,
				Description:        d.Get("description").(string),
				Name:               d.Get("name").(string),
				Regions:            regions,
				Tags:               expandTags(d.Get("tags").(*schema.Set).List()),
			}))
	if err != nil {
		return diag.FromErr(err)
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    resourceCloudAccountVMCStateRefreshFunc(*apiClient, *updateResp.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 5 * time.Second,
	}
	if _, err := stateChangeFunc.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	return resourceCloudAccountVMCRead(ctx, d, m)
}

func resourceCloudAccountVMCDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	if _, _, err := apiClient.CloudAccount.DeleteCloudAccount(cloud_account.NewDeleteCloudAccountParams().WithID(id)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func resourceCloudAccountVMCStateRefreshFunc(apiClient client.API, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ret, err := apiClient.Request.GetRequestTracker(request.NewGetRequestTrackerParams().WithID(id))
		if err != nil {
			return "", models.RequestTrackerStatusFAILED, err
		}

		status := ret.Payload.Status
		switch *status {
		case models.RequestTrackerStatusFAILED:
			return []string{""}, *status, errors.New(ret.Payload.Message)
		case models.RequestTrackerStatusINPROGRESS:
			return [...]string{id}, *status, nil
		case models.RequestTrackerStatusFINISHED:
			cloudAccountIDs := make([]string, len(ret.Payload.Resources))
			for i, r := range ret.Payload.Resources {
				cloudAccountIDs[i] = strings.TrimPrefix(r, "/iaas/api/cloud-accounts/")
			}
			return cloudAccountIDs, *status, nil
		default:
			return [...]string{id}, ret.Payload.Message, fmt.Errorf("esourceCloudAccountVMCStateRefreshFunc: unknown status %v", *status)
		}
	}
}
