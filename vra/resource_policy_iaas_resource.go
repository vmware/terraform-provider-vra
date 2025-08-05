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

func resourcePolicyIaaSResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyIaaSResourceCreate,
		ReadContext:   resourcePolicyIaaSResourceRead,
		UpdateContext: resourcePolicyIaaSResourceUpdate,
		DeleteContext: resourcePolicyIaaSResourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"enforcement_type": {
				Type:         schema.TypeString,
				Description:  "The type of enforcement for the policy.",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"HARD", "SOFT"}, true),
			},
			"failure_policy": {
				Type:         schema.TypeString,
				Description:  "Failure policy to apply when the policy fails.",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Fail", "Ignore"}, true),
			},
			"name": {
				Type:        schema.TypeString,
				Description: "A human-friendly name used as an identifier for the policy instance.",
				Required:    true,
			},
			"resource_rules": {
				Type:        schema.TypeSet,
				Description: "Resource Rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_groups": {
							Type:        schema.TypeSet,
							Description: "List of API groups the resources belong to.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							MinItems: 1,
							Required: true,
						},
						"api_versions": {
							Type:        schema.TypeSet,
							Description: "List of API Versions the resources belong to.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							MinItems: 1,
							Required: true,
						},
						"operations": {
							Type:        schema.TypeSet,
							Description: "List of Operations the admission hook cares about.",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"CREATE", "UPDATE", "DELETE"}, true),
							},
							Required: true,
						},
						"resources": {
							Type:        schema.TypeSet,
							Description: "List of Resources this rule applies to.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							MinItems: 1,
							Required: true,
						},
					},
				},
				Required: true,
			},
			"validation_actions": {
				Type:        schema.TypeSet,
				Description: "List of validation actions.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"Deny", "Warn"}, true),
				},
				MinItems: 1,
				Required: true,
			},
			"validations": {
				Type:        schema.TypeSet,
				Description: "List of CEL expressions which are used to validate admission requests.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expression": {
							Type:        schema.TypeString,
							Description: "Expression which will be evaluated by CEL.",
							Required:    true,
						},
						"message": {
							Type:        schema.TypeString,
							Description: "Message displayed when validation fails.",
							Optional:    true,
						},
						"message_expression": {
							Type:        schema.TypeString,
							Description: "CEL expression that evaluates to the validation failure message that is returned when this rule fails.",
							Optional:    true,
						},
						"reason": {
							Type:        schema.TypeString,
							Description: "Machine-readable description of why this validation failed.",
							Optional:    true,
						},
					},
				},
				Required: true,
				MinItems: 1,
				MaxItems: 3,
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
			"exclude_resource_rules": {
				Type:        schema.TypeSet,
				Description: "Exclude Resource Rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_groups": {
							Type:        schema.TypeSet,
							Description: "List of API groups the resources belong to.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							MinItems: 1,
							Required: true,
						},
						"api_versions": {
							Type:        schema.TypeSet,
							Description: "List of API Versions the resources belong to.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							MinItems: 1,
							Required: true,
						},
						"operations": {
							Type:        schema.TypeSet,
							Description: "List of Operations the admission hook cares about.",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"CREATE", "UPDATE", "DELETE"}, true),
							},
							Required: true,
						},
						"resources": {
							Type:        schema.TypeSet,
							Description: "List of Resources this rule applies to.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							MinItems: 1,
							Required: true,
						},
					},
				},
				Optional: true,
			},
			"match_conditions": {
				Type:        schema.TypeSet,
				Description: "List of conditions that must be met for a request to be validated.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expression": {
							Type:        schema.TypeString,
							Description: "Expression which will be evaluated by CEL.",
							Required:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Identifier for this match condition.",
							Required:    true,
						},
					},
				},
				MaxItems: 3,
				Optional: true,
			},
			"match_expressions": {
				Type:        schema.TypeSet,
				Description: "List of label selector requirements that must be met for an object to be validated.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Description: "The label key that the selector applies to.",
							Required:    true,
						},
						"operator": {
							Type:         schema.TypeString,
							Description:  "A key's relationship to a set of values.",
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"In", "NotIn", "Exists", "DoesNotExist"}, true),
						},
						"values": {
							Type:        schema.TypeSet,
							Description: "An array of string values.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required: true,
						},
					},
				},
				Optional: true,
			},
			"match_labels": {
				Type:        schema.TypeMap,
				Description: "Map of {key,value} pairs that must be met for an object to be validated.",
				Optional:    true,
			},
			"match_policy": {
				Type:         schema.TypeString,
				Description:  "Match policy.",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Exact", "Equivalent"}, true),
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

func resourcePolicyIaaSResourceCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	validationActionsList := d.Get("validation_actions").(*schema.Set).List()
	if !compareUnique(validationActionsList) {
		return diag.Errorf("`validation_actions` must be unique")
	}
	validationActions := expandStringList(validationActionsList)

	definition := &PolicyIaaSResourceDefinition{
		AutomationPolicy: PolicyIaaSResourceAutomationPolicy{
			FailurePolicy:   d.Get("failure_policy").(string),
			MatchConditions: expandPolicyIaaSResourceAutomationPolicyMatchConditions(d.Get("match_conditions").(*schema.Set).List()),
			MatchConstraints: PolicyIaaSResourceAutomationPolicyMatchConstraints{
				ExcludeResourceRules: expandPolicyIaaSResourceAutomationPolicyExcludeResourceRule(d.Get("exclude_resource_rules").(*schema.Set).List()),
				MatchPolicy:          withString(d.Get("match_policy").(string)),
				ResourceRules:        expandPolicyIaaSResourceAutomationPolicyResourceRule(d.Get("resource_rules").(*schema.Set).List()),
			},
			ValidationActions: validationActions,
			Validations:       expandPolicyIaaSResourceAutomationPolicyValidations(d.Get("validations").(*schema.Set).List()),
		},
	}

	_, matchExpressionsOK := d.GetOk("match_expressions")
	_, matchLabelsOK := d.GetOk("match_labels")
	if matchExpressionsOK || matchLabelsOK {
		definition.AutomationPolicy.MatchConstraints.ObjectSelector = &PolicyIaaSResourceAutomationPolicyObjectSelector{
			MatchExpressions: expandPolicyIaaSResourceAutomationPolicyObjectSelectorMatchExpressions(d.Get("match_expressions").(*schema.Set).List()),
			MatchLabels:      expandPolicyIaaSResourceAutomationPolicyObjectSelectorMatchLabels(d.Get("match_labels").(map[string]any)),
		}
	}

	_, createdResp, err := apiClient.Policies.CreatePolicyUsingPOST1(policies.NewCreatePolicyUsingPOST1Params().WithPolicy(&models.Policy{
		Criteria:        expandPolicyCriteria(d.Get("criteria").(*schema.Set).List()),
		Definition:      definition,
		Description:     d.Get("description").(string),
		EnforcementType: d.Get("enforcement_type").(string),
		Name:            d.Get("name").(string),
		ProjectID:       d.Get("project_id").(string),
		ScopeCriteria:   expandPolicyCriteria(d.Get("project_criteria").(*schema.Set).List()),
		TypeID:          withString(PolicyIaaSResourceTypeID),
	}))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdResp.Payload.ID.String())

	return resourcePolicyIaaSResourceRead(ctx, d, m)
}

