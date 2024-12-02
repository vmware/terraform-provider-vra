// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"testing"
)

func TestExpandInputs(t *testing.T) {
	type Message struct {
		Name string
		Body string
		Time int64
	}
	m := Message{"Foo", "Bar", 1294706395881547000}

	inputs := make(map[string]interface{})
	inputs["whole"] = 2
	inputs["text"] = "This is string"
	inputs["flag"] = true
	inputs["fraction"] = 2.4
	inputs["message"] = m

	expandedInputs := expandInputs(inputs)

	if expandedInputs["whole"] != 2 {
		t.Errorf("int type input is not expanded correctly.")
	}

	if expandedInputs["text"] != "This is string" {
		t.Errorf("string type input is not expanded correctly.")
	}

	if expandedInputs["flag"] != true {
		t.Errorf("bool type input is not expanded correctly.")
	}

	if expandedInputs["fraction"] != 2.4 {
		t.Errorf("float type input is not expanded correctly.")
	}

	if expandedInputs["message"] != m {
		t.Errorf("object type input is not expanded correctly.")
	}
}
