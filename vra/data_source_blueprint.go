package vra

import (
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/vmware/vra-sdk-go/pkg/client/blueprint"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceBlueprint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBlueprintRead,

		Schema: map[string]*schema.Schema{
			"content": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_source_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_source_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_source_sync_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_source_sync_messages": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"content_source_sync_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_source_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_scope_org": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag to indicate blueprint can be requested from any project in org",
			},
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"total_released_versions": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_versions": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"valid": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"validation_messages": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metadata": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBlueprintRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")
	projectID, projectIDOk := d.GetOk("project_id")

	if !idOk && !nameOk {
		return fmt.Errorf("one of id or name must be assigned")
	}

	var resp *blueprint.ListBlueprintsUsingGET1OK
	var err error
	projects := make([]string, 1)

	if projectIDOk {
		projects = append(projects, projectID.(string))
		resp, err = apiClient.Blueprint.ListBlueprintsUsingGET1(
			blueprint.NewListBlueprintsUsingGET1Params().WithName(withString(name.(string))).WithProjects(projects))
	} else {
		resp, err = apiClient.Blueprint.ListBlueprintsUsingGET1(
			blueprint.NewListBlueprintsUsingGET1Params().WithName(withString(name.(string))))
	}

	if err != nil {
		return err
	}

	if resp.GetPayload().NumberOfElements == 0 {
		return fmt.Errorf("blueprint %s not found", name)
	}

	if resp.GetPayload().NumberOfElements > 1 {
		return fmt.Errorf("more than one blueprint found with the same name, try to narrow filter by project_id")
	}

	setFields := func(blueprint *models.Blueprint) {
		d.SetId(blueprint.ID)
		d.Set("content", blueprint.Content)
		d.Set("content_source_id", blueprint.ContentSourceID)
		d.Set("content_source_pat", blueprint.ContentSourcePath)
		d.Set("content_source_sync_at", blueprint.ContentSourceSyncAt)
		d.Set("content_source_sync_messages", blueprint.ContentSourceSyncMessages)
		d.Set("content_source_sync_status", blueprint.ContentSourceSyncStatus)
		d.Set("content_source_type", blueprint.ContentSourceType)
		d.Set("created_at", blueprint.CreatedAt)
		d.Set("created_by", blueprint.CreatedBy)
		d.Set("description", blueprint.Description)
		d.Set("name", blueprint.Name)
		d.Set("org_id", blueprint.OrgID)
		d.Set("project_id", blueprint.ProjectID)
		d.Set("project_name", blueprint.ProjectName)
		d.Set("request_scope_org", blueprint.RequestScopeOrg)
		d.Set("self_link", blueprint.SelfLink)
		d.Set("status", blueprint.Status)
		d.Set("total_released_versions", blueprint.TotalReleasedVersions)
		d.Set("total_versions", blueprint.TotalVersions)
		d.Set("updated_at", blueprint.UpdatedAt)
		d.Set("updated_by", blueprint.UpdatedBy)
		d.Set("valid", blueprint.Valid)
	}

	for _, bp := range resp.Payload.Content {
		if (idOk && bp.ID == id) || (nameOk && bp.Name == name.(string)) {
			bpDetails, err := apiClient.Blueprint.GetBlueprintUsingGET1(
				blueprint.NewGetBlueprintUsingGET1Params().WithBlueprintID(strfmt.UUID(bp.ID)))
			if err != nil {
				return err
			}
			setFields(bpDetails.GetPayload())

			return nil
		}
	}

	return fmt.Errorf("blueprint %s not found", name)
}
