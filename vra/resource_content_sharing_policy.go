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
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func resourceContentSharingPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContentSharingPolicyCreate,
		ReadContext:   resourceContentSharingPolicyRead,
		UpdateContext: resourceContentSharingPolicyUpdate,
		DeleteContext: resourceContentSharingPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Required arguments
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The policy name.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the project to which the policy belongs.",
			},

			// Optional arguments
			"catalog_item_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of catalog item ids to share.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				AtLeastOneOf: []string{"catalog_source_ids"},
			},
			"catalog_source_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of catalog source ids to share.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				AtLeastOneOf: []string{"catalog_item_ids"},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The policy description.",
			},

			// Computed attributes
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
			"last_updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Most recent policy update timestamp.",
			},
			"last_updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Most recent policy editor..",
			},
			"org_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the organization to which the policy belongs.",
			},
		},
	}
}

func resourceContentSharingPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	catalogItemsIDs := []string{}
	if v, ok := d.GetOk("catalog_item_ids"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.Errorf("catalog_item_ids are not unique")
		}
		catalogItemsIDs = expandStringList(v.(*schema.Set).List())
	}
	catalogSourceIDs := []string{}
	if v, ok := d.GetOk("catalog_source_ids"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.Errorf("catalog_source_ids are not unique")
		}
		catalogSourceIDs = expandStringList(v.(*schema.Set).List())
	}
	definition := buildContentSharingPolicyDefinition(catalogItemsIDs, catalogSourceIDs, d.Get("project_id").(string))
	policy := models.Policy{
		Definition:      definition,
		Description:     d.Get("description").(string),
		EnforcementType: EnforcementTypeHard,
		Name:            d.Get("name").(string),
		ProjectID:       d.Get("project_id").(string),
		TypeID:          withString(CatalogEntitlementTypeID),
	}
	_, createResp, err := apiClient.Policies.DryRunPolicyUsingPOST2(policies.NewDryRunPolicyUsingPOST2Params().WithPolicy(&policy))
	if err != nil {
		return diag.FromErr(err)
	}

	id := createResp.GetPayload().ID.String()
	d.SetId(id)

	return resourceContentSharingPolicyRead(ctx, d, m)
}

func resourceContentSharingPolicyRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	resp, err := apiClient.Policies.GetPolicyUsingGET5(policies.NewGetPolicyUsingGET5Params().WithID(strfmt.UUID(id)))
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

func resourceContentSharingPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	catalogItemsIDs := []string{}
	if v, ok := d.GetOk("catalog_item_ids"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.Errorf("catalog_item_ids are not unique")
		}
		catalogItemsIDs = expandStringList(v.(*schema.Set).List())
	}
	catalogSourceIDs := []string{}
	if v, ok := d.GetOk("catalog_source_ids"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.Errorf("catalog_source_ids are not unique")
		}
		catalogSourceIDs = expandStringList(v.(*schema.Set).List())
	}
	definition := buildContentSharingPolicyDefinition(catalogItemsIDs, catalogSourceIDs, d.Get("project_id").(string))
	policy := models.Policy{
		Definition:      definition,
		Description:     d.Get("description").(string),
		EnforcementType: EnforcementTypeHard,
		ID:              strfmt.UUID(id),
		Name:            d.Get("name").(string),
		ProjectID:       d.Get("project_id").(string),
		TypeID:          withString(CatalogEntitlementTypeID),
	}
	if _, _, err := apiClient.Policies.DryRunPolicyUsingPOST2(policies.NewDryRunPolicyUsingPOST2Params().WithPolicy(&policy)); err != nil {
		return diag.FromErr(err)
	}

	return resourceContentSharingPolicyRead(ctx, d, m)
}

func resourceContentSharingPolicyDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	if _, err := apiClient.Policies.DeletePolicyUsingDELETE5(policies.NewDeletePolicyUsingDELETE5Params().WithID(strfmt.UUID(id))); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
