// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func TestDeploymentResourceActionID(t *testing.T) {
	deploymentUUID := strfmt.UUID("d8f9b3a1-0000-4000-8000-000000000001")
	resourceUUID := strfmt.UUID("d8f9b3a1-0000-4000-8000-000000000002")

	id := deploymentResourceActionID(deploymentUUID, resourceUUID, "Cloud.vSphere.Machine.Snapshot.Create")
	expected := "d8f9b3a1-0000-4000-8000-000000000001/d8f9b3a1-0000-4000-8000-000000000002/Cloud.vSphere.Machine.Snapshot.Create"
	if id != expected {
		t.Errorf("expected %q, got %q", expected, id)
	}
}

func TestParseDeploymentResourceActionID(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		deploymentID string
		resourceID   string
		actionID     string
		expectError  bool
	}{
		{
			name:         "valid id",
			id:           "11111111-1111-4111-8111-111111111111/22222222-2222-4222-8222-222222222222/PowerOff",
			deploymentID: "11111111-1111-4111-8111-111111111111",
			resourceID:   "22222222-2222-4222-8222-222222222222",
			actionID:     "PowerOff",
		},
		{
			name:         "action id containing separators",
			id:           "11111111-1111-4111-8111-111111111111/22222222-2222-4222-8222-222222222222/Custom/Action/Name",
			deploymentID: "11111111-1111-4111-8111-111111111111",
			resourceID:   "22222222-2222-4222-8222-222222222222",
			actionID:     "Custom/Action/Name",
		},
		{
			name:        "too few segments",
			id:          "11111111-1111-4111-8111-111111111111/22222222-2222-4222-8222-222222222222",
			expectError: true,
		},
		{
			name:        "empty action id",
			id:          "11111111-1111-4111-8111-111111111111/22222222-2222-4222-8222-222222222222/",
			expectError: true,
		},
		{
			name:        "empty string",
			id:          "",
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			deploymentUUID, resourceUUID, actionID, err := parseDeploymentResourceActionID(test.id)
			if test.expectError && err == nil {
				t.Fatalf("expected an error for the id %q", test.id)
			}
			if test.expectError {
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			actual := []string{deploymentUUID.String(), resourceUUID.String(), actionID}
			expected := []string{test.deploymentID, test.resourceID, test.actionID}
			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("expected parsed id %#v, got %#v", expected, actual)
			}
		})
	}
}

func TestFlattenResourceActionSchemaProperties(t *testing.T) {
	tests := []struct {
		name         string
		actionSchema interface{}
		expectedKeys []string
	}{
		{
			name: "schema with properties",
			actionSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"count": map[string]interface{}{"type": "integer"},
					"name":  map[string]interface{}{"type": "string"},
				},
			},
			expectedKeys: []string{"count", "name"},
		},
		{
			name:         "nil schema",
			actionSchema: nil,
			expectedKeys: []string{},
		},
		{
			name:         "schema is not a map",
			actionSchema: "not-a-map",
			expectedKeys: []string{},
		},
		{
			name:         "schema without properties",
			actionSchema: map[string]interface{}{"type": "object"},
			expectedKeys: []string{},
		},
		{
			name:         "properties is not a map",
			actionSchema: map[string]interface{}{"properties": []interface{}{"a"}},
			expectedKeys: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			properties := flattenResourceActionSchemaProperties(test.actionSchema)
			if properties == nil {
				t.Fatal("expected a non-nil map")
			}
			if len(properties) != len(test.expectedKeys) {
				t.Fatalf("expected %d properties, got %d", len(test.expectedKeys), len(properties))
			}
			for _, key := range test.expectedKeys {
				if _, ok := properties[key]; !ok {
					t.Errorf("expected the property %q to be present", key)
				}
			}
		})
	}
}

func TestFlattenResourceAction(t *testing.T) {
	action := &models.ResourceAction{
		Description: "Creates a snapshot.",
		DisplayName: "Create Snapshot",
		ID:          "Cloud.vSphere.Machine.Snapshot.Create",
		Name:        "Snapshot",
		Schema: map[string]interface{}{
			"properties": map[string]interface{}{
				"name": map[string]interface{}{"type": "string"},
			},
		},
		Valid: true,
	}

	flattened := flattenResourceAction(action)

	if flattened["id"] != action.ID {
		t.Errorf("expected the id %q, got %q", action.ID, flattened["id"])
	}
	if flattened["display_name"] != action.DisplayName {
		t.Errorf("expected the display name %q, got %q", action.DisplayName, flattened["display_name"])
	}
	if flattened["valid"] != true {
		t.Errorf("expected the action to be valid")
	}
	expectedSchema := `{"properties":{"name":{"type":"string"}}}`
	if flattened["schema_json"] != expectedSchema {
		t.Errorf("expected the schema %q, got %q", expectedSchema, flattened["schema_json"])
	}
}

func TestFlattenResourceActionWithoutSchema(t *testing.T) {
	flattened := flattenResourceAction(&models.ResourceAction{ID: "PowerOff"})

	if flattened["schema_json"] != "" {
		t.Errorf("expected an empty schema, got %q", flattened["schema_json"])
	}
}

