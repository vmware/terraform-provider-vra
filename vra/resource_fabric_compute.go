// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"
	"errors"

	"github.com/vmware/vra-sdk-go/pkg/client/fabric_compute"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFabricCompute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFabricComputeCreate,
		ReadContext:   resourceFabricComputeRead,
		UpdateContext: resourceFabricComputeUpdate,
		DeleteContext: resourceFabricComputeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was created. The date is in ISO 8601 and UTC.",
			},
			"custom_properties": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "A list of key value pair of custom properties for the fabric compute resource.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly description.",
			},
			"external_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the external entity on the provider side.",
			},
			"external_region_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The external region id of the fabric compute.",
			},
			"external_zone_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The external zone id of the fabric compute.",
			},
			"lifecycle_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Lifecycle status of the compute instance.",
			},
			"links": linksSchema(),
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly name used as an identifier for the fabric compute resource instance.",
			},
			"org_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the organization this entity belongs to.",
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email of the user that owns the entity.",
			},
			"power_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Power state of fabric compute instance.",
			},
			"tags": tagsSchema(),
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the fabric compute instance.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC.",
			},
		},
	}
}

func resourceFabricComputeCreate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return diag.FromErr(errors.New("vra_fabric_compute resources are only importable"))
}

func resourceFabricComputeRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	getResp, err := apiClient.FabricCompute.GetFabricCompute(fabric_compute.NewGetFabricComputeParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *fabric_compute.GetFabricComputeNotFound:
			return diag.Errorf("fabric compute '%s' not found", id)
		default:
			// nop
		}
		return diag.FromErr(err)
	}

	fabricCompute := getResp.GetPayload()
	d.SetId(*fabricCompute.ID)
	d.Set("created_at", fabricCompute.CreatedAt)
	d.Set("custom_properties", fabricCompute.CustomProperties)
	d.Set("description", fabricCompute.Description)
	d.Set("external_id", fabricCompute.ExternalID)
	d.Set("external_region_id", fabricCompute.ExternalRegionID)
	d.Set("external_zone_id", fabricCompute.ExternalZoneID)
	d.Set("lifecycle_state", fabricCompute.LifecycleState)
	d.Set("name", fabricCompute.Name)
	d.Set("org_id", fabricCompute.OrgID)
	d.Set("owner", fabricCompute.Owner)
	d.Set("power_state", fabricCompute.PowerState)
	d.Set("type", fabricCompute.Type)
	d.Set("updated_at", fabricCompute.UpdatedAt)

	if err := d.Set("links", flattenLinks(fabricCompute.Links)); err != nil {
		return diag.Errorf("error setting zone links - error: %#v", err)
	}

	if err := d.Set("tags", flattenTags(fabricCompute.Tags)); err != nil {
		return diag.Errorf("error setting zone tags - error: %v", err)
	}

	return nil
}

func resourceFabricComputeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	fabricComputeSpecification := &models.FabricComputeSpecification{
		Tags: expandTags(d.Get("tags").(*schema.Set).List()),
	}

	if _, err := apiClient.FabricCompute.UpdateFabricCompute(fabric_compute.NewUpdateFabricComputeParams().WithID(id).WithBody(fabricComputeSpecification)); err != nil {
		return diag.FromErr(err)
	}

	return resourceFabricComputeRead(ctx, d, m)
}

func resourceFabricComputeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
