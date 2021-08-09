package vra

import (
	"context"
	"errors"
	"strconv"

	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
				Type:     schema.TypeString,
				Required: true,
			},
			"nsx_hostname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"regions": {
				Type:     schema.TypeList,
				Required: true,
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
			"region_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudAccountVMCCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var regions []string

	apiClient := m.(*Client).apiClient

	tags := expandTags(d.Get("tags").(*schema.Set).List())
	if v, ok := d.GetOk("regions"); ok {
		if !compareUnique(v.([]interface{})) {
			return diag.FromErr(errors.New("specified regions are not unique"))
		}
		regions = expandStringList(v.([]interface{}))
	}

	cloudAccountProperties := make(map[string]string)
	cloudAccountProperties["acceptSelfSignedCertificate"] = strconv.FormatBool(d.Get("accept_self_signed_cert").(bool))
	cloudAccountProperties["apiKey"] = d.Get("api_token").(string)
	cloudAccountProperties["dcId"] = d.Get("dc_id").(string)
	cloudAccountProperties["hostName"] = d.Get("vcenter_hostname").(string)
	cloudAccountProperties["nsxHostName"] = d.Get("nsx_hostname").(string)
	cloudAccountProperties["sddcId"] = d.Get("sddc_name").(string)

	createResp, err := apiClient.CloudAccount.CreateCloudAccount(
		cloud_account.NewCreateCloudAccountParams().
			WithTimeout(IncreasedTimeOut).
			WithBody(&models.CloudAccountSpecification{
				AssociatedCloudAccountIds: []string{},
				CloudAccountProperties:    cloudAccountProperties,
				CloudAccountType:          withString("vmc"),
				CreateDefaultZones:        false,
				Description:               d.Get("description").(string),
				Name:                      withString(d.Get("name").(string)),
				PrivateKey:                withString(d.Get("vcenter_password").(string)),
				PrivateKeyID:              withString(d.Get("vcenter_username").(string)),
				RegionIds:                 regions,
				Tags:                      tags,
			}))

	if err != nil {
		return diag.FromErr(err)
	}

	// The returned EnabledRegionIds and Hrefs containing the region ids can be in a different order than the request order.
	// Call a routine to normalize the order to correspond with the users region order.
	regionsIds, err := flattenAndNormalizeCloudAccountRegionIds(regions, createResp.Payload)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("region_ids", regionsIds)

	if err := d.Set("tags", flattenTags(tags)); err != nil {
		return diag.Errorf("error setting cloud account tags - error: %#v", err)
	}
	d.SetId(*createResp.Payload.ID)

	return resourceCloudAccountVMCRead(ctx, d, m)
}

func resourceCloudAccountVMCRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	ret, err := apiClient.CloudAccount.GetCloudAccount(cloud_account.NewGetCloudAccountParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *cloud_account.GetCloudAccountNotFound:
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	vmcAccount := *ret.Payload
	regions := vmcAccount.EnabledRegionIds

	d.Set("created_at", vmcAccount.CreatedAt)
	d.Set("dc_id", vmcAccount.CloudAccountProperties["dcId"])
	d.Set("description", vmcAccount.Description)
	d.Set("name", vmcAccount.Name)
	d.Set("nsx_hostname", vmcAccount.CloudAccountProperties["nsxHostName"])
	d.Set("org_id", vmcAccount.OrgID)
	d.Set("owner", vmcAccount.Owner)
	d.Set("regions", regions)
	d.Set("sddc_name", vmcAccount.CloudAccountProperties["sddcId"])
	d.Set("updated_at", vmcAccount.UpdatedAt)
	d.Set("vcenter_hostname", vmcAccount.CloudAccountProperties["hostName"])
	d.Set("vcenter_username", vmcAccount.CloudAccountProperties["privateKeyId"])

	if err := d.Set("links", flattenLinks(vmcAccount.Links)); err != nil {
		return diag.Errorf("error setting cloud_account_vmc links - error: %#v", err)
	}

	// The returned EnabledRegionIds and Hrefs containing the region ids can be in a different order than the request order.
	// Call a routine to normalize the order to correspond with the users region order.
	regionsIds, err := flattenAndNormalizeCloudAccountRegionIds(regions, &vmcAccount)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("region_ids", regionsIds)

	if err := d.Set("tags", flattenTags(vmcAccount.Tags)); err != nil {
		return diag.Errorf("error setting cloud account tags - error: %#v", err)
	}

	return nil
}

func resourceCloudAccountVMCUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var regions []string

	apiClient := m.(*Client).apiClient

	id := d.Id()

	if v, ok := d.GetOk("regions"); ok {
		if !compareUnique(v.([]interface{})) {
			return diag.FromErr(errors.New("specified regions are not unique"))
		}
		regions = expandStringList(v.([]interface{}))
	}

	_, err := apiClient.CloudAccount.UpdateCloudAccount(cloud_account.NewUpdateCloudAccountParams().WithID(id).WithBody(&models.UpdateCloudAccountSpecification{
		CreateDefaultZones: false,
		Description:        d.Get("description").(string),
		RegionIds:          regions,
		Tags:               expandTags(d.Get("tags").(*schema.Set).List()),
	}))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCloudAccountVMCRead(ctx, d, m)
}

func resourceCloudAccountVMCDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	_, err := apiClient.CloudAccount.DeleteCloudAccount(cloud_account.NewDeleteCloudAccountParams().WithID(id))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
