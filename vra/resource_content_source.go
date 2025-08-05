// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/vmware/vra-sdk-go/pkg/client/content_source"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func resourceContentSource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContentSourceCreate,
		ReadContext:   resourceContentSourceRead,
		DeleteContext: resourceContentSourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required arguments
			"config": {
				Type:        schema.TypeSet,
				Description: "The content source custom configuration.",
				ForceNew:    true,
				MaxItems:    1,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"branch": {
							Type:        schema.TypeString,
							Description: "The content source branch name.",
							Required:    true,
						},
						"content_type": {
							Type:         schema.TypeString,
							Description:  "The content source type.",
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"BLUEPRINT", "IMAGE", "ABX_SCRIPTS", "TERRAFORM_CONFIGURATION"}, true),
						},
						"integration_id": {
							Type:        schema.TypeString,
							Description: "The content source integration id as seen integrations.",
							Required:    true,
						},
						"path": {
							Type:        schema.TypeString,
							Description: "Path to refer to in the content source repository and branch.",
							Optional:    true,
						},
						"project_name": {
							Type:        schema.TypeString,
							Description: "The name of the project.",
							Computed:    true,
						},
						"repository": {
							Type:        schema.TypeString,
							Description: "The content source repository.",
							Required:    true,
						},
					},
				},
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the content source instance.",
				ForceNew:    true,
				Required:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "The id of the project this entity belongs to.",
				ForceNew:    true,
				Required:    true,
			},
			"sync_enabled": {
				Type:        schema.TypeBool,
				Description: "Wether or not sync is enabled for this content source.",
				ForceNew:    true,
				Required:    true,
			},
			"type_id": {
				Type:         schema.TypeString,
				Description:  "The type of this content source.",
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"com.gitlab", "com.github", "org.bitbucket"}, true),
			},

			// Optional arguments
			"description": {
				Type:        schema.TypeString,
				Description: "A human-friendly description for the catalog source instance.",
				ForceNew:    true,
				Optional:    true,
			},

			// Computed attributes
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was created. The date is in ISO 8601 and UTC.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user the entity was created by.",
			},
			"last_updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC.",
			},
			"last_updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user the entity was last updated by.",
			},
			"org_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the organization this entity belongs to.",
			},
		},
	}
}

func resourceContentSourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	config := expandContentSourceRepositoryConfig(d.Get("config").(*schema.Set).List())

	contentSourceSpecification := models.ContentSource{
		Config:      config[0],
		Description: d.Get("description").(string),
		Name:        withString(d.Get("name").(string)),
		TypeID:      withString(d.Get("type_id").(string)),
		SyncEnabled: d.Get("sync_enabled").(bool),
		ProjectID:   withString(d.Get("project_id").(string)),
	}
	resp, err := apiClient.ContentSource.CreateContentSourceUsingPOST(content_source.NewCreateContentSourceUsingPOSTParams().WithSource(&contentSourceSpecification))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.GetPayload().ID.String())

	return resourceContentSourceRead(ctx, d, m)
}

func resourceContentSourceRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	resp, err := apiClient.ContentSource.GetContentSourceUsingGET(content_source.NewGetContentSourceUsingGETParams().WithID(strfmt.UUID(id)))
	if err != nil {
		switch err.(type) {
		case *content_source.GetContentSourceUsingGETNotFound:
			return diag.Errorf("content source with id `%s` not found", id)
		default:
			// nop
		}
		return diag.FromErr(err)
	}

	contentSource := *resp.Payload
	d.SetId(contentSource.ID.String())
	_ = d.Set("created_at", contentSource.CreatedAt.String())
	_ = d.Set("created_by", contentSource.CreatedBy)
	_ = d.Set("description", contentSource.Description)
	_ = d.Set("name", contentSource.Name)
	_ = d.Set("last_updated_at", contentSource.LastUpdatedAt.String())
	_ = d.Set("last_updated_by", contentSource.LastUpdatedBy)
	_ = d.Set("org_id", contentSource.OrgID)
	_ = d.Set("project_id", contentSource.ProjectID)
	_ = d.Set("sync_enabled", contentSource.SyncEnabled)
	_ = d.Set("type_id", contentSource.TypeID)

	config, err := flattenContentsourceRepositoryConfig(contentSource.Config)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("config", config)

	return nil
}

func resourceContentSourceDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	if _, err := apiClient.ContentSource.DeleteContentSourceUsingDELETE(content_source.NewDeleteContentSourceUsingDELETEParams().WithID(strfmt.UUID(id))); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
