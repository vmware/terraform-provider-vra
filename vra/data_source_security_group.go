// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"

	"github.com/vmware/vra-sdk-go/pkg/client/security_group"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSecurityGroupRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeString,
				Required: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
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
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": linksSchema(),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rules": rulesSchema(false),
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client).apiClient

	filter := d.Get("filter").(string)

	getResp, err := apiClient.SecurityGroup.GetSecurityGroups(security_group.NewGetSecurityGroupsParams().WithDollarFilter(withString(filter)))
	if err != nil {
		return err
	}

	securityGroups := getResp.Payload
	if len(securityGroups.Content) > 1 {
		return fmt.Errorf("vra_security_group must filter to a single security group")
	}
	if len(securityGroups.Content) == 0 {
		return fmt.Errorf("as_security_group filter did not match any security groups")
	}

	securityGroup := securityGroups.Content[0]
	d.SetId(*securityGroup.ID)
	d.Set("created_at", securityGroup.CreatedAt)
	d.Set("description", securityGroup.Description)
	d.Set("external_id", securityGroup.ExternalID)
	d.Set("external_region_id", securityGroup.ExternalRegionID)
	d.Set("name", securityGroup.Name)
	d.Set("organization_id", securityGroup.OrgID)
	d.Set("owner", securityGroup.Owner)
	d.Set("rules", securityGroup.Rules)
	d.Set("updated_at", securityGroup.UpdatedAt)

	return nil
}
