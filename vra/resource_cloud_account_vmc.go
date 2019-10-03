package vra

import (
	"fmt"
	"strconv"

	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCloudAccountVMC() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudAccountVMCCreate,
		Read:   resourceCloudAccountVMCRead,
		Update: resourceCloudAccountVMCUpdate,
		Delete: resourceCloudAccountVMCDelete,

		Schema: map[string]*schema.Schema{
			"accept_self_signed_cert": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"api_token": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"dc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"nsx_hostname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"regions": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"region_ids": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sddc_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": tagsSchema(),
			"vcenter_hostname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vcenter_password": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"vcenter_username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceCloudAccountVMCCreate(d *schema.ResourceData, m interface{}) error {
	var regions []string

	apiClient := m.(*Client).apiClient

	tags := expandTags(d.Get("tags").(*schema.Set).List())
	if v, ok := d.GetOk("regions"); ok {
		if !compareUnique(v.([]interface{})) {
			return fmt.Errorf("specified regions are not unique")
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

	createResp, err := apiClient.CloudAccount.CreateCloudAccount(cloud_account.NewCreateCloudAccountParams().WithBody(&models.CloudAccountSpecification{
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
		return err
	}

	// The returned EnabledRegionIds and Hrefs containing the region ids can be in a different order than the request order.
	// Call a routine to normalize the order to correspond with the users region order.
	regionsIds, err := flattenAndNormalizeCloudAccountRegionIds(regions, createResp.Payload)
	if err != nil {
		return err
	}
	d.Set("region_ids", regionsIds)

	if err := d.Set("tags", flattenTags(tags)); err != nil {
		return fmt.Errorf("error setting cloud account tags - error: %#v", err)
	}
	d.SetId(*createResp.Payload.ID)

	return resourceCloudAccountVMCRead(d, m)
}

func resourceCloudAccountVMCRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	ret, err := apiClient.CloudAccount.GetCloudAccount(cloud_account.NewGetCloudAccountParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *cloud_account.GetCloudAccountNotFound:
			d.SetId("")
			return nil
		}
		return err
	}
	vmcAccount := *ret.Payload

	d.Set("dc_id", vmcAccount.CloudAccountProperties["dcId"])
	d.Set("description", vmcAccount.Description)
	d.Set("name", vmcAccount.Name)
	d.Set("nsx_hostname", vmcAccount.CloudAccountProperties["nsxHostName"])
	d.Set("sddc_name", vmcAccount.CloudAccountProperties["sddcId"])
	d.Set("vcenter_hostname", vmcAccount.CloudAccountProperties["hostName"])
	d.Set("vcenter_username", vmcAccount.CloudAccountProperties["privateKeyId"])

	regions := vmcAccount.EnabledRegionIds
	d.Set("regions", regions)

	// The returned EnabledRegionIds and Hrefs containing the region ids can be in a different order than the request order.
	// Call a routine to normalize the order to correspond with the users region order.
	regionsIds, err := flattenAndNormalizeCloudAccountRegionIds(regions, &vmcAccount)
	if err != nil {
		return err
	}
	d.Set("region_ids", regionsIds)

	if err := d.Set("tags", flattenTags(vmcAccount.Tags)); err != nil {
		return fmt.Errorf("error setting cloud account tags - error: %#v", err)
	}

	return nil
}

func resourceCloudAccountVMCUpdate(d *schema.ResourceData, m interface{}) error {
	var regions []string

	apiClient := m.(*Client).apiClient

	id := d.Id()

	if v, ok := d.GetOk("regions"); ok {
		if !compareUnique(v.([]interface{})) {
			return fmt.Errorf("specified regions are not unique")
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
		return err
	}

	return resourceCloudAccountVMCRead(d, m)
}

func resourceCloudAccountVMCDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	_, err := apiClient.CloudAccount.DeleteCloudAccount(cloud_account.NewDeleteCloudAccountParams().WithID(id))
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
