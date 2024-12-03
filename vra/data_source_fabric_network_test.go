// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
)

func TestAccDataSourceVRAFabricNetworkBasic(t *testing.T) {
	rInt := acctest.RandInt()
	dataSourceName := "data.vra_fabric_network.my-fabric-network"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAFabricNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRAFabricNetworkNoneConfig(rInt),
				ExpectError: regexp.MustCompile("vra_fabric_network filter did not match any fabric network"),
			},
			{
				Config: testAccDataSourceVRAFabricNetworkOneConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "external_region_id", "us-east-1"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "appnet-isolated-dev"),
				),
			},
		},
	})
}

func testAccDataSourceVRAFabricNetworkBaseConfig(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	id := os.Getenv("VRA_AWS_ACCESS_KEY_ID")
	secret := os.Getenv("VRA_AWS_SECRET_ACCESS_KEY")
	return fmt.Sprintf(`
	resource "vra_cloud_account_aws" "my-cloud-account" {
		name = "my-cloud-account-%d"
		description = "test cloud account"
		access_key = "%s"
		secret_key = "%s"
		regions = ["us-east-1"]
	}
	`, rInt, id, secret)
}

func testAccDataSourceVRAFabricNetworkNoneConfig(rInt int) string {
	return testAccDataSourceVRAFabricNetworkBaseConfig(rInt) + `
		data "vra_fabric_network" "my-fabric-network" {
			filter = "name eq 'foobar appnet-isolated-dev'"
		}`
}

func testAccDataSourceVRAFabricNetworkOneConfig(rInt int) string {
	return testAccDataSourceVRAFabricNetworkBaseConfig(rInt) + `
	data "vra_fabric_network" "my-fabric-network" {
			filter = "name eq 'appnet-isolated-dev'"
		}`
}

func testAccCheckVRAFabricNetworkDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vra_cloud_account_aws" {
			continue
		}

		_, err := apiClient.CloudAccount.GetAwsCloudAccount(cloud_account.NewGetAwsCloudAccountParams().WithID(rs.Primary.ID))
		if err == nil {
			return err
		}
	}

	return nil
}
