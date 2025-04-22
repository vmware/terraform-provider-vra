// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"encoding/json"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

const (
	PolicyApprovalTypeID     string = "com.vmware.policy.approval"
	PolicyDay2ActionTypeID   string = "com.vmware.policy.deployment.action"
	PolicyIaaSResourceTypeID string = "com.vmware.policy.supervisor.iaas"
	PolicyLeaseTypeID        string = "com.vmware.policy.deployment.lease"
)

type PolicyApprovalDefinition struct {
	Actions              []string `json:"actions"`
	Approvers            []string `json:"approvers"`
	ApprovalLevel        int      `json:"level"`
	ApprovalMode         string   `json:"approvalMode"`
	ApprovalType         string   `json:"approverType"`
	AutoApprovalDecision string   `json:"autoApprovalDecision"`
	AutoApprovalExpiry   int      `json:"autoApprovalExpiry"`
}

type PolicyDay2ActionDefinition struct {
	AllowedActions []PolicyDay2ActionAllowedAction `json:"allowedActions"`
}

type PolicyDay2ActionAllowedAction struct {
	Actions     []string `json:"actions"`
	Authorities []string `json:"authorities"`
}

type PolicyIaaSResourceDefinition struct {
	AutomationPolicy PolicyIaaSResourceAutomationPolicy `json:"automationPolicy"`
}

type PolicyIaaSResourceAutomationPolicy struct {
	FailurePolicy     string                                               `json:"failurePolicy,omitempty"`
	MatchConditions   []*PolicyIaaSResourceAutomationPolicyMatchConditions `json:"matchConditions,omitempty"`
	MatchConstraints  PolicyIaaSResourceAutomationPolicyMatchConstraints   `json:"matchConstraints"`
	ValidationActions []string                                             `json:"validationActions,omitempty"`
	Validations       []PolicyIaaSResourceAutomationPolicyValidations      `json:"validations"`
}

type PolicyIaaSResourceAutomationPolicyMatchConstraints struct {
	ExcludeResourceRules []*PolicyIaaSResourceAutomationPolicyResourceRule `json:"excludeResourceRules,omitempty"`
	MatchPolicy          *string                                           `json:"matchPolicy,omitempty"`
	ObjectSelector       *PolicyIaaSResourceAutomationPolicyObjectSelector `json:"objectSelector,omitempty"`
	ResourceRules        []PolicyIaaSResourceAutomationPolicyResourceRule  `json:"resourceRules"`
}

type PolicyIaaSResourceAutomationPolicyObjectSelector struct {
	MatchExpressions []*PolicyIaaSResourceAutomationPolicyMatchExpressions `json:"matchExpressions,omitempty"`
	MatchLabels      *map[string]string                                    `json:"matchLabels,omitempty"`
}

type PolicyIaaSResourceAutomationPolicyMatchExpressions struct {
	Key      string   `json:"key"`
	Operator string   `json:"operator"`
	Values   []string `json:"values"`
}
type PolicyIaaSResourceAutomationPolicyResourceRule struct {
	APIGroups   []string `json:"apiGroups"`
	APIVersions []string `json:"apiVersions"`
	Operations  []string `json:"operations"`
	Resources   []string `json:"resources"`
}

type PolicyIaaSResourceAutomationPolicyMatchConditions struct {
	Expression string `json:"expression"`
	Name       string `json:"name"`
}
type PolicyIaaSResourceAutomationPolicyValidations struct {
	Expression        string  `json:"expression"`
	Message           *string `json:"message,omitempty"`
	MessageExpression *string `json:"messageExpression,omitempty"`
	Reason            *string `json:"reason,omitempty"`
}
type PolicyLeaseDefinition struct {
	LeaseGrace        *int `json:"leaseGrace,omitempty"`
	LeaseTermMax      int  `json:"leaseTermMax"`
	LeaseTotalTermMax int  `json:"leaseTotalTermMax"`
}

func policyDefinitionConvert(genericDefinition any, castedDefinition any) error {
	definitionJSON, err := json.Marshal(genericDefinition)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(definitionJSON, &castedDefinition); err != nil {
		return err
	}

	return nil
}

func expandPolicyIaaSResourceAutomationPolicyMatchConditions(matchConditionsMap []any) []*PolicyIaaSResourceAutomationPolicyMatchConditions {
	matchConditions := make([]*PolicyIaaSResourceAutomationPolicyMatchConditions, 0, len(matchConditionsMap))

	for _, matchCondition := range matchConditionsMap {
		matchConditionMap := matchCondition.(map[string]any)
		helper := &PolicyIaaSResourceAutomationPolicyMatchConditions{
			Expression: matchConditionMap["expression"].(string),
			Name:       matchConditionMap["name"].(string),
		}
		matchConditions = append(matchConditions, helper)
	}

	return matchConditions
}

