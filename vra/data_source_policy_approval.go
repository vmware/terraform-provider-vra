// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/vmware/vra-sdk-go/pkg/client/policies"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePolicyApproval() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyApprovalRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Computed:      true,
				ConflictsWith: []string{"search"},
				Description:   "The id of the policy instance.",
				Optional:      true,
			},
			"search": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"id"},
				Description:   "Search criteria to narrow down the policy instance.",
				Optional:      true,
			},
			"actions": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of actions to trigger approval.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"approval_level": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The level defines the order in which the policy is enforced. Level 1 approvals are applied first, followed by level 2 approvals, and so on.",
			},
			"approval_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Who must approve the request.",
			},
			"approval_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Approval Type.",
			},
			"approvers": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of approvers of the policy.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"auto_approval_decision": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Automatically approve or reject a request after the number of days specified in the Auto expiry trigger field.",
			},
			"auto_approval_expiry": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of days the approvers have to respond before the Auto action is triggered.",
			},
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
			"criteria": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "The policy criteria.",
				Elem: &schema.Schema{
					Type:        schema.TypeMap,
					Description: "The policy criteria clause.",
				},
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly description for the policy instance.",
			},
			"enforcement_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of enforcement for the policy.",
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
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly name used as an identifier for the policy instance.",
			},
			"org_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the organization this entity belongs to.",
			},
			"project_criteria": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "The project based criteria.",
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the project this entity belongs to.",
			},
		},
	}
}

func dataSourcePolicyApprovalRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(*Client).apiClient

	id, idOk := d.GetOk("id")
	search, searchOk := d.GetOk("search")
	if !idOk && !searchOk {
		return diag.Errorf("one of `id` or `search` is required")
	}

	var policy *models.Policy
	if id != "" {
		getResp, err := apiClient.Policies.GetPolicyUsingGET5(policies.NewGetPolicyUsingGET5Params().WithID(strfmt.UUID(id.(string))))
		if err != nil {
			switch err.(type) {
			case *policies.GetPolicyUsingGET5NotFound:
				return diag.Errorf("policy with id `%s` not found", id)
			default:
				// nop
			}
			return diag.FromErr(err)
		}

		policy = getResp.GetPayload()
		if *policy.TypeID != PolicyApprovalTypeID {
			return diag.Errorf("policy with id `%s` is not an approval policy", id)
		}
	} else {
		getResp, err := apiClient.Policies.GetPoliciesUsingGET5(policies.NewGetPoliciesUsingGET5Params().WithTypeID(PolicyApprovalTypeID).WithExpandDefinition(withBool(true)).WithSearch(withString(search.(string))))
		if err != nil {
			return diag.FromErr(err)
		}

		policies := getResp.Payload
		if len(policies.Content) == 0 {
			return diag.Errorf("vra_policy_approval `search` criteria did not match any policy")
		}
		if len(policies.Content) > 1 {
			return diag.Errorf("vra_policy_approval `search` criteria must filter to a single policy")
		}

		policy = policies.Content[0]
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

	var definition PolicyApprovalDefinition
	if err := policyDefinitionConvert(policy.Definition, &definition); err != nil {
		return diag.FromErr(err)
	}

	d.Set("actions", definition.Actions)
	d.Set("approval_level", definition.ApprovalLevel)
	d.Set("approval_mode", definition.ApprovalMode)
	d.Set("approval_type", definition.ApprovalType)
	d.Set("approvers", definition.Approvers)
	d.Set("auto_approval_decision", definition.AutoApprovalDecision)
	d.Set("auto_approval_expiry", definition.AutoApprovalExpiry)

	return nil
}
