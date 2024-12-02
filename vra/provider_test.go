// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProviderVRA *schema.Provider

func init() {
	testAccProviderVRA = Provider()
	testAccProviders = map[string]*schema.Provider{
		"vra": testAccProviderVRA,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(_ *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("VRA_URL"); v == "" {
		t.Fatal("VRA_URL must be set for acceptance tests")
	}

	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}
}

func testAccPreCheckMachine(t *testing.T) {
	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}

	envVars := [...]string{
		"VRA_URL",
		"VRA_AWS_CLOUD_ACCOUNT_NAME",
		"VRA_FLAVOR_1",
		"VRA_FLAVOR_2",
		"VRA_IMAGE_1",
		"VRA_IMAGE_2",
		"VRA_REGION",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckIntegration(t *testing.T) {
	if v := os.Getenv("VRA_URL"); v == "" {
		t.Fatal("VRA_URL must be set for acceptance tests")
	}

	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}

	if v := os.Getenv("VRA_GITHUB_TOKEN"); v == "" {
		t.Fatal("VRA_GITHUB_TOKEN must be set for acceptance tests")
	}
}

func testAccPreCheckContentSharingPolicy(t *testing.T) {
	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}

	envVars := [...]string{
		"VRA_URL",
		"VRA_PROJECT_ID",
		"VRA_CATALOG_ITEM_ID",
		"VRA_CATALOG_SOURCE_ID",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckLoadBalancer(t *testing.T) {
	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}

	if os.Getenv("VRA_IMAGE") == "" && os.Getenv("VRA_IMAGE_REF") == "" {
		t.Fatal("VRA_IMAGE or VRA_IMAGE_REF must be set for acceptance tests")
	}

	envVars := [...]string{
		"VRA_URL",
		"VRA_AWS_CLOUD_ACCOUNT_NAME",
		"VRA_FLAVOR",
		"VRA_REGION",
		"VRA_FABRIC_NETWORK",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckBlockDevice(t *testing.T) {
	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}

	envVars := [...]string{
		"VRA_URL",
		"VRA_AWS_CLOUD_ACCOUNT_NAME",
		"VRA_REGION",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckBlockDeviceSnapshotResource(t *testing.T) {
	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}

	envVars := [...]string{
		"VRA_URL",
		"VRA_VSPHERE_CLOUD_ACCOUNT_NAME",
		"VRA_REGION",
		"VRA_PROJECT",
		"VRA_VSPHERE_DATASTORE_NAME",
		"VRA_VSPHERE_STORAGE_POLICY_NAME",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckImageProfile(t *testing.T) {
	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}

	if os.Getenv("VRA_IMAGE") == "" && os.Getenv("VRA_IMAGE_REF") == "" {
		t.Fatal("VRA_IMAGE or VRA_IMAGE_REF must be set for acceptance tests")
	}

	envVars := [...]string{
		"VRA_URL",
		"VRA_AWS_CLOUD_ACCOUNT_NAME",
		"VRA_REGION",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckStorageProfile(t *testing.T) {
	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}

	envVars := [...]string{
		"VRA_URL",
		"VRA_AWS_CLOUD_ACCOUNT_NAME",
		"VRA_AWS_REGION_NAME",
		"VRA_ARM_CLOUD_ACCOUNT_NAME",
		"VRA_ARM_REGION_NAME",
		"VRA_ARM_FABRIC_STORAGE_ACCOUNT_NAME",
		"VRA_VSPHERE_CLOUD_ACCOUNT_NAME",
		"VRA_VSPHERE_REGION_NAME",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckAWS(t *testing.T) {
	if v := os.Getenv("VRA_URL"); v == "" {
		t.Fatal("VRA_URL must be set for acceptance tests")
	}

	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("VRA_AWS_ACCESS_KEY_ID"); v == "" {
		t.Fatal("VRA_AWS_ACCESS_KEY_ID must be set for acceptance tests")
	}
	if v := os.Getenv("VRA_AWS_SECRET_ACCESS_KEY"); v == "" {
		t.Fatal("VRA_AWS_SECRET_ACCESS_KEY must be set for acceptance tests")
	}
}

func testAccPreCheckAzure(t *testing.T) {
	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}

	envVars := [...]string{
		"VRA_URL",
		"VRA_ARM_SUBSCRIPTION_ID",
		"VRA_ARM_TENANT_ID",
		"VRA_ARM_CLIENT_APP_ID",
		"VRA_ARM_CLIENT_APP_KEY",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckStorageProfileAws(t *testing.T) {
	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}

	envVars := [...]string{
		"VRA_URL",
		"VRA_AWS_CLOUD_ACCOUNT_NAME",
		"VRA_AWS_REGION_NAME",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckStorageProfileAzure(t *testing.T) {
	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}

	envVars := [...]string{
		"VRA_URL",
		"VRA_ARM_SUBSCRIPTION_ID",
		"VRA_ARM_TENANT_ID",
		"VRA_ARM_CLIENT_APP_ID",
		"VRA_ARM_CLIENT_APP_KEY",
		"VRA_ARM_REGION",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckStorageProfileVsphere(t *testing.T) {
	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}

	envVars := [...]string{
		"VRA_URL",
		"VRA_VSPHERE_REGION",
		"VRA_VSPHERE_CLOUD_ACCOUNT_NAME",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckVsphere(t *testing.T) {
	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}

	envVars := [...]string{
		"VRA_URL",
		"VRA_VSPHERE_USERNAME",
		"VRA_VSPHERE_PASSWORD",
		"VRA_VSPHERE_DATACOLLECTOR_NAME",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckVsphereForDataStore(t *testing.T) {
	testAccPreCheckVra(t)

	// The vCenter should have already been added into the vRA
	// Set the VRA_VSPHERE_DATASTORE_NAME to an existing unique datastore name
	envVars := [...]string{
		"VRA_VSPHERE_DATASTORE_NAME",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckVsphereForStoragePolicy(t *testing.T) {
	testAccPreCheckVra(t)

	// The vCenter should have already been added into the vRA
	// Set the VRA_VSPHERE_STORAGE_POLICY_NAME to an existing unique data policy name
	envVars := [...]string{
		//"VRA_URL",
		"VRA_VSPHERE_STORAGE_POLICY_NAME",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckBlockDeviceSnapshot(t *testing.T) {
	testAccPreCheckVra(t)

	// The vCenter should have already been added into the vRA
	// Set the VRA_BLOCK_DEVICE_ID to an existing block device id which has snapshots created
	envVars := [...]string{
		"VRA_BLOCK_DEVICE_ID",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckGCP(t *testing.T) {
	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}

	envVars := [...]string{
		"VRA_URL",
		"VRA_GCP_CLIENT_EMAIL",
		"VRA_GCP_PRIVATE_KEY_ID",
		"VRA_GCP_PRIVATE_KEY",
		"VRA_GCP_PROJECT_ID",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckVMC(t *testing.T) {
	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}

	envVars := [...]string{
		"VRA_URL",
		"VRA_VMC_API_TOKEN",
		"VRA_VMC_SDDC_NAME",
		"VRA_VMC_VCENTER_HOSTNAME",
		"VRA_VMC_VCENTER_USERNAME",
		"VRA_VMC_VCENTER_PASSWORD",
		"VRA_VMC_NSX_HOSTNAME",
		"VRA_VMC_DATA_COLLECTOR_NAME",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckNSXV(t *testing.T) {
	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}

	envVars := [...]string{
		"VRA_URL",
		"VRA_NSXV_HOSTNAME",
		"VRA_NSXV_PASSWORD",
		"VRA_NSXV_USERNAME",
		"VRA_NSXV_DATA_COLLECTOR_NAME",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckNSXT(t *testing.T) {
	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}

	envVars := [...]string{
		"VRA_URL",
		"VRA_NSXT_HOSTNAME",
		"VRA_NSXT_PASSWORD",
		"VRA_NSXT_USERNAME",
		"VRA_NSXT_DATA_COLLECTOR_NAME",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckCatalogItem(t *testing.T) {
	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}

	envVars := [...]string{
		"VRA_URL",
		"VRA_CATALOG_ITEM_NAME",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckDeployment(t *testing.T) {
	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}

	envVars := [...]string{
		"VRA_URL",
		"VRA_CATALOG_ITEM_NAME",
		"VRA_PROJECT_NAME",
		"VRA_BLUEPRINT_ID",
		"VRA_BLUEPRINT_VERSION",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckDeploymentDataSource(t *testing.T) {
	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}

	envVars := [...]string{
		"VRA_URL",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckBlueprint(t *testing.T) {
	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}

	envVars := [...]string{
		"VRA_URL",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckVra(t *testing.T) {
	if v := os.Getenv("VRA_URL"); v == "" {
		t.Fatal("VRA_URL must be set for acceptance tests")
	}

	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}
}

func testAccPreCheckContentSource(t *testing.T) {
	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}

	envVars := [...]string{
		"VRA_URL",
		"VRA_INTEGRATION_ID",
		"VRA_CONTENT_SOURCE_PATH",
		"VRA_CONTENT_SOURCE_BRANCH",
		"VRA_CONTENT_SOURCE_REPO",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckFabricStorageAccountAzure(t *testing.T) {
	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}

	envVars := [...]string{
		"VRA_URL",
		"VRA_ARM_CLOUD_ACCOUNT_NAME",
		"VRA_ARM_REGION_NAME",
		"VRA_ARM_FABRIC_STORAGE_ACCOUNT_NAME",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckFabricCompute(t *testing.T) {
	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}

	envVars := [...]string{
		"VRA_URL",
		"VRA_FABRIC_COMPUTE_NAME",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}

func testAccPreCheckFabricDatastoreVsphere(t *testing.T) {
	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}

	envVars := [...]string{
		"VRA_URL",
		"VRA_FABRIC_DATASTORE_VSPHERE_NAME",
	}

	for _, name := range envVars {
		if v := os.Getenv(name); v == "" {
			t.Fatalf("%s must be set for acceptance tests\n", name)
		}
	}
}
