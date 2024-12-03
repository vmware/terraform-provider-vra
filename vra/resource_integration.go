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
	"github.com/vmware/vra-sdk-go/pkg/client/integration"
	"github.com/vmware/vra-sdk-go/pkg/client/request"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIntegration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIntegrationCreate,
		ReadContext:   resourceIntegrationRead,
		UpdateContext: resourceIntegrationUpdate,
		DeleteContext: resourceIntegrationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"associated_cloud_account_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Ids of the cloud accounts to associate with this integration.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Certificate to be used to connect to the integration.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was created. The date is in ISO 8601 and UTC.",
			},
			"custom_properties": {
				Type:        schema.TypeMap,
				Computed:    true,
				Optional:    true,
				Description: "Additional custom properties that may be used to extend the Integration.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A human-friendly description.",
			},
			"integration_properties": {
				Type:        schema.TypeMap,
				Required:    true,
				Description: "Integration specific properties supplied in as name value pairs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"integration_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Integration type.",
			},
			"links": linksSchema(),
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the integration.",
			},
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
			"private_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Secret access key or password to be used to authenticate with the integration.",
			},
			"private_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Access key id or username to be used to authenticate with the integration.",
			},
			"tags": tagsSchema(),
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC.",
			},
		},
	}
}

func resourceIntegrationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	var associatedCloudAccountIDs []string
	if v, ok := d.GetOk("associated_cloud_account_ids"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.FromErr(errors.New("specified associated cloud account ids are not unique"))
		}
		associatedCloudAccountIDs = expandStringList(v.(*schema.Set).List())
	}

	customProperties := expandCustomProperties(d.Get("custom_properties").(map[string]interface{}))

	integrationProperties := make(map[string]string)
	for key, value := range d.Get("integration_properties").(map[string]interface{}) {
		integrationProperties[key] = value.(string)
	}

	createResp, err := apiClient.Integration.CreateIntegrationAsync(
		integration.NewCreateIntegrationAsyncParams().
			WithAPIVersion(IaaSAPIVersion).
			WithBody(&models.IntegrationSpecification{
				AssociatedCloudAccountIds: associatedCloudAccountIDs,
				CertificateInfo: &models.CertificateInfoSpecification{
					Certificate: withString(d.Get("certificate").(string)),
				},
				CustomProperties:      customProperties,
				Description:           d.Get("description").(string),
				IntegrationProperties: integrationProperties,
				IntegrationType:       withString(d.Get("integration_type").(string)),
				Name:                  withString(d.Get("name").(string)),
				PrivateKey:            d.Get("private_key").(string),
				PrivateKeyID:          d.Get("private_key_id").(string),
				Tags:                  expandTags(d.Get("tags").(*schema.Set).List()),
			}))
	if err != nil {
		return diag.FromErr(err)
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    resourceIntegrationStateRefreshFunc(*apiClient, *createResp.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
	}

	resourceIDs, err := stateChangeFunc.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	integration := (resourceIDs.([]string))[0]

	d.SetId(integration)

	return resourceIntegrationRead(ctx, d, m)
}

func resourceIntegrationRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	ret, err := apiClient.Integration.GetIntegration(integration.NewGetIntegrationParams().WithAPIVersion(IaaSAPIVersion).WithID(id))
	if err != nil {
		switch err.(type) {
		case *integration.GetIntegrationNotFound:
			d.SetId("")
			return diag.Errorf("integration '%s' not found", id)
		}
		return diag.FromErr(err)
	}

	integration := *ret.Payload
	d.Set("created_at", integration.CreatedAt)
	d.Set("custom_properties", integration.CustomProperties)
	d.Set("description", integration.Description)
	d.Set("integration_properties", integration.IntegrationProperties)
	d.Set("integration_type", integration.IntegrationType)
	d.Set("name", integration.Name)
	d.Set("org_id", integration.OrgID)
	d.Set("owner", integration.Owner)
	d.Set("updated_at", integration.UpdatedAt)

	if err := d.Set("links", flattenLinks(integration.Links)); err != nil {
		return diag.Errorf("error setting integration links - error: %#v", err)
	}

	if err := d.Set("tags", flattenTags(integration.Tags)); err != nil {
		return diag.Errorf("Error setting integration tags - error: %#v", err)
	}

	return nil
}

func resourceIntegrationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	var associatedCloudAccountIDs []string
	if v, ok := d.GetOk("associated_cloud_account_ids"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.FromErr(errors.New("specified associated cloud account ids are not unique"))
		}
		associatedCloudAccountIDs = expandStringList(v.(*schema.Set).List())
	}

	customProperties := expandCustomProperties(d.Get("custom_properties").(map[string]interface{}))

	integrationProperties := make(map[string]string)
	for key, value := range d.Get("integration_properties").(map[string]interface{}) {
		integrationProperties[key] = value.(string)
	}

	id := d.Id()
	updateResp, err := apiClient.Integration.UpdateIntegrationAsync(
		integration.NewUpdateIntegrationAsyncParams().
			WithAPIVersion(IaaSAPIVersion).
			WithID(id).
			WithBody(&models.UpdateIntegrationSpecification{
				AssociatedCloudAccountIds: associatedCloudAccountIDs,
				CertificateInfo: &models.CertificateInfoSpecification{
					Certificate: withString(d.Get("certificate").(string)),
				},
				CustomProperties:      customProperties,
				Description:           d.Get("description").(string),
				IntegrationProperties: integrationProperties,
				PrivateKey:            d.Get("private_key").(string),
				PrivateKeyID:          d.Get("private_key_id").(string),
				Tags:                  expandTags(d.Get("tags").(*schema.Set).List()),
			}))
	if err != nil {
		return diag.FromErr(err)
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    resourceIntegrationStateRefreshFunc(*apiClient, *updateResp.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 5 * time.Second,
	}
	if _, err := stateChangeFunc.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	return resourceIntegrationRead(ctx, d, m)
}

func resourceIntegrationDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	if _, _, err := apiClient.Integration.DeleteIntegration(integration.NewDeleteIntegrationParams().WithAPIVersion(IaaSAPIVersion).WithID(id)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func resourceIntegrationStateRefreshFunc(apiClient client.API, id string) retry.StateRefreshFunc {
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
			integrationIDs := make([]string, len(ret.Payload.Resources))
			for i, r := range ret.Payload.Resources {
				integrationIDs[i] = strings.TrimPrefix(r, "/iaas/api/integrations/")
			}
			return integrationIDs, *status, nil
		default:
			return [...]string{id}, ret.Payload.Message, fmt.Errorf("resourceIntegrationStateRefreshFunc: unknown status %v", *status)
		}
	}
}
