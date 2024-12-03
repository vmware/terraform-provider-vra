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

	"github.com/vmware/vra-sdk-go/pkg/client/network"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client"
	"github.com/vmware/vra-sdk-go/pkg/client/request"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkCreate,
		ReadContext:   resourceNetworkRead,
		UpdateContext: resourceNetworkUpdate,
		DeleteContext: resourceNetworkDelete,
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
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"constraints": constraintsSchema(),
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
			"outbound_access": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": tagsSchema(),
			"cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_id": {
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
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to create vra_network resource")
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

	if v, ok := d.GetOk("deployment_id"); ok {
		networkSpecification.DeploymentID = v.(string)
	}

	if v, ok := d.GetOk("outbound_access"); ok {
		networkSpecification.OutboundAccess = v.(bool)
	}
	log.Printf("[DEBUG] create network: %#v", networkSpecification)
	createNetworkCreated, err := apiClient.Network.CreateNetwork(network.NewCreateNetworkParams().WithBody(&networkSpecification))
	if err != nil {
		return diag.FromErr(err)
	}
	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    networkStateRefreshFunc(*apiClient, *createNetworkCreated.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
	}
	resourceIDs, err := stateChangeFunc.WaitForStateContext(ctx)
	log.Printf("Waitforstate returned: %T %+v %+v\n", resourceIDs, resourceIDs, err)

	if err != nil {
		return diag.FromErr(err)
	}

	networkIDs := resourceIDs.([]string)
	d.SetId(networkIDs[0])
	log.Printf("Finished to create vra_network resource with name %s", d.Get("name"))

	return resourceNetworkRead(ctx, d, m)
}

func networkStateRefreshFunc(apiClient client.API, id string) retry.StateRefreshFunc {
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

func resourceNetworkRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Reading the vra_network resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	resp, err := apiClient.Network.GetNetwork(network.NewGetNetworkParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *network.GetNetworkNotFound:
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	network := *resp.Payload
	d.Set("cidr", network.Cidr)
	d.Set("custom_properties", network.CustomProperties)
	d.Set("description", network.Description)
	d.Set("deployment_id", network.DeploymentID)
	d.Set("external_id", network.ExternalID)
	d.Set("external_zone_id", network.ExternalZoneID)
	d.Set("name", network.Name)
	d.Set("organization_id", network.OrgID)
	d.Set("owner", network.Owner)
	d.Set("project_id", network.ProjectID)
	d.Set("updated_at", network.UpdatedAt)

	if err := d.Set("tags", flattenTags(network.Tags)); err != nil {
		return diag.Errorf("error setting network tags - error: %v", err)
	}

	if err := d.Set("links", flattenLinks(network.Links)); err != nil {
		return diag.Errorf("error setting network links - error: %#v", err)
	}

	log.Printf("Finished reading the vra_network resource with name %s", d.Get("name"))
	return nil
}

func resourceNetworkUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return diag.FromErr(errors.New("Updating a network resource is not allowed"))
}

func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to delete the vra_network resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	deleteNetworkAccepted, err := apiClient.Network.DeleteNetwork(network.NewDeleteNetworkParams().WithID(id))
	if err != nil {
		return diag.FromErr(err)
	}
	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    networkStateRefreshFunc(*apiClient, *deleteNetworkAccepted.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 5 * time.Second,
	}

	if _, err = stateChangeFunc.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("Finished deleting the vra_network resource with name %s", d.Get("name"))
	return nil
}
