package cas

import (
	"github.com/vmware/terraform-provider-cas/sdk"

	"github.com/hashicorp/terraform/helper/schema"
)

// constraintsSchema returns the schema to use for the constraints property
func constraintsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
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

func expandConstraints(configConstraints []interface{}) []tango.Constraint {
	constraints := make([]tango.Constraint, 0, len(configConstraints))

	for _, configConstraint := range configConstraints {
		constraintMap := configConstraint.(map[string]interface{})

		constraint := tango.Constraint{
			Mandatory:  constraintMap["mandatory"].(bool),
			Expression: constraintMap["expression"].(string),
		}

		constraints = append(constraints, constraint)
	}

	return constraints
}

func flattenConstraints(constraints []tango.Constraint) []interface{} {
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
