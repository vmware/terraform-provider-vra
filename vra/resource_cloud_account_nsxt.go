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

func resourceCloudAccountNSXT() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudAccountNSXTCreate,
		ReadContext:   resourceCloudAccountNSXTRead,
		UpdateContext: resourceCloudAccountNSXTUpdate,
		DeleteContext: resourceCloudAccountNSXTDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required arguments
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Host name for the NSX-T endpoint.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A human-friendly name used as an identifier in APIs that support this option.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password for the user used to authenticate with the cloud Account.",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username to authenticate with the cloud account.",
			},
			// Optional arguments
			"accept_self_signed_cert": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Accept self signed certificate when connecting.",
			},
			"dc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identifier of a data collector vm deployed in the on premise infrastructure. Refer to the data-collector API to create or list data collectors.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A human-friendly description.",
			},
			"manager_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Create NSX-T cloud account in Manager (legacy) mode. When set to true, NSX-T cloud account is created in Manager mode. Mode cannot be changed after cloud account is created. Default value is false.",
			},
			"tags": tagsSchema(),
			// Computed attributes
			"associated_cloud_account_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
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

func resourceCloudAccountNSXTCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	createResp, err := apiClient.CloudAccount.CreateNsxTCloudAccountAsync(
		cloud_account.NewCreateNsxTCloudAccountAsyncParams().
			WithAPIVersion(IaaSAPIVersion).
			WithBody(&models.CloudAccountNsxTSpecification{
				AcceptSelfSignedCertificate: d.Get("accept_self_signed_cert").(bool),
				Dcid:                        withString(d.Get("dc_id").(string)),
				Description:                 d.Get("description").(string),
				HostName:                    withString(d.Get("hostname").(string)),
				ManagerMode:                 d.Get("manager_mode").(bool),
				Name:                        withString(d.Get("name").(string)),
				Password:                    withString(d.Get("password").(string)),
				Tags:                        expandTags(d.Get("tags").(*schema.Set).List()),
				Username:                    withString(d.Get("username").(string)),
			}))
	if err != nil {
		return diag.FromErr(err)
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    resourceCloudAccountNSXTStateRefreshFunc(*apiClient, *createResp.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
	}

	resourceIDs, err := stateChangeFunc.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	cloudAccountNSXT := (resourceIDs.([]string))[0]

	d.SetId(cloudAccountNSXT)

	return resourceCloudAccountNSXTRead(ctx, d, m)
}

func resourceCloudAccountNSXTRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	ret, err := apiClient.CloudAccount.GetNsxTCloudAccount(cloud_account.NewGetNsxTCloudAccountParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *cloud_account.GetNsxTCloudAccountNotFound:
			d.SetId("")
			return diag.Errorf("nsx-t cloud account '%s' not found", id)
		}
		return diag.FromErr(err)
	}
	nsxtAccount := *ret.Payload
	d.Set("associated_cloud_account_ids", flattenAssociatedCloudAccountIDs(nsxtAccount.Links))
	d.Set("created_at", nsxtAccount.CreatedAt)
	d.Set("dc_id", nsxtAccount.Dcid)
	d.Set("description", nsxtAccount.Description)
	d.Set("manager_mode", nsxtAccount.ManagerMode)
	d.Set("name", nsxtAccount.Name)
	d.Set("org_id", nsxtAccount.OrgID)
	d.Set("owner", nsxtAccount.Owner)
	d.Set("updated_at", nsxtAccount.UpdatedAt)
	d.Set("username", nsxtAccount.Username)

	if err := d.Set("links", flattenLinks(nsxtAccount.Links)); err != nil {
		return diag.Errorf("error setting cloud_account_nsxt links - error: %#v", err)
	}

	if err := d.Set("tags", flattenTags(nsxtAccount.Tags)); err != nil {
		return diag.Errorf("error setting cloud_account_nsxt tags - error: %#v", err)
	}

	return nil
}

func resourceCloudAccountNSXTUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	updateResp, err := apiClient.CloudAccount.UpdateNsxTCloudAccountAsync(
		cloud_account.NewUpdateNsxTCloudAccountAsyncParams().
			WithAPIVersion(IaaSAPIVersion).
			WithID(id).
			WithBody(&models.UpdateCloudAccountNsxTSpecification{
				AcceptSelfSignedCertificate: d.Get("accept_self_signed_cert").(bool),
				Dcid:                        withString(d.Get("dc_id").(string)),
				Description:                 d.Get("description").(string),
				HostName:                    withString(d.Get("hostname").(string)),
				Name:                        d.Get("name").(string),
				Password:                    withString(d.Get("password").(string)),
				Tags:                        expandTags(d.Get("tags").(*schema.Set).List()),
				Username:                    withString(d.Get("username").(string)),
			}))
	if err != nil {
		return diag.FromErr(err)
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    resourceCloudAccountNSXTStateRefreshFunc(*apiClient, *updateResp.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 5 * time.Second,
	}
	if _, err := stateChangeFunc.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	return resourceCloudAccountNSXTRead(ctx, d, m)
}

func resourceCloudAccountNSXTDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	if _, _, err := apiClient.CloudAccount.DeleteCloudAccountNsxT(cloud_account.NewDeleteCloudAccountNsxTParams().WithID(id)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func resourceCloudAccountNSXTStateRefreshFunc(apiClient client.API, id string) retry.StateRefreshFunc {
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
			return [...]string{id}, ret.Payload.Message, fmt.Errorf("resourceCloudAccountNSXTStateRefreshFunc: unknown status %v", *status)
		}
	}
}
