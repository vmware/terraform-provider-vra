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

func dataSourceRegionEnumerationGCP() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRegionEnumerationGCPRead,

		Schema: map[string]*schema.Schema{
			"client_email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "GCP Client email.",
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
				Computed:    true,
				Description: "A set of region ids that can be enabled for this cloud account.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceRegionEnumerationGCPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*Client).apiClient

	enumResp, err := apiClient.CloudAccount.EnumerateGcpRegionsAsync(
		cloud_account.NewEnumerateGcpRegionsAsyncParams().
			WithAPIVersion(IaaSAPIVersion).
			WithTimeout(IncreasedTimeOut).
			WithBody(&models.CloudAccountGcpRegionEnumerationSpecification{
				ClientEmail:  d.Get("client_email").(string),
				PrivateKey:   d.Get("private_key").(string),
				PrivateKeyID: d.Get("private_key_id").(string),
				ProjectID:    d.Get("project_id").(string),
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
	d.SetId(d.Get("private_key_id").(string))

	return nil
}
