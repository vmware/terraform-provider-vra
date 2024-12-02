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
)

func TestAccDataSourceVRACloudAccountNSXT(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName1 := "vra_cloud_account_nsxt.this"
	dataSourceName1 := "data.vra_cloud_account_nsxt.this"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckNSXT(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRACloudAccountNSXTNotFound(rInt),
				ExpectError: regexp.MustCompile("cloud account foobar not found"),
			},
			{
				Config: testAccDataSourceVRACloudAccountNSXT(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "dc_id", dataSourceName1, "dc_id"),
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "hostname", dataSourceName1, "hostname"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
					resource.TestCheckResourceAttrPair(resourceName1, "tags", dataSourceName1, "tags"),
					resource.TestCheckResourceAttrPair(resourceName1, "username", dataSourceName1, "username"),
				),
			},
		},
	})
}

func testAccDataSourceVRACloudAccountNSXTBase(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	hostname := os.Getenv("VRA_NSXT_HOSTNAME")
	password := os.Getenv("VRA_NSXT_PASSWORD")
	username := os.Getenv("VRA_NSXT_USERNAME")
	dcName := os.Getenv("VRA_NSXT_DATA_COLLECTOR_NAME")
	return fmt.Sprintf(`
	data "vra_data_collector" "dc" {
      name = "%s"
	}

	resource "vra_cloud_account_nsxt" "this" {
	  name        = "my-nsxt-account-%d"
	  description = "my nsx-t cloud account"
	  dc_id        = data.vra_data_collector.dc.id
	  hostname    = "%s"
	  password    = "%s"
	  username    = "%s"

	  accept_self_signed_cert = true

	  tags {
		key   = "foo"
		value = "bar"
	  }

	  tags {
		key = "where"
		value = "waldo"
	  }
	}`, dcName, rInt, hostname, password, username)
}

func testAccDataSourceVRACloudAccountNSXTNotFound(rInt int) string {
	return testAccDataSourceVRACloudAccountNSXTBase(rInt) + `
	data "vra_cloud_account_nsxt" "this" {
		name = "foobar"
	}`
}

func testAccDataSourceVRACloudAccountNSXT(rInt int) string {
	return testAccDataSourceVRACloudAccountNSXTBase(rInt) + `
	data "vra_cloud_account_nsxt" "this" {
		name = "${vra_cloud_account_nsxt.this.name}"
	}`
}
