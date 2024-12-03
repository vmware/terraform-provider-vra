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

func resourceCloudAccountAWS() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudAccountAWSCreate,
		ReadContext:   resourceCloudAccountAWSRead,
		UpdateContext: resourceCloudAccountAWSUpdate,
		DeleteContext: resourceCloudAccountAWSDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required arguments
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Aws Access key ID.",
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
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Aws Secret Access Key.",
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

func resourceCloudAccountAWSCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var regions []*models.RegionSpecification

	apiClient := m.(*Client).apiClient

	if v, ok := d.GetOk("regions"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.FromErr(errors.New("specified regions are not unique"))
		}
		regions = expandRegionSpecificationList(v.(*schema.Set).List())
	}

	createResp, err := apiClient.CloudAccount.CreateAwsCloudAccountAsync(
		cloud_account.NewCreateAwsCloudAccountAsyncParams().
			WithAPIVersion(IaaSAPIVersion).
			WithBody(&models.CloudAccountAwsSpecification{
				AccessKeyID:        withString(d.Get("access_key").(string)),
				CreateDefaultZones: false,
				Description:        d.Get("description").(string),
				Name:               withString(d.Get("name").(string)),
				Regions:            regions,
				SecretAccessKey:    withString(d.Get("secret_key").(string)),
				Tags:               expandTags(d.Get("tags").(*schema.Set).List()),
			}))
	if err != nil {
		return diag.FromErr(err)
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    resourceCloudAccountAWSStateRefreshFunc(*apiClient, *createResp.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
	}

	resourceIDs, err := stateChangeFunc.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	cloudAccountAWS := (resourceIDs.([]string))[0]

	d.SetId(cloudAccountAWS)

	return resourceCloudAccountAWSRead(ctx, d, m)
}

func resourceCloudAccountAWSRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	ret, err := apiClient.CloudAccount.GetAwsCloudAccount(cloud_account.NewGetAwsCloudAccountParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *cloud_account.GetAwsCloudAccountNotFound:
			d.SetId("")
			return diag.Errorf("aws cloud account '%s' not found", id)
		}
		return diag.FromErr(err)
	}

	awsAccount := *ret.Payload
	d.Set("access_key", awsAccount.AccessKeyID)
	d.Set("created_at", awsAccount.CreatedAt)
	d.Set("description", awsAccount.Description)
	d.Set("name", awsAccount.Name)
	d.Set("org_id", awsAccount.OrgID)
	d.Set("owner", awsAccount.Owner)
	d.Set("updated_at", awsAccount.UpdatedAt)

	if err := d.Set("links", flattenLinks(awsAccount.Links)); err != nil {
		return diag.Errorf("error setting cloud_account_aws links - error: %#v", err)
	}

	if err := d.Set("regions", extractIDsFromRegion(awsAccount.EnabledRegions)); err != nil {
		return diag.Errorf("error setting cloud_account_aws regions - error: %#v", err)
	}

	if err := d.Set("tags", flattenTags(awsAccount.Tags)); err != nil {
		return diag.Errorf("Error setting cloud_account_aws tags - error: %#v", err)
	}

	return nil
}

func resourceCloudAccountAWSUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var regions []*models.RegionSpecification

	apiClient := m.(*Client).apiClient

	if v, ok := d.GetOk("regions"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.FromErr(errors.New("specified regions are not unique"))
		}
		regions = expandRegionSpecificationList(v.(*schema.Set).List())
	}

	id := d.Id()
	updateResp, err := apiClient.CloudAccount.UpdateAWSCloudAccountAsync(
		cloud_account.NewUpdateAWSCloudAccountAsyncParams().
			WithAPIVersion(IaaSAPIVersion).
			WithID(id).
			WithBody(&models.UpdateCloudAccountAwsSpecification{
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
		Refresh:    resourceCloudAccountAWSStateRefreshFunc(*apiClient, *updateResp.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 5 * time.Second,
	}

	if _, err := stateChangeFunc.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	return resourceCloudAccountAWSRead(ctx, d, m)
}

func resourceCloudAccountAWSDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	if _, _, err := apiClient.CloudAccount.DeleteAwsCloudAccount(cloud_account.NewDeleteAwsCloudAccountParams().WithID(id)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func resourceCloudAccountAWSStateRefreshFunc(apiClient client.API, id string) retry.StateRefreshFunc {
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
			return [...]string{id}, ret.Payload.Message, fmt.Errorf("resourceCloudAccountAWSStateRefreshFunc: unknown status %v", *status)
		}
	}
}
