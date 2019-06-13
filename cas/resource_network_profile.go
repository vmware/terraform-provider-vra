package cas

import (
	"fmt"
	"log"

	"github.com/vmware/cas-sdk-go/pkg/client/network_profile"
	"github.com/vmware/cas-sdk-go/pkg/models"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNetworkProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkProfileCreate,
		Read:   resourceNetworkProfileRead,
		Update: resourceNetworkProfileUpdate,
		Delete: resourceNetworkProfileDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"region_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"custom_properties": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"fabric_network_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"isolated_network_cidr_prefix": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"isolation_external_fabric_network_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"isolation_network_domain_cidr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"isolation_network_domain_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"isolation_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
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

func resourceNetworkProfileCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("Starting to create cas_network_profile resource")
	apiClient := m.(*Client).apiClient

	name := d.Get("name").(string)
	regionID := d.Get("region_id").(string)

	networkProfileSpecification := models.NetworkProfileSpecification{
		IsolationType:                    d.Get("isolation_type").(string),
		IsolationNetworkDomainID:         d.Get("isolation_network_domain_id").(string),
		IsolationNetworkDomainCIDR:       d.Get("isolation_network_domain_cidr").(string),
		IsolationExternalFabricNetworkID: d.Get("isolation_external_fabric_network_id").(string),
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
		if !compareUnique(v.([]interface{})) {
			return fmt.Errorf("Specified fabric network ids are not unique")
		}
		networkProfileSpecification.FabricNetworkIds = expandStringList(v.([]interface{}))
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		if !compareUnique(v.([]interface{})) {
			return fmt.Errorf("Specified security group ids are not unique")
		}
		networkProfileSpecification.SecurityGroupIds = expandStringList(v.([]interface{}))
	}

	log.Printf("[DEBUG] create network profile: %#v", networkProfileSpecification)
	createNetworkProfileCreated, err := apiClient.NetworkProfile.CreateNetworkProfile(network_profile.NewCreateNetworkProfileParams().WithBody(&networkProfileSpecification))
	if err != nil {
		return err
	}

	d.SetId(*createNetworkProfileCreated.Payload.ID)
	log.Printf("Finished to create cas_network_profile resource with name %s", d.Get("name"))

	return resourceNetworkProfileRead(d, m)
}

func resourceNetworkProfileRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("Reading the cas_network_profile resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	resp, err := apiClient.NetworkProfile.GetNetworkProfile(network_profile.NewGetNetworkProfileParams().WithID(id))
	if err != nil {
		return err
	}

	networkProfile := *resp.Payload
	d.Set("custom_properties", networkProfile.CustomProperties)
	d.Set("description", networkProfile.Description)
	d.Set("external_region_id", networkProfile.ExternalRegionID)
	d.Set("isolation_type", networkProfile.IsolationType)
	d.Set("isolation_network_domain_cidr", networkProfile.IsolationNetworkDomainCIDR)
	d.Set("isolated_network_cidr_prefix", networkProfile.IsolatedNetworkCIDRPrefix)
	d.Set("name", networkProfile.Name)
	d.Set("organization_id", networkProfile.OrganizationID)
	d.Set("owner", networkProfile.Owner)
	d.Set("updated_at", networkProfile.UpdatedAt)

	if err := d.Set("tags", flattenTags(networkProfile.Tags)); err != nil {
		return fmt.Errorf("error setting network profile tags - error: %v", err)
	}

	if err := d.Set("links", flattenLinks(networkProfile.Links)); err != nil {
		return fmt.Errorf("error setting network profile links - error: %#v", err)
	}

	log.Printf("Finished reading the cas_network_profile resource with name %s", d.Get("name"))
	return nil
}

func resourceNetworkProfileUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	name := d.Get("name").(string)
	regionID := d.Get("region_id").(string)

	networkProfileSpecification := models.NetworkProfileSpecification{
		IsolationType:                    d.Get("isolation_type").(string),
		IsolationNetworkDomainID:         d.Get("isolation_network_domain_id").(string),
		IsolationNetworkDomainCIDR:       d.Get("isolation_network_domain_cidr").(string),
		IsolationExternalFabricNetworkID: d.Get("isolation_external_fabric_network_id").(string),
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
		if !compareUnique(v.([]interface{})) {
			return fmt.Errorf("Specified fabric network ids are not unique")
		}
		networkProfileSpecification.FabricNetworkIds = expandStringList(v.([]interface{}))
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		if !compareUnique(v.([]interface{})) {
			return fmt.Errorf("Specified security group ids are not unique")
		}
		networkProfileSpecification.SecurityGroupIds = expandStringList(v.([]interface{}))
	}

	_, err := apiClient.NetworkProfile.UpdateNetworkProfile(network_profile.NewUpdateNetworkProfileParams().WithID(id).WithBody(&networkProfileSpecification))
	if err != nil {
		return err
	}

	return resourceNetworkProfileRead(d, m)
}

func resourceNetworkProfileDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("Starting to delete the cas_network_profile resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	_, err := apiClient.NetworkProfile.DeleteNetworkProfile(network_profile.NewDeleteNetworkProfileParams().WithID(id))
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("Finished deleting the cas_network_profile resource with name %s", d.Get("name"))
	return nil
}
