// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/vmware/vra-sdk-go/pkg/client/policies"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func resourceContentSharingPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContentSharingPolicyCreate,
		ReadContext:   resourceContentSharingPolicyRead,
		UpdateContext: resourceContentSharingPolicyUpdate,
		DeleteContext: resourceContentSharingPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Required arguments
			"catalog_item_ids": {
				Type:        schema.TypeSet,
				Description: "List of catalog item ids to share.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:     true,
				AtLeastOneOf: []string{"catalog_source_ids"},
			},
			"catalog_source_ids": {
				Type:        schema.TypeSet,
				Description: "List of catalog source ids to share.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:     true,
				AtLeastOneOf: []string{"catalog_item_ids"},
			},
			"name": {
				Type:        schema.TypeString,
				Description: "A human-friendly name used as an identifier for the policy instance.",
				Required:    true,
			},

			// Optional arguments
			"description": {
				Type:        schema.TypeString,
				Description: "A human-friendly description for the policy instance.",
				Optional:    true,
			},
			"entitlement_type": {
				Type:         schema.TypeString,
				Description:  "Entitlement type.",
				Optional:     true,
				RequiredWith: []string{"principals"},
				AtLeastOneOf: []string{"project_id"},
				ValidateFunc: validation.StringInSlice([]string{"USER", "ROLE"}, true),
			},
			"principals": {
				Type:        schema.TypeSet,
				Description: "List of users or roles that can share content.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reference_id": {
							Type:        schema.TypeString,
							Description: "The reference ID of the principal.",
							Optional:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "The type of the principal.",
							Required:    true,
						},
					},
				},
				Optional:     true,
				RequiredWith: []string{"entitlement_type"},
			},
			"project_criteria": {
				Type:          schema.TypeSet,
				ConflictsWith: []string{"project_id"},
				Description:   "The project based criteria.",
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
				ForceNew: true,
				Optional: true,
			},
			"project_id": {
				Type:         schema.TypeString,
				Description:  "The id of the project this entity belongs to.",
				ForceNew:     true,
				Optional:     true,
				AtLeastOneOf: []string{"entitlement_type"},
			},

			// Computed attributes
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
			"enforcement_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of enforcement for the policy.",
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
		},
	}
}

func resourceContentSharingPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	items := make([]PolicyContentSharingItem, 0)
	if v, ok := d.GetOk("catalog_item_ids"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.Errorf("catalog_item_ids must be unique")
		}
		for _, catalogItemID := range expandStringList(v.(*schema.Set).List()) {
			items = append(items, PolicyContentSharingItem{
				ID:   catalogItemID,
				Type: CatalogItemIdentifier,
			})
		}
	}
	if v, ok := d.GetOk("catalog_source_ids"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.Errorf("catalog_source_ids must be unique")
		}
		for _, catalogSourceID := range expandStringList(v.(*schema.Set).List()) {
			items = append(items, PolicyContentSharingItem{
				ID:   catalogSourceID,
				Type: CatalogSourceIdentifier,
			})
		}
	}

	var entitledUsers []PolicyContentSharingEntitledUser
	_, projectIDOk := d.GetOk("project_id")
	_, entitlementTypeOk := d.GetOk("entitlement_type")
	if !projectIDOk && !entitlementTypeOk {
		return diag.Errorf("`entitlement_type` or `project_id` must be specified")
	} else if projectIDOk && !entitlementTypeOk {
		// For backwards compatibility, we will share the content with all users and groups
		// in the project if the project_id is specified and entitlement_type is not specified.
		// Warning: A ptoperties drift will be detected every time the user runs an apply command.
		entitledUsers = []PolicyContentSharingEntitledUser{
			{
				UserType: "USER",
				Items:    items,
				Principals: []PolicyContentSharingPrincipal{
					{
						ReferenceID: "",
						Type:        "PROJECT",
					},
				},
			},
		}
	} else {
		entitledUsers = []PolicyContentSharingEntitledUser{
			{
				UserType:   d.Get("entitlement_type").(string),
				Items:      items,
				Principals: expandPolicyContentSharingPrincipal(d.Get("principals").(*schema.Set).List()),
			},
		}
	}

	_, createdResp, err := apiClient.Policies.CreatePolicyUsingPOST1(policies.NewCreatePolicyUsingPOST1Params().WithPolicy(&models.Policy{
		Definition: &PolicyContentSharingDefinition{
			EntitledUsers: entitledUsers,
		},
		Description:     d.Get("description").(string),
		EnforcementType: EnforcementTypeHard,
		Name:            d.Get("name").(string),
		ProjectID:       d.Get("project_id").(string),
		ScopeCriteria:   expandPolicyCriteria(d.Get("project_criteria").(*schema.Set).List()),
		TypeID:          withString(PolicyCatalogEntitlementTypeID),
	}))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdResp.Payload.ID.String())

	return resourceContentSharingPolicyRead(ctx, d, m)
}

func resourceContentSharingPolicyRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	getResp, err := apiClient.Policies.GetPolicyUsingGET5(policies.NewGetPolicyUsingGET5Params().WithID(strfmt.UUID(id)))
	if err != nil {
		switch err.(type) {
		case *policies.GetPolicyUsingGET5NotFound:
			return diag.Errorf("policy with id `%s` not found", id)
		default:
			// nop
		}
		return diag.FromErr(err)
	}

	policy := getResp.GetPayload()
	if *policy.TypeID != PolicyCatalogEntitlementTypeID {
		return diag.Errorf("policy with id `%s` is not a content sharing policy", id)
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

func resourceContentSharingPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	items := make([]PolicyContentSharingItem, 0)
	if v, ok := d.GetOk("catalog_item_ids"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.Errorf("catalog_item_ids must be unique")
		}
		for _, catalogItemID := range expandStringList(v.(*schema.Set).List()) {
			items = append(items, PolicyContentSharingItem{
				ID:   catalogItemID,
				Type: CatalogItemIdentifier,
			})
		}
	}
	if v, ok := d.GetOk("catalog_source_ids"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.Errorf("catalog_source_ids must be unique")
		}
		for _, catalogSourceID := range expandStringList(v.(*schema.Set).List()) {
			items = append(items, PolicyContentSharingItem{
				ID:   catalogSourceID,
				Type: CatalogSourceIdentifier,
			})
		}
	}
	var entitledUsers []PolicyContentSharingEntitledUser
	_, projectIDOk := d.GetOk("project_id")
	_, entitlementTypeOk := d.GetOk("entitlement_type")
	if !projectIDOk && !entitlementTypeOk {
		return diag.Errorf("`entitlement_type` or `project_id` must be specified")
	} else if projectIDOk && !entitlementTypeOk {
		// For backwards compatibility, we will share the content with all users and groups
		// in the project if the project_id is specified and entitlement_type is not specified.
		entitledUsers = []PolicyContentSharingEntitledUser{
			{
				UserType: "USER",
				Items:    items,
				Principals: []PolicyContentSharingPrincipal{
					{
						ReferenceID: "",
						Type:        "PROJECT",
					},
				},
			},
		}
	} else {
		entitledUsers = []PolicyContentSharingEntitledUser{
			{
				UserType:   d.Get("entitlement_type").(string),
				Items:      items,
				Principals: expandPolicyContentSharingPrincipal(d.Get("principals").(*schema.Set).List()),
			},
		}
	}

	_, _, err := apiClient.Policies.CreatePolicyUsingPOST1(policies.NewCreatePolicyUsingPOST1Params().WithPolicy(&models.Policy{
		Definition: &PolicyContentSharingDefinition{
			EntitledUsers: entitledUsers,
		},
		Description:     d.Get("description").(string),
		EnforcementType: EnforcementTypeHard,
		ID:              strfmt.UUID(id),
		Name:            d.Get("name").(string),
		ProjectID:       d.Get("project_id").(string),
		ScopeCriteria:   expandPolicyCriteria(d.Get("project_criteria").(*schema.Set).List()),
		TypeID:          withString(PolicyCatalogEntitlementTypeID),
	}))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceContentSharingPolicyRead(ctx, d, m)
}

func resourceContentSharingPolicyDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	if _, err := apiClient.Policies.DeletePolicyUsingDELETE5(policies.NewDeletePolicyUsingDELETE5Params().WithID(strfmt.UUID(id))); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
