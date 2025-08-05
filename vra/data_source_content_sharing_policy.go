// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/vmware/vra-sdk-go/pkg/client/policies"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceContentSharingPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContentSharingPolicyRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Input attributes
			"id": {
				Type:          schema.TypeString,
				Computed:      true,
				ConflictsWith: []string{"name"},
				Description:   "The id of the policy instance.",
				Optional:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Computed:      true,
				ConflictsWith: []string{"id"},
				Description:   "The name of the policy instance.",
				Optional:      true,
			},

			// Computed attributes
			"catalog_item_ids": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of catalog item ids to share.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"catalog_source_ids": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of catalog source ids to share.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was created. The date is in ISO 8601 and UTC.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user the entity was created by.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly description for the policy instance.",
			},
			"enforcement_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of enforcement for the policy.",
			},
			"entitlement_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Entitlement type.",
			},
			"last_updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC.",
			},
			"last_updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user the entity was last updated by.",
			},
			"org_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the organization this entity belongs to.",
			},
			"principals": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of users or roles that can share content.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reference_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The reference ID of the principal.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the principal.",
						},
					},
				},
			},
			"project_criteria": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "The project based criteria.",
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the project this entity belongs to.",
			},
		},
	}
}

func dataSourceContentSharingPolicyRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")
	if !idOk && !nameOk {
		return diag.Errorf("one of id or name must be provided")
	}

	var policy *models.Policy
	if id != "" {
		getResp, err := apiClient.Policies.GetPolicyUsingGET5(policies.NewGetPolicyUsingGET5Params().WithID(strfmt.UUID(id.(string))))
		if err != nil {
			switch err.(type) {
			case *policies.GetPolicyUsingGET5NotFound:
				return diag.Errorf("policy with id `%s` not found", id)
			default:
				// nop
			}
			return diag.FromErr(err)
		}

		policy = getResp.GetPayload()
		if *policy.TypeID != PolicyCatalogEntitlementTypeID {
			return diag.Errorf("policy with id `%s` is not a content sharing policy", id)
		}
	} else {
		getResp, err := apiClient.Policies.GetPoliciesUsingGET5(policies.NewGetPoliciesUsingGET5Params().WithTypeID(PolicyCatalogEntitlementTypeID).WithExpandDefinition(withBool(true)).WithSearch(withString(name.(string))))
		if err != nil {
			return diag.FromErr(err)
		}

		policies := getResp.Payload
		if len(policies.Content) == 0 {
			return diag.Errorf("vra_content_sharing_policy `name` criteria did not match any policy")
		}
		if len(policies.Content) > 1 {
			return diag.Errorf("vra_content_sharing_policy `name` criteria must filter to a single policy")
		}

		policy = policies.Content[0]
	}

	d.SetId(policy.ID.String())
	_ = d.Set("created_at", policy.CreatedAt.String())
	_ = d.Set("created_by", policy.CreatedBy)
	_ = d.Set("enforcement_type", policy.EnforcementType)
	_ = d.Set("last_updated_at", policy.LastUpdatedAt.String())
	_ = d.Set("last_updated_by", policy.LastUpdatedBy)
	_ = d.Set("name", policy.Name)
	_ = d.Set("org_id", policy.OrgID)

	if policy.Description != "" {
		_ = d.Set("description", policy.Description)
	}
	if policy.ScopeCriteria != nil {
		_ = d.Set("project_criteria", flattenPolicyCriteria(*policy.ScopeCriteria))
	}
	if policy.ProjectID != "" {
		_ = d.Set("project_id", policy.ProjectID)
	}

	var definition PolicyContentSharingDefinition
	if err := policyDefinitionConvert(policy.Definition, &definition); err != nil {
		return diag.FromErr(err)
	}

	catalogItemIDs := make([]string, 0)
	catalogSourceIDs := make([]string, 0)
	principalsMap := make([]any, 0)
	for _, entitledUser := range definition.EntitledUsers {
		_ = d.Set("entitlement_type", entitledUser.UserType)

		for _, item := range entitledUser.Items {
			if item.Type == CatalogItemIdentifier {
				catalogItemIDs = append(catalogItemIDs, item.ID)
			}
			if item.Type == CatalogSourceIdentifier {
				catalogSourceIDs = append(catalogSourceIDs, item.ID)
			}
		}

		for _, principal := range entitledUser.Principals {
			helper := make(map[string]any)
			helper["reference_id"] = principal.ReferenceID
			helper["type"] = principal.Type
			principalsMap = append(principalsMap, helper)
		}
	}
	if err := d.Set("catalog_item_ids", catalogItemIDs); err != nil {
		return diag.Errorf("error setting catalog_item_ids: %s", err.Error())
	}
	if err := d.Set("catalog_source_ids", catalogSourceIDs); err != nil {
		return diag.Errorf("error setting catalog_source_ids: %s", err.Error())
	}
	if err := d.Set("principals", principalsMap); err != nil {
		return diag.Errorf("error setting principals: %s", err.Error())
	}

	return nil
}
