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

func resourceCloudAccountGCP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudAccountGCPCreate,
		ReadContext:   resourceCloudAccountGCPRead,
		UpdateContext: resourceCloudAccountGCPUpdate,
		DeleteContext: resourceCloudAccountGCPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required arguments
			"client_email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "GCP Client email.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of this resource instance.",
			},
			"private_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "GCP Private key.",
			},
			"private_key_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "GCP Private key ID.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "GCP Project ID.",
			},
			"regions": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The set of region ids that will be enabled for this cloud account.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			// Optional arguments
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

func resourceCloudAccountGCPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var regions []*models.RegionSpecification

	apiClient := m.(*Client).apiClient

	if v, ok := d.GetOk("regions"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.FromErr(errors.New("specified regions are not unique"))
		}
		regions = expandRegionSpecificationList(v.(*schema.Set).List())
	}

	createResp, err := apiClient.CloudAccount.CreateGcpCloudAccountAsync(
		cloud_account.NewCreateGcpCloudAccountAsyncParams().
			WithAPIVersion(IaaSAPIVersion).
			WithBody(&models.CloudAccountGcpSpecification{
				ClientEmail:        withString(d.Get("client_email").(string)),
				CreateDefaultZones: false,
				Description:        d.Get("description").(string),
				Name:               withString(d.Get("name").(string)),
				PrivateKey:         withString(d.Get("private_key").(string)),
				PrivateKeyID:       withString(d.Get("private_key_id").(string)),
				ProjectID:          withString(d.Get("project_id").(string)),
				Regions:            regions,
				Tags:               expandTags(d.Get("tags").(*schema.Set).List()),
			}))
	if err != nil {
		return diag.FromErr(err)
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    resourceCloudAccountGCPStateRefreshFunc(*apiClient, *createResp.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
	}

	resourceIDs, err := stateChangeFunc.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	cloudAccountGCP := (resourceIDs.([]string))[0]

	d.SetId(cloudAccountGCP)

	return resourceCloudAccountGCPRead(ctx, d, m)
}

func resourceCloudAccountGCPRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	ret, err := apiClient.CloudAccount.GetGcpCloudAccount(cloud_account.NewGetGcpCloudAccountParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *cloud_account.GetGcpCloudAccountNotFound:
			d.SetId("")
			return diag.Errorf("gcp cloud account '%s' not found", id)
		}
		return diag.FromErr(err)
	}

	gcpAccount := *ret.Payload
	d.Set("client_email", gcpAccount.ClientEmail)
	d.Set("created_at", gcpAccount.CreatedAt)
	d.Set("description", gcpAccount.Description)
	d.Set("name", gcpAccount.Name)
	d.Set("org_id", gcpAccount.OrgID)
	d.Set("owner", gcpAccount.Owner)
	d.Set("private_key_id", gcpAccount.PrivateKeyID)
	d.Set("project_id", gcpAccount.ProjectID)
	d.Set("updated_at", gcpAccount.UpdatedAt)

	if err := d.Set("links", flattenLinks(gcpAccount.Links)); err != nil {
		return diag.Errorf("error setting cloud_account_gcp links - error: %#v", err)
	}

	if err := d.Set("regions", extractIDsFromRegion(gcpAccount.EnabledRegions)); err != nil {
		return diag.Errorf("error setting cloud_account_gcp regions - error: %#v", err)
	}

	if err := d.Set("tags", flattenTags(gcpAccount.Tags)); err != nil {
		return diag.Errorf("error setting cloud_account_gcp tags - error: %#v", err)
	}

	return nil
}

func resourceCloudAccountGCPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var regions []*models.RegionSpecification

	apiClient := m.(*Client).apiClient

	if v, ok := d.GetOk("regions"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.FromErr(errors.New("specified regions are not unique"))
		}
		regions = expandRegionSpecificationList(v.(*schema.Set).List())
	}

	id := d.Id()
	updateResp, err := apiClient.CloudAccount.UpdateGcpCloudAccountAsync(
		cloud_account.NewUpdateGcpCloudAccountAsyncParams().
			WithAPIVersion(IaaSAPIVersion).
			WithID(id).
			WithBody(&models.UpdateCloudAccountGcpSpecification{
				CreateDefaultZones: false,
				Description:        d.Get("description").(string),
				Regions:            regions,
				Tags:               expandTags(d.Get("tags").(*schema.Set).List()),
			}))
	if err != nil {
		return diag.FromErr(err)
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    resourceCloudAccountGCPStateRefreshFunc(*apiClient, *updateResp.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 5 * time.Second,
	}
	if _, err := stateChangeFunc.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	return resourceCloudAccountGCPRead(ctx, d, m)
}

func resourceCloudAccountGCPDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	if _, _, err := apiClient.CloudAccount.DeleteGcpCloudAccount(cloud_account.NewDeleteGcpCloudAccountParams().WithID(id)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func resourceCloudAccountGCPStateRefreshFunc(apiClient client.API, id string) retry.StateRefreshFunc {
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
			return [...]string{id}, ret.Payload.Message, fmt.Errorf("resourceCloudAccountGCPStateRefreshFunc: unknown status %v", *status)
		}
	}
}
