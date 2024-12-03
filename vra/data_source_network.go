// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"strings"

	"github.com/vmware/vra-sdk-go/pkg/client/network"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					return !strings.HasPrefix(new, old)
				},
			},
			"cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deployment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_properties": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"constraints": constraintsSchema(),
			"tags":        tagsSchema(),
			"outbound_access": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"external_zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": linksSchema(),
		},
	}
}

func dataSourceNetworkRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if !idOk && !nameOk {
		return fmt.Errorf("One of id or name must be assigned")
	}

	getResp, err := apiClient.Network.GetNetworks(network.NewGetNetworksParams())
	if err != nil {
		return err
	}

	setFields := func(network *models.Network) {
		d.SetId(*network.ID)
		d.Set("cidr", network.Cidr)
		d.Set("custom_properties", network.CustomProperties)
		d.Set("description", network.Description)
		d.Set("deployment_id", network.DeploymentID)
		d.Set("external_id", network.ExternalID)
		d.Set("external_zone_id", network.ExternalZoneID)
		d.Set("name", network.Name)
		d.Set("organization_id", network.OrgID)
		d.Set("owner", network.Owner)
		d.Set("project_id", network.ProjectID)
		d.Set("tags", network.Tags)
		d.Set("updated_at", network.UpdatedAt)
	}
	for _, network := range getResp.Payload.Content {
		if idOk && network.ID == id {
			setFields(network)
			return nil
		}
		if nameOk && network.Name == name {
			setFields(network)
			return nil
		}
	}

	return fmt.Errorf("network %s not found", name)
}
