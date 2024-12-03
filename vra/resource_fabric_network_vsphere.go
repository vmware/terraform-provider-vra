// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"
	"errors"
	"log"

	"github.com/vmware/vra-sdk-go/pkg/client/fabric_network"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFabricNetworkVsphere() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFabricNetworkVsphereCreate,
		ReadContext:   resourceFabricNetworkVsphereRead,
		UpdateContext: resourceFabricNetworkVsphereUpdate,
		DeleteContext: resourceFabricNetworkVsphereDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"cidr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_account_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"default_gateway": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"default_ipv6_gateway": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dns_search_domains": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"dns_server_addresses": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"external_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"is_public": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
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

func resourceFabricNetworkVsphereCreate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	log.Printf("Starting to create vra_fabric_network resource")
	return diag.FromErr(errors.New("vra_fabric_network_vsphere resources are only importable"))
}

func resourceFabricNetworkVsphereRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Reading the vra_fabric_network_vsphere resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	resp, err := apiClient.FabricNetwork.GetVsphereFabricNetwork(fabric_network.NewGetVsphereFabricNetworkParams().WithID(id))
	if err != nil {
		return diag.FromErr(err)
	}

	VsphereFabricNetwork := *resp.Payload
	d.Set("cidr", VsphereFabricNetwork.Cidr)
	d.Set("created_at", VsphereFabricNetwork.CreatedAt)
	d.Set("cloud_account_ids", VsphereFabricNetwork.CloudAccountIds)
	d.Set("default_gateway", VsphereFabricNetwork.DefaultGateway)
	d.Set("default_ipv6_gateway", VsphereFabricNetwork.DefaultIPV6Gateway)
	d.Set("dns_search_domains", VsphereFabricNetwork.DNSSearchDomains)
	d.Set("dns_server_addresses", VsphereFabricNetwork.DNSServerAddresses)
	d.Set("domain", VsphereFabricNetwork.Domain)
	d.Set("ipv6_cidr", VsphereFabricNetwork.IPV6Cidr)
	d.Set("is_default", VsphereFabricNetwork.IsDefault)
	d.Set("is_public", VsphereFabricNetwork.IsPublic)
	d.Set("org_id", VsphereFabricNetwork.OrgID)
	d.Set("name", VsphereFabricNetwork.Name)
	d.Set("owner", VsphereFabricNetwork.Owner)
	d.Set("updated_at", VsphereFabricNetwork.UpdatedAt)
	d.Set("custom_properties", VsphereFabricNetwork.CustomProperties)

	if err := d.Set("tags", flattenTags(VsphereFabricNetwork.Tags)); err != nil {
		return diag.Errorf("error setting network ip range tags - error: %v", err)
	}

	if err := d.Set("links", flattenLinks(VsphereFabricNetwork.Links)); err != nil {
		return diag.Errorf("error setting network ip range links - error: %#v", err)
	}

	log.Printf("Finished reading the vra_fabric_network_vsphere resource with name %s", d.Get("name"))

	return nil
}

func resourceFabricNetworkVsphereUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to Update vra_fabric_network resource")
	var dnsSearchDomains []string
	var dnsServerAddresses []string
	apiClient := m.(*Client).apiClient

	id := d.Id()

	if v, ok := d.GetOk("dns_search_domains"); ok {
		if !compareUnique(v.([]interface{})) {
			return diag.FromErr(errors.New("specified dns search domains are not unique"))
		}
		dnsSearchDomains = expandStringList(v.([]interface{}))
	}
	if v, ok := d.GetOk("dns_server_addresses"); ok {
		if !compareUnique(v.([]interface{})) {
			return diag.FromErr(errors.New("specified dns server addresses are not unique"))
		}
		dnsServerAddresses = expandStringList(v.([]interface{}))
	}

	isDefault := d.Get("is_default").(bool)
	isPublic := d.Get("is_public").(bool)

	VsphereFabricNetworkSpecification := models.FabricNetworkVsphereSpecification{
		Cidr:               d.Get("cidr").(string),
		DefaultGateway:     d.Get("default_gateway").(string),
		DefaultIPV6Gateway: d.Get("default_ipv6_gateway").(string),
		DNSSearchDomains:   dnsSearchDomains,
		DNSServerAddresses: dnsServerAddresses,
		Domain:             d.Get("domain").(string),
		IPV6Cidr:           d.Get("ipv6_cidr").(string),
		IsDefault:          withBool(isDefault),
		IsPublic:           withBool(isPublic),
		Tags:               expandTags(d.Get("tags").(*schema.Set).List()),
	}

	_, err := apiClient.FabricNetwork.UpdatevSphereFabricNetwork(fabric_network.NewUpdatevSphereFabricNetworkParams().WithID(id).WithBody(&VsphereFabricNetworkSpecification))
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("finished Updating vra_fabric_network_vsphere resource")
	return resourceFabricNetworkVsphereRead(ctx, d, m)

}

func resourceFabricNetworkVsphereDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	log.Printf("Starting to delete the vra_fabric_network_vsphere resource with name %s", d.Get("name"))
	return nil
}
