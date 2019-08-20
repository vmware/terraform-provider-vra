package vra

import (
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
)

// Provider represents the VRA provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VRA_URL", nil),
				Description: "The base url for API operations.",
			},
			"refresh_token": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"access_token"},
				DefaultFunc:   schema.EnvDefaultFunc("VRA_REFRESH_TOKEN", nil),
				Description:   "The refresh token for API operations.",
			},
			"access_token": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"refresh_token"},
				DefaultFunc:   schema.EnvDefaultFunc("VRA_ACCESS_TOKEN", nil),
				Description:   "The access token for API operations.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"vra_cloud_account_aws":   dataSourceCloudAccountAWS(),
			"vra_cloud_account_azure": dataSourceCloudAccountAzure(),
			"vra_fabric_network":      dataSourceFabricNetwork(),
			"vra_image":               dataSourceImage(),
			"vra_network":             dataSourceNetwork(),
			"vra_network_domain":      dataSourceNetworkDomain(),
			"vra_project":             dataSourceProject(),
			"vra_region":              dataSourceRegion(),
			"vra_security_group":      dataSourceSecurityGroup(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"vra_block_device":          resourceBlockDevice(),
			"vra_cloud_account_aws":     resourceCloudAccountAWS(),
			"vra_cloud_account_azure":   resourceCloudAccountAzure(),
			"vra_flavor_profile":        resourceFlavorProfile(),
			"vra_image_profile":         resourceImageProfile(),
			"vra_load_balancer":         resourceLoadBalancer(),
			"vra_machine":               resourceMachine(),
			"vra_network":               resourceNetwork(),
			"vra_network_profile":       resourceNetworkProfile(),
			"vra_project":               resourceProject(),
			"vra_storage_profile":       resourceStorageProfile(),
			"vra_storage_profile_aws":   resourceStorageProfileAws(),
			"vra_storage_profile_azure": resourceStorageProfileAzure(),
			"vra_zone":                  resourceZone(),
		},

		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	url := d.Get("url").(string)
	refreshToken := ""
	accessToken := ""

	if v, ok := d.GetOk("refresh_token"); ok {
		refreshToken = v.(string)
	}

	if v, ok := d.GetOk("access_token"); ok {
		accessToken = v.(string)
	}

	if accessToken == "" && refreshToken == "" {
		return nil, errors.New("refresh_token or access_token required")
	}

	if accessToken != "" {
		return NewClientFromAccessToken(url, accessToken)
	}

	return NewClientFromRefreshToken(url, refreshToken)
}
