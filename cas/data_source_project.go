package cas

import (
	"fmt"

	"github.com/vmware/cas-sdk-go/pkg/client/project"
	"github.com/vmware/cas-sdk-go/pkg/models"

	"github.com/hashicorp/terraform/helper/schema"
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

	if idOk == false && nameOk == false {
		return fmt.Errorf("One of id or name must be assigned")
	}

	getResp, err := apiClient.Project.GetProjects(project.NewGetProjectsParams())
	if err != nil {
		return err
	}

	setFields := func(project *models.Project) {
		d.SetId(*project.ID)
		d.Set("description", project.Description)
		d.Set("name", project.Name)
	}
	for _, project := range getResp.Payload.Content {
		if idOk && project.ID == id {
			setFields(project)
			return nil
		}
		if nameOk && project.Name == name {
			setFields(project)
			return nil
		}
	}

	return fmt.Errorf("project %s not found", name)
}
