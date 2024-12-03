// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/fabric_vsphere_storage_policies"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func dataSourceFabricStoragePolicyVsphere() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFabricStoragePolicyVsphereRead,
		Schema: map[string]*schema.Schema{
			"cloud_account_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filter": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"filter"},
			},
			"links": linksSchema(),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceFabricStoragePolicyVsphereRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading the vra vsphere storage policies data source with filter %s", d.Get("filter"))
	apiClient := meta.(*Client).apiClient

	var storagePolicy *models.FabricVsphereStoragePolicy

	id := d.Get("id").(string)
	filter := d.Get("filter").(string)

	if id == "" && filter == "" {
		return fmt.Errorf("one of id or filter is required")
	}

	if id != "" {
		log.Printf("Reading fabric vSphere storage policy data source with id: %s", id)
		getResp, err := apiClient.FabricvSphereStoragePolicies.GetFabricVSphereStoragePolicy(
			fabric_vsphere_storage_policies.NewGetFabricVSphereStoragePolicyParams().WithID(id))

		if err != nil {
			return err
		}
		storagePolicy = getResp.GetPayload()
	} else {
		log.Printf("Reading fabric vSphere storage policies data source with filter: %s", filter)
		getResp, err := apiClient.FabricvSphereStoragePolicies.GetFabricVSphereStoragePolicies(
			fabric_vsphere_storage_policies.NewGetFabricVSphereStoragePoliciesParams().WithDollarFilter(withString(filter)))
		if err != nil {
			return err
		}

		storagePolicies := *getResp.Payload
		if len(storagePolicies.Content) > 1 {
			return fmt.Errorf("fabric vSphere storage policies must filter to one storage policy")
		}
		if len(storagePolicies.Content) == 0 {
			return fmt.Errorf("fabric vSphere storage policies filter doesn't match to any storage policy")
		}

		storagePolicy = storagePolicies.Content[0]
	}

	d.SetId(*storagePolicy.ID)
	d.Set("cloud_account_ids", storagePolicy.CloudAccountIds)
	d.Set("created_at", storagePolicy.CreatedAt)
	d.Set("external_id", storagePolicy.ExternalID)
	d.Set("external_region_id", storagePolicy.ExternalRegionID)
	d.Set("name", storagePolicy.Name)
	d.Set("org_id", storagePolicy.OrgID)
	d.Set("updated_at", storagePolicy.UpdatedAt)

	if err := d.Set("links", flattenLinks(storagePolicy.Links)); err != nil {
		return fmt.Errorf("error setting storage policy links - error: %#v", err)
	}

	log.Println("Finished reading the storage policy data source")
	return nil
}
