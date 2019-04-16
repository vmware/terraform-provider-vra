package cas

import (
	"errors"

	tango "github.com/vmware/terraform-provider-cas/sdk"

	"github.com/hashicorp/terraform/helper/schema"
)

// Provider represents the tango provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TANGO_URL", nil),
				Description: "The base url for API operations.",
			},
			"refresh_token": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"access_token"},
				DefaultFunc:   schema.EnvDefaultFunc("TANGO_REFRESH_TOKEN", nil),
				Description:   "The refresh token for API operations.",
			},
			"access_token": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"refresh_token"},
				DefaultFunc:   schema.EnvDefaultFunc("TANGO_ACCESS_TOKEN", nil),
				Description:   "The access token for API operations.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TANGO_PROJECT_ID", nil),
				Description: "The project id to use for this template.",
			},
			"deployment_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TANGO_DEPLOYMENT_ID", nil),
				Description: "The deployment id to use for this template.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"cas_cloud_account_aws":   dataSourceCloudAccountAWS(),
			"cas_cloud_account_azure": dataSourceCloudAccountAzure(),
			"cas_image":               dataSourceImage(),
			"cas_project":             dataSourceProject(),
			"cas_region":              dataSourceRegion(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"cas_block_device":        resourceBlockDevice(),
			"cas_cloud_account_aws":   resourceCloudAccountAWS(),
			"cas_cloud_account_azure": resourceCloudAccountAzure(),
			"cas_flavor":              resourceFlavor(),
			"cas_image_profile":       resourceImageProfile(),
			"cas_load_balancer":       resourceLoadBalancer(),
			"cas_machine":             resourceMachine(),
			"cas_network":             resourceNetwork(),
			"cas_project":             resourceProject(),
			"cas_zone":                resourceZone(),
		},

		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	url := d.Get("url").(string)
	projectID := d.Get("project_id").(string)
	deploymentID := d.Get("deployment_id").(string)
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
		return tango.NewClientFromAccessToken(url, accessToken, projectID, deploymentID)
	}

	return tango.NewClientFromRefreshToken(url, refreshToken, projectID, deploymentID)
}
