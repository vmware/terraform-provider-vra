// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/client/data_collector"
)

func dataSourceDataCollector() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDataCollectorRead,

		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDataCollectorRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	name := d.Get("name").(string)

	getResp, err := apiClient.DataCollector.GetDataCollectors(data_collector.NewGetDataCollectorsParams())
	if err != nil {
		return err
	}

	dataCollectors := getResp.Payload
	if dataCollectors.TotalElements == 0 {
		return fmt.Errorf("No vra_data_collectors found")
	}

	for _, dc := range dataCollectors.Content {
		if *dc.Name == name {
			d.Set("ip_address", dc.IPAddress)
			d.Set("hostname", dc.HostName)
			d.Set("status", dc.Status)
			d.SetId(*dc.DcID)
			return nil
		}
	}

	return fmt.Errorf("vra_data_collector \"%s\" not found", name)
}