func flattenPolicyIaaSResourceAutomationPolicyMatchConditions(matchConditions []*PolicyIaaSResourceAutomationPolicyMatchConditions) []any {
	if len(matchConditions) == 0 {
		return make([]any, 0)
	}

	matchConditionsMap := make([]any, 0, len(matchConditions))

	for _, matchCondition := range matchConditions {
		helper := make(map[string]any)
		helper["expression"] = matchCondition.Expression
		helper["name"] = matchCondition.Name

		matchConditionsMap = append(matchConditionsMap, helper)
	}

	return matchConditionsMap
}

func expandPolicyIaaSResourceAutomationPolicyValidations(validationsMap []any) []PolicyIaaSResourceAutomationPolicyValidations {
	validations := make([]PolicyIaaSResourceAutomationPolicyValidations, 0, len(validationsMap))

	for _, validation := range validationsMap {
		validationMap := validation.(map[string]any)
		helper := PolicyIaaSResourceAutomationPolicyValidations{
			Expression: validationMap["expression"].(string),
		}
		if message, ok := validationMap["message"]; ok {
			helper.Message = withString(message.(string))
		}
		if messageExpression, ok := validationMap["message_expression"]; ok {
			helper.MessageExpression = withString(messageExpression.(string))
		}
		if reason, ok := validationMap["reason"]; ok {
			helper.Reason = withString(reason.(string))
		}
		validations = append(validations, helper)
	}

	return validations
}

func flattenPolicyIaaSResourceAutomationPolicyValidations(validations []PolicyIaaSResourceAutomationPolicyValidations) []any {
	if len(validations) == 0 {
		return make([]any, 0)
	}

	validationsMap := make([]any, 0, len(validations))

	for _, validation := range validations {
		helper := make(map[string]any)
		helper["expression"] = validation.Expression
		if validation.Message != nil {
			helper["message"] = *validation.Message
		}
		if validation.MessageExpression != nil {
			helper["message_expression"] = *validation.MessageExpression
		}
		if validation.Reason != nil {
			helper["reason"] = *validation.Reason
		}

		validationsMap = append(validationsMap, helper)
	}

	return validationsMap
}

func expandPolicyIaaSResourceAutomationPolicyResourceRule(resourceRulesMap []any) []PolicyIaaSResourceAutomationPolicyResourceRule {
	resourceRules := make([]PolicyIaaSResourceAutomationPolicyResourceRule, 0, len(resourceRulesMap))

	for _, resourceRule := range resourceRulesMap {
		resourceRuleMap := resourceRule.(map[string]any)
		helper := PolicyIaaSResourceAutomationPolicyResourceRule{
			APIGroups:   expandStringList(resourceRuleMap["api_groups"].(*schema.Set).List()),
			APIVersions: expandStringList(resourceRuleMap["api_versions"].(*schema.Set).List()),
			Operations:  expandStringList(resourceRuleMap["operations"].(*schema.Set).List()),
			Resources:   expandStringList(resourceRuleMap["resources"].(*schema.Set).List()),
		}

		resourceRules = append(resourceRules, helper)
	}

	return resourceRules
}

func flattenPolicyIaaSResourceAutomationPolicyResourceRule(resourceRules []PolicyIaaSResourceAutomationPolicyResourceRule) []any {
	if len(resourceRules) == 0 {
		return make([]any, 0)
	}

	resourceRulesMap := make([]any, 0, len(resourceRules))

	for _, resourceRule := range resourceRules {
		helper := make(map[string]any)
		helper["api_groups"] = resourceRule.APIGroups
		helper["api_versions"] = resourceRule.APIVersions
		helper["operations"] = resourceRule.Operations
		helper["resources"] = resourceRule.Resources

		resourceRulesMap = append(resourceRulesMap, helper)
	}

	return resourceRulesMap
}

