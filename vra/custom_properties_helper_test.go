// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"sort"
	"testing"
)

func TestRemovedCustomPropertyKeys(t *testing.T) {
	cases := []struct {
		name     string
		old      map[string]interface{}
		new      map[string]interface{}
		expected []string
	}{
		{
			name:     "one key removed",
			old:      map[string]interface{}{"prop_1": "val1", "prop_2": "val2"},
			new:      map[string]interface{}{"prop_1": "val1"},
			expected: []string{"prop_2"},
		},
		{
			name:     "all keys removed",
			old:      map[string]interface{}{"prop_1": "val1", "prop_2": "val2"},
			new:      map[string]interface{}{},
			expected: []string{"prop_1", "prop_2"},
		},
		{
			name:     "no change",
			old:      map[string]interface{}{"prop_1": "val1"},
			new:      map[string]interface{}{"prop_1": "val1"},
			expected: []string{},
		},
		{
			name:     "value changed",
			old:      map[string]interface{}{"prop_1": "val1"},
			new:      map[string]interface{}{"prop_1": "val2"},
			expected: []string{},
		},
		{
			name:     "key added",
			old:      map[string]interface{}{},
			new:      map[string]interface{}{"prop_1": "val1"},
			expected: []string{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			removed := removedCustomPropertyKeys(tc.old, tc.new)
			sort.Strings(removed)

			if len(removed) != len(tc.expected) {
				t.Fatalf("expected %v, got %v", tc.expected, removed)
			}
			for i := range removed {
				if removed[i] != tc.expected[i] {
					t.Fatalf("expected %v, got %v", tc.expected, removed)
				}
			}
		})
	}
}
