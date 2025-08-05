// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"log"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/blueprint"
)

func dataSourceBlueprintVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBlueprintVersionRead,

		Schema: map[string]*schema.Schema{
			"blueprint_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"blueprint_description": {
				Type:     schema.TypeString,
				Computed: true,
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
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Required: true,
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
				Computed: true,
			},
			"version_change_log": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBlueprintVersionRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("Reading the vra_blueprint_version data source")
	apiClient := m.(*Client).apiClient

	id := d.Get("id").(string)
	bpUUID := strfmt.UUID(d.Get("blueprint_id").(string))

	resp, err := apiClient.Blueprint.GetBlueprintVersionUsingGET1(
		blueprint.NewGetBlueprintVersionUsingGET1Params().
			WithBlueprintID(bpUUID).
			WithVersion(id).
			WithDollarSelect([]string{"*"}))

	if err != nil {
		switch err.(type) {
		case *blueprint.GetBlueprintVersionUsingGET1NotFound:
			return fmt.Errorf("blueprint version '%v' is not found", id)
		}
		return err
	}

	blueprintVersion := *resp.Payload
	d.SetId(blueprintVersion.ID)
	_ = d.Set("blueprint_id", blueprintVersion.BlueprintID)
	_ = d.Set("blueprint_description", blueprintVersion.Description)
	_ = d.Set("change_log", blueprintVersion.VersionChangeLog)
	_ = d.Set("content", blueprintVersion.Content)
	_ = d.Set("created_at", blueprintVersion.CreatedAt)
	_ = d.Set("created_by", blueprintVersion.CreatedBy)
	_ = d.Set("description", blueprintVersion.VersionDescription)
	_ = d.Set("name", blueprintVersion.Name)
	_ = d.Set("org_id", blueprintVersion.OrgID)
	_ = d.Set("project_id", blueprintVersion.ProjectID)
	_ = d.Set("project_name", blueprintVersion.ProjectName)
	_ = d.Set("status", blueprintVersion.Status)
	_ = d.Set("updated_at", blueprintVersion.UpdatedAt)
	_ = d.Set("updated_by", blueprintVersion.UpdatedBy)
	_ = d.Set("valid", blueprintVersion.Valid)
	_ = d.Set("version", blueprintVersion.Version)

	log.Printf("finished reading vra_blueprint_version data source '%v'", id)
	return nil
}
