package cas

import (
	"fmt"
	"log"
	"strings"

	"github.com/vmware/terraform-provider-cas/sdk"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkCreate,
		Read:   resourceNetworkRead,
		Update: resourceNetworkUpdate,
		Delete: resourceNetworkDelete,

		Schema: map[string]*schema.Schema{
			"cidr": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !strings.HasPrefix(new, old)
				},
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"custom_properties": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Optional: true,
			},
			"constraints": constraintsSchema(),
			"tags":        tagsSchema(),
			"outbound_access": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
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

func resourceNetworkCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)

	networkSpecification := tango.NetworkSpecification{
		Name:             d.Get("name").(string),
		ProjectID:        client.GetProjectID(),
		Constraints:      expandConstraints(d.Get("constraints").([]interface{})),
		Tags:             expandTags(d.Get("tags").([]interface{})),
		CustomProperties: expandCustomProperties(d.Get("custom_properties").(map[string]interface{})),
	}

	networkSpecification.CustomProperties["__composition_context_id"] = client.GetDeploymentID()

	if v, ok := d.GetOk("description"); ok {
		networkSpecification.Description = v.(string)
	}

	if v, ok := d.GetOk("outbound_access"); ok {
		networkSpecification.OutboundAccess = v.(bool)
	}

	log.Printf("[DEBUG] record create network: %#v", networkSpecification)
	resourceObject, err := client.CreateResource(networkSpecification)
	if err != nil {
		return err
	}

	networkObject := resourceObject.(*tango.Network)

	d.SetId(networkObject.ID)
	d.Set("name", networkObject.Name)
	d.Set("cidr", networkObject.CIDR)
	d.Set("external_zone_id", networkObject.ExternalZoneID)
	d.Set("external_id", networkObject.ExternalID)
	d.Set("self_link", networkObject.SelfLink)
	d.Set("updated_at", networkObject.UpdatedAt)
	d.Set("owner", networkObject.Owner)
	d.Set("organization_id", networkObject.OrganizationID)
	d.Set("custom_properties", networkObject.CustomProperties)

	if err := d.Set("tags", flattenTags(networkObject.Tags)); err != nil {
		return fmt.Errorf("Error setting Network tags - error: %#v", err)
	}

	if err := d.Set("links", flattenLinks(networkObject.Links)); err != nil {
		return fmt.Errorf("Error setting Network links - error: %#v", err)
	}

	return nil
}

func resourceNetworkRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)

	resourceObject, err := client.ReadResource(getSelfLink(d.Get("links").([]interface{})))
	if err != nil {
		d.SetId("")
		return nil
	}

	networkObject := resourceObject.(*tango.Network)

	d.Set("cidr", networkObject.CIDR)
	d.Set("project_id", networkObject.ProjectID)
	d.Set("external_zone_id", networkObject.ExternalZoneID)
	d.Set("external_id", networkObject.ExternalID)
	d.Set("name", networkObject.Name)
	d.Set("description", networkObject.Description)
	d.Set("self_link", networkObject.SelfLink)
	d.Set("updated_at", networkObject.UpdatedAt)
	d.Set("owner", networkObject.Owner)
	d.Set("organization_id", networkObject.OrganizationID)
	d.Set("custom_properties", networkObject.CustomProperties)

	if err := d.Set("tags", flattenTags(networkObject.Tags)); err != nil {
		return fmt.Errorf("Error setting Network tags - error: %#v", err)
	}

	if err := d.Set("links", flattenLinks(networkObject.Links)); err != nil {
		return fmt.Errorf("Error setting Network links - error: %#v", err)
	}

	return nil
}

func resourceNetworkUpdate(d *schema.ResourceData, m interface{}) error {

	return nil
}

func resourceNetworkDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)
	err := client.DeleteResource(getSelfLink(d.Get("links").([]interface{})))

	if err != nil && strings.Contains(err.Error(), "404") { // already deleted
		return nil
	}

	return err
}
