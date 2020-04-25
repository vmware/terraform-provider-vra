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
			"administrators": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"members": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"shared_resources": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"zone_assignments": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
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

	setFields := func(project *models.Project) {
		d.SetId(*project.ID)
		d.Set("administrators", flattenUserList(project.Administrators))
		d.Set("description", project.Description)
		d.Set("members", flattenUserList(project.Members))
		d.Set("name", project.Name)
		d.Set("shared_resources", project.SharedResources)
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
