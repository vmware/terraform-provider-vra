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

func resourceCloudAccountAzure() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudAccountAzureCreate,
		ReadContext:   resourceCloudAccountAzureRead,
		UpdateContext: resourceCloudAccountAzureUpdate,
		DeleteContext: resourceCloudAccountAzureDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required arguments
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Azure Client Application ID.",
			},
			"application_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Azure Client Application Secret Key.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of this resource instance.",
			},
			"regions": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The set of region ids that will be enabled for this cloud account.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Azure Subscription ID.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Azure Tenant ID.",
			},

			// Optional arguments
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A human-friendly description.",
			},
			"tags": tagsSchema(),

			//Computed attributes
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

func resourceCloudAccountAzureCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var regions []*models.RegionSpecification

	apiClient := m.(*Client).apiClient

	if v, ok := d.GetOk("regions"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.FromErr(errors.New("specified regions are not unique"))
		}
		regions = expandRegionSpecificationList(v.(*schema.Set).List())
	}

	createResp, err := apiClient.CloudAccount.CreateAzureCloudAccountAsync(
		cloud_account.NewCreateAzureCloudAccountAsyncParams().
			WithAPIVersion(IaaSAPIVersion).
			WithBody(&models.CloudAccountAzureSpecification{
				ClientApplicationID:        withString(d.Get("application_id").(string)),
				ClientApplicationSecretKey: withString(d.Get("application_key").(string)),
				CreateDefaultZones:         false,
				Description:                d.Get("description").(string),
				Name:                       withString(d.Get("name").(string)),
				Regions:                    regions,
				SubscriptionID:             withString(d.Get("subscription_id").(string)),
				Tags:                       expandTags(d.Get("tags").(*schema.Set).List()),
				TenantID:                   withString(d.Get("tenant_id").(string)),
			}))
	if err != nil {
		return diag.FromErr(err)
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    resourceCloudAccountAzureStateRefreshFunc(*apiClient, *createResp.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
	}

	resourceIDs, err := stateChangeFunc.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	cloudAccountAzure := (resourceIDs.([]string))[0]

	d.SetId(cloudAccountAzure)

	return resourceCloudAccountAzureRead(ctx, d, m)
}

func resourceCloudAccountAzureRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	ret, err := apiClient.CloudAccount.GetAzureCloudAccount(cloud_account.NewGetAzureCloudAccountParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *cloud_account.GetAzureCloudAccountNotFound:
			d.SetId("")
			return diag.Errorf("azure cloud account '%s' not found", id)
		}
		return diag.FromErr(err)
	}

	azureAccount := *ret.Payload
	d.Set("application_id", azureAccount.ClientApplicationID)
	d.Set("created_at", azureAccount.CreatedAt)
	d.Set("description", azureAccount.Description)
	d.Set("name", azureAccount.Name)
	d.Set("org_id", azureAccount.OrgID)
	d.Set("owner", azureAccount.Owner)
	d.Set("subscription_id", azureAccount.SubscriptionID)
	d.Set("tenant_id", azureAccount.TenantID)
	d.Set("updated_at", azureAccount.UpdatedAt)

	if err := d.Set("links", flattenLinks(azureAccount.Links)); err != nil {
		return diag.Errorf("error setting cloud_account_azure links - error: %#v", err)
	}

	if err := d.Set("regions", extractIDsFromRegion(azureAccount.EnabledRegions)); err != nil {
		return diag.Errorf("error setting cloud_account_azure regions - error: %#v", err)
	}

	if err := d.Set("tags", flattenTags(azureAccount.Tags)); err != nil {
		return diag.Errorf("Error setting cloud_account_azure tags - error: %#v", err)
	}

	return nil
}

func resourceCloudAccountAzureUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var regions []*models.RegionSpecification

	apiClient := m.(*Client).apiClient

	if v, ok := d.GetOk("regions"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.FromErr(errors.New("specified regions are not unique"))
		}
		regions = expandRegionSpecificationList(v.(*schema.Set).List())
	}

	id := d.Id()
	updateResp, err := apiClient.CloudAccount.UpdateAzureCloudAccountAsync(
		cloud_account.NewUpdateAzureCloudAccountAsyncParams().
			WithAPIVersion(IaaSAPIVersion).
			WithID(id).
			WithBody(&models.UpdateCloudAccountAzureSpecification{
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
		Refresh:    resourceCloudAccountAzureStateRefreshFunc(*apiClient, *updateResp.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 5 * time.Second,
	}
	if _, err := stateChangeFunc.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	return resourceCloudAccountAzureRead(ctx, d, m)
}

func resourceCloudAccountAzureDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	if _, _, err := apiClient.CloudAccount.DeleteAzureCloudAccount(cloud_account.NewDeleteAzureCloudAccountParams().WithID(id)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func resourceCloudAccountAzureStateRefreshFunc(apiClient client.API, id string) retry.StateRefreshFunc {
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
			return [...]string{id}, ret.Payload.Message, fmt.Errorf("resourceCloudAccountAzureStateRefreshFunc: unknown status %v", *status)
		}
	}
}
