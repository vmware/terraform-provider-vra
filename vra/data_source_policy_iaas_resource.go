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

func dataSourcePolicyIaaSResource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyIaaSResourceRead,

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
			"exclude_resource_rules": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Exclude Resource Rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_groups": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "List of API groups the resources belong to.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"api_versions": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "List of API Versions the resources belong to.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"operations": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "List of Operations the admission hook cares about.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"resources": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "List of Resources this rule applies to.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"failure_policy": {
				Type:        schema.TypeString,
				Description: "Failure policy to apply when the policy fails.",
				Computed:    true,
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
			"match_conditions": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of conditions that must be met for a request to be validated.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expression": {
							Type:        schema.TypeString,
							Description: "Expression which will be evaluated by CEL.",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Identifier for this match condition.",
							Computed:    true,
						},
					},
				},
			},
			"match_expressions": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of label selector requirements that must be met for an object to be validated.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The label key that the selector applies to.",
						},
						"operator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A key's relationship to a set of values.",
						},
						"values": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "An array of string values.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"match_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Map of {key,value} pairs that must be met for an object to be validated.",
			},
			"match_policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Match policy.",
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
			"resource_rules": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Resource Rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_groups": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "List of API groups the resources belong to.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"api_versions": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "List of API Versions the resources belong to.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"operations": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "List of Operations the admission hook cares about.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"resources": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "List of Resources this rule applies to.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"validation_actions": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of validation actions.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"validations": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of CEL expressions which are used to validate admission requests.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expression": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Expression which will be evaluated by CEL.",
						},
						"message": {
							Type:        schema.TypeString,
							Description: "Message displayed when validation fails.",
							Computed:    true,
						},
						"message_expression": {
							Type:        schema.TypeString,
							Description: "CEL expression that evaluates to the validation failure message that is returned when this rule fails.",
							Computed:    true,
						},
						"reason": {
							Type:        schema.TypeString,
							Description: "Machine-readable description of why this validation failed.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePolicyIaaSResourceRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		if *policy.TypeID != PolicyIaaSResourceTypeID {
			return diag.Errorf("policy with id `%s` is not an IaaS resource policy", id)
		}
	} else {
		getResp, err := apiClient.Policies.GetPoliciesUsingGET5(policies.NewGetPoliciesUsingGET5Params().WithTypeID(PolicyIaaSResourceTypeID).WithExpandDefinition(withBool(true)).WithSearch(withString(search.(string))))
		if err != nil {
			return diag.FromErr(err)
		}

		policies := getResp.Payload
		if len(policies.Content) == 0 {
			return diag.Errorf("vra_policy_iaas_resource `search` criteria did not match any policy")
		}
		if len(policies.Content) > 1 {
			return diag.Errorf("vra_policy_iaas_resource `search` criteria must filter to a single policy")
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

	var definition PolicyIaaSResourceDefinition
	if err := policyDefinitionConvert(policy.Definition, &definition); err != nil {
		return diag.FromErr(err)
	}

	if definition.AutomationPolicy.MatchConstraints.ExcludeResourceRules != nil {
		_ = d.Set("exclude_resource_rules", flattenPolicyIaaSResourceAutomationPolicyExcludeResourceRule(definition.AutomationPolicy.MatchConstraints.ExcludeResourceRules))
	}
	_ = d.Set("failure_policy", definition.AutomationPolicy.FailurePolicy)
	if len(definition.AutomationPolicy.MatchConditions) > 0 {
		_ = d.Set("match_conditions", flattenPolicyIaaSResourceAutomationPolicyMatchConditions(definition.AutomationPolicy.MatchConditions))
	}
	if definition.AutomationPolicy.MatchConstraints.ObjectSelector != nil && definition.AutomationPolicy.MatchConstraints.ObjectSelector.MatchExpressions != nil {
		_ = d.Set("match_expressions", flattenPolicyIaaSResourceAutomationPolicyObjectSelectorMatchExpressions(definition.AutomationPolicy.MatchConstraints.ObjectSelector.MatchExpressions))
	}
	if definition.AutomationPolicy.MatchConstraints.ObjectSelector != nil && definition.AutomationPolicy.MatchConstraints.ObjectSelector.MatchLabels != nil {
		_ = d.Set("match_labels", definition.AutomationPolicy.MatchConstraints.ObjectSelector.MatchLabels)
	}
	if definition.AutomationPolicy.MatchConstraints.MatchPolicy != nil {
		_ = d.Set("match_policy", *definition.AutomationPolicy.MatchConstraints.MatchPolicy)
	}
	_ = d.Set("resource_rules", flattenPolicyIaaSResourceAutomationPolicyResourceRule(definition.AutomationPolicy.MatchConstraints.ResourceRules))
	_ = d.Set("validation_actions", definition.AutomationPolicy.ValidationActions)
	_ = d.Set("validations", flattenPolicyIaaSResourceAutomationPolicyValidations(definition.AutomationPolicy.Validations))

	return nil
}
