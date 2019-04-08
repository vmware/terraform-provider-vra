package cas

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var testAccProviderCAS *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviderCAS = Provider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"tango": testAccProvider,
		"cas":   testAccProviderCAS,
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
	if v := os.Getenv("TANGO_URL"); v == "" {
		t.Fatal("TANGO_URL must be set for acceptance tests")
	}

	if os.Getenv("TANGO_REFRESH_TOKEN") == "" && os.Getenv("TANGO_ACCESS_TOKEN") == "" {
		t.Fatal("TANGO_REFRESH_TOKEN or TANGO_ACCESS_TOKEN must be set for acceptance tests")
	}

	if v := os.Getenv("TANGO_PROJECT_ID"); v == "" {
		t.Fatal("TANGO_PROJECT_ID must be set for acceptance tests")
	}

	if v := os.Getenv("TANGO_DEPLOYMENT_ID"); v == "" {
		t.Fatal("TANGO_DEPLOYMENT_ID must be set for acceptance tests")
	}
}

func testAccPreCheckAWS(t *testing.T) {
	if v := os.Getenv("TANGO_URL"); v == "" {
		t.Fatal("TANGO_URL must be set for acceptance tests")
	}

	if os.Getenv("TANGO_REFRESH_TOKEN") == "" && os.Getenv("TANGO_ACCESS_TOKEN") == "" {
		t.Fatal("TANGO_REFRESH_TOKEN or TANGO_ACCESS_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("CAS_AWS_ACCESS_KEY_ID"); v == "" {
		t.Fatal("CAS_AWS_ACCESS_KEY_ID must be set for acceptance tests")
	}
	if v := os.Getenv("CAS_AWS_SECRET_ACCESS_KEY"); v == "" {
		t.Fatal("TANCAS_AWS_SECRET_ACCESS_KEYGO_URL must be set for acceptance tests")
	}
}
