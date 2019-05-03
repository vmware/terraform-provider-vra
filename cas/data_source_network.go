package cas

import (
	"fmt"
	"strings"

	"github.com/vmware/cas-sdk-go/pkg/client/network"
	"github.com/vmware/cas-sdk-go/pkg/models"
	tango "github.com/vmware/terraform-provider-cas/sdk"

	"github.com/hashicorp/terraform/helper/schema"
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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !strings.HasPrefix(new, old)
				},
			},
			"cidr": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_properties": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"constraints": constraintsSDKSchema(),
			"tags":        tagsSDKSchema(),
			"outbound_access": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"external_zone_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": linksSchema(),
		},
	}
}

func dataSourceNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*tango.Client)
	apiClient := client.GetAPIClient()

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk == false && nameOk == false {
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
		d.Set("external_id", network.ExternalID)
		d.Set("external_zone_id", network.ExternalZoneID)
		d.Set("name", network.Name)
		d.Set("organization_id", network.OrganizationID)
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
