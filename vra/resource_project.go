// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/project"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"administrators": {
				Type:        schema.TypeSet,
				Optional:    true,
				Deprecated:  "To specify the type of principal, please use `administrator_roles` instead.",
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
				Description: "A human-friendly description.",
			},
			"machine_naming_template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The naming template to be used for resources provisioned in this project.",
			},
			"placement_policy": {
				Type:        schema.TypeString,
				Default:     "DEFAULT",
				Description: "The placement policy that will be applied when selecting a cloud zone for provisioning.",
				Optional:    true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "DEFAULT" && value != "SPREAD" {
						errors = append(errors, fmt.Errorf(
							"%q must be one of 'DEFAULT', 'SPREAD'", k))
					}
					return
				},
			},
			"members": {
				Type:        schema.TypeSet,
				Optional:    true,
				Deprecated:  "To specify the type of principal, please use `member_roles` instaed.",
				Description: "List of member users associated with the project.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"member_roles": userSchema("List of member roles associated with the project."),
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A human-friendly name used as an identifier in APIs that support this option.",
			},
			"operation_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The timeout that should be used for Blueprint operations and Provisioning tasks. The timeout is in seconds.",
			},
			"shared_resources": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Specifies whether the resources in this projects are shared or not. If not set default will be used.",
			},
			"supervisor_roles": userSchema("List of supervisor roles associated with the project."),
			"viewers": {
				Type:        schema.TypeSet,
				Optional:    true,
				Deprecated:  "To specify the type of principal, please use `viewer_roles` instead.",
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
							Optional:    true,
							Description: "The maximum amount of cpus that can be used by this cloud zone. Default is 0 (unlimited cpu).",
						},
						"max_instances": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The maximum number of instances that can be provisioned in this cloud zone. Default is 0 (unlimited instances)",
						},
						"memory_limit_mb": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The maximum amount of memory that can be used by this cloud zone. Default is 0 (unlimited memory).",
						},
						"priority": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The priority of this zone in the current project. Lower numbers mean higher priority. Default is 0 (highest)",
						},
						"storage_limit_gb": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Upper limit on storage that can be requested from a cloud zone which is part of this project. Default is 0 (unlimited storage). Supported only for vSphere cloud zones.",
						},
						"zone_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The Cloud Zone Id",
						},
					},
				},
			},
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	administrators := expandUserListAndNewUserList(d.Get("administrators").(*schema.Set).List(),
		d.Get("administrator_roles").(*schema.Set).List())
	constraints := expandProjectConstraints(d.Get("constraints").(*schema.Set).List())
	customProperties := expandCustomProperties(d.Get("custom_properties").(map[string]interface{}))
	description := d.Get("description").(string)
	machineNamingTemplate := d.Get("machine_naming_template").(string)
	members := expandUserListAndNewUserList(d.Get("members").(*schema.Set).List(), d.Get("member_roles").(*schema.Set).List())
	name := d.Get("name").(string)
	operationTimeout := int64(d.Get("operation_timeout").(int))
	placementPolicy := d.Get("placement_policy").(string)
	sharedResources := d.Get("shared_resources").(bool)
	supervisors := expandUsers(d.Get("supervisor_roles").(*schema.Set).List())
	viewers := expandUserListAndNewUserList(d.Get("viewers").(*schema.Set).List(), d.Get("viewer_roles").(*schema.Set).List())
	zoneAssignment := expandZoneAssignment(d.Get("zone_assignments").(*schema.Set).List())

	createResp, err := apiClient.Project.CreateProject(project.NewCreateProjectParams().WithBody(&models.IaaSProjectSpecification{
		Administrators:               administrators,
		Constraints:                  constraints,
		CustomProperties:             customProperties,
		Description:                  description,
		MachineNamingTemplate:        machineNamingTemplate,
		Members:                      members,
		Name:                         &name,
		OperationTimeout:             &operationTimeout,
		PlacementPolicy:              placementPolicy,
		SharedResources:              withBool(sharedResources),
		Supervisors:                  supervisors,
		Viewers:                      viewers,
		ZoneAssignmentConfigurations: zoneAssignment,
	}))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*createResp.Payload.ID)

	return resourceProjectRead(ctx, d, m)
}

func resourceProjectRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	ret, err := apiClient.Project.GetProject(project.NewGetProjectParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *project.GetProjectNotFound:
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	project := *ret.Payload
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

	return nil
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	administrators := expandUserListAndNewUserList(d.Get("administrators").(*schema.Set).List(),
		d.Get("administrator_roles").(*schema.Set).List())
	constraints := expandProjectConstraints(d.Get("constraints").(*schema.Set).List())
	customProperties := expandCustomProperties(d.Get("custom_properties").(map[string]interface{}))
	description := d.Get("description").(string)
	machineNamingTemplate := d.Get("machine_naming_template").(string)
	members := expandUserListAndNewUserList(d.Get("members").(*schema.Set).List(), d.Get("member_roles").(*schema.Set).List())
	supervisors := expandUsers(d.Get("supervisor_roles").(*schema.Set).List())
	viewers := expandUserListAndNewUserList(d.Get("viewers").(*schema.Set).List(), d.Get("viewer_roles").(*schema.Set).List())
	name := d.Get("name").(string)
	operationTimeout := int64(d.Get("operation_timeout").(int))
	placementPolicy := d.Get("placement_policy").(string)
	sharedResources := d.Get("shared_resources").(bool)
	zoneAssignment := expandZoneAssignment(d.Get("zone_assignments").(*schema.Set).List())

	_, err := apiClient.Project.UpdateProject(project.NewUpdateProjectParams().WithID(id).WithBody(&models.IaaSProjectSpecification{
		Administrators:               administrators,
		Constraints:                  constraints,
		CustomProperties:             customProperties,
		Description:                  description,
		MachineNamingTemplate:        machineNamingTemplate,
		Members:                      members,
		Name:                         &name,
		OperationTimeout:             &operationTimeout,
		PlacementPolicy:              placementPolicy,
		SharedResources:              withBool(sharedResources),
		Supervisors:                  supervisors,
		Viewers:                      viewers,
		ZoneAssignmentConfigurations: zoneAssignment,
	}))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceProjectRead(ctx, d, m)
}

func resourceProjectDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	// Workaround an issue where the cloud regions need to be removed before the project can be deleted.
	_, err := apiClient.Project.UpdateProject(project.NewUpdateProjectParams().WithID(id).WithBody(&models.IaaSProjectSpecification{
		ZoneAssignmentConfigurations: []*models.ZoneAssignmentSpecification{},
	}))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = apiClient.Project.DeleteProject(project.NewDeleteProjectParams().WithID(id))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func expandUserListAndNewUserList(userList []interface{}, configUsers []interface{}) []*models.User {
	users := make([]*models.User, 0, len(userList)+len(configUsers))

	for _, email := range userList {
		user := models.User{
			Email: withString(email.(string)),
		}
		users = append(users, &user)
	}

	users = append(users, expandUsers(configUsers)...)

	return users
}

func expandZoneAssignment(configZoneAssignments []interface{}) []*models.ZoneAssignmentSpecification {
	zoneAssignments := make([]*models.ZoneAssignmentSpecification, 0, len(configZoneAssignments))

	for _, configZone := range configZoneAssignments {
		configZoneAssignment := configZone.(map[string]interface{})

		za := models.ZoneAssignmentSpecification{
			CPULimit:           int64(configZoneAssignment["cpu_limit"].(int)),
			MaxNumberInstances: int64(configZoneAssignment["max_instances"].(int)),
			MemoryLimitMB:      int64(configZoneAssignment["memory_limit_mb"].(int)),
			Priority:           int32(configZoneAssignment["priority"].(int)),
			StorageLimitGB:     int64(configZoneAssignment["storage_limit_gb"].(int)),
			ZoneID:             configZoneAssignment["zone_id"].(string),
		}

		zoneAssignments = append(zoneAssignments, &za)
	}

	return zoneAssignments
}

func flattenZoneAssignment(list []*models.ZoneAssignment) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, zoneAssignment := range list {
		l := map[string]interface{}{
			"cpu_limit":        zoneAssignment.CPULimit,
			"max_instances":    zoneAssignment.MaxNumberInstances,
			"memory_limit_mb":  zoneAssignment.MemoryLimitMB,
			"priority":         zoneAssignment.Priority,
			"storage_limit_gb": zoneAssignment.StorageLimitGB,
			"zone_id":          zoneAssignment.ZoneID,
		}

		result = append(result, l)
	}
	return result
}

func expandProjectConstraints(configProjectConstraints []interface{}) map[string][]models.Constraint {
	projectConstraints := map[string][]models.Constraint{
		"extensibility": {},
		"network":       {},
		"storage":       {},
	}

	for _, configProjectConstraint := range configProjectConstraints {
		configConstraints := configProjectConstraint.(map[string]interface{})

		if v, ok := configConstraints["extensibility"]; ok {
			projectConstraints["extensibility"] = expandConstraintsForProject(v.(*schema.Set).List())
		}

		if v, ok := configConstraints["network"]; ok {
			projectConstraints["network"] = expandConstraintsForProject(v.(*schema.Set).List())
		}

		if v, ok := configConstraints["storage"]; ok {
			projectConstraints["storage"] = expandConstraintsForProject(v.(*schema.Set).List())
		}
	}

	return projectConstraints
}

func flattenProjectConstraints(projectConstraints map[string][]models.Constraint) []map[string]interface{} {
	if len(projectConstraints) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, 1)

	helper := make(map[string]interface{})
	if v, ok := projectConstraints["extensibility"]; ok && len(v) > 0 {
		helper["extensibility"] = flattenConstraints(v)
	}

	if v, ok := projectConstraints["network"]; ok && len(v) > 0 {
		helper["network"] = flattenConstraints(v)
	}

	if v, ok := projectConstraints["storage"]; ok && len(v) > 0 {
		helper["storage"] = flattenConstraints(v)
	}

	result = append(result, helper)
	return result
}
