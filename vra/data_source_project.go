package vra

import (
	"fmt"

	"github.com/vmware/vra-sdk-go/pkg/client/project"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceProjectRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

	setFields := func(project *models.Project) {
		d.SetId(*project.ID)
		d.Set("description", project.Description)
		d.Set("name", project.Name)
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
	} else {
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
}
