package vra

import (
	"context"

	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
				Type:     schema.TypeList,
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

	tags := expandTags(d.Get("tags").(*schema.Set).List())

	createResp, err := apiClient.CloudAccount.CreateNsxVCloudAccount(
		cloud_account.NewCreateNsxVCloudAccountParams().
			WithTimeout(IncreasedTimeOut).
			WithBody(&models.CloudAccountNsxVSpecification{
				AcceptSelfSignedCertificate: d.Get("accept_self_signed_cert").(bool),
				Dcid:                        withString(d.Get("dc_id").(string)),
				Description:                 d.Get("description").(string),
				HostName:                    withString(d.Get("hostname").(string)),
				Name:                        withString(d.Get("name").(string)),
				Password:                    withString(d.Get("password").(string)),
				Tags:                        tags,
				Username:                    withString(d.Get("username").(string)),
			}))

	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("tags", flattenTags(tags)); err != nil {
		return diag.Errorf("error setting cloud account tags - error: %#v", err)
	}
	d.SetId(*createResp.Payload.ID)

	return resourceCloudAccountNSXVRead(ctx, d, m)
}

func resourceCloudAccountNSXVRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	ret, err := apiClient.CloudAccount.GetNsxVCloudAccount(cloud_account.NewGetNsxVCloudAccountParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *cloud_account.GetNsxVCloudAccountNotFound:
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	nsxvAccount := *ret.Payload
	d.Set("associated_cloud_account_ids", flattenAssociatedCloudAccountIds(nsxvAccount.Links))
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
		return diag.Errorf("error setting cloud account tags - error: %#v", err)
	}

	return nil
}

func resourceCloudAccountNSXVUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	_, err := apiClient.CloudAccount.UpdateNsxVCloudAccount(cloud_account.NewUpdateNsxVCloudAccountParams().WithID(id).WithBody(&models.UpdateCloudAccountNsxVSpecification{
		Description: d.Get("description").(string),
		Tags:        expandTags(d.Get("tags").(*schema.Set).List()),
	}))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCloudAccountNSXVRead(ctx, d, m)
}

func resourceCloudAccountNSXVDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	_, err := apiClient.CloudAccount.DeleteCloudAccountNsxV(cloud_account.NewDeleteCloudAccountNsxVParams().WithID(id))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
