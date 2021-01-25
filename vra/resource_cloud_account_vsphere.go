package vra

import (
	"fmt"

	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCloudAccountVsphere() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudAccountVsphereCreate,
		Read:   resourceCloudAccountVsphereRead,
		Update: resourceCloudAccountVsphereUpdate,
		Delete: resourceCloudAccountVsphereDelete,
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
			"regions": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
			"associated_cloud_account_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"dcid": {
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
			"custom_properties": {
				Type:     schema.TypeMap,
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

func resourceCloudAccountVsphereCreate(d *schema.ResourceData, m interface{}) error {
	var regions, associatedCloudAccountIds []string

	apiClient := m.(*Client).apiClient

	tags := expandTags(d.Get("tags").(*schema.Set).List())
	if v, ok := d.GetOk("regions"); ok {
		if !compareUnique(v.([]interface{})) {
			return fmt.Errorf("Specified regions are not unique")
		}
		regions = expandStringList(v.([]interface{}))
	}

	if v, ok := d.GetOk("associated_cloud_account_ids"); ok {
		if !compareUnique(v.([]interface{})) {
			return fmt.Errorf("specified associated cloud account ids are not unique")
		}
		associatedCloudAccountIds = expandStringList(v.([]interface{}))
	}

	createResp, err := apiClient.CloudAccount.CreateVSphereCloudAccount(
		cloud_account.NewCreateVSphereCloudAccountParams().
			WithTimeout(IncreasedTimeOut).
			WithBody(&models.CloudAccountVsphereSpecification{
				AcceptSelfSignedCertificate: d.Get("accept_self_signed_cert").(bool),
				AssociatedCloudAccountIds:   associatedCloudAccountIds,
				CreateDefaultZones:          false,
				Dcid:                        d.Get("dcid").(string),
				Description:                 d.Get("description").(string),
				HostName:                    withString(d.Get("hostname").(string)),
				Name:                        withString(d.Get("name").(string)),
				Password:                    withString(d.Get("password").(string)),
				RegionIds:                   regions,
				Tags:                        tags,
				Username:                    withString(d.Get("username").(string)),
			}))

	if err != nil {
		return err
	}

	// The returned EnabledRegionIds and Hrefs containing the region ids can be in a different order than the request order.
	// Call a routine to normalize the order to correspond with the users region order.
	regionsIds, err := flattenAndNormalizeCloudAccountVsphereRegionIds(regions, createResp.Payload)
	if err != nil {
		return err
	}
	d.Set("region_ids", regionsIds)

	if err := d.Set("tags", flattenTags(tags)); err != nil {
		return fmt.Errorf("Error setting cloud account tags - error: %#v", err)
	}
	d.SetId(*createResp.Payload.ID)

	return resourceCloudAccountVsphereRead(d, m)
}

func resourceCloudAccountVsphereRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	ret, err := apiClient.CloudAccount.GetVSphereCloudAccount(cloud_account.NewGetVSphereCloudAccountParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *cloud_account.GetVSphereCloudAccountNotFound:
			d.SetId("")
			return nil
		}
		return err
	}
	vsphereAccount := *ret.Payload
	regions := vsphereAccount.EnabledRegionIds

	d.Set("associated_cloud_account_ids", flattenAssociatedCloudAccountIds(vsphereAccount.Links))
	d.Set("created_at", vsphereAccount.CreatedAt)
	d.Set("custom_properties", vsphereAccount.CustomProperties)
	d.Set("dcid", vsphereAccount.Dcid)
	d.Set("description", vsphereAccount.Description)
	d.Set("hostname", vsphereAccount.HostName)
	d.Set("name", vsphereAccount.Name)
	d.Set("org_id", vsphereAccount.OrgID)
	d.Set("owner", vsphereAccount.Owner)
	d.Set("regions", regions)
	d.Set("updated_at", vsphereAccount.UpdatedAt)
	d.Set("username", vsphereAccount.Username)

	if err := d.Set("links", flattenLinks(vsphereAccount.Links)); err != nil {
		return fmt.Errorf("error setting cloud_account_vsphere links - error: %#v", err)
	}

	// The returned EnabledRegionIds and Hrefs containing the region ids can be in a different order than the request order.
	// Call a routine to normalize the order to correspond with the users region order.
	regionsIds, err := flattenAndNormalizeCloudAccountVsphereRegionIds(regions, &vsphereAccount)
	if err != nil {
		return err
	}
	d.Set("region_ids", regionsIds)

	if err := d.Set("tags", flattenTags(vsphereAccount.Tags)); err != nil {
		return fmt.Errorf("Error setting cloud account tags - error: %#v", err)
	}

	return nil
}

func resourceCloudAccountVsphereUpdate(d *schema.ResourceData, m interface{}) error {
	var regions []string

	apiClient := m.(*Client).apiClient

	id := d.Id()

	if v, ok := d.GetOk("regions"); ok {
		if !compareUnique(v.([]interface{})) {
			return fmt.Errorf("Specified regions are not unique")
		}
		regions = expandStringList(v.([]interface{}))
	}
	_, err := apiClient.CloudAccount.UpdateVSphereCloudAccount(cloud_account.NewUpdateVSphereCloudAccountParams().WithID(id).WithBody(&models.UpdateCloudAccountVsphereSpecification{
		CreateDefaultZones: false,
		Description:        d.Get("description").(string),
		RegionIds:          regions,
		Tags:               expandTags(d.Get("tags").(*schema.Set).List()),
	}))
	if err != nil {
		return err
	}

	return resourceCloudAccountVsphereRead(d, m)
}

func resourceCloudAccountVsphereDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	_, err := apiClient.CloudAccount.DeleteVSphereCloudAccount(cloud_account.NewDeleteVSphereCloudAccountParams().WithID(id))
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
