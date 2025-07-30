// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import "encoding/json"

// ContentSourceRepositoryConfig - Config fields for linking an SCM integration
// with a repository and project not included in the swagger model so we're hand crafting it
type ContentSourceRepositoryConfig struct {
	Path          string `json:"path,omitempty"`
	Branch        string `json:"branch,omitempty"`
	Repository    string `json:"repository,omitempty"`
	ContentType   string `json:"contentType,omitempty"`
	ProjectName   string `json:"projectName,omitempty"`
	IntegrationID string `json:"integrationId,omitempty"`
}

// treating the config elem as an array rather than a singleton
func expandContentSourceRepositoryConfig(repoConfigs []interface{}) []*ContentSourceRepositoryConfig {
	configs := make([]*ContentSourceRepositoryConfig, 0, len(repoConfigs))

	for _, repo := range repoConfigs {
		config := repo.(map[string]interface{})

		cfg := ContentSourceRepositoryConfig{
			Path:          config["path"].(string),
			Branch:        config["branch"].(string),
			Repository:    config["repository"].(string),
			ContentType:   config["content_type"].(string),
			ProjectName:   config["project_name"].(string),
			IntegrationID: config["integration_id"].(string),
		}

		configs = append(configs, &cfg)
	}
	return configs
}

func flattenContentsourceRepositoryConfig(configSpec interface{}) ([]map[string]interface{}, error) {
	if configSpec == nil {
		return make([]map[string]interface{}, 0), nil
	}

	configJSON, err := json.Marshal(configSpec)
	if err != nil {
		return nil, err
	}

	var config ContentSourceRepositoryConfig
	if err := json.Unmarshal(configJSON, &config); err != nil {
		return nil, err
	}

	helper := make(map[string]interface{})
	helper["path"] = config.Path
	helper["branch"] = config.Branch
	helper["repository"] = config.Repository
	helper["content_type"] = config.ContentType
	helper["project_name"] = config.ProjectName
	helper["integration_id"] = config.IntegrationID

	return []map[string]interface{}{helper}, nil
}