func TestGetResourceActionInputTypesMapFromSchema(t *testing.T) {
	tests := []struct {
		name          string
		schema        map[string]interface{}
		expected      map[string]string
		errorContains string
	}{
		{
			name: "valid schema",
			schema: map[string]interface{}{
				"enabled": map[string]interface{}{"type": "boolean"},
				"ports":   map[string]interface{}{"type": "array"},
			},
			expected: map[string]string{"enabled": "boolean", "ports": "array"},
		},
		{
			name:          "property is not an object",
			schema:        map[string]interface{}{"enabled": "boolean"},
			errorContains: `input "enabled" has an invalid schema: expected an object`,
		},
		{
			name:          "property has no type",
			schema:        map[string]interface{}{"enabled": map[string]interface{}{"default": true}},
			errorContains: `input "enabled" has an invalid schema: missing type`,
		},
		{
			name:          "property type is not a string",
			schema:        map[string]interface{}{"enabled": map[string]interface{}{"type": []interface{}{"boolean", "null"}}},
			errorContains: `input "enabled" has an invalid schema: type must be a non-empty string`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := getResourceActionInputTypesMapFromSchema(test.schema)
			if test.errorContains != "" {
				if err == nil || !strings.Contains(err.Error(), test.errorContains) {
					t.Fatalf("expected error containing %q, got %v", test.errorContains, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(actual, test.expected) {
				t.Fatalf("expected %#v, got %#v", test.expected, actual)
			}
		})
	}
}

func TestGetInputsByType(t *testing.T) {
	inputs := map[string]interface{}{
		"enabled":  "true",
		"count":    "3",
		"ratio":    "1.5",
		"ports":    `[80,443]`,
		"metadata": `{"owner":"terraform"}`,
		"name":     "example",
	}
	types := map[string]string{
		"enabled":  "boolean",
		"count":    "integer",
		"ratio":    "number",
		"ports":    "array",
		"metadata": "object",
		"name":     "string",
	}

	actual, err := getInputsByType(inputs, types)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]interface{}{
		"enabled":  true,
		"count":    3,
		"ratio":    1.5,
		"ports":    []interface{}{float64(80), float64(443)},
		"metadata": map[string]interface{}{"owner": "terraform"},
		"name":     "example",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %#v, got %#v", expected, actual)
	}
}

func TestDeploymentResourceActionWaitForCompletion(t *testing.T) {
	resource := resourceDeploymentResourceAction()
	waitForCompletion, ok := resource.Schema["wait_for_completion"]
	if !ok {
		t.Fatal("wait_for_completion must remain a supported resource argument")
	}
	if waitForCompletion.Default != true {
		t.Fatal("wait_for_completion must default to true")
	}
}

func TestDeploymentResourceActionProviderSchema(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("provider schema is invalid: %v", err)
	}
}

func TestLockDeployment(t *testing.T) {
	first := strfmt.UUID("d8f9b3a1-0000-4000-8000-00000000000a")
	second := strfmt.UUID("d8f9b3a1-0000-4000-8000-00000000000b")

	// A different deployment must not be blocked by the one that is already held.
	release, err := lockDeployment(context.Background(), first)
	if err != nil {
		t.Fatalf("unexpected error locking first deployment: %v", err)
	}
	done := make(chan error, 1)
	go func() {
		releaseSecond, err := lockDeployment(context.Background(), second)
		if err == nil {
			releaseSecond()
		}
		done <- err
	}()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("unexpected error locking second deployment: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("locking a different deployment blocked")
	}
	release()

	// The same deployment must be served by the same lock.
	release, err = lockDeployment(context.Background(), first)
	if err != nil {
		t.Fatalf("unexpected error locking first deployment again: %v", err)
	}
	locked := make(chan error, 1)
	go func() {
		releaseFirst, err := lockDeployment(context.Background(), first)
		if err == nil {
			releaseFirst()
		}
		locked <- err
	}()
	select {
	case err := <-locked:
		if err != nil {
			t.Fatalf("unexpected error waiting for first deployment: %v", err)
		}
		t.Fatal("the same deployment was locked twice at the same time")
	case <-time.After(100 * time.Millisecond):
	}
	release()
	if err := <-locked; err != nil {
		t.Fatalf("unexpected error after releasing first deployment: %v", err)
	}

	// Cancellation must stop a waiter before it can acquire the deployment later.
	release, err = lockDeployment(context.Background(), first)
	if err != nil {
		t.Fatalf("unexpected error locking first deployment for cancellation test: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancelled := make(chan error, 1)
	go func() {
		_, err := lockDeployment(ctx, first)
		cancelled <- err
	}()
	cancel()
	if err := <-cancelled; !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context cancellation, got %v", err)
	}
	release()
}

func TestDeploymentResourceActionPendingStatuses(t *testing.T) {
	for _, status := range []string{
		models.RequestStatusCREATED,
		models.RequestStatusPENDING,
		models.RequestStatusINITIALIZATION,
		models.RequestStatusCHECKINGAPPROVAL,
		models.RequestStatusAPPROVALPENDING,
		models.RequestStatusUSERINTERACTIONPENDING,
		models.RequestStatusINPROGRESS,
		models.RequestStatusCOMPLETION,
	} {
		if !isDeploymentResourceActionPending(status) {
			t.Errorf("expected %s to be pending", status)
		}
	}

	for _, status := range []string{
		models.RequestStatusSUCCESSFUL,
		models.RequestStatusFAILED,
		models.RequestStatusAPPROVALREJECTED,
		models.RequestStatusABORTED,
	} {
		if isDeploymentResourceActionPending(status) {
			t.Errorf("expected %s not to be pending", status)
		}
	}
}
