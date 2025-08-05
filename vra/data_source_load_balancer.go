// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/load_balancer"
)

func dataSourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLoadBalancerRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nics": nicsSchema(false),
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"routes": routesSchema(false),
			"custom_properties": {
				Type:     schema.TypeMap,
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
			"internet_facing": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"tags":    tagsSchema(),
			"targets": LoadBalancerTargetSchema(),
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_zone_id": {
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
		},
	}
}

func dataSourceLoadBalancerRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Get("id").(string)
	resp, err := apiClient.LoadBalancer.GetLoadBalancer(load_balancer.NewGetLoadBalancerParams().WithID(id))
	if err != nil {
		return diag.FromErr(err)
	}

	loadBalancer := *resp.Payload
	d.SetId(id)
	_ = d.Set("address", loadBalancer.Address)
	_ = d.Set("created_at", loadBalancer.CreatedAt)
	_ = d.Set("custom_properties", loadBalancer.CustomProperties)
	_ = d.Set("description", loadBalancer.Description)
	_ = d.Set("deployment_id", loadBalancer.DeploymentID)
	_ = d.Set("external_id", loadBalancer.ExternalID)
	_ = d.Set("external_region_id", loadBalancer.ExternalRegionID)
	_ = d.Set("external_zone_id", loadBalancer.ExternalZoneID)
	_ = d.Set("name", loadBalancer.Name)
	_ = d.Set("organization_id", loadBalancer.OrgID)
	_ = d.Set("owner", loadBalancer.Owner)
	_ = d.Set("project_id", loadBalancer.ProjectID)
	_ = d.Set("updated_at", loadBalancer.UpdatedAt)

	if err := d.Set("tags", flattenTags(loadBalancer.Tags)); err != nil {
		return diag.Errorf("error setting load balancer tags - error: %v", err)
	}
	if err := d.Set("routes", flattenRoutes(loadBalancer.Routes)); err != nil {
		return diag.Errorf("error setting load balancer routes - error: %v", err)
	}

	if err := d.Set("links", flattenLinks(loadBalancer.Links)); err != nil {
		return diag.Errorf("error setting load balancer links - error: %#v", err)
	}

	log.Printf("Finished reading the vra_load_balancer data source with id %s", d.Get("id"))
	return nil
}
