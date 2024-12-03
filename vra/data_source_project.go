// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"

	"github.com/vmware/vra-sdk-go/pkg/client/project"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceProjectRead,

		Schema: map[string]*schema.Schema{
			"administrators": {
				Type:        schema.TypeSet,
				Computed:    true,
				Optional:    true,
				Deprecated:  "Please use `administrator_roles` instead.",
				Description: "List of administrator users associated with the project. Only administrators can manage project's configuration.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"administrator_roles": userSchema("List of administrator roles associated with the project. Only administrators can manage project's configuration."),
			"constraints": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Description: "List of storage, network and extensibility constraints to be applied when provisioning through this project.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"extensibility": constraintsSchema(),
						"network":       constraintsSchema(),
						"storage":       constraintsSchema(),
					},
				},
			},
			"custom_properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The project custom properties which are added to all requests in this project",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A human-friendly description.",
			},
			"machine_naming_template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The naming template to be used for resources provisioned in this project.",
			},
			"members": {
				Type:        schema.TypeSet,
				Computed:    true,
				Optional:    true,
				Deprecated:  "Please use `member_roles` instead.",
				Description: "List of member users associated with the project.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"member_roles": userSchema("List of member roles associated with the project."),
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A human-friendly name used as an identifier in APIs that support this option.",
			},
			"operation_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The timeout that should be used for Blueprint operations and Provisioning tasks. The timeout is in seconds.",
			},
			"placement_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The placement policy that will be applied when selecting a cloud zone for provisioning.",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The id of this project resource.",
			},
			"shared_resources": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Specifies whether the resources in this projects are shared or not. If not set default will be used.",
			},
			"supervisor_roles": userSchema("List of supervisor roles associated with the project."),
			"viewers": {
				Type:        schema.TypeSet,
				Computed:    true,
				Optional:    true,
				Deprecated:  "Please use `viewer_roles` instead.",
				Description: "List of viewer users associated with the project.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"viewer_roles": userSchema("List of viewer roles associated with the project."),
			"zone_assignments": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of configurations for zone assignment to a project.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							Description: "The maximum amount of cpus that can be used by this cloud zone. Default is 0 (unlimited cpu).",
						},
						"max_instances": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							Description: "The maximum number of instances that can be provisioned in this cloud zone. Default is 0 (unlimited instances)",
						},
						"memory_limit_mb": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							Description: "The maximum amount of memory that can be used by this cloud zone. Default is 0 (unlimited memory).",
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							Description: "The priority of this zone in the current project. Lower numbers mean higher priority. Default is 0 (highest)",
						},
						"storage_limit_gb": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							Description: "Upper limit on storage that can be requested from a cloud zone which is part of this project. Default is 0 (unlimited storage). Supported only for vSphere cloud zones.",
						},
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Cloud Zone Id",
						},
					},
				},
			},
		},
	}
}

func dataSourceProjectRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if !idOk && !nameOk {
		return fmt.Errorf("one of id or name must be provided")
	}

	setFields := func(project *models.IaaSProject) {
		d.SetId(*project.ID)
		d.Set("administrators", flattenUsers(project.Administrators))
		d.Set("administrator_roles", flattenUsers(project.Administrators))
		d.Set("constraints", flattenProjectConstraints(project.Constraints))
		d.Set("custom_properties", project.CustomProperties)
		d.Set("description", project.Description)
		d.Set("machine_naming_template", project.MachineNamingTemplate)
		d.Set("members", flattenUsers(project.Members))
		d.Set("member_roles", flattenUsers(project.Members))
		d.Set("name", project.Name)
		d.Set("operation_timeout", project.OperationTimeout)
		d.Set("placement_policy", project.PlacementPolicy)
		d.Set("shared_resources", project.SharedResources)
		d.Set("supervisor_roles", flattenUsers(project.Supervisors))
		d.Set("viewers", flattenUsers(project.Viewers))
		d.Set("viewer_roles", flattenUsers(project.Viewers))
		d.Set("zone_assignments", flattenZoneAssignment(project.Zones))
	}

	if idOk {
		getResp, err := apiClient.Project.GetProject(project.NewGetProjectParams().WithID(id.(string)))

		if err != nil {
			switch err.(type) {
			case *project.GetProjectNotFound:
				return fmt.Errorf("project %s not found", name)
			default:
				return err
			}
		}

		setFields(getResp.GetPayload())
		return nil
	}

	filter := fmt.Sprintf("name eq '%s'", name)
	getResp, err := apiClient.Project.GetProjects(project.NewGetProjectsParams().WithDollarFilter(withString(filter)))

	if err != nil {
		return err
	}

	projects := getResp.GetPayload()
	if len(projects.Content) > 1 {
		return fmt.Errorf("vra_project must filter to only one project")
	}
	if len(projects.Content) == 0 {
		return fmt.Errorf("project %s not found", name)
	}

	setFields(projects.Content[0])
	return nil
}
