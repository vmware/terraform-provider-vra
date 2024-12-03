// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"

	"github.com/vmware/vra-sdk-go/pkg/client/blueprint"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"log"
)

func resourceBlueprint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBlueprintCreate,
		ReadContext:   resourceBlueprintRead,
		UpdateContext: resourceBlueprintUpdate,
		DeleteContext: resourceBlueprintDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"content": {
				Type:     schema.TypeString,
				Optional: true,
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
				Type:     schema.TypeSet,
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
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_scope_org": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
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

func resourceBlueprintCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to create vra_blueprint resource")
	apiClient := m.(*Client).apiClient

	blueprintSpecification := models.Blueprint{
		Content:         d.Get("content").(string),
		Name:            d.Get("name").(string),
		ProjectID:       d.Get("project_id").(string),
		RequestScopeOrg: d.Get("request_scope_org").(bool),
	}

	if v, ok := d.GetOk("description"); ok {
		blueprintSpecification.Description = v.(string)
	}

	resp, err := apiClient.Blueprint.CreateBlueprintUsingPOST1(blueprint.NewCreateBlueprintUsingPOST1Params().WithBlueprint(&blueprintSpecification))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.GetPayload().ID)
	log.Printf("Finished to create vra_blueprint resource with name %s", d.Get("name"))

	return resourceBlueprintRead(ctx, d, m)
}

func resourceBlueprintRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Reading the vra_blueprint resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	bpUUID := strfmt.UUID(id)

	resp, err := apiClient.Blueprint.GetBlueprintUsingGET1(blueprint.NewGetBlueprintUsingGET1Params().WithBlueprintID(bpUUID))

	if err != nil {
		switch err.(type) {
		case *blueprint.GetBlueprintUsingGET1NotFound:
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	blueprint := *resp.Payload
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

	if err := d.Set("validation_messages", flattenValidationMessages(blueprint.ValidationMessages)); err != nil {
		return diag.Errorf("error setting validation_messages in blueprint - error: %#v", err)
	}

	log.Printf("Finished reading the vra_blueprint resource with name %s", d.Get("name"))
	return nil
}

func resourceBlueprintUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to update the vra_blueprint resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	bpUUID := strfmt.UUID(id)
	blueprintSpecification := models.Blueprint{
		Content:         d.Get("content").(string),
		Description:     d.Get("description").(string),
		Name:            d.Get("name").(string),
		ProjectID:       d.Get("project_id").(string),
		RequestScopeOrg: d.Get("request_scope_org").(bool),
	}

	_, err := apiClient.Blueprint.UpdateBlueprintUsingPUT1(
		blueprint.NewUpdateBlueprintUsingPUT1Params().WithBlueprintID(bpUUID).WithBlueprint(&blueprintSpecification))

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("Finished updating the vra_blueprint resource with name %s", d.Get("name"))
	return resourceBlueprintRead(ctx, d, m)
}

func resourceBlueprintDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to delete the vra_blueprint resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	bpUUID := strfmt.UUID(id)
	_, err := apiClient.Blueprint.DeleteBlueprintUsingDELETE1(
		blueprint.NewDeleteBlueprintUsingDELETE1Params().WithBlueprintID(bpUUID))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("Finished deleting the vra_blueprint resource with name %s", d.Get("name"))
	return nil
}

func flattenValidationMessages(blueprintValidationMessages []*models.BlueprintValidationMessage) []map[string]interface{} {
	if len(blueprintValidationMessages) == 0 {
		return make([]map[string]interface{}, 0)
	}

	validationMsgs := make([]map[string]interface{}, 0, len(blueprintValidationMessages))

	for _, validationMsg := range blueprintValidationMessages {
		helper := make(map[string]interface{})
		helper["message"] = validationMsg.Message
		helper["metadata"] = validationMsg.Metadata
		helper["path"] = validationMsg.Path
		helper["resource_name"] = validationMsg.ResourceName
		helper["type"] = validationMsg.Type

		validationMsgs = append(validationMsgs, helper)
	}

	return validationMsgs
}
