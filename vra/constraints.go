// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constraintsSchema returns the schema to use for the constraints property
func constraintsSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Constraints that are used to drive placement policies for entities such as image, network, storage, etc. Constraint expressions are matched against tags on existing placement targets.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"mandatory": {
					Type:        schema.TypeBool,
					Required:    true,
					Description: "Indicates whether this constraint should be strictly enforced or not.",
				},
				"expression": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "An expression of the form \"[!]tag-key[:[tag-value]]\", used to indicate a constraint match on keys and values of tags.",
				},
			},
		},
	}
}

func expandConstraints(configConstraints []interface{}) []*models.Constraint {
	constraints := make([]*models.Constraint, 0, len(configConstraints))

	for _, configConstraint := range configConstraints {
		constraintMap := configConstraint.(map[string]interface{})

		constraint := models.Constraint{
			Mandatory:  withBool(constraintMap["mandatory"].(bool)),
			Expression: withString(constraintMap["expression"].(string)),
		}

		constraints = append(constraints, &constraint)
	}

	return constraints
}

func expandConstraintsForProject(configConstraints []interface{}) []models.Constraint {
	constraints := make([]models.Constraint, 0, len(configConstraints))

	for _, configConstraint := range configConstraints {
		constraintMap := configConstraint.(map[string]interface{})

		constraint := models.Constraint{
			Mandatory:  withBool(constraintMap["mandatory"].(bool)),
			Expression: withString(constraintMap["expression"].(string)),
		}

		constraints = append(constraints, constraint)
	}

	return constraints
}

func flattenConstraints(constraints []models.Constraint) []interface{} {
	if len(constraints) == 0 {
		return make([]interface{}, 0)
	}

	configConstraints := make([]interface{}, 0, len(constraints))

	for _, constraint := range constraints {
		helper := make(map[string]interface{})
		helper["mandatory"] = *constraint.Mandatory
		helper["expression"] = *constraint.Expression

		configConstraints = append(configConstraints, helper)
	}

	return configConstraints
}
