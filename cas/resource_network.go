package cas

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/vmware/cas-sdk-go/pkg/client/network"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/cas-sdk-go/pkg/client"
	"github.com/vmware/cas-sdk-go/pkg/client/request"
	"github.com/vmware/cas-sdk-go/pkg/models"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkCreate,
		Read:   resourceNetworkRead,
		Update: resourceNetworkUpdate,
		Delete: resourceNetworkDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !strings.HasPrefix(new, old)
				},
			},
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"constraints": constraintsSchema(),
			"custom_properties": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"outbound_access": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": tagsSchema(),
			"cidr": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_zone_id": &schema.Schema{
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
			"self_link": &schema.Schema{
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

func resourceNetworkCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("Starting to create cas_network resource")
	apiClient := m.(*Client).apiClient

	name := d.Get("name").(string)
	projectID := d.Get("project_id").(string)
	constraints := expandConstraints(d.Get("constraints").(*schema.Set).List())
	tags := expandTags(d.Get("tags").(*schema.Set).List())
	customProperties := expandCustomProperties(d.Get("custom_properties").(map[string]interface{}))

	networkSpecification := models.NetworkSpecification{
		Name:             &name,
		ProjectID:        &projectID,
		Constraints:      constraints,
		Tags:             tags,
		CustomProperties: customProperties,
	}

	if v, ok := d.GetOk("description"); ok {
		networkSpecification.Description = v.(string)
	}

	if v, ok := d.GetOk("outbound_access"); ok {
		networkSpecification.OutboundAccess = v.(bool)
	}
	log.Printf("[DEBUG] create network: %#v", networkSpecification)
	createNetworkCreated, err := apiClient.Network.CreateNetwork(network.NewCreateNetworkParams().WithBody(&networkSpecification))
	if err != nil {
		return err
	}
	stateChangeFunc := resource.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    networkStateRefreshFunc(*apiClient, *createNetworkCreated.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    5 * time.Minute,
		MinTimeout: 5 * time.Second,
	}
	resourceIds, err := stateChangeFunc.WaitForState()
	log.Printf("Waitforstate returned: %T %+v %+v\n", resourceIds, resourceIds, err)

	if err != nil {
		return err
	}

	networkIDs := resourceIds.([]string)
	d.SetId(networkIDs[0])
	log.Printf("Finished to create cas_network resource with name %s", d.Get("name"))

	return resourceNetworkRead(d, m)
}

func networkStateRefreshFunc(apiClient client.MulticloudIaaS, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ret, err := apiClient.Request.GetRequestTracker(request.NewGetRequestTrackerParams().WithID(id))
		if err != nil {
			return "", models.RequestTrackerStatusFAILED, err
		}

		status := ret.Payload.Status
		switch *status {
		case models.RequestTrackerStatusFAILED:
			return []string{""}, *status, fmt.Errorf(ret.Payload.Message)
		case models.RequestTrackerStatusINPROGRESS:
			return [...]string{id}, *status, nil
		case models.RequestTrackerStatusFINISHED:
			networkIDs := make([]string, len(ret.Payload.Resources))
			for i, r := range ret.Payload.Resources {
				networkIDs[i] = strings.TrimPrefix(r, "/iaas/api/networks/")
			}
			return networkIDs, *status, nil
		default:
			return [...]string{id}, ret.Payload.Message, fmt.Errorf("networkStateRefreshFunc: unknown status %v", *status)
		}
	}
}

func resourceNetworkRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("Reading the cas_network resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	resp, err := apiClient.Network.GetNetwork(network.NewGetNetworkParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *network.GetNetworkNotFound:
			d.SetId("")
			return nil
		}
		return err
	}

	network := *resp.Payload
	d.Set("cidr", network.Cidr)
	d.Set("custom_properties", network.CustomProperties)
	d.Set("description", network.Description)
	d.Set("external_id", network.ExternalID)
	d.Set("external_zone_id", network.ExternalZoneID)
	d.Set("name", network.Name)
	d.Set("organization_id", network.OrganizationID)
	d.Set("owner", network.Owner)
	d.Set("project_id", network.ProjectID)
	d.Set("updated_at", network.UpdatedAt)

	if err := d.Set("tags", flattenTags(network.Tags)); err != nil {
		return fmt.Errorf("error setting network tags - error: %v", err)
	}

	if err := d.Set("links", flattenLinks(network.Links)); err != nil {
		return fmt.Errorf("error setting network links - error: %#v", err)
	}

	log.Printf("Finished reading the cas_network resource with name %s", d.Get("name"))
	return nil
}

func resourceNetworkUpdate(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("Updating a network resource is not allowed")
}

func resourceNetworkDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("Starting to delete the cas_network resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	deleteNetwork, err := apiClient.Network.DeleteNetwork(network.NewDeleteNetworkParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *network.DeleteNetworkNotFound:
			d.SetId("")
			return nil
		}
		return err
	}
	stateChangeFunc := resource.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    networkStateRefreshFunc(*apiClient, *deleteNetwork.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    5 * time.Minute,
		MinTimeout: 5 * time.Second,
	}

	_, err = stateChangeFunc.WaitForState()
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("Finished deleting the cas_network resource with name %s", d.Get("name"))
	return nil
}
