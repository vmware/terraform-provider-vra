package cas

import (
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
)

// Provider represents the CAS provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CAS_URL", nil),
				Description: "The base url for API operations.",
			},
			"refresh_token": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"access_token"},
				DefaultFunc:   schema.EnvDefaultFunc("CAS_REFRESH_TOKEN", nil),
				Description:   "The refresh token for API operations.",
			},
			"access_token": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"refresh_token"},
				DefaultFunc:   schema.EnvDefaultFunc("CAS_ACCESS_TOKEN", nil),
				Description:   "The access token for API operations.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"cas_cloud_account_aws":   dataSourceCloudAccountAWS(),
			"cas_cloud_account_azure": dataSourceCloudAccountAzure(),
			"cas_fabric_network":      dataSourceFabricNetwork(),
			"cas_image":               dataSourceImage(),
			"cas_network":             dataSourceNetwork(),
			"cas_network_domain":      dataSourceNetworkDomain(),
			"cas_project":             dataSourceProject(),
			"cas_region":              dataSourceRegion(),
			"cas_security_group":      dataSourceSecurityGroup(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"cas_block_device":        resourceBlockDevice(),
			"cas_cloud_account_aws":   resourceCloudAccountAWS(),
			"cas_cloud_account_azure": resourceCloudAccountAzure(),
			"cas_flavor_profile":      resourceFlavorProfile(),
			"cas_image_profile":       resourceImageProfile(),
			"cas_load_balancer":       resourceLoadBalancer(),
			"cas_machine":             resourceMachine(),
			"cas_network":             resourceNetwork(),
			"cas_network_profile":     resourceNetworkProfile(),
			"cas_project":             resourceProject(),
			"cas_storage_profile":     resourceStorageProfile(),
			"cas_storage_profile_aws": resourceStorageProfileAws(),
			"cas_zone":                resourceZone(),
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
