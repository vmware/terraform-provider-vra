package tango

import (
	"fmt"
	"log"
	"strings"

	"tango-terraform-provider/tango/client"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoadBalancerCreate,
		Read:   resourceLoadBalancerRead,
		Update: resourceLoadBalancerUpdate,
		Delete: resourceLoadBalancerDelete,

		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
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
			"tags": tagsSchema(),
			"custom_properties": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Optional: true,
			},
			"target_links": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"internet_facing": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"nics": nicsSchema(true),
			"routes": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"port": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"member_protocol": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"member_port": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"health_check_configuration": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"port": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"url_path": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"interval_seconds": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
									},
									"timeout_seconds": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
									},
									"unhealthy_threshold": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
									},
									"healthy_threshold": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"external_zone_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_region_id": &schema.Schema{
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
			"created_at": &schema.Schema{
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

func resourceLoadBalancerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)

	loadBalancerSpecification := tango.LoadBalancerSpecification{
		Name:             d.Get("name").(string),
		ProjectID:        client.GetProjectID(),
		Nics:             expandNics(d.Get("nics").([]interface{})),
		Routes:           expandRoutes(d.Get("routes").([]interface{})),
		CustomProperties: expandCustomProperties(d.Get("custom_properties").(map[string]interface{})),
		Tags:             expandTags(d.Get("tags").([]interface{})),
	}

	loadBalancerSpecification.CustomProperties["__composition_context_id"] = client.GetDeploymentID()

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

	log.Printf("[DEBUG] record create load balancer: %#v", loadBalancerSpecification)
	resourceObject, err := client.CreateResource(loadBalancerSpecification)
	if err != nil {
		return err
	}

	loadBalancerObject := resourceObject.(*tango.LoadBalancer)

	d.SetId(loadBalancerObject.ID)
	d.Set("address", loadBalancerObject.Address)
	d.Set("name", loadBalancerObject.Name)
	d.Set("external_zone_id", loadBalancerObject.ExternalZoneID)
	d.Set("external_region_id", loadBalancerObject.ExternalRegionID)
	d.Set("external_id", loadBalancerObject.ExternalID)
	d.Set("self_link", loadBalancerObject.SelfLink)
	d.Set("created_at", loadBalancerObject.CreatedAt)
	d.Set("updated_at", loadBalancerObject.UpdatedAt)
	d.Set("owner", loadBalancerObject.Owner)
	d.Set("organization_id", loadBalancerObject.OrganizationID)
	d.Set("custom_properties", loadBalancerObject.CustomProperties)

	if err := d.Set("tags", flattenTags(loadBalancerObject.Tags)); err != nil {
		return fmt.Errorf("Error setting Load Balancer tags - error: %#v", err)
	}

	if err := d.Set("routes", flattenRoutes(loadBalancerObject.Routes)); err != nil {
		return fmt.Errorf("Error setting Load Balancer routes - error: %#v", err)
	}

	if err := d.Set("links", flattenLinks(loadBalancerObject.Links)); err != nil {
		return fmt.Errorf("Error setting Load Balancer links - error: %#v", err)
	}

	return nil
}

func expandRoutes(configRoutes []interface{}) []tango.Route {
	routes := make([]tango.Route, 0, len(configRoutes))

	for _, configRoute := range configRoutes {
		routeMap := configRoute.(map[string]interface{})

		route := tango.Route{
			Protocol:       routeMap["protocol"].(string),
			Port:           routeMap["port"].(string),
			MemberProtocol: routeMap["member_protocol"].(string),
			MemberPort:     routeMap["member_port"].(string),
		}

		if v, ok := routeMap["health_check_configuration"].([]interface{}); ok && len(v) == 1 {
			configHCC := v[0].(map[string]interface{})

			healthCheckConfiguration := &tango.HealthCheckConfiguration{
				Protocol: configHCC["protocol"].(string),
				Port:     configHCC["port"].(string),
			}

			if v, ok := configHCC["url_path"].(string); ok && v != "" {
				healthCheckConfiguration.URLPath = v
			}

			if v, ok := configHCC["interval_seconds"].(int); ok && v != 0 {
				healthCheckConfiguration.IntervalSeconds = v
			}

			if v, ok := configHCC["timeout_seconds"].(int); ok && v != 0 {
				healthCheckConfiguration.TimeoutSeconds = v
			}

			if v, ok := configHCC["unhealthy_threshold"].(int); ok && v != 0 {
				healthCheckConfiguration.UnhealthThreshold = v
			}

			if v, ok := configHCC["healthy_threshold"].(int); ok && v != 0 {
				healthCheckConfiguration.HealthThreshold = v
			}

			route.HCC = healthCheckConfiguration
		}

		routes = append(routes, route)
	}

	return routes
}

