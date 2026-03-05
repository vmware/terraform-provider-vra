// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceCloudAccountVCF() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudAccountVCFRead,

		Schema: map[string]*schema.Schema{
			// Optional arguments
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
				Description:   "The id of this resource instance.",
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"id"},
				Description:   "The name of this resource instance.",
			},

			// Computed attributes
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was created. The date is in ISO 8601 and UTC.",
			},
			"dc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identifier of a data collector vm deployed in the on premise infrastructure.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly description.",
			},
			"enabled_regions": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "The set of regions that are enabled for this cloud account.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique identifier of the region on the provider side.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique identifier of the region.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the region on the provider side.",
						},
					},
				},
			},
			"links": linksSchema(),
			"nsx_hostname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address or FQDN of the NSX Manager Server in the specified workload domain.",
			},
			"nsx_username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Username used to authenticate to the NSX Manager in the specified workload domain.",
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
			"regions": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "The set of region ids that are enabled for this cloud account.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sddc_manager_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SDDC manager integration id.",
			},
			"tags": tagsSchema(),
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC.",
			},
			"vcenter_hostname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address or FQDN of the vCenter Server in the specified workload domain.",
			},
			"vcenter_username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Username used to authenticate to the vCenter Server in the specified workload domain.",
			},
			"workload_domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Id of the workload domain.",
			},
			"workload_domain_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the workload domain.",
			},
			"vsphere_cloud_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the vSphere cloud account associated with this VCF cloud account.",
			},
		},
	}
}

func dataSourceCloudAccountVCFRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	id := d.Get("id").(string)
	name := d.Get("name").(string)

	if id == "" && name == "" {
		return fmt.Errorf("one of 'id' or 'name' must be set")
	}

	var cloudAccountVCF *models.CloudAccountVcf
	if id != "" {
		getResp, err := apiClient.CloudAccount.GetVcfCloudAccount(cloud_account.NewGetVcfCloudAccountParams().WithID(id))
		if err != nil {
			switch err.(type) {
			case *cloud_account.GetVcfCloudAccountNotFound:
				return fmt.Errorf("vcf cloud account with id '%s' not found", id)
			default:
				// nop
			}
			return err
		}

		cloudAccountVCF = getResp.GetPayload()
	} else {
		getResp, err := apiClient.CloudAccount.GetVcfCloudAccounts(cloud_account.NewGetVcfCloudAccountsParams())
		if err != nil {
			return err
		}

		for _, account := range getResp.Payload.Content {
			if account.Name == name {
				cloudAccountVCF = account
			}
		}

		if cloudAccountVCF == nil {
			return fmt.Errorf("vcf cloud account with name '%s' not found", name)
		}
	}

	d.SetId(*cloudAccountVCF.ID)
	d.Set("created_at", cloudAccountVCF.CreatedAt)
	d.Set("dc_id", cloudAccountVCF.CustomProperties["dcId"])
	d.Set("description", cloudAccountVCF.Description)
	d.Set("name", cloudAccountVCF.Name)
	d.Set("nsx_hostname", cloudAccountVCF.NsxHostName)
	d.Set("nsx_username", cloudAccountVCF.NsxUsername)
	d.Set("org_id", cloudAccountVCF.OrgID)
	d.Set("owner", cloudAccountVCF.Owner)
	d.Set("sddc_manager_id", cloudAccountVCF.SddcManagerID)
	d.Set("updated_at", cloudAccountVCF.UpdatedAt)
	d.Set("vcenter_hostname", cloudAccountVCF.VcenterHostName)
	d.Set("vcenter_username", cloudAccountVCF.VcenterUsername)
	d.Set("workload_domain_id", cloudAccountVCF.VcfDomainID)
	d.Set("workload_domain_name", cloudAccountVCF.VcfDomainName)

	vsphereLink := cloudAccountVCF.CustomProperties["vsphere"]
	parts := strings.Split(vsphereLink, "/")
	if len(parts) > 0 {
		d.Set("vsphere_cloud_account_id", parts[len(parts)-1])
	} else {
		log.Printf("unexpected format for vsphere cloud account link: %s", vsphereLink)
		d.Set("vsphere_cloud_account_id", "")
	}

	if err := d.Set("links", flattenLinks(cloudAccountVCF.Links)); err != nil {
		return fmt.Errorf("error setting cloud_account_vcf links - error: %#v", err)
	}

	if err := d.Set("enabled_regions", flattenEnabledRegions(cloudAccountVCF.EnabledRegions)); err != nil {
		return fmt.Errorf("error setting cloud_account_vcf enabled_regions - error: %#v", err)
	}

	if err := d.Set("regions", extractIDsFromRegion(cloudAccountVCF.EnabledRegions)); err != nil {
		return fmt.Errorf("error setting cloud_account_vcf regions - error: %#v", err)
	}

	if err := d.Set("tags", flattenTags(cloudAccountVCF.Tags)); err != nil {
		return fmt.Errorf("error setting cloud_account_vcf tags - error: %#v", err)
	}

	return nil
}
