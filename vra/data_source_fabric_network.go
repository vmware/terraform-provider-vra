// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/fabric_network"
)

func dataSourceFabricNetwork() *schema.Resource {
	return &schema.Resource{
		Read: resourceFabricNetworkRead,

		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cloud_account_ids": {
				Type:     schema.TypeSet,
				Computed: true,
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
			"custom_properties": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func resourceFabricNetworkRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading the vra_fabric_network data source with name %s", d.Get("name"))
	apiClient := meta.(*Client).apiClient

	filter := d.Get("filter").(string)

	getResp, err := apiClient.FabricNetwork.GetFabricNetworks(fabric_network.NewGetFabricNetworksParams().WithDollarFilter(withString(filter)))
	if err != nil {
		return err
	}

	fabricNetworks := getResp.Payload
	if len(fabricNetworks.Content) > 1 {
		return fmt.Errorf("vra_fabric_network must filter to a fabric network")
	}
	if len(fabricNetworks.Content) == 0 {
		return fmt.Errorf("vra_fabric_network filter did not match any fabric network")
	}

	fabricNetwork := fabricNetworks.Content[0]
	d.SetId(*fabricNetwork.ID)
	_ = d.Set("cidr", fabricNetwork.Cidr)
	_ = d.Set("cloud_account_ids", fabricNetwork.CloudAccountIds)
	_ = d.Set("created_at", fabricNetwork.CreatedAt)
	_ = d.Set("description", fabricNetwork.Description)
	_ = d.Set("external_id", fabricNetwork.ExternalID)
	_ = d.Set("external_region_id", fabricNetwork.ExternalRegionID)
	_ = d.Set("is_default", fabricNetwork.IsDefault)
	_ = d.Set("is_public", fabricNetwork.IsPublic)
	_ = d.Set("name", fabricNetwork.Name)
	_ = d.Set("organization_id", fabricNetwork.OrgID)
	_ = d.Set("owner", fabricNetwork.Owner)
	_ = d.Set("updated_at", fabricNetwork.UpdatedAt)
	_ = d.Set("custom_properties", fabricNetwork.CustomProperties)

	if err := d.Set("tags", flattenTags(fabricNetwork.Tags)); err != nil {
		return fmt.Errorf("error getting network profile tags - error: %v", err)
	}

	if err := d.Set("links", flattenLinks(fabricNetwork.Links)); err != nil {
		return fmt.Errorf("error getting network profile links - error: %#v", err)
	}

	log.Printf("Finished reading the vra_fabric_network data source with name %s", d.Get("name"))
	return nil
}
