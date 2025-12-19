// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"errors"
	"fmt"

	"github.com/vmware/vra-sdk-go/pkg/client/network"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"name", "filter"},
				Optional:      true,
				Computed:      true,
				Description:   "The id of the network instance",
			},
			"name": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"id", "filter"},
				Optional:      true,
				Computed:      true,
				Description:   "The human-friendly name of the network instance",
			},
			"filter": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"id", "name"},
				Optional:      true,
				Description:   "The search criteria to narrow down the network instance",
			},
			"return_first": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Return the first matching network instance when set to true",
			},
			"cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IPv4 address range of the network in CIDR format",
			},
			"cloud_account_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Set of ids of the cloud accounts this resource belongs to",
			},
			"custom_properties": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Additional properties that may be used to extend the base resource",
			},
			"deployment_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Deployment id that is associated with this resource",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly description",
			},
			"external_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "External entity Id on the provider side",
			},
			"external_region_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The external regionId of the resource",
			},
			"external_zone_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The external zoneId of the resource",
			},
			"links": linksSchema(),
			"organization_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the organization this entity belongs to",
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email of the user or display name of the group that owns the entity",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the project this resource belongs to",
			},
			"tags": tagsSchema(),
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC",
			},
		},
	}
}

func dataSourceNetworkRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")
	filter, filterOK := d.GetOk("filter")
	returnFirst, _ := d.Get("return_first").(bool)

	if !idOk && !nameOk && !filterOK {
		return errors.New("one of id, name or filter must be assigned")
	}

	var net *models.Network
	if idOk {
		getResp, err := apiClient.Network.GetNetwork(network.NewGetNetworkParams().WithID(id.(string)))
		if err != nil {
			switch err.(type) {
			case *network.GetNetworkNotFound:
				return fmt.Errorf("network with id '%s' not found", id)
			default:
				// nop
			}
			return err
		}

		net = getResp.GetPayload()
	} else {
		var filterParam string
		if nameOk {
			filterParam = fmt.Sprintf("name eq '%s'", name.(string))
		} else {
			filterParam = filter.(string)
		}
		getResp, err := apiClient.Network.GetNetworks(network.NewGetNetworksParams().WithDollarFilter(&filterParam))
		if err != nil {
			return err
		}
		networks := getResp.GetPayload()
		if len(networks.Content) > 1 && !returnFirst {
			if nameOk {
				return fmt.Errorf("there are more than one network with name '%s'", name)
			}
			return errors.New("must filter to one network")
		}
		if len(networks.Content) == 0 {
			if nameOk {
				return fmt.Errorf("network with name '%s' not found", name)
			}
			return fmt.Errorf("filter doesn't match to any network")
		}
		net = networks.Content[0]
	}

	d.SetId(*net.ID)
	d.Set("cidr", net.Cidr)
	d.Set("cloud_account_ids", net.CloudAccountIds)
	d.Set("created_at", net.CreatedAt)
	d.Set("custom_properties", net.CustomProperties)
	d.Set("deployment_id", net.DeploymentID)
	d.Set("description", net.Description)
	d.Set("external_id", net.ExternalID)
	d.Set("external_region_id", net.ExternalRegionID)
	d.Set("external_zone_id", net.ExternalZoneID)
	d.Set("name", net.Name)
	d.Set("organization_id", net.OrgID)
	d.Set("owner", net.Owner)
	d.Set("project_id", net.ProjectID)
	d.Set("tags", net.Tags)
	d.Set("updated_at", net.UpdatedAt)

	if err := d.Set("links", flattenLinks(net.Links)); err != nil {
		return fmt.Errorf("error setting network links - error: %#v", err)
	}

	return nil
}
