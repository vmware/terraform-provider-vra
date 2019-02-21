package tango

import (
	"errors"
	"tango-terraform-provider/tango/client"

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
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TANGO_PROJECT_ID", nil),
				Description: "The project id to use for this template.",
			},
			"deployment_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TANGO_DEPLOYMENT_ID", nil),
				Description: "The deployment id to use for this template.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"tango_machine":       resourceMachine(),
			"tango_network":       resourceNetwork(),
			"tango_block_device":  resourceBlockDevice(),
			"tango_load_balancer": resourceLoadBalancer(),
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
