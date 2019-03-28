package cas

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/cas-sdk-go/pkg/client/project"
	"github.com/vmware/cas-sdk-go/pkg/models"

	tango "github.com/vmware/terraform-provider-cas/sdk"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"administrators": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"members": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"zone_assignments": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"max_instances": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)
	apiClient := client.GetAPIClient()

	administrators := expandUserList(d.Get("administrators").(*schema.Set).List())
	description := d.Get("description").(string)
	members := expandUserList(d.Get("members").(*schema.Set).List())
	name := d.Get("name").(string)
	zoneAssignment := expandZoneAssignment(d.Get("zone_assignments").(*schema.Set).List())

	createResp, err := apiClient.Project.CreateProject(project.NewCreateProjectParams().WithBody(&models.ProjectSpecification{
		Administrators:               administrators,
		Description:                  description,
		Members:                      members,
		Name:                         &name,
		ZoneAssignmentConfigurations: zoneAssignment,
	}))
	if err != nil {
		return err
	}

	d.SetId(*createResp.Payload.ID)

	return resourceProjectRead(d, m)
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)
	apiClient := client.GetAPIClient()

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
	d.Set("members", flattenUserList(Project.Members))
	d.Set("name", Project.Name)
	d.Set("zone_assignments", flattenZoneAssignment(Project.Zones))

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)
	apiClient := client.GetAPIClient()

	id := d.Id()
	administrators := expandUserList(d.Get("administrators").(*schema.Set).List())
	description := d.Get("description").(string)
	members := expandUserList(d.Get("members").(*schema.Set).List())
	name := d.Get("name").(string)
	zoneAssignment := expandZoneAssignment(d.Get("zone_assignments").(*schema.Set).List())

	_, err := apiClient.Project.UpdateProject(project.NewUpdateProjectParams().WithID(id).WithBody(&models.ProjectSpecification{
		Administrators:               administrators,
		Description:                  description,
		Members:                      members,
		Name:                         &name,
		ZoneAssignmentConfigurations: zoneAssignment,
	}))
	if err != nil {
		return err
	}

	return resourceProjectRead(d, m)
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*tango.Client)
	apiClient := client.GetAPIClient()

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
			MaxNumberInstances: int64(configZoneAssignment["max_instances"].(int)),
			Priority:           int32(configZoneAssignment["priority"].(int)),
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
			"max_instances": zoneAssignment.MaxNumberInstances,
			"priority":      zoneAssignment.Priority,
			"zone_id":       zoneAssignment.ZoneID,
		}

		result = append(result, l)
	}
	return result
}
