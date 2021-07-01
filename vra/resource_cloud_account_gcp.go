package vra

import (
	"context"
	"errors"

	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"private_key_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Optional arguments
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"regions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

func resourceCloudAccountGCPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var regions []string

	apiClient := m.(*Client).apiClient

	if v, ok := d.GetOk("regions"); ok {
		if !compareUnique(v.([]interface{})) {
			return diag.FromErr(errors.New("specified regions are not unique"))
		}
		regions = expandStringList(v.([]interface{}))
	}

	createResp, err := apiClient.CloudAccount.CreateGcpCloudAccount(cloud_account.NewCreateGcpCloudAccountParams().WithBody(&models.CloudAccountGcpSpecification{
		Description:        d.Get("description").(string),
		Name:               withString(d.Get("name").(string)),
		ClientEmail:        withString(d.Get("client_email").(string)),
		PrivateKey:         withString(d.Get("private_key").(string)),
		PrivateKeyID:       withString(d.Get("private_key_id").(string)),
		ProjectID:          withString(d.Get("project_id").(string)),
		CreateDefaultZones: false,
		RegionIds:          regions,
		Tags:               expandTags(d.Get("tags").(*schema.Set).List()),
	}))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*createResp.Payload.ID)

	return resourceCloudAccountGCPRead(ctx, d, m)
}

func resourceCloudAccountGCPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	ret, err := apiClient.CloudAccount.GetGcpCloudAccount(cloud_account.NewGetGcpCloudAccountParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *cloud_account.GetGcpCloudAccountNotFound:
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	gcpAccount := *ret.Payload
	regions := gcpAccount.EnabledRegionIds

	d.Set("client_email", gcpAccount.ClientEmail)
	d.Set("created_at", gcpAccount.CreatedAt)
	d.Set("description", gcpAccount.Description)
	d.Set("name", gcpAccount.Name)
	d.Set("org_id", gcpAccount.OrgID)
	d.Set("owner", gcpAccount.Owner)
	d.Set("private_key_id", gcpAccount.PrivateKeyID)
	d.Set("project_id", gcpAccount.ProjectID)
	d.Set("regions", regions)
	d.Set("updated_at", gcpAccount.UpdatedAt)

	if err := d.Set("links", flattenLinks(gcpAccount.Links)); err != nil {
		return diag.Errorf("error setting cloud_account_gcp links - error: %#v", err)
	}

	// The returned EnabledRegionIds and Hrefs containing the region ids can be in a different order than the request order.
	// Call a routine to normalize the order to correspond with the users region order.
	regionsIds, err := flattenAndNormalizeCLoudAccountGcpRegionIds(regions, &gcpAccount)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("region_ids", regionsIds)

	if err := d.Set("tags", flattenTags(gcpAccount.Tags)); err != nil {
		return diag.Errorf("error setting cloud account tags - error: %#v", err)
	}

	return nil
}

func resourceCloudAccountGCPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var regions []string

	apiClient := m.(*Client).apiClient

	id := d.Id()

	if v, ok := d.GetOk("regions"); ok {
		if !compareUnique(v.([]interface{})) {
			return diag.FromErr(errors.New("specified regions are not unique"))
		}
		regions = expandStringList(v.([]interface{}))
	}
	tags := expandTags(d.Get("tags").(*schema.Set).List())

	_, err := apiClient.CloudAccount.UpdateGcpCloudAccount(cloud_account.NewUpdateGcpCloudAccountParams().WithID(id).WithBody(&models.UpdateCloudAccountGcpSpecification{
		Description:        d.Get("description").(string),
		CreateDefaultZones: false,
		RegionIds:          regions,
		Tags:               tags,
	}))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCloudAccountGCPRead(ctx, d, m)
}

func resourceCloudAccountGCPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	_, err := apiClient.CloudAccount.DeleteGcpCloudAccount(cloud_account.NewDeleteGcpCloudAccountParams().WithID(id))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
