// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/fabric_vsphere_datastore"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceFabricDatastoreVsphere() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVsphereDatastoreRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Computed:      true,
				ConflictsWith: []string{"filter"},
				Description:   "The id of the vSphere fabric datastore resource instance.",
				Optional:      true,
			},
			"filter": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"id"},
				Description:   "Search criteria to narrow down the vSphere fabric datastore resource instance.",
				Optional:      true,
			},
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

func dataSourceVsphereDatastoreRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading the vSphere fabric datastore data source with id %s or filter %s", d.Get("id"), d.Get("filter"))
	apiClient := meta.(*Client).apiClient

	id := d.Get("id").(string)
	filter := d.Get("filter").(string)

	if id == "" && filter == "" {
		return fmt.Errorf("one of id or filter is required")
	}

	var fabricVsphereDatastore *models.FabricVsphereDatastore
	if id != "" {
		log.Printf("Reading vSphere fabric datastore data source with id: %s", id)
		getResp, err := apiClient.FabricvSphereDatastore.GetFabricVSphereDatastore(
			fabric_vsphere_datastore.NewGetFabricVSphereDatastoreParams().WithID(id))
		if err != nil {
			switch err.(type) {
			case *fabric_vsphere_datastore.GetFabricVSphereDatastoreNotFound:
				return fmt.Errorf("vSphere fabric datastore '%s' not found", id)
			default:
				// nop
			}
			return err
		}

		fabricVsphereDatastore = getResp.GetPayload()
	} else {
		log.Printf("Reading vSphere fabric datastore data source with filter: %s", filter)
		getResp, err := apiClient.FabricvSphereDatastore.GetFabricVSphereDatastores(
			fabric_vsphere_datastore.NewGetFabricVSphereDatastoresParams().WithDollarFilter(withString(filter)))
		if err != nil {
			return err
		}

		datastores := *getResp.Payload
		if len(datastores.Content) > 1 {
			return fmt.Errorf("must filter to one vSphere fabric datastore")
		}
		if len(datastores.Content) == 0 {
			return fmt.Errorf("filter doesn't match to any vSphere fabric datastore")
		}

		fabricVsphereDatastore = datastores.Content[0]
	}

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
		return fmt.Errorf("error setting vSphere fabric datastore links - error: %#v", err)
	}

	if err := d.Set("tags", flattenTags(fabricVsphereDatastore.Tags)); err != nil {
		return fmt.Errorf("error setting vSphere fabric datastore tags - error: %v", err)
	}

	log.Println("Finished reading the vSphere fabric datastore data source")
	return nil
}
