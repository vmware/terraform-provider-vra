// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"encoding/json"
	"testing"
)

// Constants for repeated literal values used across test cases.
const (
	testProviderID = "prov-123"
	testUUID       = "uuid-456"
	testPlatformID = "prov-abc"
)

// requireNoErr is a test helper that marks the calling test as failed if err is non-nil.
func requireNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// computeChangedKeys
// ---------------------------------------------------------------------------

func TestComputeChangedKeys(t *testing.T) {
	cases := []struct {
		name string
		old  interface{}
		new  interface{}
		want map[string]bool
	}{
		{
			name: "added key detected",
			old:  map[string]interface{}{"a": "1"},
			new:  map[string]interface{}{"a": "1", "b": "2"},
			want: map[string]bool{"b": true},
		},
		{
			name: "removed key detected",
			old:  map[string]interface{}{"a": "1", "b": "2"},
			new:  map[string]interface{}{"a": "1"},
			want: map[string]bool{"b": true},
		},
		{
			name: "modified value detected",
			old:  map[string]interface{}{"port": "9091"},
			new:  map[string]interface{}{"port": "9094"},
			want: map[string]bool{"port": true},
		},
		{
			name: "no change returns empty set",
			old:  map[string]interface{}{"a": "1"},
			new:  map[string]interface{}{"a": "1"},
			want: map[string]bool{},
		},
		{
			name: "nil inputs return empty set",
			old:  nil,
			new:  nil,
			want: map[string]bool{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := computeChangedKeys(tc.old, tc.new)
			for k := range tc.want {
				if !got[k] {
					t.Errorf("expected %q in changed keys, got %v", k, got)
				}
			}
			for k := range got {
				if !tc.want[k] {
					t.Errorf("unexpected key %q in changed keys", k)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// isLiteralDefault
// ---------------------------------------------------------------------------

func TestIsLiteralDefault(t *testing.T) {
	cases := []struct {
		name string
		v    interface{}
		want bool
	}{
		{"nil is not a literal", nil, false},
		{"map is not a literal (server-side descriptor)", map[string]interface{}{"bind": "_resource.id"}, false},
		{"string is a literal", "hello", true},
		{"integer is a literal", 42, true},
		{"bool is a literal", true, true},
		{"empty array is a literal", []interface{}{}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isLiteralDefault(tc.v); got != tc.want {
				t.Errorf("isLiteralDefault(%v) = %v, want %v", tc.v, got, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// nativePlatformID
// ---------------------------------------------------------------------------

func TestNativePlatformID(t *testing.T) {
	cases := []struct {
		name  string
		props map[string]interface{}
		want  string
	}{
		{
			name:  "providerId preferred over uuid",
			props: map[string]interface{}{"providerId": testProviderID, "uuid": testUUID},
			want:  testProviderID,
		},
		{
			name:  "uuid used when providerId absent",
			props: map[string]interface{}{"uuid": testUUID},
			want:  testUUID,
		},
		{
			name:  "empty props return empty string",
			props: map[string]interface{}{},
			want:  "",
		},
		{
			name:  "non-string providerId is ignored",
			props: map[string]interface{}{"providerId": 12345},
			want:  "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := nativePlatformID(tc.props); got != tc.want {
				t.Errorf("nativePlatformID() = %q, want %q", got, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// overlayLastRequestInputs
// ---------------------------------------------------------------------------

func TestOverlayLastRequestInputs(t *testing.T) {
	t.Run("updates existing deployment input key", func(t *testing.T) {
		all := map[string]interface{}{"port": "9091"}
		overlayLastRequestInputs(all, map[string]interface{}{"port": "9094"}, nil)
		if all["port"] != "9094" {
			t.Errorf("existing key should be updated, got %v", all["port"])
		}
	})

	t.Run("updates user-managed key absent from deployment inputs", func(t *testing.T) {
		all := map[string]interface{}{}
		overlayLastRequestInputs(all,
			map[string]interface{}{"customKey": "newVal"},
			map[string]bool{"customKey": true},
		)
		if all["customKey"] != "newVal" {
			t.Errorf("user-managed key should be overlaid, got %v", all["customKey"])
		}
	})

	t.Run("action-internal keys do not leak into state", func(t *testing.T) {
		all := map[string]interface{}{"site": "fre"}
		overlayLastRequestInputs(all,
			map[string]interface{}{"internalWorkflowToken": "secret"},
			nil,
		)
		if _, leaked := all["internalWorkflowToken"]; leaked {
			t.Error("action-internal key must not be added to state")
		}
	})

	t.Run("nil userManagedKeys does not panic", func(t *testing.T) {
		all := map[string]interface{}{"a": "old"}
		overlayLastRequestInputs(all,
			map[string]interface{}{"a": "new", "unknown": "x"},
			nil,
		)
		if all["a"] != "new" {
			t.Errorf("existing key should be updated, got %v", all["a"])
		}
		if _, leaked := all["unknown"]; leaked {
			t.Error("unknown key should not have been added")
		}
	})
}

// ---------------------------------------------------------------------------
// pickFieldValue
// ---------------------------------------------------------------------------

func TestPickFieldValue(t *testing.T) {
	cases := []struct {
		name         string
		key          string
		fieldType    string
		typed        map[string]interface{}
		props        map[string]interface{}
		defVal       interface{}
		platformID   string
		selfRefIDKey string
		want         interface{}
	}{
		{
			name:       "deployment input takes priority",
			key:        "port",
			fieldType:  "integer",
			typed:      map[string]interface{}{"port": 9094},
			props:      map[string]interface{}{"port": 9091},
			platformID: testProviderID,
			want:       9094,
		},
		{
			name:      "resource property used when no deployment input",
			key:       "site",
			fieldType: "string",
			typed:     map[string]interface{}{},
			props:     map[string]interface{}{"site": "fre"},
			want:      "fre",
		},
		{
			name:      "literal schema default used when no input or property",
			key:       "enabled",
			fieldType: "boolean",
			typed:     map[string]interface{}{},
			props:     map[string]interface{}{},
			defVal:    true,
			want:      true,
		},
		{
			name:         "map-typed default skipped, self-ref Id heuristic applies",
			key:          "resourceId",
			fieldType:    "string",
			typed:        map[string]interface{}{},
			props:        map[string]interface{}{},
			defVal:       map[string]interface{}{"bind": "_resource.id"},
			platformID:   testPlatformID,
			selfRefIDKey: "resourceId",
			want:         testPlatformID,
		},
		{
			name:         "self-ref Id heuristic returns platform identifier when key matches selfRefIDKey",
			key:          "virtualServiceId",
			fieldType:    "string",
			typed:        map[string]interface{}{},
			props:        map[string]interface{}{},
			platformID:   testPlatformID,
			selfRefIDKey: "virtualServiceId",
			want:         testPlatformID,
		},
		{
			name:         "self-ref Id heuristic does NOT apply for a different resource's Id field",
			key:          "networkId",
			fieldType:    "string",
			typed:        map[string]interface{}{},
			props:        map[string]interface{}{},
			platformID:   testPlatformID,
			selfRefIDKey: "virtualServiceId", // resource is VirtualService, not Network
			want:         nil,
		},
		{
			name:         "self-ref Id heuristic does not apply to non-string fields",
			key:          "virtualServiceId",
			fieldType:    "integer",
			typed:        map[string]interface{}{},
			props:        map[string]interface{}{},
			platformID:   testPlatformID,
			selfRefIDKey: "virtualServiceId",
			want:         nil,
		},
		{
			name:      "nil returned when nothing matches",
			key:       "unknownField",
			fieldType: "string",
			typed:     map[string]interface{}{},
			props:     map[string]interface{}{},
			want:      nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := pickFieldValue(tc.key, tc.fieldType, tc.typed, tc.props, tc.defVal, tc.platformID, tc.selfRefIDKey)
			if got != tc.want {
				t.Errorf("pickFieldValue() = %v, want %v", got, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// getInputTypesMapFromSchema
// ---------------------------------------------------------------------------

func TestGetInputTypesMapFromSchema(t *testing.T) {
	t.Run("normal schema returns type map", func(t *testing.T) {
		schema := map[string]interface{}{
			"port": map[string]interface{}{"type": "integer"},
			"name": map[string]interface{}{"type": "string"},
		}
		got, err := getInputTypesMapFromSchema(schema)
		requireNoErr(t, err)
		if got["port"] != "integer" || got["name"] != "string" {
			t.Errorf("unexpected result: %v", got)
		}
	})

	t.Run("non-map property skipped without panic", func(t *testing.T) {
		schema := map[string]interface{}{
			"good": map[string]interface{}{"type": "string"},
			"bad":  "not-a-map",
		}
		got, err := getInputTypesMapFromSchema(schema)
		requireNoErr(t, err)
		if _, exists := got["bad"]; exists {
			t.Error("non-map property should be skipped")
		}
		if got["good"] != "string" {
			t.Errorf("valid property should still be present, got %v", got)
		}
	})

	t.Run("property without type field is skipped", func(t *testing.T) {
		schema := map[string]interface{}{
			"refProp": map[string]interface{}{"$ref": "#/definitions/Foo"},
		}
		got, err := getInputTypesMapFromSchema(schema)
		requireNoErr(t, err)
		if _, exists := got["refProp"]; exists {
			t.Error("property without type field should be skipped")
		}
	})

	t.Run("non-string type field is skipped", func(t *testing.T) {
		schema := map[string]interface{}{
			"weird": map[string]interface{}{"type": []string{"string", "null"}},
		}
		got, err := getInputTypesMapFromSchema(schema)
		requireNoErr(t, err)
		if _, exists := got["weird"]; exists {
			t.Error("property with non-string type field should be skipped")
		}
	})

	t.Run("empty schema returns empty map", func(t *testing.T) {
		got, err := getInputTypesMapFromSchema(map[string]interface{}{})
		requireNoErr(t, err)
		if len(got) != 0 {
			t.Errorf("expected empty map, got %v", got)
		}
	})
}

// ---------------------------------------------------------------------------
// updateUserInputs
// ---------------------------------------------------------------------------

func TestUpdateUserInputs(t *testing.T) {
	t.Run("nil allInputs returns nil", func(t *testing.T) {
		if got := updateUserInputs(nil, map[string]interface{}{"a": "1"}, nil); got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})

	t.Run("nil userInputs returns nil", func(t *testing.T) {
		if got := updateUserInputs(map[string]interface{}{"a": "1"}, nil, nil); got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})

	t.Run("key absent from deployment inputs preserves user value", func(t *testing.T) {
		// Key exists in userInputs but not in deployment.Inputs (action-only key).
		// The user-configured value must be preserved, not collapsed to nil.
		got := updateUserInputs(
			map[string]interface{}{},
			map[string]interface{}{"actionOnlyKey": "myvalue"},
			map[string]string{},
		)
		if got["actionOnlyKey"] != "myvalue" {
			t.Errorf("user value must be preserved when key absent from deployment inputs, got %v", got["actionOnlyKey"])
		}
	})

	t.Run("platform value is always used as ground truth", func(t *testing.T) {
		// Platform reports "9091"; Terraform state (userInputs) has a corrupted or
		// stale value. The platform value must win unconditionally so that
		// terraform apply -refresh-only corrects any corrupted state.
		got := updateUserInputs(
			map[string]interface{}{"port": "9091"},
			map[string]interface{}{"port": "9091-9095"},
			map[string]string{"port": "string"},
		)
		if got["port"] != "9091" {
			t.Errorf("platform value should always be used as ground truth, got %v", got["port"])
		}
	})

	t.Run("matching value returns decoded platform value", func(t *testing.T) {
		got := updateUserInputs(
			map[string]interface{}{"env": "prod"},
			map[string]interface{}{"env": "prod"},
			map[string]string{"env": "string"},
		)
		if got["env"] != "prod" {
			t.Errorf("matching value should be preserved, got %v", got["env"])
		}
	})
}

// ---------------------------------------------------------------------------
// decodeInputValue
// ---------------------------------------------------------------------------

func TestDecodeInputValue(t *testing.T) {
	t.Run("string type returns string", func(t *testing.T) {
		got := decodeInputValue("env", "prod", map[string]string{"env": "string"})
		if got != "prod" {
			t.Errorf("expected 'prod', got %v", got)
		}
	})

	t.Run("integer type returns string representation", func(t *testing.T) {
		got := decodeInputValue("port", float64(9092), map[string]string{"port": "integer"})
		if got != "9092" {
			t.Errorf("expected '9092', got %v", got)
		}
	})

	t.Run("array type returns JSON string", func(t *testing.T) {
		got := decodeInputValue("ports", []interface{}{"9092", "9093"}, map[string]string{"ports": "array"})
		if got != `["9092","9093"]` {
			t.Errorf("expected JSON array string, got %v", got)
		}
	})

	t.Run("object type returns JSON string", func(t *testing.T) {
		got := decodeInputValue("config", map[string]interface{}{"a": "b"}, map[string]string{"config": "object"})
		if got != `{"a":"b"}` {
			t.Errorf("expected JSON object string, got %v", got)
		}
	})

	t.Run("auto-detect array when no type info", func(t *testing.T) {
		got := decodeInputValue("ports", []interface{}{"9092"}, map[string]string{})
		if got != `["9092"]` {
			t.Errorf("expected JSON array, got %v", got)
		}
	})

	t.Run("auto-detect map when no type info", func(t *testing.T) {
		got := decodeInputValue("obj", map[string]interface{}{"x": 1}, map[string]string{})
		str, ok := got.(string)
		if !ok {
			t.Fatalf("expected string, got %T", got)
		}
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(str), &parsed); err != nil {
			t.Fatalf("expected valid JSON, got %v", str)
		}
	})

	t.Run("plain value with no type info returns Sprint", func(t *testing.T) {
		got := decodeInputValue("count", 42, map[string]string{})
		if got != "42" {
			t.Errorf("expected '42', got %v", got)
		}
	})
}

// ---------------------------------------------------------------------------
// resolveUnmatchedArrayFields
// ---------------------------------------------------------------------------

func TestResolveUnmatchedArrayFields(t *testing.T) {
	t.Run("matches array-of-object by structure", func(t *testing.T) {
		schemaMap := map[string]interface{}{
			"datagrid_abc": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"port":           map[string]interface{}{"type": "integer"},
						"port_range_end": map[string]interface{}{"type": "integer"},
					},
				},
			},
		}
		inputTypesMap := map[string]string{"datagrid_abc": "array"}
		rawInputs := map[string]interface{}{
			"portRangeObj": `[{"port":9092,"port_range_end":9092}]`,
		}
		result := map[string]interface{}{}
		resolveUnmatchedArrayFields(schemaMap, inputTypesMap, rawInputs, result)
		if result["datagrid_abc"] == nil {
			t.Fatal("expected datagrid_abc to be matched")
		}
		arr, ok := result["datagrid_abc"].([]interface{})
		if !ok || len(arr) != 1 {
			t.Fatalf("expected parsed array with 1 item, got %v", result["datagrid_abc"])
		}
	})

	t.Run("does not match when key counts differ", func(t *testing.T) {
		schemaMap := map[string]interface{}{
			"field1": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"a": map[string]interface{}{"type": "string"},
						"b": map[string]interface{}{"type": "string"},
						"c": map[string]interface{}{"type": "string"},
					},
				},
			},
		}
		inputTypesMap := map[string]string{"field1": "array"}
		rawInputs := map[string]interface{}{
			"input1": `[{"a":"1","b":"2"}]`,
		}
		result := map[string]interface{}{}
		resolveUnmatchedArrayFields(schemaMap, inputTypesMap, rawInputs, result)
		if result["field1"] != nil {
			t.Error("should not match when key counts differ")
		}
	})

	t.Run("skips already-matched fields", func(t *testing.T) {
		schemaMap := map[string]interface{}{
			"field1": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{"a": map[string]interface{}{"type": "string"}},
				},
			},
		}
		inputTypesMap := map[string]string{"field1": "array"}
		rawInputs := map[string]interface{}{"something": `[{"a":"x"}]`}
		result := map[string]interface{}{"field1": "already-set"}

		resolveUnmatchedArrayFields(schemaMap, inputTypesMap, rawInputs, result)
		if result["field1"] != "already-set" {
			t.Error("should not overwrite already-matched field")
		}
	})

	t.Run("skips non-array types", func(t *testing.T) {
		schemaMap := map[string]interface{}{
			"field1": map[string]interface{}{"type": "string"},
		}
		inputTypesMap := map[string]string{"field1": "string"}
		result := map[string]interface{}{}
		resolveUnmatchedArrayFields(schemaMap, inputTypesMap, map[string]interface{}{}, result)
		if len(result) != 0 {
			t.Error("non-array types should be skipped")
		}
	})

	t.Run("skips non-JSON string inputs", func(t *testing.T) {
		schemaMap := map[string]interface{}{
			"field1": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{"a": map[string]interface{}{"type": "string"}},
				},
			},
		}
		inputTypesMap := map[string]string{"field1": "array"}
		rawInputs := map[string]interface{}{"input1": "not-json"}
		result := map[string]interface{}{}
		resolveUnmatchedArrayFields(schemaMap, inputTypesMap, rawInputs, result)
		if len(result) != 0 {
			t.Error("non-JSON inputs should be skipped")
		}
	})
}

// ---------------------------------------------------------------------------
// resolveByTitle
// ---------------------------------------------------------------------------

func TestResolveByTitle(t *testing.T) {
	t.Run("matches by normalized title", func(t *testing.T) {
		schemaMap := map[string]interface{}{
			"textField_abc": map[string]interface{}{
				"type":            "string",
				"title":           "Site",
				"$dynamicDefault": map[string]interface{}{},
			},
		}
		inputTypesMap := map[string]string{"textField_abc": "string"}
		rawInputs := map[string]interface{}{"site": "fre"}
		result := map[string]interface{}{}

		resolveByTitle(schemaMap, rawInputs, inputTypesMap, result)
		if result["textField_abc"] != "fre" {
			t.Errorf("expected 'fre', got %v", result["textField_abc"])
		}
	})

	t.Run("multi-word title normalized correctly", func(t *testing.T) {
		schemaMap := map[string]interface{}{
			"textField_xyz": map[string]interface{}{
				"type":            "string",
				"title":           "Tier 1 Gateway",
				"$dynamicDefault": map[string]interface{}{},
			},
		}
		inputTypesMap := map[string]string{"textField_xyz": "string"}
		rawInputs := map[string]interface{}{"tier1gateway": "gw-01"}
		result := map[string]interface{}{}

		resolveByTitle(schemaMap, rawInputs, inputTypesMap, result)
		if result["textField_xyz"] != "gw-01" {
			t.Errorf("expected 'gw-01', got %v", result["textField_xyz"])
		}
	})

	t.Run("skips fields without dynamicDefault or data", func(t *testing.T) {
		schemaMap := map[string]interface{}{
			"field1": map[string]interface{}{
				"type":  "string",
				"title": "Site",
			},
		}
		inputTypesMap := map[string]string{"field1": "string"}
		rawInputs := map[string]interface{}{"site": "fre"}
		result := map[string]interface{}{}

		resolveByTitle(schemaMap, rawInputs, inputTypesMap, result)
		if result["field1"] != nil {
			t.Error("should not match fields without $dynamicDefault or $data")
		}
	})

	t.Run("skips readOnly fields", func(t *testing.T) {
		schemaMap := map[string]interface{}{
			"field1": map[string]interface{}{
				"type":            "string",
				"title":           "Site",
				"readOnly":        true,
				"$dynamicDefault": map[string]interface{}{},
			},
		}
		inputTypesMap := map[string]string{"field1": "string"}
		rawInputs := map[string]interface{}{"site": "fre"}
		result := map[string]interface{}{}

		resolveByTitle(schemaMap, rawInputs, inputTypesMap, result)
		if result["field1"] != nil {
			t.Error("should not match readOnly fields")
		}
	})

	t.Run("skips already-matched fields", func(t *testing.T) {
		schemaMap := map[string]interface{}{
			"field1": map[string]interface{}{
				"type":            "string",
				"title":           "Site",
				"$dynamicDefault": map[string]interface{}{},
			},
		}
		inputTypesMap := map[string]string{"field1": "string"}
		rawInputs := map[string]interface{}{"site": "fre"}
		result := map[string]interface{}{"field1": "already-set"}

		resolveByTitle(schemaMap, rawInputs, inputTypesMap, result)
		if result["field1"] != "already-set" {
			t.Error("should not overwrite already-matched field")
		}
	})

	t.Run("boolean type conversion", func(t *testing.T) {
		schemaMap := map[string]interface{}{
			"toggle_abc": map[string]interface{}{
				"type":            "boolean",
				"title":           "Enabled",
				"$dynamicDefault": map[string]interface{}{},
			},
		}
		inputTypesMap := map[string]string{"toggle_abc": "boolean"}
		rawInputs := map[string]interface{}{"enabled": "true"}
		result := map[string]interface{}{}

		resolveByTitle(schemaMap, rawInputs, inputTypesMap, result)
		if result["toggle_abc"] != true {
			t.Errorf("expected true, got %v", result["toggle_abc"])
		}
	})

	t.Run("integer type conversion", func(t *testing.T) {
		schemaMap := map[string]interface{}{
			"num_abc": map[string]interface{}{
				"type":            "integer",
				"title":           "Count",
				"$dynamicDefault": map[string]interface{}{},
			},
		}
		inputTypesMap := map[string]string{"num_abc": "integer"}
		rawInputs := map[string]interface{}{"count": "42"}
		result := map[string]interface{}{}

		resolveByTitle(schemaMap, rawInputs, inputTypesMap, result)
		if result["num_abc"] != 42 {
			t.Errorf("expected 42, got %v", result["num_abc"])
		}
	})

	t.Run("array type JSON deserialization", func(t *testing.T) {
		schemaMap := map[string]interface{}{
			"arr_abc": map[string]interface{}{
				"type":  "array",
				"title": "Ports",
				"$data": map[string]interface{}{},
			},
		}
		inputTypesMap := map[string]string{"arr_abc": "array"}
		rawInputs := map[string]interface{}{"ports": `["9092","9093"]`}
		result := map[string]interface{}{}

		resolveByTitle(schemaMap, rawInputs, inputTypesMap, result)
		arr, ok := result["arr_abc"].([]interface{})
		if !ok || len(arr) != 2 {
			t.Fatalf("expected parsed array with 2 items, got %v", result["arr_abc"])
		}
	})
}

// ---------------------------------------------------------------------------
// resolveFromResourceProperties / extractValuesFromProperties
// ---------------------------------------------------------------------------

func TestResolveFromResourceProperties(t *testing.T) {
	t.Run("extracts values from resource properties by title keyword", func(t *testing.T) {
		schemaMap := map[string]interface{}{
			"arr_field": map[string]interface{}{
				"type":  "array",
				"title": "Current Service Ports",
				"items": map[string]interface{}{"type": "string"},
			},
		}
		inputTypesMap := map[string]string{"arr_field": "array"}
		currentProps := map[string]interface{}{
			"services": []interface{}{
				map[string]interface{}{"port": 9092, "protocol": "TCP"},
				map[string]interface{}{"port": 9093, "protocol": "TCP"},
			},
		}
		result := map[string]interface{}{}

		resolveFromResourceProperties(schemaMap, inputTypesMap, currentProps, result)
		arr, ok := result["arr_field"].([]interface{})
		if !ok || len(arr) != 2 {
			t.Fatalf("expected 2 values, got %v", result["arr_field"])
		}
		if arr[0] != "9092" || arr[1] != "9093" {
			t.Errorf("expected [9092 9093], got %v", arr)
		}
	})

	t.Run("skips already-matched fields", func(t *testing.T) {
		schemaMap := map[string]interface{}{
			"arr_field": map[string]interface{}{
				"type":  "array",
				"title": "Current Ports",
				"items": map[string]interface{}{"type": "string"},
			},
		}
		inputTypesMap := map[string]string{"arr_field": "array"}
		result := map[string]interface{}{"arr_field": "already-set"}

		resolveFromResourceProperties(schemaMap, inputTypesMap, map[string]interface{}{}, result)
		if result["arr_field"] != "already-set" {
			t.Error("should not overwrite already-matched field")
		}
	})

	t.Run("skips fields with dynamicDefault", func(t *testing.T) {
		schemaMap := map[string]interface{}{
			"arr_field": map[string]interface{}{
				"type":            "array",
				"title":           "Ports",
				"$dynamicDefault": map[string]interface{}{},
				"items":           map[string]interface{}{"type": "string"},
			},
		}
		inputTypesMap := map[string]string{"arr_field": "array"}
		currentProps := map[string]interface{}{
			"services": []interface{}{
				map[string]interface{}{"port": 9092},
			},
		}
		result := map[string]interface{}{}

		resolveFromResourceProperties(schemaMap, inputTypesMap, currentProps, result)
		if result["arr_field"] != nil {
			t.Error("should not match fields with $dynamicDefault")
		}
	})

	t.Run("skips array-of-object types", func(t *testing.T) {
		schemaMap := map[string]interface{}{
			"arr_field": map[string]interface{}{
				"type":  "array",
				"title": "Ports",
				"items": map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
				},
			},
		}
		inputTypesMap := map[string]string{"arr_field": "array"}
		result := map[string]interface{}{}

		resolveFromResourceProperties(schemaMap, inputTypesMap, map[string]interface{}{}, result)
		if result["arr_field"] != nil {
			t.Error("should not match array-of-object fields")
		}
	})

	t.Run("no match returns nothing", func(t *testing.T) {
		schemaMap := map[string]interface{}{
			"arr_field": map[string]interface{}{
				"type":  "array",
				"title": "Unrelated Title",
				"items": map[string]interface{}{"type": "string"},
			},
		}
		inputTypesMap := map[string]string{"arr_field": "array"}
		currentProps := map[string]interface{}{
			"services": []interface{}{
				map[string]interface{}{"port": 9092},
			},
		}
		result := map[string]interface{}{}

		resolveFromResourceProperties(schemaMap, inputTypesMap, currentProps, result)
		if result["arr_field"] != nil {
			t.Error("should not match when title doesn't contain any property field name")
		}
	})
}

func TestExtractValuesFromProperties(t *testing.T) {
	t.Run("extracts values when field name appears in title", func(t *testing.T) {
		props := map[string]interface{}{
			"services": []interface{}{
				map[string]interface{}{"port": 9092, "name": "svc1"},
				map[string]interface{}{"port": 9093, "name": "svc2"},
			},
		}
		values := extractValuesFromProperties("Current Service Ports", props)
		if len(values) != 2 || values[0] != "9092" || values[1] != "9093" {
			t.Errorf("expected [9092 9093], got %v", values)
		}
	})

	t.Run("returns nil when no array properties exist", func(t *testing.T) {
		props := map[string]interface{}{
			"name": "myresource",
		}
		values := extractValuesFromProperties("Some Title", props)
		if values != nil {
			t.Errorf("expected nil, got %v", values)
		}
	})

	t.Run("returns nil when no field name matches title", func(t *testing.T) {
		props := map[string]interface{}{
			"items": []interface{}{
				map[string]interface{}{"count": 5},
			},
		}
		values := extractValuesFromProperties("Unrelated Title", props)
		if values != nil {
			t.Errorf("expected nil, got %v", values)
		}
	})

	t.Run("skips non-object array items", func(t *testing.T) {
		props := map[string]interface{}{
			"names": []interface{}{"a", "b", "c"},
		}
		values := extractValuesFromProperties("Names List", props)
		if values != nil {
			t.Errorf("expected nil for non-object array, got %v", values)
		}
	})

	t.Run("empty properties return nil", func(t *testing.T) {
		values := extractValuesFromProperties("Title", map[string]interface{}{})
		if values != nil {
			t.Errorf("expected nil for empty props, got %v", values)
		}
	})

	t.Run("rejects substring match in middle of word", func(t *testing.T) {
		props := map[string]interface{}{
			"services": []interface{}{
				map[string]interface{}{"port": 9092, "name": "svc1"},
			},
		}
		// "port" should NOT match "transport" (field not at word start)
		values := extractValuesFromProperties("Transport Config", props)
		if values != nil {
			t.Errorf("expected nil for substring-only match, got %v", values)
		}
	})
}

// ---------------------------------------------------------------------------
// fieldNameMatchesTitle
// ---------------------------------------------------------------------------

func TestFieldNameMatchesTitle(t *testing.T) {
	t.Run("exact word match", func(t *testing.T) {
		if !fieldNameMatchesTitle("port", "Current Port") {
			t.Error("expected match for exact word")
		}
	})

	t.Run("plural word match", func(t *testing.T) {
		if !fieldNameMatchesTitle("port", "Current Service Ports") {
			t.Error("expected match for plural")
		}
	})

	t.Run("case insensitive", func(t *testing.T) {
		if !fieldNameMatchesTitle("Port", "current port") {
			t.Error("expected case-insensitive match")
		}
	})

	t.Run("rejects mid-word substring", func(t *testing.T) {
		if fieldNameMatchesTitle("port", "Transport Config") {
			t.Error("should not match mid-word substring")
		}
	})

	t.Run("rejects end-of-word substring", func(t *testing.T) {
		if fieldNameMatchesTitle("id", "Grid Settings") {
			t.Error("should not match end-of-word substring")
		}
	})

	t.Run("short field name at word start matches", func(t *testing.T) {
		if !fieldNameMatchesTitle("id", "ID Field") {
			t.Error("expected match when field is full word")
		}
	})
}
