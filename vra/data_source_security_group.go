package vra

import (
	"fmt"

	"github.com/vmware/vra-sdk-go/pkg/client/security_group"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSecurityGroupRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeString,
				Required: true,
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"egress": rulesSchema(false),
			"external_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_region_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ingress": rulesSchema(false),
			"links":   linksSchema(),
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": &schema.Schema{
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
	d.Set("egress", securityGroup.Egress)
	d.Set("external_id", securityGroup.ExternalID)
	d.Set("external_region_id", securityGroup.ExternalRegionID)
	d.Set("ingress", securityGroup.Ingress)
	d.Set("name", securityGroup.Name)
	d.Set("organization_id", securityGroup.OrganizationID)
	d.Set("owner", securityGroup.Owner)
	d.Set("updated_at", securityGroup.UpdatedAt)

	return nil
}
