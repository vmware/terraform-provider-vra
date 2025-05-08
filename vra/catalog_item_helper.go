// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import "encoding/json"

const (
	CatalogItemVMImageTypeID     string = "com.vmw.vmimage"
	CatalogItemVroWorkflowTypeID string = "com.vmw.vro.workflow"
)

type CatalogItemVMImagePublishSpec struct {
	CloudConfig *string `json:"cloudConfig,omitempty"`
	ImageName   string  `json:"imageName"`
	SelectZone  *bool   `json:"selectZone,omitempty"`
}

type CatalogItemVroWorkflowPublishSpec struct {
	WorkflowID string `json:"workflowId"`
}

func catalogItemSpecConvert(genericSpec any, castedSpec any) error {
	specJSON, err := json.Marshal(genericSpec)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(specJSON, &castedSpec); err != nil {
		return err
	}

	return nil
}
