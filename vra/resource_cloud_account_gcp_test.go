// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
)

func TestAccVRACloudAccountGCP_Basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckGCP(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRACloudAccountGCPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRACloudAccountGCPConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRACloudAccountGCPExists("vra_cloud_account_gcp.my-cloud-account"),
					resource.TestMatchResourceAttr(
						"vra_cloud_account_gcp.my-cloud-account", "name", regexp.MustCompile("^my-cloud-account-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_gcp.my-cloud-account", "description", "test cloud account"),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_gcp.my-cloud-account", "client_email", os.Getenv("VRA_GCP_CLIENT_EMAIL")),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_gcp.my-cloud-account", "project_id", os.Getenv("VRA_GCP_PROJECT_ID")),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_gcp.my-cloud-account", "private_key_id", os.Getenv("VRA_GCP_PRIVATE_KEY_ID")),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_gcp.my-cloud-account", "tags.#", "2"),
				),
			},
			{
				Config: testAccCheckVRACloudAccountGCPUpdateConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRACloudAccountGCPExists("vra_cloud_account_gcp.my-cloud-account"),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_gcp.my-cloud-account", "description", "your test cloud account"),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_gcp.my-cloud-account", "client_email", os.Getenv("VRA_GCP_CLIENT_EMAIL")),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_gcp.my-cloud-account", "project_id", os.Getenv("VRA_GCP_PROJECT_ID")),
					resource.TestCheckResourceAttr(
						"vra_cloud_account_gcp.my-cloud-account", "private_key_id", os.Getenv("VRA_GCP_PRIVATE_KEY_ID")),
				),
			},
		},
	})
}
func TestAccVRACloudAccountGCP_Duplicate(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckGCP(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRACloudAccountGCPDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckVRACloudAccountGCPConfigDuplicateRegion(rInt),
				ExpectError: regexp.MustCompile("specified regions are not unique"),
			},
		},
	})
}

func testAccCheckVRACloudAccountGCPExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no cloud account ID is set")
		}

		return nil
	}
}

func testAccCheckVRACloudAccountGCPDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vra_cloud_account_gcp" {
			continue
		}

		_, err := apiClient.CloudAccount.GetGcpCloudAccount(cloud_account.NewGetGcpCloudAccountParams().WithID(rs.Primary.ID))
		if err == nil {
			return err
		}
	}

	return nil
}

func testAccCheckVRACloudAccountGCPConfig(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	clientEmail := os.Getenv("VRA_GCP_CLIENT_EMAIL")
	projectID := os.Getenv("VRA_GCP_PROJECT_ID")
	privateKeyID := os.Getenv("VRA_GCP_PRIVATE_KEY_ID")
	privateKey := os.Getenv("VRA_GCP_PRIVATE_KEY")
	return fmt.Sprintf(`
resource "vra_cloud_account_gcp" "my-cloud-account" {
	name = "my-cloud-account-%d"
	description = "test cloud account"
	client_email = "%s"
	project_id = "%s"
	private_key_id = "%s"
	private_key = "%s"
	regions = ["us-west2"]
	tags {
		key = "foo"
		value = "bar"
	}
	tags {
		key = "where"
		value = "waldo"
	}
 }`, rInt, clientEmail, projectID, privateKeyID, privateKey)
}

func testAccCheckVRACloudAccountGCPUpdateConfig(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	clientEmail := os.Getenv("VRA_GCP_CLIENT_EMAIL")
	projectID := os.Getenv("VRA_GCP_PROJECT_ID")
	privateKeyID := os.Getenv("VRA_GCP_PRIVATE_KEY_ID")
	privateKey := os.Getenv("VRA_GCP_PRIVATE_KEY")
	return fmt.Sprintf(`
resource "vra_cloud_account_gcp" "my-cloud-account" {
	name = "my-cloud-account-%d"
	description = "your test cloud account"
	client_email = "%s"
	project_id = "%s"
	private_key_id = "%s"
	private_key = "%s"
	regions = ["us-west2"]
 }`, rInt, clientEmail, projectID, privateKeyID, privateKey)
}

func testAccCheckVRACloudAccountGCPConfigDuplicateRegion(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	clientEmail := os.Getenv("VRA_GCP_CLIENT_EMAIL")
	projectID := os.Getenv("VRA_GCP_PROJECT_ID")
	privateKeyID := os.Getenv("VRA_GCP_PRIVATE_KEY_ID")
	privateKey := os.Getenv("VRA_GCP_PRIVATE_KEY")
	return fmt.Sprintf(`
resource "vra_cloud_account_gcp" "my-cloud-account" {
	name = "my-cloud-account-%d"
	description = "test cloud account"
	client_email = "%s"
	project_id = "%s"
	private_key_id = "%s"
	private_key = "%s"
	regions = ["us-west2", "us-west2"]
 }`, rInt, clientEmail, projectID, privateKeyID, privateKey)
}
