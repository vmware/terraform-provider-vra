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

func TestAccVRACloudAccountNSXT_Basic(t *testing.T) {
	rInt := acctest.RandInt()
	const nsxtAccount = "vra_cloud_account_nsxt.this"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckNSXT(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRACloudAccountNSXTDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRACloudAccountNSXTConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRACloudAccountNSXTExists(nsxtAccount),
					resource.TestMatchResourceAttr(
						nsxtAccount, "name", regexp.MustCompile("^my-nsxt-account-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						nsxtAccount, "description", "my nsx-t cloud account"),
					resource.TestCheckResourceAttr(
						nsxtAccount, "hostname", os.Getenv("VRA_NSXT_HOSTNAME")),
					resource.TestCheckResourceAttr(
						nsxtAccount, "username", os.Getenv("VRA_NSXT_USERNAME")),
					resource.TestCheckResourceAttrPair(
						nsxtAccount, "dc_id", "data.vra_data_collector.dc", "id"),
					resource.TestCheckResourceAttr(
						nsxtAccount, "tags.#", "2"),
				),
			},
			{
				Config: testAccCheckVRACloudAccountNSXTUpdateConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRACloudAccountNSXTExists(nsxtAccount),
					resource.TestMatchResourceAttr(
						nsxtAccount, "name", regexp.MustCompile("^my-nsxt-account-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						nsxtAccount, "description", "your nsx-t cloud account"),
					resource.TestCheckResourceAttr(
						nsxtAccount, "hostname", os.Getenv("VRA_NSXT_HOSTNAME")),
					resource.TestCheckResourceAttr(
						nsxtAccount, "username", os.Getenv("VRA_NSXT_USERNAME")),
					resource.TestCheckResourceAttrPair(
						nsxtAccount, "dc_id", "data.vra_data_collector.dc", "id"),
					resource.TestCheckResourceAttr(
						nsxtAccount, "tags.#", "2"),
				),
			},
		},
	})
}

func testAccCheckVRACloudAccountNSXTExists(n string) resource.TestCheckFunc {
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

func testAccCheckVRACloudAccountNSXTDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vra_cloud_account_nsxt" {
			continue
		}

		_, err := apiClient.CloudAccount.GetNsxTCloudAccount(cloud_account.NewGetNsxTCloudAccountParams().WithID(rs.Primary.ID))
		if err == nil {
			return err
		}
	}

	return nil
}

func testAccCheckVRACloudAccountNSXTConfig(rInt int) string {
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

func testAccCheckVRACloudAccountNSXTUpdateConfig(rInt int) string {
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
	  description = "your nsx-t cloud account"
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
