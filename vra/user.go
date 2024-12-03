// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

// userSchema returns the schema to use for the administrator / member / viewer property in Project
func userSchema(description string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: description,
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:        schema.TypeString,
					Default:     "user",
					Optional:    true,
					Description: "Type of the principal. Currently supported ‘user’ (default) and 'group’.",
				},
				"email": {
					Type:        schema.TypeString,
					Description: "The email of the user or name of the group.",
					Required:    true,
				},
			},
		},
	}
}

func expandUsers(configUsers []interface{}) []*models.User {
	users := make([]*models.User, 0, len(configUsers))

	for _, configUser := range configUsers {
		userMap := configUser.(map[string]interface{})

		var user models.User

		if v, found := userMap["type"]; found && v != nil {
			user.Type = v.(string)
		}

		if v, found := userMap["email"]; found && v != nil {
			user.Email = withString(v.(string))
		}

		users = append(users, &user)
	}

	return users
}

func flattenUsers(users []*models.User) []interface{} {
	if len(users) == 0 {
		return make([]interface{}, 0)
	}

	configUsers := make([]interface{}, 0, len(users))

	for _, user := range users {
		helper := make(map[string]interface{})
		helper["email"] = *user.Email
		helper["type"] = user.Type

		configUsers = append(configUsers, helper)
	}

	return configUsers
}
