package vra

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProviderVRA *schema.Provider

func init() {
	testAccProviderVRA = Provider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"vra": testAccProviderVRA,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
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

	if os.Getenv("VRA_IMAGE") == "" && os.Getenv("VRA_IMAGE_REF") == "" {
		t.Fatal("VRA_IMAGE or VRA_IMAGE_REF must be set for acceptance tests")
	}

	envVars := [...]string{
		"VRA_URL",
		"VRA_AWS_CLOUD_ACCOUNT_NAME",
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

func testAccPreCheckVra(t *testing.T) {
	if v := os.Getenv("VRA_URL"); v == "" {
		t.Fatal("VRA_URL must be set for acceptance tests")
	}

	if os.Getenv("VRA_REFRESH_TOKEN") == "" && os.Getenv("VRA_ACCESS_TOKEN") == "" {
		t.Fatal("VRA_REFRESH_TOKEN or VRA_ACCESS_TOKEN must be set for acceptance tests")
	}
}
