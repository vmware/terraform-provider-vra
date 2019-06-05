package cas

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/cas-sdk-go/pkg/client/network"
)

func dataSourceNetworkDomain() *schema.Resource {
	return &schema.Resource{
		Read: resourceNetworkDomainRead,

		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"cidr": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_properties": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"cloud_account_ids": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": tagsSchema(),
			"external_region_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": linksSchema(),
			"organization_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNetworkDomainRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading the cas_network_domain data source with name %s", d.Get("name"))
	apiClient := meta.(*Client).apiClient

	filter := d.Get("filter").(string)

	getResp, err := apiClient.Network.GetNetworkDomains(network.NewGetNetworkDomainsParams().WithDollarFilter(withString(filter)))
	if err != nil {
		return err
	}

	NetworkDomains := getResp.Payload
	if len(NetworkDomains.Content) > 1 {
		return fmt.Errorf("cas_network_domain must filter to a network domain")
	}
	if len(NetworkDomains.Content) == 0 {
		return fmt.Errorf("cas_network_domain filter did not match any network domain")
	}

	NetworkDomain := NetworkDomains.Content[0]
	d.SetId(*NetworkDomain.ID)
	d.Set("cidr", NetworkDomain.Cidr)
	d.Set("cloud_account_ids", NetworkDomain.CloudAccountIds)
	d.Set("created_at", NetworkDomain.CreatedAt)
	d.Set("description", NetworkDomain.Description)
	d.Set("external_id", NetworkDomain.ExternalID)
	d.Set("external_region_id", NetworkDomain.ExternalRegionID)
	d.Set("custom_properties", NetworkDomain.CustomProperties)
	d.Set("name", NetworkDomain.Name)
	d.Set("organization_id", NetworkDomain.OrganizationID)
	d.Set("owner", NetworkDomain.Owner)
	d.Set("updated_at", NetworkDomain.UpdatedAt)

	if err := d.Set("tags", flattenTags(NetworkDomain.Tags)); err != nil {
		return fmt.Errorf("error getting network domain tags - error: %v", err)
	}

	if err := d.Set("links", flattenLinks(NetworkDomain.Links)); err != nil {
		return fmt.Errorf("error getting network domain links - error: %#v", err)
	}

	log.Printf("Finished reading the cas_network_domain data source with name %s", d.Get("name"))
	return nil
}
