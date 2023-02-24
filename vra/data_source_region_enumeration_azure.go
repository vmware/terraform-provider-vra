package vra

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceRegionEnumerationAzure() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRegionEnumerationAzureRead,

		Schema: map[string]*schema.Schema{
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

func dataSourceRegionEnumerationAzureRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*Client).apiClient

	enumResp, err := apiClient.CloudAccount.EnumerateAzureRegionsAsync(
		cloud_account.NewEnumerateAzureRegionsAsyncParams().
			WithAPIVersion(IaaSAPIVersion).
			WithTimeout(IncreasedTimeOut).
			WithBody(&models.CloudAccountAzureRegionEnumerationSpecification{
				ClientApplicationID:        d.Get("application_id").(string),
				ClientApplicationSecretKey: d.Get("application_key").(string),
				SubscriptionID:             d.Get("subscription_id").(string),
				TenantID:                   d.Get("tenant_id").(string),
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
			WithAPIVersion(IaaSAPIVersion).
			WithTimeout(IncreasedTimeOut).
			WithID(enumID))
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("regions", extractIdsFromRegionSpecification(getResp.Payload.ExternalRegions))
	d.SetId(d.Get("application_id").(string))

	return nil
}
