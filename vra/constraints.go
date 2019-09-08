package vra

import (
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform/helper/schema"
)

// constraintsSchema returns the schema to use for the constraints property
func constraintsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"mandatory": {
					Type:     schema.TypeBool,
					Required: true,
				},
				"expression": {
					Type:     schema.TypeString,
					Required: true,
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

/*
func flattenConstraints(constraints []*models.Constraint) []interface{} {
	if len(constraints) == 0 {
		return make([]interface{}, 0)
	}

	configConstraints := make([]interface{}, 0, len(constraints))

	for _, constraint := range constraints {
		helper := make(map[string]interface{})
		helper["mandatory"] = constraint.Mandatory
		helper["expression"] = constraint.Expression

		configConstraints = append(configConstraints, helper)
	}

	return configConstraints
}
*/
