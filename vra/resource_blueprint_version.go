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

func resourceBlueprintVersion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBlueprintVersionCreate,
		ReadContext:   resourceBlueprintVersionRead,
		UpdateContext: resourceBlueprintVersionUpdate,
		DeleteContext: resourceBlueprintVersionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"blueprint_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"blueprint_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"change_log": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"content": {
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
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"release": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceBlueprintVersionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to create vra_blueprint_version resource")
	apiClient := m.(*Client).apiClient

	blueprintVersionRequestSpecification := models.BlueprintVersionRequest{
		ChangeLog: d.Get("change_log").(string),
		Release:   d.Get("release").(bool),
		Version:   withString(d.Get("version").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		blueprintVersionRequestSpecification.Description = v.(string)
	}

	resp, err := apiClient.Blueprint.CreateBlueprintVersionUsingPOST1(
		blueprint.NewCreateBlueprintVersionUsingPOST1Params().
			WithBlueprintID(strfmt.UUID(d.Get("blueprint_id").(string))).
			WithVersionRequest(&blueprintVersionRequestSpecification))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.GetPayload().ID)
	log.Printf("Finished to create vra_blueprint_version resource with blueprint_id %s version %s", d.Get("blueprint_id"), d.Get("version"))

	return resourceBlueprintVersionRead(ctx, d, m)
}

func resourceBlueprintVersionRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Reading the vra_blueprint_version resource with blueprint_id %s and version %s", d.Get("blueprint_id"), d.Get("version"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	bpUUID := strfmt.UUID(d.Get("blueprint_id").(string))

	resp, err := apiClient.Blueprint.GetBlueprintVersionUsingGET1(
		blueprint.NewGetBlueprintVersionUsingGET1Params().
			WithBlueprintID(bpUUID).
			WithVersion(id).
			WithDollarSelect([]string{"*"}))

	if err != nil {
		switch err.(type) {
		case *blueprint.GetBlueprintVersionUsingGET1NotFound:
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	blueprintVersion := *resp.Payload
	d.SetId(blueprintVersion.ID)
	d.Set("blueprint_id", blueprintVersion.BlueprintID)
	d.Set("blueprint_description", blueprintVersion.Description)
	d.Set("change_log", blueprintVersion.VersionChangeLog)
	d.Set("content", blueprintVersion.Content)
	d.Set("created_at", blueprintVersion.CreatedAt)
	d.Set("created_by", blueprintVersion.CreatedBy)
	d.Set("description", blueprintVersion.VersionDescription)
	d.Set("name", blueprintVersion.Name)
	d.Set("org_id", blueprintVersion.OrgID)
	d.Set("project_id", blueprintVersion.ProjectID)
	d.Set("project_name", blueprintVersion.ProjectName)
	d.Set("status", blueprintVersion.Status)
	d.Set("updated_at", blueprintVersion.UpdatedAt)
	d.Set("updated_by", blueprintVersion.UpdatedBy)
	d.Set("valid", blueprintVersion.Valid)
	d.Set("version", blueprintVersion.Version)

	log.Printf("Finished reading the vra_blueprint_version resource with blueprint_id %s and version %s", d.Get("blueprint_id"), d.Get("version"))
	return nil
}

func resourceBlueprintVersionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Starting to update the vra_blueprint_version resource with blueprint_id %s and version %s", d.Get("blueprint_id"), d.Get("version"))
	apiClient := m.(*Client).apiClient

	if d.HasChange("release") {
		bpUUID := strfmt.UUID(d.Get("blueprint_id").(string))
		version := d.Get("version").(string)
		if d.Get("release").(bool) {
			_, err := apiClient.Blueprint.ReleaseBlueprintVersionUsingPOST1(
				blueprint.NewReleaseBlueprintVersionUsingPOST1Params().
					WithBlueprintID(bpUUID).
					WithVersion(version))
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			_, err := apiClient.Blueprint.UnReleaseBlueprintVersionUsingPOST1(
				blueprint.NewUnReleaseBlueprintVersionUsingPOST1Params().
					WithBlueprintID(bpUUID).
					WithVersion(version))
			if err != nil {
				return diag.FromErr(err)
			}
		}
		log.Printf("Finished updating the vra_blueprint resource with name %s", d.Get("name"))
	} else {
		log.Printf("only changes supported on vra_blueprint_version resource are to release flag")
	}

	return resourceBlueprintVersionRead(ctx, d, m)
}

func resourceBlueprintVersionDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	log.Printf("Starting to delete the vra_blueprint_version resource with blueprint_id %s and version %s", d.Get("blueprint_id"), d.Get("version"))
	log.Printf("vra_blueprint_version cannot be deleted in vRA. It can only  be unreleased. Removing local state")
	d.SetId("")
	log.Printf("Finished deleting the vra_blueprint resource with name %s", d.Get("name"))
	return nil
}
