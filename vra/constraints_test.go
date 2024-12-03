// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"github.com/vmware/vra-sdk-go/pkg/models"

	"testing"
)

func TestExpandConstraints(t *testing.T) {
	c1 := map[string]interface{}{"mandatory": true, "expression": "Foo:Bar"}
	c2 := map[string]interface{}{"mandatory": false, "expression": "Env:Test"}

	constraints := make([]interface{}, 0)
	expandedConstraints := expandConstraints(constraints)

	if len(expandedConstraints) != 0 {
		t.Errorf("error while expanding when there are no constraints")
	}

	constraints = append(constraints, c1)
	constraints = append(constraints, c2)

	expandedConstraints = expandConstraints(constraints)

	if len(expandedConstraints) != 2 {
		t.Errorf("not all constraints are expanded correctly")
	}

	if *expandedConstraints[0].Expression != c1["expression"] || *expandedConstraints[0].Mandatory != c1["mandatory"] {
		t.Errorf("constraint %#v is not expanded correctly", c1)
	}

	if *expandedConstraints[1].Expression != c2["expression"] || *expandedConstraints[1].Mandatory != c2["mandatory"] {
		t.Errorf("constraint %#v is not expanded correctly", c2)
	}
}

func TestExpandConstraintsForProject(t *testing.T) {
	c1 := map[string]interface{}{"mandatory": true, "expression": "Foo:Bar"}
	c2 := map[string]interface{}{"mandatory": false, "expression": "Env:Test"}

	constraints := make([]interface{}, 0)
	expandedConstraints := expandConstraintsForProject(constraints)

	if len(expandedConstraints) != 0 {
		t.Errorf("error while expanding when there are no constraints")
	}

	constraints = append(constraints, c1)
	constraints = append(constraints, c2)

	expandedConstraints = expandConstraintsForProject(constraints)

	if len(expandedConstraints) != 2 {
		t.Errorf("not all constraints are expanded correctly")
	}

	if *expandedConstraints[0].Expression != c1["expression"] || *expandedConstraints[0].Mandatory != c1["mandatory"] {
		t.Errorf("constraint %#v is not expanded correctly", c1)
	}

	if *expandedConstraints[1].Expression != c2["expression"] || *expandedConstraints[1].Mandatory != c2["mandatory"] {
		t.Errorf("constraint %#v is not expanded correctly", c2)
	}
}

func TestFlattenConstraints(t *testing.T) {
	constraints := make([]models.Constraint, 0)
	flattenedConstraints := flattenConstraints(constraints)

	if len(flattenedConstraints) != 0 {
		t.Errorf("error while flattening when there are no constraints")
	}

	constraint1 := models.Constraint{Expression: withString("Foo:Bar"), Mandatory: withBool(true)}
	constraint2 := models.Constraint{Expression: withString("Env:Test"), Mandatory: withBool(false)}

	constraints = append(constraints, constraint1)
	constraints = append(constraints, constraint2)

	flattenedConstraints = flattenConstraints(constraints)

	if len(flattenedConstraints) != 2 {
		t.Errorf("not all constraints are flattened correctly")
	}

	fc1 := flattenedConstraints[0].(map[string]interface{})
	if fc1["expression"].(string) != *constraint1.Expression || fc1["mandatory"].(bool) != *constraint1.Mandatory {
		t.Errorf("constraint %#v, %#v is not flattened correctly", *constraint2.Expression, *constraint2.Mandatory)
	}

	fc2 := flattenedConstraints[1].(map[string]interface{})
	if fc2["expression"].(string) != *constraint2.Expression || fc2["mandatory"].(bool) != *constraint2.Mandatory {
		t.Errorf("constraint %#v, %#v is not flattened correctly", *constraint2.Expression, *constraint2.Mandatory)
	}
}
