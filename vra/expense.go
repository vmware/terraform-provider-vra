// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

// expenseSchema returns the schema to use for the expense property
func expenseSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"additional_expense": {
					Type:     schema.TypeFloat,
					Optional: true,
					Computed: true,
				},
				"code": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"compute_expense": {
					Type:     schema.TypeFloat,
					Optional: true,
					Computed: true,
				},
				"last_update_time": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"message": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"network_expense": {
					Type:     schema.TypeFloat,
					Optional: true,
					Computed: true,
				},
				"storage_expense": {
					Type:     schema.TypeFloat,
					Optional: true,
					Computed: true,
				},
				"total_expense": {
					Type:     schema.TypeFloat,
					Optional: true,
					Computed: true,
				},
				"unit": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func flattenExpense(expense *models.Expense) []interface{} {
	if expense == nil {
		return make([]interface{}, 0)
	}

	helper := make(map[string]interface{})

	helper["additional_expense"] = expense.AdditionalExpense
	helper["code"] = expense.Code
	helper["compute_expense"] = expense.ComputeExpense
	helper["last_update_time"] = expense.LastUpdatedTime.String()
	helper["message"] = expense.Message
	helper["network_expense"] = expense.NetworkExpense
	helper["storage_expense"] = expense.StorageExpense
	helper["total_expense"] = expense.TotalExpense
	helper["unit"] = expense.Unit

	return []interface{}{helper}
}
