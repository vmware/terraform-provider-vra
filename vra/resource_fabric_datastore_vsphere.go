// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"
	"errors"

	"github.com/vmware/vra-sdk-go/pkg/client/fabric_vsphere_datastore"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFabricDatastoreVsphere() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFabricDatastoreVsphereCreate,
		ReadContext:   resourceFabricDatastoreVsphereRead,
		UpdateContext: resourceFabricDatastoreVsphereUpdate,
		DeleteContext: resourceFabricDatastoreVsphereDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"cloud_account_ids": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Set of ids of the cloud accounts this entity belongs to.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was created. The date is in ISO 8601 and UTC.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly description.",
			},
			"external_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "External entity Id on the provider side.",
			},
			"external_region_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Id of datacenter in which the datastore is present.",
			},
			"free_size_gb": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates free size available in datastore.",
			},
			"links": linksSchema(),
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly name used as an identifier for the vSphere fabric datastore resource instance.",
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
			"tags": tagsSchema(),
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of datastore.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC.",
			},
		},
	}
}

func resourceFabricDatastoreVsphereCreate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return diag.FromErr(errors.New("vra_fabric_datastore_vsphere resources are only importable"))
}

func resourceFabricDatastoreVsphereRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	resp, err := apiClient.FabricvSphereDatastore.GetFabricVSphereDatastore(fabric_vsphere_datastore.NewGetFabricVSphereDatastoreParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *fabric_vsphere_datastore.GetFabricVSphereDatastoreNotFound:
			return diag.Errorf("vSphere fabric datastore '%s' not found", id)
		default:
			// nop
		}
		return diag.FromErr(err)
	}

	fabricVsphereDatastore := resp.GetPayload()
	d.SetId(*fabricVsphereDatastore.ID)
	d.Set("cloud_account_ids", fabricVsphereDatastore.CloudAccountIds)
	d.Set("created_at", fabricVsphereDatastore.CreatedAt)
	d.Set("description", fabricVsphereDatastore.Description)
	d.Set("external_id", fabricVsphereDatastore.ExternalID)
	d.Set("external_region_id", fabricVsphereDatastore.ExternalRegionID)
	d.Set("free_size_gb", fabricVsphereDatastore.FreeSizeGB)
	d.Set("name", fabricVsphereDatastore.Name)
	d.Set("org_id", fabricVsphereDatastore.OrgID)
	d.Set("owner", fabricVsphereDatastore.Owner)
	d.Set("type", fabricVsphereDatastore.Type)
	d.Set("updated_at", fabricVsphereDatastore.UpdatedAt)

	if err := d.Set("links", flattenLinks(fabricVsphereDatastore.Links)); err != nil {
		return diag.Errorf("error setting vSphere fabric datastore links - error: %#v", err)
	}

	if err := d.Set("tags", flattenTags(fabricVsphereDatastore.Tags)); err != nil {
		return diag.Errorf("error setting vSphere fabric datastore tags - error: %v", err)
	}

	return nil
}

func resourceFabricDatastoreVsphereUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()

	fabricVsphereDatastoreSpecification := &models.FabricVsphereDatastoreSpecification{
		Tags: expandTags(d.Get("tags").(*schema.Set).List()),
	}

	_, err := apiClient.FabricvSphereDatastore.UpdateFabricVsphereDatastore(fabric_vsphere_datastore.NewUpdateFabricVsphereDatastoreParams().WithID(id).WithBody(fabricVsphereDatastoreSpecification))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceFabricDatastoreVsphereRead(ctx, d, m)
}

func resourceFabricDatastoreVsphereDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
