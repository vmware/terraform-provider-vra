// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/vmware/vra-sdk-go/pkg/client/network_profile"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetworkProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkProfileCreate,
		ReadContext:   resourceNetworkProfileRead,
		UpdateContext: resourceNetworkProfileUpdate,
		DeleteContext: resourceNetworkProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"cloud_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the cloud account this entity belongs to.",
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"custom_properties": {
				Type:     schema.TypeMap,
				Computed: true,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fabric_network_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"isolated_network_cidr_prefix": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"isolated_external_fabric_network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"isolated_network_domain_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"isolated_network_domain_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"isolation_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": tagsSchema(),
			"external_region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": linksSchema(),
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization_id": {
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

func resourceNetworkProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to create vra_network_profile resource")
	apiClient := m.(*Client).apiClient

	name := d.Get("name").(string)
	regionID := d.Get("region_id").(string)

	networkProfileSpecification := models.NetworkProfileSpecification{
		IsolationType:                    d.Get("isolation_type").(string),
		IsolationNetworkDomainID:         d.Get("isolated_network_domain_id").(string),
		IsolationNetworkDomainCIDR:       d.Get("isolated_network_domain_cidr").(string),
		IsolationExternalFabricNetworkID: d.Get("isolated_external_fabric_network_id").(string),
		IsolatedNetworkCIDRPrefix:        int32(d.Get("isolated_network_cidr_prefix").(int)),
		Name:                             &name,
		RegionID:                         &regionID,
		Tags:                             expandTags(d.Get("tags").(*schema.Set).List()),
		CustomProperties:                 expandCustomProperties(d.Get("custom_properties").(map[string]interface{})),
	}

	if v, ok := d.GetOk("description"); ok {
		networkProfileSpecification.Description = v.(string)
	}

	if v, ok := d.GetOk("fabric_network_ids"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.FromErr(errors.New("specified fabric network ids are not unique"))
		}
		networkProfileSpecification.FabricNetworkIds = expandStringList(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.FromErr(errors.New("specified security group ids are not unique"))
		}
		networkProfileSpecification.SecurityGroupIds = expandStringList(v.(*schema.Set).List())
	}

	log.Printf("[DEBUG] create network profile: %#v", networkProfileSpecification)
	createNetworkProfileCreated, err := apiClient.NetworkProfile.CreateNetworkProfile(network_profile.NewCreateNetworkProfileParams().WithBody(&networkProfileSpecification))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*createNetworkProfileCreated.Payload.ID)
	log.Printf("Finished to create vra_network_profile resource with name %s", d.Get("name"))

	return resourceNetworkProfileRead(ctx, d, m)
}

func resourceNetworkProfileRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Reading the vra_network_profile resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	resp, err := apiClient.NetworkProfile.GetNetworkProfile(network_profile.NewGetNetworkProfileParams().WithID(id))
	if err != nil {
		return diag.FromErr(err)
	}

	networkProfile := *resp.Payload
	d.Set("cloud_account_id", networkProfile.CloudAccountID)
	d.Set("created_at", networkProfile.CreatedAt)
	d.Set("custom_properties", networkProfile.CustomProperties)
	d.Set("description", networkProfile.Description)
	d.Set("external_region_id", networkProfile.ExternalRegionID)
	d.Set("isolation_type", networkProfile.IsolationType)
	d.Set("isolated_network_domain_cidr", networkProfile.IsolationNetworkDomainCIDR)
	d.Set("isolated_network_cidr_prefix", networkProfile.IsolatedNetworkCIDRPrefix)
	d.Set("name", networkProfile.Name)
	d.Set("org_id", networkProfile.OrgID)
	d.Set("organization_id", networkProfile.OrgID)
	d.Set("owner", networkProfile.Owner)
	d.Set("updated_at", networkProfile.UpdatedAt)

	if err := d.Set("tags", flattenTags(networkProfile.Tags)); err != nil {
		return diag.Errorf("error setting network profile tags - error: %v", err)
	}

	if err := d.Set("links", flattenLinks(networkProfile.Links)); err != nil {
		return diag.Errorf("error setting network profile links - error: %#v", err)
	}

	if regionLink, ok := networkProfile.Links["region"]; ok {
		if regionLink.Href != "" {
			d.Set("region_id", strings.TrimPrefix(regionLink.Href, "/iaas/api/regions/"))
		}
	}

	if fabricNetworkLinks, ok := networkProfile.Links["fabric-networks"]; ok {
		if len(fabricNetworkLinks.Hrefs) != 0 {
			var networkIDs []string
			for i, link := range fabricNetworkLinks.Hrefs {
				networkIDs = append(networkIDs, strings.TrimPrefix(link, "/iaas/api/fabric-networks/"))
				log.Printf("Appending network profile link %s on index %d", link, i)
			}
			d.Set("fabric_network_ids", networkIDs)
		}
	}
	log.Printf("Finished reading the vra_network_profile data source with filter %s", d.Get("filter"))
	return nil
}

func resourceNetworkProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	name := d.Get("name").(string)
	regionID := d.Get("region_id").(string)

	networkProfileSpecification := models.NetworkProfileSpecification{
		IsolationType:                    d.Get("isolation_type").(string),
		IsolationNetworkDomainID:         d.Get("isolated_network_domain_id").(string),
		IsolationNetworkDomainCIDR:       d.Get("isolated_network_domain_cidr").(string),
		IsolationExternalFabricNetworkID: d.Get("isolated_external_fabric_network_id").(string),
		IsolatedNetworkCIDRPrefix:        int32(d.Get("isolated_network_cidr_prefix").(int)),
		Name:                             &name,
		RegionID:                         &regionID,
		Tags:                             expandTags(d.Get("tags").(*schema.Set).List()),
		CustomProperties:                 expandCustomProperties(d.Get("custom_properties").(map[string]interface{})),
	}

	if v, ok := d.GetOk("description"); ok {
		networkProfileSpecification.Description = v.(string)
	}

	if v, ok := d.GetOk("fabric_network_ids"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.FromErr(errors.New("specified fabric network ids are not unique"))
		}
		networkProfileSpecification.FabricNetworkIds = expandStringList(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.FromErr(errors.New("specified security group ids are not unique"))
		}
		networkProfileSpecification.SecurityGroupIds = expandStringList(v.(*schema.Set).List())
	}

	_, err := apiClient.NetworkProfile.UpdateNetworkProfile(network_profile.NewUpdateNetworkProfileParams().WithID(id).WithBody(&networkProfileSpecification))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetworkProfileRead(ctx, d, m)
}

func resourceNetworkProfileDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to delete the vra_network_profile resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	_, err := apiClient.NetworkProfile.DeleteNetworkProfile(network_profile.NewDeleteNetworkProfileParams().WithID(id))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("Finished deleting the vra_network_profile resource with name %s", d.Get("name"))
	return nil
}
