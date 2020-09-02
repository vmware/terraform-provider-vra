package vra

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/project"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"administrators": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of administrator users associated with the project. Only administrators can manage project's configuration.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
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
			"members": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of member users associated with the project.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
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
			"viewers": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of viewer users associated with the project.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
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

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	administrators := expandUserList(d.Get("administrators").(*schema.Set).List())
	constraints := expandProjectConstraints(d.Get("constraints").(*schema.Set).List())
	description := d.Get("description").(string)
	machineNamingTemplate := d.Get("machine_naming_template").(string)
	members := expandUserList(d.Get("members").(*schema.Set).List())
	name := d.Get("name").(string)
	operationTimeout := d.Get("operation_timeout").(int)
	sharedResources := d.Get("shared_resources").(bool)
	viewers := expandUserList(d.Get("viewers").(*schema.Set).List())
	zoneAssignment := expandZoneAssignment(d.Get("zone_assignments").(*schema.Set).List())

	createResp, err := apiClient.Project.CreateProject(project.NewCreateProjectParams().WithBody(&models.ProjectSpecification{
		Administrators:               administrators,
		Constraints:                  constraints,
		Description:                  description,
		MachineNamingTemplate:        machineNamingTemplate,
		Members:                      members,
		Name:                         &name,
		OperationTimeout:             int64(operationTimeout),
		SharedResources:              withBool(sharedResources),
		Viewers:                      viewers,
		ZoneAssignmentConfigurations: zoneAssignment,
	}))
	if err != nil {
		return err
	}

	d.SetId(*createResp.Payload.ID)

	return resourceProjectRead(d, m)
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	ret, err := apiClient.Project.GetProject(project.NewGetProjectParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *project.GetProjectNotFound:
			d.SetId("")
			return nil
		}
		return err
	}
	project := *ret.Payload
	d.Set("administrators", flattenUserList(project.Administrators))
	d.Set("constraints", flattenProjectConstraints(project.Constraints))
	d.Set("description", project.Description)
	d.Set("machine_naming_template", project.MachineNamingTemplate)
	d.Set("members", flattenUserList(project.Members))
	d.Set("name", project.Name)
	d.Set("operation_timeout", project.OperationTimeout)
	d.Set("shared_resources", project.SharedResources)
	d.Set("viewers", flattenUserList(project.Viewers))
	d.Set("zone_assignments", flattenZoneAssignment(project.Zones))

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	administrators := expandUserList(d.Get("administrators").(*schema.Set).List())
	constraints := expandProjectConstraints(d.Get("constraints").(*schema.Set).List())
	description := d.Get("description").(string)
	machineNamingTemplate := d.Get("machine_naming_template").(string)
	members := expandUserList(d.Get("members").(*schema.Set).List())
	viewers := expandUserList(d.Get("viewers").(*schema.Set).List())
	name := d.Get("name").(string)
	operationTimeout := d.Get("operation_timeout").(int)
	sharedResources := d.Get("shared_resources").(bool)
	zoneAssignment := expandZoneAssignment(d.Get("zone_assignments").(*schema.Set).List())

	_, err := apiClient.Project.UpdateProject(project.NewUpdateProjectParams().WithID(id).WithBody(&models.ProjectSpecification{
		Administrators:               administrators,
		Constraints:                  constraints,
		Description:                  description,
		MachineNamingTemplate:        machineNamingTemplate,
		Members:                      members,
		Name:                         &name,
		OperationTimeout:             int64(operationTimeout),
		SharedResources:              withBool(sharedResources),
		Viewers:                      viewers,
		ZoneAssignmentConfigurations: zoneAssignment,
	}))
	if err != nil {
		return err
	}

	return resourceProjectRead(d, m)
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	// Workaround an issue where the cloud regions need to be removed before the project can be deleted.
	_, err := apiClient.Project.UpdateProject(project.NewUpdateProjectParams().WithID(id).WithBody(&models.ProjectSpecification{
		ZoneAssignmentConfigurations: []*models.ZoneAssignmentConfig{},
	}))
	if err != nil {
		return err
	}

	_, err = apiClient.Project.DeleteProject(project.NewDeleteProjectParams().WithID(id))
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func expandUserList(userList []interface{}) []*models.User {
	users := make([]*models.User, 0, len(userList))

	for _, email := range userList {
		user := models.User{
			Email: withString(email.(string)),
		}
		users = append(users, &user)
	}

	return users
}

func flattenUserList(userList []*models.User) []*string {
	result := make([]*string, 0, len(userList))

	for _, user := range userList {
		result = append(result, user.Email)
	}

	return result
}

func expandZoneAssignment(configZoneAssignments []interface{}) []*models.ZoneAssignmentConfig {
	zoneAssignments := make([]*models.ZoneAssignmentConfig, 0, len(configZoneAssignments))

	for _, configZone := range configZoneAssignments {
		configZoneAssignment := configZone.(map[string]interface{})

		za := models.ZoneAssignmentConfig{
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

func flattenZoneAssignment(list []*models.ZoneAssignmentConfig) []map[string]interface{} {
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
	projectConstraints := make(map[string][]models.Constraint)

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
