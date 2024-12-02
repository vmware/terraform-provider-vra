// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/vmware/vra-sdk-go/pkg/client/network_ip_range"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNetworkIPRange() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkIPRangeCreate,
		ReadContext:   resourceNetworkIPRangeRead,
		UpdateContext: resourceNetworkIPRangeUpdate,
		DeleteContext: resourceNetworkIPRangeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was created. The date is in ISO 8601 and UTC.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A human-friendly description.",
			},
			"end_ip_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "End IP address of the range.",
				// Do we need to validate?
			},
			"external_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "External entity Id on the provider side.",
			},
			"fabric_network_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Id of the fabric network.",
				Deprecated:  "Please use `fabric_network_ids` instead.",
			},
			"fabric_network_ids": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"fabric_network_id"},
				Description:   "The Ids of the fabric networks.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ip_version": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "IP address version: IPv4 or IPv6.",
				ValidateFunc: validation.StringInSlice([]string{"IPv4", "IPv6"}, true),
			},
			"links": linksSchema(),
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the network IP range.",
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
			"start_ip_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Start IP address of the range.",
				// Do we need to validate?
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC.",
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceNetworkIPRangeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to create vra_network_ip_range resource")
	apiClient := m.(*Client).apiClient

	name := d.Get("name").(string)
	endIPAddress := d.Get("end_ip_address").(string)
	startIPAddress := d.Get("start_ip_address").(string)
	fabricNetworkIDs := []string{}
	if v, ok := d.GetOk("fabric_network_ids"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.FromErr(errors.New("specified fabric_network_ids are not unique"))
		}
		fabricNetworkIDs = expandStringList(v.(*schema.Set).List())
	} else if v, ok := d.GetOk("fabric_network_id"); ok {
		fabricNetworkIDs = append(fabricNetworkIDs, v.(string))
	}

	networkIPRangeSpecification := models.NetworkIPRangeSpecification{
		EndIPAddress:     &endIPAddress,
		FabricNetworkIds: fabricNetworkIDs,
		IPVersion:        d.Get("ip_version").(string),
		Name:             &name,
		StartIPAddress:   &startIPAddress,
		Tags:             expandTags(d.Get("tags").(*schema.Set).List()),
	}

	if v, ok := d.GetOk("description"); ok {
		networkIPRangeSpecification.Description = v.(string)
	}

	log.Printf("[DEBUG] Creating vra_network_ip_range with specification: %#v", networkIPRangeSpecification)

	createNetworkIPRangeCreated, err := apiClient.NetworkIPRange.CreateInternalNetworkIPRange(network_ip_range.NewCreateInternalNetworkIPRangeParams().WithBody(&networkIPRangeSpecification))
	if err != nil {
		return diag.FromErr(err)
	}
	if createNetworkIPRangeCreated != nil {
		d.SetId(*createNetworkIPRangeCreated.Payload.ID)
	}
	log.Printf("Finished creating vra_network_ip_range resource with name %s", d.Get("name"))

	return resourceNetworkIPRangeRead(ctx, d, m)
}

func resourceNetworkIPRangeRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Reading the vra_network_ip_range resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	resp, err := apiClient.NetworkIPRange.GetInternalNetworkIPRange(network_ip_range.NewGetInternalNetworkIPRangeParams().WithID(id))
	if err != nil {
		return diag.FromErr(err)
	}

	networkIPRange := *resp.Payload
	d.Set("created_at", networkIPRange.CreatedAt)
	d.Set("description", networkIPRange.Description)
	d.Set("end_ip_address", networkIPRange.EndIPAddress)
	d.Set("external_id", networkIPRange.ExternalID)
	d.Set("ip_version", networkIPRange.IPVersion)
	d.Set("name", networkIPRange.Name)
	d.Set("org_id", networkIPRange.OrgID)
	d.Set("owner", networkIPRange.Owner)
	d.Set("start_ip_address", networkIPRange.StartIPAddress)
	d.Set("updated_at", networkIPRange.UpdatedAt)

	if fabricNetworkLinks, ok := networkIPRange.Links["fabric-networks"]; ok {
		fabricNetworkIDs := make([]string, 0, len(fabricNetworkLinks.Hrefs))
		for _, fabricNetworkLink := range fabricNetworkLinks.Hrefs {
			fabricNetworkIDs = append(fabricNetworkIDs, strings.TrimPrefix(fabricNetworkLink, "/iaas/api/fabric-networks/"))
		}
		d.Set("fabric_network_ids", fabricNetworkIDs)
	} else if fabricNetworkLink, ok := networkIPRange.Links["fabric-network"]; ok {
		fabricNetworkIDs := []string{strings.TrimPrefix(fabricNetworkLink.Href, "/iaas/api/fabric-networks/")}
		d.Set("fabric_network_ids", fabricNetworkIDs)
	}

	if err := d.Set("links", flattenLinks(networkIPRange.Links)); err != nil {
		return diag.Errorf("error setting network ip range links - error: %#v", err)
	}

	if err := d.Set("tags", flattenTags(networkIPRange.Tags)); err != nil {
		return diag.Errorf("error setting network ip range tags - error: %v", err)
	}

	log.Printf("Finished reading the vra_network_ip_range resource with name %s", d.Get("name"))

	return nil
}

func resourceNetworkIPRangeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to update the vra_network_ip_range resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	name := d.Get("name").(string)
	endIPAddress := d.Get("end_ip_address").(string)
	startIPAddress := d.Get("start_ip_address").(string)
	fabricNetworkIDs := []string{}
	if v, ok := d.GetOk("fabric_network_ids"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.FromErr(errors.New("specified fabric_network_ids are not unique"))
		}
		fabricNetworkIDs = expandStringList(v.(*schema.Set).List())
	} else if v, ok := d.GetOk("fabric_network_id"); ok {
		fabricNetworkIDs = append(fabricNetworkIDs, v.(string))
	}

	networkIPRangeSpecification := models.NetworkIPRangeSpecification{
		EndIPAddress:     &endIPAddress,
		FabricNetworkIds: fabricNetworkIDs,
		IPVersion:        d.Get("ip_version").(string),
		Name:             &name,
		StartIPAddress:   &startIPAddress,
		Tags:             expandTags(d.Get("tags").(*schema.Set).List()),
	}

	if v, ok := d.GetOk("description"); ok {
		networkIPRangeSpecification.Description = v.(string)
	}
	log.Printf("[DEBUG] Updating vra_network_ip_range resource with specification: %#v", networkIPRangeSpecification)

	_, err := apiClient.NetworkIPRange.UpdateInternalNetworkIPRange(network_ip_range.NewUpdateInternalNetworkIPRangeParams().WithID(id).WithBody(&networkIPRangeSpecification))
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("Finished updating the vra_network_ip_range resource with name %s", d.Get("name"))
	return resourceNetworkIPRangeRead(ctx, d, m)

}

func resourceNetworkIPRangeDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to delete the vra_network_ip_range resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	_, err := apiClient.NetworkIPRange.DeleteInternalNetworkIPRange(network_ip_range.NewDeleteInternalNetworkIPRangeParams().WithID(id))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("Finished deleting the vra_network_ip_range resource with name %s", d.Get("name"))
	return nil
}
