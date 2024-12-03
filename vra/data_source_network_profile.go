// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/network_profile"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"log"
	"strings"
)

func dataSourceNetworkProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkProfileRead,

		Schema: map[string]*schema.Schema{
			"custom_properties": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fabric_network_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"isolated_network_cidr_prefix": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"isolated_external_fabric_network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"isolated_network_domain_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"isolated_network_domain_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"isolation_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": linksSchema(),
			"name": {
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
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": tagsSchema(),
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNetworkProfileRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading the vra_network_profile data source with filter %s", d.Get("filter"))
	apiClient := meta.(*Client).apiClient

	var networkProfile *models.NetworkProfile

	id := d.Get("id").(string)
	filter := d.Get("filter").(string)

	if id == "" && filter == "" {
		return fmt.Errorf("one of id or filter is required")
	}

	if id != "" {
		getResp, err := apiClient.NetworkProfile.GetNetworkProfile(network_profile.NewGetNetworkProfileParams().WithID(id))

		if err != nil {
			return err
		}
		networkProfile = getResp.GetPayload()
	} else {
		getResp, err := apiClient.NetworkProfile.GetNetworkProfiles(network_profile.NewGetNetworkProfilesParams().WithDollarFilter(withString(filter)))
		if err != nil {
			return err
		}

		networkProfiles := *getResp.Payload
		if len(networkProfiles.Content) > 1 {
			return fmt.Errorf("vra_network_profile must filter to a network profile")
		}
		if len(networkProfiles.Content) == 0 {
			return fmt.Errorf("vra_network_profile filter did not match any network profile")
		}

		networkProfile = networkProfiles.Content[0]
	}

	d.SetId(*networkProfile.ID)
	d.Set("custom_properties", networkProfile.CustomProperties)
	d.Set("description", networkProfile.Description)
	d.Set("external_region_id", networkProfile.ExternalRegionID)
	d.Set("isolation_type", networkProfile.IsolationType)
	d.Set("isolated_network_domain_cidr", networkProfile.IsolationNetworkDomainCIDR)
	d.Set("isolated_network_cidr_prefix", networkProfile.IsolatedNetworkCIDRPrefix)
	d.Set("name", networkProfile.Name)
	d.Set("organization_id", networkProfile.OrgID)
	d.Set("owner", networkProfile.Owner)
	d.Set("updated_at", networkProfile.UpdatedAt)

	if err := d.Set("tags", flattenTags(networkProfile.Tags)); err != nil {
		return fmt.Errorf("error setting network profile tags - error: %v", err)
	}

	if err := d.Set("links", flattenLinks(networkProfile.Links)); err != nil {
		return fmt.Errorf("error setting network profile links - error: %#v", err)
	}

	if fabricNetworkLinks, ok := networkProfile.Links["fabric-networks"]; ok {
		if fabricNetworkLinks.Hrefs != nil {
			fabricNetworkIDs := make([]string, 0, len(fabricNetworkLinks.Hrefs))

			for _, link := range fabricNetworkLinks.Hrefs {
				fabricNetworkIDs = append(fabricNetworkIDs, strings.TrimPrefix(link, "/iaas/api/fabric-networks/"))
			}

			d.Set("fabric_network_ids", fabricNetworkIDs)
		}
	}

	if extFabricNetworkLink, ok := networkProfile.Links["isolated-external-fabric-networks"]; ok {
		if extFabricNetworkLink.Href != "" {
			d.Set("isolated_external_fabric_network_id", strings.TrimPrefix(extFabricNetworkLink.Href, "/iaas/api/fabric-networks/"))
		}
	}

	if networkDomainLink, ok := networkProfile.Links["network-domains"]; ok {
		if networkDomainLink.Href != "" {
			d.Set("isolated_network_domain_id", strings.TrimPrefix(networkDomainLink.Href, "/iaas/api/network-domains/"))
		}
	}

	if securityGroupLinks, ok := networkProfile.Links["security-groups"]; ok {
		if securityGroupLinks.Hrefs != nil {
			securityGroupIDs := make([]string, 0, len(securityGroupLinks.Hrefs))

			for _, link := range securityGroupLinks.Hrefs {
				securityGroupIDs = append(securityGroupIDs, strings.TrimPrefix(link, "/iaas/api/security-groups/"))
			}

			d.Set("security_group_ids", securityGroupIDs)
		}
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
