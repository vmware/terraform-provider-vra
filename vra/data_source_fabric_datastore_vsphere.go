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
			"cloud_account_ids": {
				Type:     schema.TypeList,
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
			"free_size_gb": {
				Type:     schema.TypeString,
				Computed: true,
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
			"type": {
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

func dataSourceVsphereDatastoreRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading the vra vsphere datastore data source with filter %s", d.Get("filter"))
	apiClient := meta.(*Client).apiClient

	var datastore *models.FabricVsphereDatastore

	id := d.Get("id").(string)
	filter := d.Get("filter").(string)

	if id == "" && filter == "" {
		return fmt.Errorf("one of id or filter is required")
	}

	if id != "" {
		log.Printf("Reading fabric vSphere datastore data source with id: %s", id)
		getResp, err := apiClient.FabricvSphereDatastore.GetFabricVSphereDatastore(
			fabric_vsphere_datastore.NewGetFabricVSphereDatastoreParams().WithID(id))

		if err != nil {
			return err
		}
		datastore = getResp.GetPayload()
	} else {
		log.Printf("Reading fabric vSphere datastore data source with filter: %s", filter)
		getResp, err := apiClient.FabricvSphereDatastore.GetFabricVSphereDatastores(
			fabric_vsphere_datastore.NewGetFabricVSphereDatastoresParams().WithDollarFilter(withString(filter)))
		if err != nil {
			return err
		}

		datastores := *getResp.Payload
		if len(datastores.Content) > 1 {
			return fmt.Errorf("fabric vSphere datastore must filter to one datastore")
		}
		if len(datastores.Content) == 0 {
			return fmt.Errorf("fabric vSphere datastore filter doesn't match to any datastore")
		}

		datastore = datastores.Content[0]
	}

	d.SetId(*datastore.ID)
	d.Set("cloud_account_ids", datastore.CloudAccountIds)
	d.Set("created_at", datastore.CreatedAt)
	d.Set("external_id", datastore.ExternalID)
	d.Set("external_region_id", datastore.ExternalRegionID)
	d.Set("free_size_gb", datastore.FreeSizeGB)
	d.Set("name", datastore.Name)
	d.Set("org_id", datastore.OrgID)
	d.Set("type", datastore.Type)
	d.Set("updated_at", datastore.UpdatedAt)

	if err := d.Set("links", flattenLinks(datastore.Links)); err != nil {
		return fmt.Errorf("error setting datastore links - error: %#v", err)
	}

	log.Println("Finished reading the datastore data source")
	return nil
}
