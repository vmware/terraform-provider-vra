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
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"machine_naming_template": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"members": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"operation_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"shared_resources": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"viewers": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"zone_assignments": {
				Type:     schema.TypeSet,
				Optional: true,
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
	Project := *ret.Payload
	d.Set("administrators", flattenUserList(Project.Administrators))
	d.Set("description", Project.Description)
	d.Set("machine_naming_template", Project.MachineNamingTemplate)
	d.Set("members", flattenUserList(Project.Members))
	d.Set("name", Project.Name)
	d.Set("operation_timeout", Project.OperationTimeout)
	d.Set("shared_resources", Project.SharedResources)
	d.Set("viewers", flattenUserList(Project.Viewers))
	d.Set("zone_assignments", flattenZoneAssignment(Project.Zones))

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	administrators := expandUserList(d.Get("administrators").(*schema.Set).List())
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
