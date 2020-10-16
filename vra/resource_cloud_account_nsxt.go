package vra

import (
	"fmt"

	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCloudAccountNSXT() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudAccountNSXTCreate,
		Read:   resourceCloudAccountNSXTRead,
		Update: resourceCloudAccountNSXTUpdate,
		Delete: resourceCloudAccountNSXTDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceCloudAccountNSXTCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	tags := expandTags(d.Get("tags").(*schema.Set).List())

	createResp, err := apiClient.CloudAccount.CreateNsxTCloudAccount(
		cloud_account.NewCreateNsxTCloudAccountParams().
			WithTimeout(IncreasedTimeOut).
			WithBody(&models.CloudAccountNsxTSpecification{
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
		return err
	}

	if err := d.Set("tags", flattenTags(tags)); err != nil {
		return fmt.Errorf("error setting cloud account tags - error: %#v", err)
	}
	d.SetId(*createResp.Payload.ID)

	return resourceCloudAccountNSXTRead(d, m)
}

func resourceCloudAccountNSXTRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	ret, err := apiClient.CloudAccount.GetNsxTCloudAccount(cloud_account.NewGetNsxTCloudAccountParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *cloud_account.GetNsxTCloudAccountNotFound:
			d.SetId("")
			return nil
		}
		return err
	}
	nsxtAccount := *ret.Payload
	d.Set("associated_cloud_account_ids", flattenAssociatedCloudAccountIds(nsxtAccount.Links))
	d.Set("created_at", nsxtAccount.CreatedAt)
	d.Set("dc_id", nsxtAccount.Dcid)
	d.Set("description", nsxtAccount.Description)
	d.Set("name", nsxtAccount.Name)
	d.Set("org_id", nsxtAccount.OrgID)
	d.Set("owner", nsxtAccount.Owner)
	d.Set("updated_at", nsxtAccount.UpdatedAt)
	d.Set("username", nsxtAccount.Username)

	if err := d.Set("links", flattenLinks(nsxtAccount.Links)); err != nil {
		return fmt.Errorf("error setting cloud_account_nsxt links - error: %#v", err)
	}

	if err := d.Set("tags", flattenTags(nsxtAccount.Tags)); err != nil {
		return fmt.Errorf("error setting cloud account tags - error: %#v", err)
	}

	return nil
}

func resourceCloudAccountNSXTUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	_, err := apiClient.CloudAccount.UpdateNsxTCloudAccount(cloud_account.NewUpdateNsxTCloudAccountParams().WithID(id).WithBody(&models.UpdateCloudAccountNsxTSpecification{
		Description: d.Get("description").(string),
		Tags:        expandTags(d.Get("tags").(*schema.Set).List()),
	}))
	if err != nil {
		return err
	}

	return resourceCloudAccountNSXTRead(d, m)
}

func resourceCloudAccountNSXTDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	_, err := apiClient.CloudAccount.DeleteCloudAccountNsxT(cloud_account.NewDeleteCloudAccountNsxTParams().WithID(id))
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