func expandPolicyIaaSResourceAutomationPolicyExcludeResourceRule(resourceRulesMap []any) []*PolicyIaaSResourceAutomationPolicyResourceRule {
	resourceRules := make([]*PolicyIaaSResourceAutomationPolicyResourceRule, 0, len(resourceRulesMap))

	for _, resourceRule := range resourceRulesMap {
		resourceRuleMap := resourceRule.(map[string]interface{})
		helper := &PolicyIaaSResourceAutomationPolicyResourceRule{
			APIGroups:   expandStringList(resourceRuleMap["api_groups"].(*schema.Set).List()),
			APIVersions: expandStringList(resourceRuleMap["api_versions"].(*schema.Set).List()),
			Operations:  expandStringList(resourceRuleMap["operations"].(*schema.Set).List()),
			Resources:   expandStringList(resourceRuleMap["resources"].(*schema.Set).List()),
		}

		resourceRules = append(resourceRules, helper)
	}

	return resourceRules
}

func flattenPolicyIaaSResourceAutomationPolicyExcludeResourceRule(resourceRules []*PolicyIaaSResourceAutomationPolicyResourceRule) []any {
	if len(resourceRules) == 0 {
		return make([]any, 0)
	}

	resourceRulesMap := make([]any, 0, len(resourceRules))

	for _, resourceRule := range resourceRules {
		helper := make(map[string]any)
		helper["api_groups"] = resourceRule.APIGroups
		helper["api_versions"] = resourceRule.APIVersions
		helper["operations"] = resourceRule.Operations
		helper["resources"] = resourceRule.Resources

		resourceRulesMap = append(resourceRulesMap, helper)
	}

	return resourceRulesMap
}

func expandPolicyIaaSResourceAutomationPolicyObjectSelectorMatchExpressions(matchExpressionsMap []any) []*PolicyIaaSResourceAutomationPolicyMatchExpressions {
	matchExpressions := make([]*PolicyIaaSResourceAutomationPolicyMatchExpressions, 0, len(matchExpressionsMap))

	for _, matchExpression := range matchExpressionsMap {
		matchExpressionMap := matchExpression.(map[string]interface{})
		helper := &PolicyIaaSResourceAutomationPolicyMatchExpressions{
			Key:      matchExpressionMap["key"].(string),
			Operator: matchExpressionMap["operator"].(string),
			Values:   expandStringList(matchExpressionMap["values"].(*schema.Set).List()),
		}

		matchExpressions = append(matchExpressions, helper)
	}

	return matchExpressions
}

func flattenPolicyIaaSResourceAutomationPolicyObjectSelectorMatchExpressions(matchExpressions []*PolicyIaaSResourceAutomationPolicyMatchExpressions) []any {
	if len(matchExpressions) == 0 {
		return make([]any, 0)
	}

	matchExpressionsMap := make([]any, 0, len(matchExpressions))

	for _, matchExpression := range matchExpressions {
		helper := make(map[string]any)
		helper["key"] = matchExpression.Key
		helper["operator"] = matchExpression.Operator
		helper["values"] = matchExpression.Values

		matchExpressionsMap = append(matchExpressionsMap, helper)
	}

	return matchExpressionsMap
}
func expandPolicyIaaSResourceAutomationPolicyObjectSelectorMatchLabels(matchLabels map[string]any) *map[string]string {
	if len(matchLabels) == 0 {
		return nil
	}

	matchLabelsMap := make(map[string]string)

	for key, value := range matchLabels {
		matchLabelsMap[key] = value.(string)
	}

	return &matchLabelsMap
}

func expandPolicyCriteria(criteria []any) *models.Criteria {
	if len(criteria) == 0 {
		return nil
	}

	matchExpression := make([]models.Clause, 0)

	for _, clause := range criteria {
		clauseMap := clause.(map[string]any)
		if len(clauseMap) > 0 {
			helper := make(map[string]any)
			for key, value := range clauseMap {
				valueType := reflect.ValueOf(value)
				if valueType.Kind() == reflect.String {
					var valueJSON json.RawMessage
					if json.Unmarshal([]byte(value.(string)), &valueJSON) != nil {
						helper[key] = value
					} else {
						helper[key] = valueJSON
					}
				} else {
					helper[key] = value
				}
			}

			matchExpression = append(matchExpression, helper)
		}
	}

	return &models.Criteria{
		MatchExpression: matchExpression,
	}
}

func flattenPolicyCriteria(criteria models.Criteria) []any {
	criteriaMap := make([]any, 0, len(criteria.MatchExpression))

	for _, expression := range criteria.MatchExpression {
		helper := make(map[string]any)

		for key, value := range expression.(map[string]any) {
			valueType := reflect.ValueOf(value)
			if valueType.Kind() == reflect.String {
				helper[key] = value
			} else if value != nil {
				valueJSON, _ := json.Marshal(value)
				helper[key] = string(valueJSON)
			}
		}

		criteriaMap = append(criteriaMap, helper)
	}

	return criteriaMap
}