func flattenRoutes(routes []tango.Route) []map[string]interface{} {
	if len(routes) == 0 {
		return make([]map[string]interface{}, 0)
	}

	configRoutes := make([]map[string]interface{}, 0, len(routes))

	for _, route := range routes {
		helper := make(map[string]interface{})
		helper["protocol"] = route.Protocol
		helper["port"] = route.Port
		helper["member_protocol"] = route.MemberProtocol
		helper["member_port"] = route.MemberPort

		if route.HCC != nil {
			hccs := [1]map[string]interface{}{}
			hccs[0] = make(map[string]interface{})
			hccs[0]["protocol"] = route.HCC.Protocol
			hccs[0]["port"] = route.HCC.Port
			hccs[0]["url_path"] = route.HCC.URLPath
			hccs[0]["interval_seconds"] = route.HCC.IntervalSeconds
			hccs[0]["timeout_seconds"] = route.HCC.TimeoutSeconds
			hccs[0]["unhealthy_threshold"] = route.HCC.UnhealthThreshold
			hccs[0]["healthy_threshold"] = route.HCC.HealthThreshold

			helper["health_check_configuration"] = hccs
		}

		configRoutes = append(configRoutes, helper)
	}

	return configRoutes
}

func resourceLoadBalancerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)

	resourceObject, err := client.ReadResource(getSelfLink(d.Get("links").([]interface{})))
	if err != nil {
		d.SetId("")
		return nil
	}

	loadBalancerObject := resourceObject.(*tango.LoadBalancer)

	d.Set("address", loadBalancerObject.Address)
	d.Set("external_zone_id", loadBalancerObject.ExternalZoneID)
	d.Set("external_region_id", loadBalancerObject.ExternalRegionID)
	d.Set("external_id", loadBalancerObject.ExternalID)
	d.Set("name", loadBalancerObject.Name)
	d.Set("description", loadBalancerObject.Description)
	d.Set("self_link", loadBalancerObject.SelfLink)
	d.Set("created_at", loadBalancerObject.CreatedAt)
	d.Set("updated_at", loadBalancerObject.UpdatedAt)
	d.Set("owner", loadBalancerObject.Owner)
	d.Set("organization_id", loadBalancerObject.OrganizationID)
	d.Set("custom_properties", loadBalancerObject.CustomProperties)

	if err := d.Set("tags", flattenTags(loadBalancerObject.Tags)); err != nil {
		return fmt.Errorf("Error setting Load Balancer tags - error: %#v", err)
	}

	if err := d.Set("routes", flattenRoutes(loadBalancerObject.Routes)); err != nil {
		return fmt.Errorf("Error setting Load Balancer routes - error: %#v", err)
	}

	if err := d.Set("links", flattenLinks(loadBalancerObject.Links)); err != nil {
		return fmt.Errorf("Error setting Load Balancer links - error: %#v", err)
	}

	return nil
}

func resourceLoadBalancerUpdate(d *schema.ResourceData, m interface{}) error {

	return nil
}

func resourceLoadBalancerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)
	err := client.DeleteResource(getSelfLink(d.Get("links").([]interface{})))

	if err != nil && strings.Contains(err.Error(), "404") { // already deleted
		return nil
	}

	return err
}
