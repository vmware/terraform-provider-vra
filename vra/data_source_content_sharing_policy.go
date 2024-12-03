// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/policies"
)

func dataSourceContentSharingPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContentSharingPolicyRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Input attributes
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Description:   "The policy ID.",
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Description:   "The policy name.",
				ConflictsWith: []string{"id"},
			},

			// Computed attributes
			"catalog_item_ids": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of catalog item ids to share.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"catalog_source_ids": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of catalog source ids to share.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy creation timestamp.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy author.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The policy description.",
			},
			"last_updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Most recent policy update timestamp.",
			},
			"last_updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Most recent policy editor.",
			},
			"org_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the organization to which the policy belongs.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the project to which the policy belongs.",
			},
		},
	}
}

func dataSourceContentSharingPolicyRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if !idOk && !nameOk {
		return diag.Errorf("one of id or name must be provided")
	}

	if nameOk {
		resp, err := apiClient.Policies.GetPoliciesUsingGET5(policies.NewGetPoliciesUsingGET5Params())
		if err != nil {
			return diag.FromErr(err)
		}
		for _, policy := range resp.GetPayload().Content {
			if policy.Name == name {
				id = policy.ID.String()
				break
			}
		}
		if id == "" {
			return diag.Errorf("content sharing policy with name '%s' not found", name)
		}
	}

	resp, err := apiClient.Policies.GetPolicyUsingGET5(policies.NewGetPolicyUsingGET5Params().WithID(strfmt.UUID(id.(string))))
	if err != nil {
		switch err.(type) {
		case *policies.GetPolicyUsingGET5NotFound:
			return diag.Errorf("content sharing policy with id '%s' not found", id)
		default:
			// nop
		}
		return diag.FromErr(err)
	}

	policy := resp.GetPayload()
	d.SetId(string(policy.ID))
	d.Set("name", policy.Name)
	d.Set("created_at", policy.CreatedAt.String())
	d.Set("created_by", policy.CreatedBy)
	d.Set("description", policy.Description)
	d.Set("last_updated_at", policy.LastUpdatedAt.String())
	d.Set("last_updated_by", policy.LastUpdatedBy)
	d.Set("org_id", policy.OrgID)
	d.Set("project_id", policy.ProjectID)

	catalogItemIDs, err := extractCatalogItemIDsFromContentSharingPolicy(policy.Definition)
	if err != nil {
		return diag.Errorf("error extracting catalog item ids from content sharing policy: %s", err.Error())
	}

	if err := d.Set("catalog_item_ids", catalogItemIDs); err != nil {
		return diag.Errorf("error setting catalog_item_ids: %s", err.Error())
	}

	catalogSourceIDs, err := extractCatalogSourceIDsFromContentSharingPolicy(policy.Definition)
	if err != nil {
		return diag.Errorf("error extracting catalog source ids from content sharing policy: %s", err.Error())
	}

	if err := d.Set("catalog_source_ids", catalogSourceIDs); err != nil {
		return diag.Errorf("error setting catalog_source_ids: %s", err.Error())
	}

	return nil
}
