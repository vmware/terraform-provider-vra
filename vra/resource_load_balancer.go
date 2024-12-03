// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/vmware/vra-sdk-go/pkg/client"
	"github.com/vmware/vra-sdk-go/pkg/client/load_balancer"
	"github.com/vmware/vra-sdk-go/pkg/client/request"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoadBalancerCreate,
		ReadContext:   resourceLoadBalancerRead,
		UpdateContext: resourceLoadBalancerUpdate,
		DeleteContext: resourceLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					return !strings.HasPrefix(new, old)
				},
			},
			"nics": nicsSchema(true),
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"routes": routesSchema(true),
			"custom_properties": {
				Type:     schema.TypeMap,
				Computed: true,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"deployment_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"internet_facing": {
				Type:     schema.TypeBool,
				Optional: true,
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

func resourceLoadBalancerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	if v, ok := d.GetOk("deployment_id"); ok {
		loadBalancerSpecification.DeploymentID = v.(string)
	}

	if v, ok := d.GetOk("internet_facing"); ok {
		loadBalancerSpecification.InternetFacing = v.(bool)
	}

	if _, ok := d.GetOk("targets"); ok {
		loadBalancerSpecification.TargetLinks = expandLoadBalancerTargets(d.Get("targets").(*schema.Set).List())
	}

	log.Printf("[DEBUG] create load lalancer: %#v", loadBalancerSpecification)
	createLoadBalancerCreated, err := apiClient.LoadBalancer.CreateLoadBalancer(load_balancer.NewCreateLoadBalancerParams().WithBody(&loadBalancerSpecification))
	if err != nil {
		return diag.FromErr(err)
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    loadBalancerStateRefreshFunc(*apiClient, *createLoadBalancerCreated.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
	}

	resourceIDs, err := stateChangeFunc.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	loadBalancerIDs := resourceIDs.([]string)
	i := strings.LastIndex(loadBalancerIDs[0], "/")
	loadBalancerID := loadBalancerIDs[0][i+1 : len(loadBalancerIDs[0])]
	d.SetId(loadBalancerID)
	log.Printf("Finished to create vra_load_balancer resource with name %s", d.Get("name"))

	return resourceLoadBalancerRead(ctx, d, m)
}

func loadBalancerStateRefreshFunc(apiClient client.API, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ret, err := apiClient.Request.GetRequestTracker(request.NewGetRequestTrackerParams().WithID(id))
		if err != nil {
			return "", models.RequestTrackerStatusFAILED, err
		}

		status := ret.Payload.Status
		switch *status {
		case models.RequestTrackerStatusFAILED:
			return []string{""}, *status, errors.New(ret.Payload.Message)
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

func resourceLoadBalancerRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		return diag.FromErr(err)
	}

	loadBalancer := *resp.Payload
	d.Set("address", loadBalancer.Address)
	d.Set("created_at", loadBalancer.CreatedAt)
	d.Set("custom_properties", loadBalancer.CustomProperties)
	d.Set("description", loadBalancer.Description)
	d.Set("deployment_id", loadBalancer.DeploymentID)
	d.Set("external_id", loadBalancer.ExternalID)
	d.Set("external_region_id", loadBalancer.ExternalRegionID)
	d.Set("external_zone_id", loadBalancer.ExternalZoneID)
	d.Set("name", loadBalancer.Name)
	d.Set("organization_id", loadBalancer.OrgID)
	d.Set("owner", loadBalancer.Owner)
	d.Set("project_id", loadBalancer.ProjectID)
	d.Set("updated_at", loadBalancer.UpdatedAt)

	if err := d.Set("tags", flattenTags(loadBalancer.Tags)); err != nil {
		return diag.Errorf("error setting load balancer tags - error: %v", err)
	}
	if err := d.Set("routes", flattenRoutes(loadBalancer.Routes)); err != nil {
		return diag.Errorf("error setting load balancer routes - error: %v", err)
	}

	if err := d.Set("links", flattenLinks(loadBalancer.Links)); err != nil {
		return diag.Errorf("error setting load balancer links - error: %#v", err)
	}

	log.Printf("Finished reading the vra_load_balancer resource with name %s", d.Get("name"))
	return nil
}

func resourceLoadBalancerUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return diag.FromErr(errors.New("Updating a load balancer resource is not allowed"))
}

func resourceLoadBalancerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to delete the vra_load_balancer resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	deleteLoadBalancer, err := apiClient.LoadBalancer.DeleteLoadBalancer(load_balancer.NewDeleteLoadBalancerParams().WithID(id))
	if err != nil {
		return diag.FromErr(err)
	}
	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    loadBalancerStateRefreshFunc(*apiClient, *deleteLoadBalancer.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 5 * time.Second,
	}

	if _, err = stateChangeFunc.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("Finished deleting the vra_load_balancer resource with name %s", d.Get("name"))
	return nil
}
