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

func dataSourcePolicyLease() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyLeaseRead,

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
					Type: schema.TypeMap,
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
			"lease_grace": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The duration in days that an expired object should be held before it is deleted.",
			},
			"lease_term_max": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum duration in days between creation (or renewal) and expiration.",
			},
			"lease_total_term_max": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum duration in days between creation and expiration. Unaffected by renewal.",
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

func dataSourcePolicyLeaseRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		if *policy.TypeID != PolicyLeaseTypeID {
			return diag.Errorf("policy with id `%s` is not a lease policy", id)
		}
	} else {
		getResp, err := apiClient.Policies.GetPoliciesUsingGET5(policies.NewGetPoliciesUsingGET5Params().WithTypeID(PolicyLeaseTypeID).WithExpandDefinition(withBool(true)).WithSearch(withString(search.(string))))
		if err != nil {
			return diag.FromErr(err)
		}

		policies := getResp.Payload
		if len(policies.Content) == 0 {
			return diag.Errorf("vra_policy_lease `search` criteria did not match any policy")
		}
		if len(policies.Content) > 1 {
			return diag.Errorf("vra_policy_lease `search` criteria must filter to a single policy")
		}

		policy = policies.Content[0]
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

	var definition PolicyLeaseDefinition
	if err := policyDefinitionConvert(policy.Definition, &definition); err != nil {
		return diag.FromErr(err)
	}

	if definition.LeaseGrace != nil {
		_ = d.Set("lease_grace", *definition.LeaseGrace)
	}
	_ = d.Set("lease_term_max", definition.LeaseTermMax)
	_ = d.Set("lease_total_term_max", definition.LeaseTotalTermMax)

	return nil
}
