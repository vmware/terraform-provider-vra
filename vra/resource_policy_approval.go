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

func resourcePolicyApproval() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyApprovalCreate,
		ReadContext:   resourcePolicyApprovalRead,
		UpdateContext: resourcePolicyApprovalUpdate,
		DeleteContext: resourcePolicyApprovalDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required arguments
			"actions": {
				Type:        schema.TypeSet,
				Description: "List of actions to trigger approval.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 1,
				Required: true,
			},
			"approval_level": {
				Type:         schema.TypeInt,
				Description:  "The level defines the order in which the policy is enforced. Level 1 approvals are applied first, followed by level 2 approvals, and so on.",
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 99),
			},
			"approval_mode": {
				Type:         schema.TypeString,
				Description:  "Who must approve the request.",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ANY_OF", "ALL_OF"}, true),
			},
			"approval_type": {
				Type:         schema.TypeString,
				Description:  "Approval Type.",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"USER", "ROLE"}, true),
			},
			"approvers": {
				Type:        schema.TypeSet,
				Description: "List of approvers of the policy.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 1,
				Required: true,
			},
			"auto_approval_decision": {
				Type:         schema.TypeString,
				Description:  "Automatically approve or reject a request after the number of days specified in the Auto expiry trigger field.",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"APPROVE", "REJECT", "NO_EXPIRY"}, true),
			},
			"auto_approval_expiry": {
				Type:         schema.TypeInt,
				Description:  "The number of days the approvers have, to respond before the Auto action is triggered.",
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 30),
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

func resourcePolicyApprovalCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	actionsList := d.Get("actions").(*schema.Set).List()
	if !compareUnique(actionsList) {
		return diag.Errorf("`actions` must be unique")
	}
	actions := expandStringList(actionsList)

	approversList := d.Get("approvers").(*schema.Set).List()
	if !compareUnique(approversList) {
		return diag.Errorf("`approvers` must be unique")
	}
	approvers := expandStringList(approversList)

	_, createdResp, err := apiClient.Policies.CreatePolicyUsingPOST1(policies.NewCreatePolicyUsingPOST1Params().WithPolicy(&models.Policy{
		Criteria: expandPolicyCriteria(d.Get("criteria").(*schema.Set).List()),
		Definition: &PolicyApprovalDefinition{
			Actions:              actions,
			ApprovalLevel:        d.Get("approval_level").(int),
			ApprovalMode:         d.Get("approval_mode").(string),
			ApprovalType:         d.Get("approval_type").(string),
			Approvers:            approvers,
			AutoApprovalExpiry:   d.Get("auto_approval_expiry").(int),
			AutoApprovalDecision: d.Get("auto_approval_decision").(string),
		},
		Description:     d.Get("description").(string),
		EnforcementType: d.Get("enforcement_type").(string),
		Name:            d.Get("name").(string),
		ProjectID:       d.Get("project_id").(string),
		ScopeCriteria:   expandPolicyCriteria(d.Get("project_criteria").(*schema.Set).List()),
		TypeID:          withString(PolicyApprovalTypeID),
	}))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdResp.Payload.ID.String())

	return resourcePolicyApprovalRead(ctx, d, m)
}

func resourcePolicyApprovalRead(_ context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
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
	if *policy.TypeID != PolicyApprovalTypeID {
		return diag.Errorf("policy with id `%s` is not an approval policy", id)
	}

	d.SetId(policy.ID.String())
	_ = d.Set("created_at", policy.CreatedAt.String())
	_ = d.Set("created_by", policy.CreatedBy)
	_ = d.Set("enforcement_type", policy.EnforcementType)
	_ = d.Set("last_updated_at", policy.LastUpdatedAt.String())
	_ = d.Set("last_updated_by", policy.LastUpdatedBy)
	_ = d.Set("name", policy.Name)
	_ = d.Set("org_id", policy.OrgID)

	if policy.Criteria != nil {
		_ = d.Set("criteria", flattenPolicyCriteria(*policy.Criteria))
	}
	if policy.Description != "" {
		_ = d.Set("description", policy.Description)
	}
	if policy.ScopeCriteria != nil {
		_ = d.Set("project_criteria", flattenPolicyCriteria(*policy.ScopeCriteria))
	}
	if policy.ProjectID != "" {
		_ = d.Set("project_id", policy.ProjectID)
	}

	var definition PolicyApprovalDefinition
	if err := policyDefinitionConvert(policy.Definition, &definition); err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("actions", definition.Actions)
	_ = d.Set("approval_level", definition.ApprovalLevel)
	_ = d.Set("approval_mode", definition.ApprovalMode)
	_ = d.Set("approval_type", definition.ApprovalType)
	_ = d.Set("approvers", definition.Approvers)
	_ = d.Set("auto_approval_decision", definition.AutoApprovalDecision)
	_ = d.Set("auto_approval_expiry", definition.AutoApprovalExpiry)

	return nil
}

func resourcePolicyApprovalUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	actionsList := d.Get("actions").(*schema.Set).List()
	if !compareUnique(actionsList) {
		return diag.Errorf("`actions` must be unique")
	}
	actions := expandStringList(actionsList)

	approversList := d.Get("approvers").(*schema.Set).List()
	if !compareUnique(approversList) {
		return diag.Errorf("`approvers` must be unique")
	}
	approvers := expandStringList(approversList)

	_, _, err := apiClient.Policies.CreatePolicyUsingPOST1(policies.NewCreatePolicyUsingPOST1Params().WithPolicy(&models.Policy{
		Criteria: expandPolicyCriteria(d.Get("criteria").(*schema.Set).List()),
		Definition: &PolicyApprovalDefinition{
			Actions:              actions,
			ApprovalLevel:        d.Get("approval_level").(int),
			ApprovalMode:         d.Get("approval_mode").(string),
			ApprovalType:         d.Get("approval_type").(string),
			Approvers:            approvers,
			AutoApprovalExpiry:   d.Get("auto_approval_expiry").(int),
			AutoApprovalDecision: d.Get("auto_approval_decision").(string),
		},
		Description:     d.Get("description").(string),
		EnforcementType: d.Get("enforcement_type").(string),
		ID:              strfmt.UUID(id),
		Name:            d.Get("name").(string),
		ProjectID:       d.Get("project_id").(string),
		ScopeCriteria:   expandPolicyCriteria(d.Get("project_criteria").(*schema.Set).List()),
		TypeID:          withString(PolicyApprovalTypeID),
	}))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePolicyApprovalRead(ctx, d, m)
}

func resourcePolicyApprovalDelete(_ context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	if _, err := apiClient.Policies.DeletePolicyUsingDELETE5(policies.NewDeletePolicyUsingDELETE5Params().WithID(strfmt.UUID(id))); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
