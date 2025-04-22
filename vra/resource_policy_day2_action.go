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
	"github.com/vmware/vra-sdk-go/pkg/client/policies"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func resourcePolicyDay2Action() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyDay2ActionCreate,
		ReadContext:   resourcePolicyDay2ActionRead,
		UpdateContext: resourcePolicyDay2ActionUpdate,
		DeleteContext: resourcePolicyDay2ActionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required arguments
			"authorities": {
				Type:        schema.TypeSet,
				Description: "List of authorities that will be allowed to perform certain actions.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 1,
				Required: true,
			},
			"enforcement_type": {
				Type:         schema.TypeString,
				Description:  "The type of enforcement for the policy.",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"HARD", "SOFT"}, true),
			},
			"name": {
				Type:        schema.TypeString,
				Description: "A human-friendly name used as an identifier for the policy instance.",
				Required:    true,
			},

			// Optional arguments
			"actions": {
				Type:        schema.TypeSet,
				Description: "List of allowed actions for authority/authorities.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"criteria": {
				Type:        schema.TypeSet,
				Description: "The policy criteria.",
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
				Optional: true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "A human-friendly description for the policy instance.",
				Optional:    true,
			},
			"project_criteria": {
				Type:          schema.TypeSet,
				ConflictsWith: []string{"project_id"},
				Description:   "The project based criteria.",
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
				ForceNew: true,
				Optional: true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "The id of the project this entity belongs to.",
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

func resourcePolicyDay2ActionCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	actions := []string{}
	if actionsList, ok := d.GetOk("actions"); ok {
		if !compareUnique(actionsList.(*schema.Set).List()) {
			return diag.Errorf("`actions` must be unique")
		}
		actions = expandStringList(actionsList.(*schema.Set).List())
	}

	authoritiesList := d.Get("authorities").(*schema.Set).List()
	if !compareUnique(authoritiesList) {
		return diag.Errorf("`authorities` must be unique")
	}
	authorities := expandStringList(authoritiesList)

	_, createdResp, err := apiClient.Policies.CreatePolicyUsingPOST1(policies.NewCreatePolicyUsingPOST1Params().WithPolicy(&models.Policy{
		Criteria: expandPolicyCriteria(d.Get("criteria").(*schema.Set).List()),
		Definition: &PolicyDay2ActionDefinition{
			AllowedActions: []PolicyDay2ActionAllowedAction{
				{
					Actions:     actions,
					Authorities: authorities,
				},
			},
		},
		Description:     d.Get("description").(string),
		EnforcementType: d.Get("enforcement_type").(string),
		Name:            d.Get("name").(string),
		ProjectID:       d.Get("project_id").(string),
		ScopeCriteria:   expandPolicyCriteria(d.Get("project_criteria").(*schema.Set).List()),
		TypeID:          withString(PolicyDay2ActionTypeID),
	}))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdResp.Payload.ID.String())

	return resourcePolicyDay2ActionRead(ctx, d, m)
}

func resourcePolicyDay2ActionRead(_ context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	getResp, err := apiClient.Policies.GetPolicyUsingGET5(policies.NewGetPolicyUsingGET5Params().WithID(strfmt.UUID(id)))
	if err != nil {
		switch err.(type) {
		case *policies.GetPolicyUsingGET5NotFound:
			return diag.Errorf("policy with id `%s` not found", id)
		default:
			// nop
		}
		return diag.FromErr(err)
	}

	policy := getResp.GetPayload()
	if *policy.TypeID != PolicyDay2ActionTypeID {
		return diag.Errorf("policy with id `%s` is not a day2 action policy", id)
	}

	d.SetId(policy.ID.String())
	d.Set("created_at", policy.CreatedAt.String())
	d.Set("created_by", policy.CreatedBy)
	d.Set("enforcement_type", policy.EnforcementType)
	d.Set("last_updated_at", policy.LastUpdatedAt.String())
	d.Set("last_updated_by", policy.LastUpdatedBy)
	d.Set("name", policy.Name)
	d.Set("org_id", policy.OrgID)

	if policy.Criteria != nil {
		d.Set("criteria", flattenPolicyCriteria(*policy.Criteria))
	}
	if policy.Description != "" {
		d.Set("description", policy.Description)
	}
	if policy.ScopeCriteria != nil {
		d.Set("project_criteria", flattenPolicyCriteria(*policy.ScopeCriteria))
	}
	if policy.ProjectID != "" {
		d.Set("project_id", policy.ProjectID)
	}

	var definition PolicyDay2ActionDefinition
	if err := policyDefinitionConvert(policy.Definition, &definition); err != nil {
		return diag.FromErr(err)
	}

	if len(definition.AllowedActions) > 0 {
		if len(definition.AllowedActions[0].Actions) > 0 {
			d.Set("actions", definition.AllowedActions[0].Actions)
		}
		d.Set("authorities", definition.AllowedActions[0].Authorities)
	}

	return nil
}

func resourcePolicyDay2ActionUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	actions := []string{}
	if actionsList, ok := d.GetOk("actions"); ok {
		if !compareUnique(actionsList.(*schema.Set).List()) {
			return diag.Errorf("`actions` must be unique")
		}
		actions = expandStringList(actionsList.(*schema.Set).List())
	}

	authoritiesList := d.Get("authorities").(*schema.Set).List()
	if !compareUnique(authoritiesList) {
		return diag.Errorf("`authorities` must be unique")
	}
	authorities := expandStringList(authoritiesList)

	_, _, err := apiClient.Policies.CreatePolicyUsingPOST1(policies.NewCreatePolicyUsingPOST1Params().WithPolicy(&models.Policy{
		Criteria: expandPolicyCriteria(d.Get("criteria").(*schema.Set).List()),
		Definition: &PolicyDay2ActionDefinition{
			AllowedActions: []PolicyDay2ActionAllowedAction{
				{
					Actions:     actions,
					Authorities: authorities,
				},
			},
		},
		Description:     d.Get("description").(string),
		EnforcementType: d.Get("enforcement_type").(string),
		ID:              strfmt.UUID(id),
		Name:            d.Get("name").(string),
		ProjectID:       d.Get("project_id").(string),
		ScopeCriteria:   expandPolicyCriteria(d.Get("project_criteria").(*schema.Set).List()),
		TypeID:          withString(PolicyDay2ActionTypeID),
	}))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePolicyDay2ActionRead(ctx, d, m)
}

func resourcePolicyDay2ActionDelete(_ context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	if _, err := apiClient.Policies.DeletePolicyUsingDELETE5(policies.NewDeletePolicyUsingDELETE5Params().WithID(strfmt.UUID(id))); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
