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

func resourceCloudAccountNSXV() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudAccountNSXVCreate,
		ReadContext:   resourceCloudAccountNSXVRead,
		UpdateContext: resourceCloudAccountNSXVUpdate,
		DeleteContext: resourceCloudAccountNSXVDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required arguments
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Optional arguments
			"accept_self_signed_cert": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"dc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
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
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudAccountNSXVCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	createResp, err := apiClient.CloudAccount.CreateNsxVCloudAccountAsync(
		cloud_account.NewCreateNsxVCloudAccountAsyncParams().
			WithAPIVersion(IaaSAPIVersion).
			WithBody(&models.CloudAccountNsxVSpecification{
				AcceptSelfSignedCertificate: d.Get("accept_self_signed_cert").(bool),
				Dcid:                        withString(d.Get("dc_id").(string)),
				Description:                 d.Get("description").(string),
				HostName:                    withString(d.Get("hostname").(string)),
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
		Refresh:    resourceCloudAccountNSXVStateRefreshFunc(*apiClient, *createResp.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
	}

	resourceIDs, err := stateChangeFunc.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	cloudAccountNSXV := (resourceIDs.([]string))[0]

	d.SetId(cloudAccountNSXV)

	return resourceCloudAccountNSXVRead(ctx, d, m)
}

func resourceCloudAccountNSXVRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	ret, err := apiClient.CloudAccount.GetNsxVCloudAccount(cloud_account.NewGetNsxVCloudAccountParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *cloud_account.GetNsxVCloudAccountNotFound:
			d.SetId("")
			return diag.Errorf("nsx-v cloud account '%s' not found", id)
		}
		return diag.FromErr(err)
	}
	nsxvAccount := *ret.Payload
	d.Set("associated_cloud_account_ids", flattenAssociatedCloudAccountIDs(nsxvAccount.Links))
	d.Set("created_at", nsxvAccount.CreatedAt)
	d.Set("dc_id", nsxvAccount.Dcid)
	d.Set("description", nsxvAccount.Description)
	d.Set("name", nsxvAccount.Name)
	d.Set("org_id", nsxvAccount.OrgID)
	d.Set("owner", nsxvAccount.Owner)
	d.Set("updated_at", nsxvAccount.UpdatedAt)
	d.Set("username", nsxvAccount.Username)

	if err := d.Set("links", flattenLinks(nsxvAccount.Links)); err != nil {
		return diag.Errorf("error setting cloud_account_nsxv links - error: %#v", err)
	}

	if err := d.Set("tags", flattenTags(nsxvAccount.Tags)); err != nil {
		return diag.Errorf("error setting cloud_account_nsxv tags - error: %#v", err)
	}

	return nil
}

func resourceCloudAccountNSXVUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	updateResp, err := apiClient.CloudAccount.UpdateNsxVCloudAccountAsync(
		cloud_account.NewUpdateNsxVCloudAccountAsyncParams().
			WithAPIVersion(IaaSAPIVersion).
			WithID(id).
			WithBody(&models.UpdateCloudAccountNsxVSpecification{
				AcceptSelfSignedCertificate: d.Get("accept_self_signed_cert").(bool),
				Dcid:                        withString(d.Get("dc_id").(string)),
				HostName:                    withString(d.Get("hostname").(string)),
				Name:                        d.Get("name").(string),
				Password:                    withString(d.Get("password").(string)),
				Username:                    withString(d.Get("username").(string)),
			}))
	if err != nil {
		return diag.FromErr(err)
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    resourceCloudAccountNSXVStateRefreshFunc(*apiClient, *updateResp.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 5 * time.Second,
	}
	if _, err := stateChangeFunc.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	return resourceCloudAccountNSXVRead(ctx, d, m)
}

func resourceCloudAccountNSXVDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	if _, _, err := apiClient.CloudAccount.DeleteCloudAccountNsxV(cloud_account.NewDeleteCloudAccountNsxVParams().WithID(id)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func resourceCloudAccountNSXVStateRefreshFunc(apiClient client.API, id string) retry.StateRefreshFunc {
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
			return [...]string{id}, ret.Payload.Message, fmt.Errorf("resourceCloudAccountNSXVStateRefreshFunc: unknown status %v", *status)
		}
	}
}
