// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"errors"
	"os"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VRA_URL", nil),
				Description: "The base url for API operations.",
			},
			"organization": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VCFA_ORGANIZATION", nil),
				Description: "Organization name (required for VCF Automation).",
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
			"insecure": {
				Type:        schema.TypeBool,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"VRA_INSECURE", "VRA7_INSECURE"}, nil),
				Optional:    true,
				Description: "Specify whether to validate TLS certificates.",
				ValidateDiagFunc: schema.SchemaValidateDiagFunc(func(_ interface{}, _ cty.Path) diag.Diagnostics {
					var diags diag.Diagnostics
					if envVar, ok := os.LookupEnv("VRA7_INSECURE"); ok && envVar != "" {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Warning,
							Summary:  "Deprecated environment variable.",
							Detail:   "'VRA7_INSECURE' is deprecated; use 'VRA_INSECURE'.",
						})
					}
					return diags
				}),
			},
			"reauthorize_timeout": {
				Type:        schema.TypeString,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"VRA_REAUTHORIZE_TIMEOUT", "VRA7_REAUTHORIZE_TIMEOUT"}, nil),
				Optional:    true,
				Description: "Specify timeout for how often to reauthorize the access token.",
				ValidateDiagFunc: schema.SchemaValidateDiagFunc(func(_ interface{}, _ cty.Path) diag.Diagnostics {
					var diags diag.Diagnostics
					if envVar, ok := os.LookupEnv("VRA7_REAUTHORIZE_TIMEOUT"); ok && envVar != "" {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Warning,
							Summary:  "Deprecated environment variable.",
							Detail:   "'VRA7_REAUTHORIZE_TIMEOUT' is deprecated; use 'VRA_REAUTHORIZE_TIMEOUT'.",
						})
					}
					return diags
				}),
			},
			"api_timeout": {
				Type:        schema.TypeInt,
				DefaultFunc: schema.EnvDefaultFunc("VRA_API_TIMEOUT", 30),
				Optional:    true,
				Description: "Specify timeout in seconds for API operations.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"vra_block_device":                  dataSourceBlockDevice(),
			"vra_block_device_snapshots":        dataSourceBlockDeviceSnapshots(),
			"vra_blueprint":                     dataSourceBlueprint(),
			"vra_blueprint_version":             dataSourceBlueprintVersion(),
			"vra_catalog_item":                  dataSourceCatalogItem(),
			"vra_catalog_item_entitlement":      dataSourceCatalogItemEntitlement(),
			"vra_catalog_source_blueprint":      dataSourceCatalogSourceBlueprint(),
			"vra_catalog_source_entitlement":    dataSourceCatalogSourceEntitlement(),
			"vra_cloud_account_aws":             dataSourceCloudAccountAWS(),
			"vra_cloud_account_azure":           dataSourceCloudAccountAzure(),
			"vra_cloud_account_gcp":             dataSourceCloudAccountGCP(),
			"vra_cloud_account_nsxt":            dataSourceCloudAccountNSXT(),
			"vra_cloud_account_nsxv":            dataSourceCloudAccountNSXV(),
			"vra_cloud_account_vmc":             dataSourceCloudAccountVMC(),
			"vra_cloud_account_vsphere":         dataSourceCloudAccountVsphere(),
			"vra_content_sharing_policy":        dataSourceContentSharingPolicy(),
			"vra_data_collector":                dataSourceDataCollector(),
			"vra_deployment":                    dataSourceDeployment(),
			"vra_fabric_compute":                dataSourceFabricCompute(),
			"vra_fabric_datastore_vsphere":      dataSourceFabricDatastoreVsphere(),
			"vra_fabric_network":                dataSourceFabricNetwork(),
			"vra_fabric_storage_account_azure":  dataSourceFabricStorageAccountAzure(),
			"vra_fabric_storage_policy_vsphere": dataSourceFabricStoragePolicyVsphere(),
			"vra_image":                         dataSourceImage(),
			"vra_image_profile":                 dataSourceImageProfile(),
			"vra_machine":                       dataSourceMachine(),
			"vra_network":                       dataSourceNetwork(),
			"vra_network_domain":                dataSourceNetworkDomain(),
			"vra_network_profile":               dataSourceNetworkProfile(),
			"vra_project":                       dataSourceProject(),
			"vra_region":                        dataSourceRegion(),
			"vra_region_enumeration":            dataSourceRegionEnumeration(),
			"vra_region_enumeration_aws":        dataSourceRegionEnumerationAWS(),
			"vra_region_enumeration_azure":      dataSourceRegionEnumerationAzure(),
			"vra_region_enumeration_gcp":        dataSourceRegionEnumerationGCP(),
			"vra_region_enumeration_vmc":        dataSourceRegionEnumerationVMC(),
			"vra_region_enumeration_vsphere":    dataSourceRegionEnumerationVsphere(),
			"vra_security_group":                dataSourceSecurityGroup(),
			"vra_storage_profile":               dataSourceStorageProfile(),
			"vra_storage_profile_aws":           datasourceStorageProfileAws(),
			"vra_storage_profile_azure":         datasourceStorageProfileAzure(),
			"vra_storage_profile_vsphere":       dataSourceStorageProfileVsphere(),
			"vra_zone":                          dataSourceZone(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"vra_block_device":               resourceBlockDevice(),
			"vra_block_device_snapshot":      resourceBlockDeviceSnapshot(),
			"vra_blueprint":                  resourceBlueprint(),
			"vra_blueprint_version":          resourceBlueprintVersion(),
			"vra_catalog_item_entitlement":   resourceCatalogItemEntitlement(),
			"vra_catalog_source_blueprint":   resourceCatalogSourceBlueprint(),
			"vra_catalog_source_entitlement": resourceCatalogSourceEntitlement(),
			"vra_cloud_account_aws":          resourceCloudAccountAWS(),
			"vra_cloud_account_azure":        resourceCloudAccountAzure(),
			"vra_cloud_account_gcp":          resourceCloudAccountGCP(),
			"vra_cloud_account_nsxt":         resourceCloudAccountNSXT(),
			"vra_cloud_account_nsxv":         resourceCloudAccountNSXV(),
			"vra_cloud_account_vmc":          resourceCloudAccountVMC(),
			"vra_cloud_account_vsphere":      resourceCloudAccountVsphere(),
			"vra_content_sharing_policy":     resourceContentSharingPolicy(),
			"vra_content_source":             resourceContentSource(),
			"vra_deployment":                 resourceDeployment(),
			"vra_fabric_compute":             resourceFabricCompute(),
			"vra_fabric_datastore_vsphere":   resourceFabricDatastoreVsphere(),
			"vra_fabric_network_vsphere":     resourceFabricNetworkVsphere(),
			"vra_flavor_profile":             resourceFlavorProfile(),
			"vra_image_profile":              resourceImageProfile(),
			"vra_integration":                resourceIntegration(),
			"vra_load_balancer":              resourceLoadBalancer(),
			"vra_machine":                    resourceMachine(),
			"vra_network":                    resourceNetwork(),
			"vra_network_profile":            resourceNetworkProfile(),
			"vra_network_ip_range":           resourceNetworkIPRange(),
			"vra_project":                    resourceProject(),
			"vra_storage_profile":            resourceStorageProfile(),
			"vra_storage_profile_aws":        resourceStorageProfileAws(),
			"vra_storage_profile_azure":      resourceStorageProfileAzure(),
			"vra_storage_profile_vsphere":    resourceStorageProfileVsphere(),
			"vra_zone":                       resourceZone(),
		},

		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	url := d.Get("url").(string)
	organization := ""
	refreshToken := ""
	accessToken := ""
	reauth := "0"
	apiTimeout := 0

	if v, ok := d.GetOk("organization"); ok {
		organization = v.(string)
	}

	if v, ok := d.GetOk("refresh_token"); ok {
		refreshToken = v.(string)
	}

	if v, ok := d.GetOk("access_token"); ok {
		accessToken = v.(string)
	}

	insecure := d.Get("insecure").(bool)

	if v, ok := d.GetOk("reauthorize_timeout"); ok {
		reauth = v.(string)
	}

	if v, ok := d.GetOk("api_timeout"); ok {
		apiTimeout = v.(int)
	}

	if accessToken == "" && refreshToken == "" {
		return nil, errors.New("refresh_token or access_token required")
	}

	if accessToken != "" {
		return NewClientFromAccessToken(url, accessToken, insecure, apiTimeout)
	}

	return NewClientFromRefreshToken(url, organization, refreshToken, insecure, reauth, apiTimeout)
}
