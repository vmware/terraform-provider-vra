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

func resourcePolicyLease() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyLeaseCreate,
		ReadContext:   resourcePolicyLeaseRead,
		UpdateContext: resourcePolicyLeaseUpdate,
		DeleteContext: resourcePolicyLeaseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required arguments
			"enforcement_type": {
				Type:         schema.TypeString,
				Description:  "The type of enforcement for the policy.",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"HARD", "SOFT"}, true),
			},
			"lease_term_max": {
				Type:         schema.TypeInt,
				Description:  "The maximum duration in days between creation (or renewal) and expiration.",
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 32767),
			},
			"lease_total_term_max": {
				Type:         schema.TypeInt,
				Description:  "The maximum duration in days between creation and expiration. Unaffected by renewal.",
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 32767),
			},
			"name": {
				Type:        schema.TypeString,
				Description: "A human-friendly name used as an identifier for the policy instance.",
				Required:    true,
			},

			// Optional arguments
			"criteria": {
				Type:        schema.TypeSet,
				Description: "The policy criteria.",
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
				Optional: true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "A human-friendly description for the policy instance.",
				Optional:    true,
			},
			"lease_grace": {
				Type:         schema.TypeInt,
				Description:  "The duration in days that an expired object should be held before it is deleted.",
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 127),
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
				Type:        schema.TypeString,
				Description: "The id of the project this entity belongs to.",
				ForceNew:    true,
				Optional:    true,
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

func resourcePolicyLeaseCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	definition := &PolicyLeaseDefinition{
		LeaseTermMax:      d.Get("lease_term_max").(int),
		LeaseTotalTermMax: d.Get("lease_total_term_max").(int),
	}
	if leaseGrace, ok := d.GetOk("lease_grace"); ok {
		definition.LeaseGrace = withInt(leaseGrace.(int))
	}

	_, createdResp, err := apiClient.Policies.CreatePolicyUsingPOST1(policies.NewCreatePolicyUsingPOST1Params().WithPolicy(&models.Policy{
		Criteria:        expandPolicyCriteria(d.Get("criteria").(*schema.Set).List()),
		Definition:      definition,
		Description:     d.Get("description").(string),
		EnforcementType: d.Get("enforcement_type").(string),
		Name:            d.Get("name").(string),
		ProjectID:       d.Get("project_id").(string),
		ScopeCriteria:   expandPolicyCriteria(d.Get("project_criteria").(*schema.Set).List()),
		TypeID:          withString(PolicyLeaseTypeID),
	}))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdResp.Payload.ID.String())

	return resourcePolicyLeaseRead(ctx, d, m)
}

func resourcePolicyLeaseRead(_ context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
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
	if *policy.TypeID != PolicyLeaseTypeID {
		return diag.Errorf("policy with id `%s` is not a lease policy", id)
	}

	d.SetId(policy.ID.String())
	_ = d.Set("created_at", policy.CreatedAt.String())
	_ = d.Set("created_by", policy.CreatedBy)
	_ = d.Set("enforcement_type", policy.EnforcementType)
	_ = d.Set("last_updated_at", policy.LastUpdatedAt.String())
	_ = d.Set("last_updated_by", policy.LastUpdatedBy)
	_ = d.Set("name", policy.Name)
	_ = d.Set("org_id", policy.OrgID)

	if policy.Criteria != nil {
		_ = d.Set("criteria", flattenPolicyCriteria(*policy.Criteria))
	}
	if policy.Description != "" {
		_ = d.Set("description", policy.Description)
	}
	if policy.ScopeCriteria != nil {
		_ = d.Set("project_criteria", flattenPolicyCriteria(*policy.ScopeCriteria))
	}
	if policy.ProjectID != "" {
		_ = d.Set("project_id", policy.ProjectID)
	}

	var definition PolicyLeaseDefinition
	if err := policyDefinitionConvert(policy.Definition, &definition); err != nil {
		return diag.FromErr(err)
	}

	if definition.LeaseGrace != nil {
		_ = d.Set("lease_grace", definition.LeaseGrace)
	}
	_ = d.Set("lease_term_max", definition.LeaseTermMax)
	_ = d.Set("lease_total_term_max", definition.LeaseTotalTermMax)

	return nil
}

func resourcePolicyLeaseUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	definition := &PolicyLeaseDefinition{
		LeaseTermMax:      d.Get("lease_term_max").(int),
		LeaseTotalTermMax: d.Get("lease_total_term_max").(int),
	}
	if leaseGrace, ok := d.GetOk("lease_grace"); ok {
		definition.LeaseGrace = withInt(leaseGrace.(int))
	}

	_, _, err := apiClient.Policies.CreatePolicyUsingPOST1(policies.NewCreatePolicyUsingPOST1Params().WithPolicy(&models.Policy{
		Criteria:        expandPolicyCriteria(d.Get("criteria").(*schema.Set).List()),
		Definition:      definition,
		Description:     d.Get("description").(string),
		EnforcementType: d.Get("enforcement_type").(string),
		ID:              strfmt.UUID(id),
		Name:            d.Get("name").(string),
		ProjectID:       d.Get("project_id").(string),
		ScopeCriteria:   expandPolicyCriteria(d.Get("project_criteria").(*schema.Set).List()),
		TypeID:          withString(PolicyLeaseTypeID),
	}))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePolicyLeaseRead(ctx, d, m)
}

func resourcePolicyLeaseDelete(_ context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	if _, err := apiClient.Policies.DeletePolicyUsingDELETE5(policies.NewDeletePolicyUsingDELETE5Params().WithID(strfmt.UUID(id))); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
