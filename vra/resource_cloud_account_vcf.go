// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"
	"errors"
	"fmt"
	"log"
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

func resourceCloudAccountVCF() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudAccountVCFCreate,
		ReadContext:   resourceCloudAccountVCFRead,
		UpdateContext: resourceCloudAccountVCFUpdate,
		DeleteContext: resourceCloudAccountVCFDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required arguments
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of this resource instance.",
			},
			"nsx_hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP address or FQDN of the NSX Manager server in the specified workload domain.",
			},
			"nsx_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password used to authenticate to the NSX Manager in the specified workload domain.",
			},
			"nsx_username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username used to authenticate to the NSX Manager in the specified workload domain.",
			},
			"regions": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Use `enabled_regions` instead.",
				Description:   "The set of region ids that will be enabled for this cloud account.",
				AtLeastOneOf:  []string{"enabled_regions"},
				ConflictsWith: []string{"enabled_regions"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enabled_regions": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				Description:   "The set of regions that will be enabled for this cloud account.",
				AtLeastOneOf:  []string{"regions"},
				ConflictsWith: []string{"regions"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_region_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Unique identifier of the region on the provider side.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique identifier of the region.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the region on the provider side.",
						},
					},
				},
			},
			"vcenter_hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP address or FQDN of the vCenter Server in the specified workload domain.",
			},
			"vcenter_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password used to authenticate to the vCenter Server in the specified workload domain.",
			},
			"vcenter_username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username used to authenticate to the vCenter Server in the specified workload domain.",
			},
			"workload_domain_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of the workload domain to add as VCF cloud account.",
			},
			"workload_domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the workload domain to add as VCF cloud account.",
			},

			// Optional arguments
			"accept_self_signed_cert": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to accept self signed certificate when connecting to NSX Manager and vCenter Server.",
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
			"sddc_manager_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "SDDC manager integration id.",
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
				Type:        schema.TypeString, //
				Computed:    true,
				Description: "Email of the user that owns the entity.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC.",
			},
			"vsphere_cloud_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the vSphere cloud account associated with this VCF cloud account.",
			},
		},
	}
}

func resourceCloudAccountVCFCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var regions []*models.RegionSpecification

	apiClient := m.(*Client).apiClient

	if v, ok := d.GetOk("enabled_regions"); ok {
		regions = expandEnabledRegions(v.(*schema.Set).List())
	} else if v, ok := d.GetOk("regions"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.FromErr(errors.New("specified regions are not unique"))
		}
		regions = expandRegionSpecificationList(v.(*schema.Set).List())
	} else {
		return diag.FromErr(errors.New("one of `regions` or `enable_regions` must be specified"))
	}

	createResp, err := apiClient.CloudAccount.CreateVcfCloudAccountAsync(
		cloud_account.NewCreateVcfCloudAccountAsyncParams().
			WithAPIVersion(IaaSAPIVersion).
			WithBody(&models.CloudAccountVcfSpecification{
				AcceptSelfSignedCertificate: d.Get("accept_self_signed_cert").(bool),
				CreateDefaultZones:          false,
				DcID:                        d.Get("dc_id").(string),
				Description:                 d.Get("description").(string),
				Name:                        withString(d.Get("name").(string)),
				NsxHostName:                 withString(d.Get("nsx_hostname").(string)),
				NsxPassword:                 withString(d.Get("nsx_password").(string)),
				NsxUsername:                 withString(d.Get("nsx_username").(string)),
				Regions:                     regions,
				SddcManagerID:               d.Get("sddc_manager_id").(string),
				Tags:                        expandTags(d.Get("tags").(*schema.Set).List()),
				VcenterHostName:             withString(d.Get("vcenter_hostname").(string)),
				VcenterPassword:             withString(d.Get("vcenter_password").(string)),
				VcenterUsername:             withString(d.Get("vcenter_username").(string)),
				WorkloadDomainID:            withString(d.Get("workload_domain_id").(string)),
				WorkloadDomainName:          withString(d.Get("workload_domain_name").(string)),
			}))
	if err != nil {
		return diag.FromErr(err)
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    resourceCloudAccountVCFStateRefreshFunc(*apiClient, *createResp.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
	}

	resourceIDs, err := stateChangeFunc.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	cloudAccountVCF := (resourceIDs.([]string))[0]

	d.SetId(cloudAccountVCF)

	return resourceCloudAccountVCFRead(ctx, d, m)
}

func resourceCloudAccountVCFRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	ret, err := apiClient.CloudAccount.GetVcfCloudAccount(cloud_account.NewGetVcfCloudAccountParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *cloud_account.GetVcfCloudAccountNotFound:
			d.SetId("")
			return diag.Errorf("vcf cloud account '%s' not found", id)
		}
		return diag.FromErr(err)
	}

	vcfAccount := *ret.Payload
	d.Set("created_at", vcfAccount.CreatedAt)
	d.Set("dc_id", vcfAccount.CustomProperties["dcId"])
	d.Set("description", vcfAccount.Description)
	d.Set("name", vcfAccount.Name)
	d.Set("nsx_hostname", vcfAccount.NsxHostName)
	d.Set("nsx_username", vcfAccount.NsxUsername)
	d.Set("org_id", vcfAccount.OrgID)
	//	d.Set("owner", vcfAccount.Owner)
	d.Set("sddc_manager_id", vcfAccount.SddcManagerID)
	d.Set("updated_at", vcfAccount.UpdatedAt)
	d.Set("vcenter_hostname", vcfAccount.VcenterHostName)
	d.Set("vcenter_username", vcfAccount.VcenterUsername)
	d.Set("workload_domain_id", vcfAccount.VcfDomainID)
	d.Set("workload_domain_name", vcfAccount.VcfDomainName)
	vsphereLink := vcfAccount.CustomProperties["vsphere"]
	parts := strings.Split(vsphereLink, "/")
	if len(parts) > 0 {
		d.Set("vsphere_cloud_account_id", parts[len(parts)-1])
	} else {
		log.Printf("unexpected format for vsphere cloud account link: %s", vsphereLink)
		d.Set("vsphere_cloud_account_id", "")
	}

	if err := d.Set("links", flattenLinks(vcfAccount.Links)); err != nil {
		return diag.Errorf("error setting cloud_account_vcf links - error: %#v", err)
	}

	if err := d.Set("enabled_regions", flattenEnabledRegions(vcfAccount.EnabledRegions)); err != nil {
		return diag.Errorf("error setting cloud_account_vcf enabled_regions - error: %#v", err)
	}

	if err := d.Set("regions", extractIDsFromRegion(vcfAccount.EnabledRegions)); err != nil {
		return diag.Errorf("error setting cloud_account_vcf regions - error: %#v", err)
	}

	if err := d.Set("tags", flattenTags(vcfAccount.Tags)); err != nil {
		return diag.Errorf("error setting cloud_account_vcf tags - error: %#v", err)
	}

	return nil
}

func resourceCloudAccountVCFUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var regions []*models.RegionSpecification

	apiClient := m.(*Client).apiClient

	if v, ok := d.GetOk("enabled_regions"); ok {
		regions = expandEnabledRegions(v.(*schema.Set).List())
	} else if v, ok := d.GetOk("regions"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.FromErr(errors.New("specified regions are not unique"))
		}
		regions = expandRegionSpecificationList(v.(*schema.Set).List())
	} else {
		return diag.FromErr(errors.New("one of `regions` or `enable_regions` must be specified"))
	}

	id := d.Id()
	updateResp, err := apiClient.CloudAccount.UpdateVcfCloudAccountAsync(
		cloud_account.NewUpdateVcfCloudAccountAsyncParams().
			WithAPIVersion(IaaSAPIVersion).
			WithID(id).
			WithBody(&models.UpdateCloudAccountVcfSpecification{
				AcceptSelfSignedCertificate: d.Get("accept_self_signed_cert").(bool),
				CreateDefaultZones:          false,
				DcID:                        d.Get("dc_id").(string),
				Description:                 d.Get("description").(string),
				Name:                        d.Get("name").(string),
				NsxHostName:                 withString(d.Get("nsx_hostname").(string)),
				NsxPassword:                 withString(d.Get("nsx_password").(string)),
				NsxUsername:                 withString(d.Get("nsx_username").(string)),
				Regions:                     regions,
				SddcManagerID:               d.Get("sddc_manager_id").(string),
				Tags:                        expandTags(d.Get("tags").(*schema.Set).List()),
				VcenterHostName:             withString(d.Get("vcenter_hostname").(string)),
				VcenterPassword:             withString(d.Get("vcenter_password").(string)),
				VcenterUsername:             withString(d.Get("vcenter_username").(string)),
				WorkloadDomainID:            withString(d.Get("workload_domain_id").(string)),
				WorkloadDomainName:          withString(d.Get("workload_domain_name").(string)),
			}))
	if err != nil {
		return diag.FromErr(err)
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    resourceCloudAccountVCFStateRefreshFunc(*apiClient, *updateResp.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 5 * time.Second,
	}
	if _, err := stateChangeFunc.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	return resourceCloudAccountVCFRead(ctx, d, m)
}

func resourceCloudAccountVCFDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	if _, _, err := apiClient.CloudAccount.DeleteVcfCloudAccount(cloud_account.NewDeleteVcfCloudAccountParams().WithID(id)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func resourceCloudAccountVCFStateRefreshFunc(apiClient client.API, id string) retry.StateRefreshFunc {
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
				cloudAccountIDs[i] = strings.TrimPrefix(r, "/iaas/api/cloud-accounts-vcf/")
				cloudAccountIDs[i] = strings.TrimPrefix(cloudAccountIDs[i], "/iaas/api/cloud-accounts/")
			}
			return cloudAccountIDs, *status, nil
		default:
			return [...]string{id}, ret.Payload.Message, fmt.Errorf("resourceCloudAccountVCFStateRefreshFunc: unknown status %v", *status)
		}
	}
}
