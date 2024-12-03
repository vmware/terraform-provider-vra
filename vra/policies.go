// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"encoding/json"
)

const (
	CatalogEntitlementTypeID string = "com.vmware.policy.catalog.entitlement"
	CatalogItemIdentifier    string = "CATALOG_ITEM_IDENTIFIER"
	CatalogSourceIdentifier  string = "CATALOG_SOURCE_IDENTIFIER"
	EnforcementTypeHard      string = "HARD"
	EnforcementTypeSoft      string = "SOFT"
)

type ContentSharingItem struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

type ContentSharingPrincipal struct {
	ReferenceID string `json:"referenceId,omitempty"`
	Type        string `json:"type,omitempty"`
}

type ContentSharingEntitledUser struct {
	UserType   string                    `json:"userType,omitempty"`
	Items      []ContentSharingItem      `json:"items,omitempty"`
	Principals []ContentSharingPrincipal `json:"principals,omitempty"`
}

type ContentSharingDefinition struct {
	EntitledUsers []ContentSharingEntitledUser `json:"entitledUsers,omitempty"`
}

func buildContentSharingPolicyDefinition(catalogItemIDs []string, catalogSourceIDs []string, projectID string) ContentSharingDefinition {
	contentSharingItems := make([]ContentSharingItem, 0, len(catalogItemIDs)+len(catalogSourceIDs))
	for _, catalogItemID := range catalogItemIDs {
		contentSharingItem := ContentSharingItem{
			ID:   catalogItemID,
			Type: CatalogItemIdentifier,
		}
		contentSharingItems = append(contentSharingItems, contentSharingItem)
	}
	for _, catalogSourceID := range catalogSourceIDs {
		contentSharingItem := ContentSharingItem{
			ID:   catalogSourceID,
			Type: CatalogSourceIdentifier,
		}
		contentSharingItems = append(contentSharingItems, contentSharingItem)
	}
	entitledUser := ContentSharingEntitledUser{
		UserType: "USER",
		Items:    contentSharingItems,
		Principals: []ContentSharingPrincipal{
			{
				ReferenceID: projectID,
				Type:        "PROJECT",
			},
		},
	}
	return ContentSharingDefinition{
		EntitledUsers: []ContentSharingEntitledUser{entitledUser},
	}
}

func extractCatalogItemIDsFromContentSharingPolicy(definition interface{}) ([]string, error) {
	catalogItemIDs := make([]string, 0)

	definitionJSON, err := json.Marshal(definition)
	if err != nil {
		return catalogItemIDs, err
	}

	var contentSharingDefinition ContentSharingDefinition
	if err := json.Unmarshal(definitionJSON, &contentSharingDefinition); err != nil {
		return catalogItemIDs, err
	}

	for _, entitledUser := range contentSharingDefinition.EntitledUsers {
		for _, item := range entitledUser.Items {
			if item.Type == CatalogItemIdentifier {
				catalogItemIDs = append(catalogItemIDs, item.ID)
			}
		}
	}

	return catalogItemIDs, nil
}

func extractCatalogSourceIDsFromContentSharingPolicy(definition interface{}) ([]string, error) {
	catalogSourceIDs := make([]string, 0)

	definitionJSON, err := json.Marshal(definition)
	if err != nil {
		return catalogSourceIDs, err
	}

	var contentSharingDefinition ContentSharingDefinition
	if err := json.Unmarshal(definitionJSON, &contentSharingDefinition); err != nil {
		return catalogSourceIDs, err
	}

	for _, entitledUser := range contentSharingDefinition.EntitledUsers {
		for _, item := range entitledUser.Items {
			if item.Type == CatalogSourceIdentifier {
				catalogSourceIDs = append(catalogSourceIDs, item.ID)
			}
		}
	}

	return catalogSourceIDs, nil
}
