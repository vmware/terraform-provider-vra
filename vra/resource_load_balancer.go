package vra

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/vmware/vra-sdk-go/pkg/client"
	"github.com/vmware/vra-sdk-go/pkg/client/load_balancer"
	"github.com/vmware/vra-sdk-go/pkg/client/request"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoadBalancerCreate,
		Read:   resourceLoadBalancerRead,
		Update: resourceLoadBalancerUpdate,
		Delete: resourceLoadBalancerDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !strings.HasPrefix(new, old)
				},
			},
			"nics": nicsSchema(true),
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"routes": routesSchema(true),
			"custom_properties": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"internet_facing": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": tagsSchema(),
			"target_links": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_region_id": &schema.Schema{
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

func resourceLoadBalancerCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("Starting to create vra_load_balancer resource")
	apiClient := m.(*Client).apiClient

	name := d.Get("name").(string)
	projectID := d.Get("project_id").(string)
	tags := expandTags(d.Get("tags").(*schema.Set).List())
	customProperties := expandCustomProperties(d.Get("custom_properties").(map[string]interface{}))
	nics := expandNics(d.Get("nics").(*schema.Set).List())
	routes := expandRoutes(d.Get("routes").(*schema.Set).List())

	loadBalancerSpecification := models.LoadBalancerSpecification{
		Name:             &name,
		ProjectID:        &projectID,
		Routes:           routes,
		Tags:             tags,
		CustomProperties: customProperties,
		Nics:             nics,
	}

	if v, ok := d.GetOk("description"); ok {
		loadBalancerSpecification.Description = v.(string)
	}

	if v, ok := d.GetOk("internet_facing"); ok {
		loadBalancerSpecification.InternetFacing = v.(bool)
	}

	if v, ok := d.GetOk("target_links"); ok {
		targetLinks := make([]string, 0)
		for _, value := range v.([]interface{}) {
			targetLinks = append(targetLinks, value.(string))
		}

		loadBalancerSpecification.TargetLinks = targetLinks
	}

	log.Printf("[DEBUG] create load lalancer: %#v", loadBalancerSpecification)
	createLoadBalancerCreated, err := apiClient.LoadBalancer.CreateLoadBalancer(load_balancer.NewCreateLoadBalancerParams().WithBody(&loadBalancerSpecification))
	if err != nil {
		return err
	}

	stateChangeFunc := resource.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    loadBalancerStateRefreshFunc(*apiClient, *createLoadBalancerCreated.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    5 * time.Minute,
		MinTimeout: 5 * time.Second,
	}

	resourceIDs, err := stateChangeFunc.WaitForState()
	if err != nil {
		return err
	}

	loadBalancerIDs := resourceIDs.([]string)
	i := strings.LastIndex(loadBalancerIDs[0], "/")
	loadBalancerID := loadBalancerIDs[0][i+1 : len(loadBalancerIDs[0])]
	d.SetId(loadBalancerID)
	log.Printf("Finished to create vra_load_balancer resource with name %s", d.Get("name"))

	return resourceLoadBalancerRead(d, m)
}

func loadBalancerStateRefreshFunc(apiClient client.MulticloudIaaS, id string) resource.StateRefreshFunc {
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
			loadBalancerIDs := make([]string, len(ret.Payload.Resources))
			for i, r := range ret.Payload.Resources {
				loadBalancerIDs[i] = strings.TrimPrefix(r, "/iaas/api/load-balancer/")
			}
			return loadBalancerIDs, *status, nil
		default:
			return [...]string{id}, ret.Payload.Message, fmt.Errorf("loadBalancerStateRefreshFunc: unknown status %v", *status)
		}
	}
}

func resourceLoadBalancerRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("Reading the vra_load_balancer resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	resp, err := apiClient.LoadBalancer.GetLoadBalancer(load_balancer.NewGetLoadBalancerParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *load_balancer.GetLoadBalancerNotFound:
			d.SetId("")
			return nil
		}
		return err
	}

	loadBalancer := *resp.Payload
	d.Set("address", loadBalancer.Address)
	d.Set("created_at", loadBalancer.CreatedAt)
	d.Set("custom_properties", loadBalancer.CustomProperties)
	d.Set("description", loadBalancer.Description)
	d.Set("external_id", loadBalancer.ExternalID)
	d.Set("external_region_id", loadBalancer.ExternalRegionID)
	d.Set("external_zone_id", loadBalancer.ExternalZoneID)
	d.Set("name", loadBalancer.Name)
	d.Set("organization_id", loadBalancer.OrganizationID)
	d.Set("owner", loadBalancer.Owner)
	d.Set("project_id", loadBalancer.ProjectID)
	d.Set("updated_at", loadBalancer.UpdatedAt)

	if err := d.Set("tags", flattenTags(loadBalancer.Tags)); err != nil {
		return fmt.Errorf("error setting machine tags - error: %v", err)
	}
	if err := d.Set("routes", flattenRoutes(loadBalancer.Routes)); err != nil {
		return fmt.Errorf("error setting machine tags - error: %v", err)
	}

	if err := d.Set("links", flattenLinks(loadBalancer.Links)); err != nil {
		return fmt.Errorf("error setting machine links - error: %#v", err)
	}

	log.Printf("Finished reading the vra_machine resource with name %s", d.Get("name"))
	return nil
}

func resourceLoadBalancerUpdate(d *schema.ResourceData, m interface{}) error {

	//Day 2 scale operation for load balancer
	log.Printf("Starting to update vra_load_balancer resource")
	apiClient := m.(*Client).apiClient

	name := d.Get("name").(string)
	projectID := d.Get("project_id").(string)
	tags := expandTags(d.Get("tags").(*schema.Set).List())
	customProperties := expandCustomProperties(d.Get("custom_properties").(map[string]interface{}))
	nics := expandNics(d.Get("nics").(*schema.Set).List())
	routes := expandRoutes(d.Get("routes").(*schema.Set).List())

	loadBalancerSpecification := models.LoadBalancerSpecification{
		Name:             &name,
		ProjectID:        &projectID,
		Routes:           routes,
		Tags:             tags,
		CustomProperties: customProperties,
		Nics:             nics,
	}

	if v, ok := d.GetOk("description"); ok {
		loadBalancerSpecification.Description = v.(string)
	}

	if v, ok := d.GetOk("internet_facing"); ok {
		loadBalancerSpecification.InternetFacing = v.(bool)
	}

	if v, ok := d.GetOk("target_links"); ok {
		targetLinks := make([]string, 0)
		for _, value := range v.([]interface{}) {
			targetLinks = append(targetLinks, value.(string))
		}

		loadBalancerSpecification.TargetLinks = targetLinks
	}

	log.Printf("[DEBUG] scale load lalancer: %#v", loadBalancerSpecification)
	scaleLoadBalancerCreated, err := apiClient.LoadBalancer.ScaleLoadBalancer(load_balancer.NewScaleLoadBalancerParams().WithBody(&loadBalancerSpecification))
	if err != nil {
		return err
	}

	stateChangeFunc := resource.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    loadBalancerStateRefreshFunc(*apiClient, *scaleLoadBalancerCreated.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    5 * time.Minute,
		MinTimeout: 5 * time.Second,
	}

	resourceIDs, err := stateChangeFunc.WaitForState()
	if err != nil {
		return err
	}

	loadBalancerIDs := resourceIDs.([]string)
	i := strings.LastIndex(loadBalancerIDs[0], "/")
	loadBalancerID := loadBalancerIDs[0][i+1 : len(loadBalancerIDs[0])]
	d.SetId(loadBalancerID)
	log.Printf("Finished to update vra_load_balancer resource with name %s", d.Get("name"))

	return resourceLoadBalancerRead(d, m)
}

func resourceLoadBalancerDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("Starting to delete the vra_load_balancer resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	deleteLoadBalancer, err := apiClient.LoadBalancer.DeleteLoadBalancer(load_balancer.NewDeleteLoadBalancerParams().WithID(id))
	if err != nil {
		return err
	}
	stateChangeFunc := resource.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    loadBalancerStateRefreshFunc(*apiClient, *deleteLoadBalancer.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    5 * time.Minute,
		MinTimeout: 5 * time.Second,
	}

	_, err = stateChangeFunc.WaitForState()
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("Finished deleting the vra_load_balancer resource with name %s", d.Get("name"))
	return nil
}
