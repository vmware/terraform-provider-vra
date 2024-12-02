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

func TestAccVRACloudAccountNSXV_Basic(t *testing.T) {
	rInt := acctest.RandInt()
	const nsxvAccount = "vra_cloud_account_nsxv.this"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckNSXV(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRACloudAccountNSXVDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRACloudAccountNSXVConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRACloudAccountNSXVExists(nsxvAccount),
					resource.TestMatchResourceAttr(
						nsxvAccount, "name", regexp.MustCompile("^my-nsxv-account-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						nsxvAccount, "description", "my nsx-v cloud account"),
					resource.TestCheckResourceAttr(
						nsxvAccount, "hostname", os.Getenv("VRA_NSXV_HOSTNAME")),
					resource.TestCheckResourceAttr(
						nsxvAccount, "username", os.Getenv("VRA_NSXV_USERNAME")),
					resource.TestCheckResourceAttrPair(
						nsxvAccount, "dc_id", "data.vra_data_collector.dc", "id"),
					resource.TestCheckResourceAttr(
						nsxvAccount, "tags.#", "2"),
				),
			},
			{
				Config: testAccCheckVRACloudAccountNSXVUpdateConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRACloudAccountNSXVExists(nsxvAccount),
					resource.TestMatchResourceAttr(
						nsxvAccount, "name", regexp.MustCompile("^my-nsxv-account-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						nsxvAccount, "description", "your nsx-v cloud account"),
					resource.TestCheckResourceAttr(
						nsxvAccount, "hostname", os.Getenv("VRA_NSXV_HOSTNAME")),
					resource.TestCheckResourceAttr(
						nsxvAccount, "username", os.Getenv("VRA_NSXV_USERNAME")),
					resource.TestCheckResourceAttrPair(
						nsxvAccount, "dc_id", "data.vra_data_collector.dc", "id"),
					resource.TestCheckResourceAttr(
						nsxvAccount, "tags.#", "2"),
				),
			},
		},
	})
}

func testAccCheckVRACloudAccountNSXVExists(n string) resource.TestCheckFunc {
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

func testAccCheckVRACloudAccountNSXVDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vra_cloud_account_nsxv" {
			continue
		}

		_, err := apiClient.CloudAccount.GetNsxTCloudAccount(cloud_account.NewGetNsxTCloudAccountParams().WithID(rs.Primary.ID))
		if err == nil {
			return err
		}
	}

	return nil
}

func testAccCheckVRACloudAccountNSXVConfig(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	hostname := os.Getenv("VRA_NSXV_HOSTNAME")
	password := os.Getenv("VRA_NSXV_PASSWORD")
	username := os.Getenv("VRA_NSXV_USERNAME")
	dcName := os.Getenv("VRA_NSXV_DATA_COLLECTOR_NAME")
	return fmt.Sprintf(`
	data "vra_data_collector" "dc" {
      name = "%s"
	}

	resource "vra_cloud_account_nsxv" "this" {
	  name        = "my-nsxv-account-%d"
	  description = "my nsx-v cloud account"
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

func testAccCheckVRACloudAccountNSXVUpdateConfig(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	hostname := os.Getenv("VRA_NSXV_HOSTNAME")
	password := os.Getenv("VRA_NSXV_PASSWORD")
	username := os.Getenv("VRA_NSXV_USERNAME")
	dcName := os.Getenv("VRA_NSXV_DATA_COLLECTOR_NAME")
	return fmt.Sprintf(`
	data "vra_data_collector" "dc" {
      name = "%s"
	}

	resource "vra_cloud_account_nsxv" "this" {
	  name        = "my-nsxv-account-%d"
	  description = "your nsx-v cloud account"
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
