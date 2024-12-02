// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/compute"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceMachine() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMachineRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"address": {
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
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_properties": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"deployment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
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
			"external_zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": linksSchema(),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"power_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMachineRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading the vra_machine data source with filter %s", d.Get("filter"))
	apiClient := meta.(*Client).apiClient

	var machine *models.Machine

	id := d.Get("id").(string)
	filter := d.Get("filter").(string)

	if id == "" && filter == "" {
		return fmt.Errorf("one of id or filter is required")
	}

	if id != "" {
		getResp, err := apiClient.Compute.GetMachine(compute.NewGetMachineParams().WithID(id))

		if err != nil {
			return err
		}
		machine = getResp.GetPayload()
	} else {
		getResp, err := apiClient.Compute.GetMachines(compute.NewGetMachinesParams().WithDollarFilter(withString(filter)))
		if err != nil {
			return err
		}

		machines := *getResp.Payload
		if len(machines.Content) > 1 {
			return fmt.Errorf("vra_machine must filter to a machine")
		}
		if len(machines.Content) == 0 {
			return fmt.Errorf("vra_machine filter did not match any machine")
		}

		machine = machines.Content[0]
	}

	d.SetId(*machine.ID)
	d.Set("address", machine.Address)
	d.Set("cloud_account_ids", machine.CloudAccountIds)
	d.Set("created_at", machine.CreatedAt)
	d.Set("custom_properties", machine.CustomProperties)
	d.Set("deployment_id", machine.DeploymentID)
	d.Set("description", machine.Description)
	d.Set("external_zone_id", machine.ExternalZoneID)
	d.Set("external_region_id", machine.ExternalRegionID)
	d.Set("external_id", machine.ExternalID)
	d.Set("name", machine.Name)
	d.Set("org_id", machine.OrgID)
	d.Set("owner", machine.Owner)
	d.Set("power_state", machine.PowerState)
	d.Set("project_id", machine.ProjectID)
	d.Set("updated_at", machine.UpdatedAt)

	if err := d.Set("tags", flattenTags(machine.Tags)); err != nil {
		return fmt.Errorf("error setting machine tags - error: %v", err)
	}

	if err := d.Set("links", flattenLinks(machine.Links)); err != nil {
		return fmt.Errorf("error setting machine links - error: %#v", err)
	}

	log.Printf("Finished reading the vra_machine data source with filter %s", d.Get("filter"))
	return nil
}