func resourcePolicyIaaSResourceRead(_ context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
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
	if *policy.TypeID != PolicyIaaSResourceTypeID {
		return diag.Errorf("policy with id `%s` is not a day2 action policy", id)
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

func resourcePolicyIaaSResourceUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	validationActionsList := d.Get("validation_actions").(*schema.Set).List()
	if !compareUnique(validationActionsList) {
		return diag.Errorf("`validation_actions` must be unique")
	}
	validationActions := expandStringList(validationActionsList)

	definition := &PolicyIaaSResourceDefinition{
		AutomationPolicy: PolicyIaaSResourceAutomationPolicy{
			FailurePolicy:   d.Get("failure_policy").(string),
			MatchConditions: expandPolicyIaaSResourceAutomationPolicyMatchConditions(d.Get("match_conditions").(*schema.Set).List()),
			MatchConstraints: PolicyIaaSResourceAutomationPolicyMatchConstraints{
				ExcludeResourceRules: expandPolicyIaaSResourceAutomationPolicyExcludeResourceRule(d.Get("exclude_resource_rules").(*schema.Set).List()),
				MatchPolicy:          withString(d.Get("match_policy").(string)),
				ResourceRules:        expandPolicyIaaSResourceAutomationPolicyResourceRule(d.Get("resource_rules").(*schema.Set).List()),
			},
			ValidationActions: validationActions,
			Validations:       expandPolicyIaaSResourceAutomationPolicyValidations(d.Get("validations").(*schema.Set).List()),
		},
	}

	_, matchExpressionsOK := d.GetOk("match_expressions")
	_, matchLabelsOK := d.GetOk("match_labels")
	if matchExpressionsOK || matchLabelsOK {
		definition.AutomationPolicy.MatchConstraints.ObjectSelector = &PolicyIaaSResourceAutomationPolicyObjectSelector{
			MatchExpressions: expandPolicyIaaSResourceAutomationPolicyObjectSelectorMatchExpressions(d.Get("match_expressions").(*schema.Set).List()),
			MatchLabels:      expandPolicyIaaSResourceAutomationPolicyObjectSelectorMatchLabels(d.Get("match_labels").(map[string]any)),
		}
	}

	_, _, err := apiClient.Policies.CreatePolicyUsingPOST1(policies.NewCreatePolicyUsingPOST1Params().WithPolicy(&models.Policy{
		Criteria:        expandPolicyCriteria(d.Get("criteria").(*schema.Set).List()),
		Definition:      definition,
		Description:     d.Get("description").(string),
		EnforcementType: d.Get("enforcement_type").(string),
		ID:              strfmt.UUID(id),
		Name:            d.Get("name").(string),
		ProjectID:       d.Get("project_id").(string),
		ScopeCriteria:   expandPolicyCriteria(d.Get("project_criteria").(*schema.Set).List()),
		TypeID:          withString(PolicyIaaSResourceTypeID),
	}))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePolicyIaaSResourceRead(ctx, d, m)
}

func resourcePolicyIaaSResourceDelete(_ context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	if _, err := apiClient.Policies.DeletePolicyUsingDELETE5(policies.NewDeletePolicyUsingDELETE5Params().WithID(strfmt.UUID(id))); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
