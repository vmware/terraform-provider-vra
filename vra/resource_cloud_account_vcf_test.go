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

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
)

func TestAccVRACloudAccountVCF_Basic(t *testing.T) {
	rInt := acctest.RandInt()
	resource1 := "vra_cloud_account_vcf.this"
	dataSource1 := "data.vra_data_collector.dc"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckVCF(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRACloudAccountVCFDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRACloudAccountVCFConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRACloudAccountVCFExists(resource1),
					resource.TestMatchResourceAttr(
						resource1, "name", regexp.MustCompile("^test-vcf-account-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						resource1, "description", "my vcf cloud account"),
					resource.TestCheckResourceAttr(
						resource1, "nsx_hostname", os.Getenv("VRA_VCF_NSX_HOSTNAME")),
					resource.TestCheckResourceAttr(
						resource1, "nsx_username", os.Getenv("VRA_VCF_NSX_USERNAME")),
					resource.TestCheckResourceAttr(
						resource1, "vcenter_hostname", os.Getenv("VRA_VCF_VCENTER_HOSTNAME")),
					resource.TestCheckResourceAttr(
						resource1, "vcenter_username", os.Getenv("VRA_VCF_VCENTER_USERNAME")),
					resource.TestCheckResourceAttr(
						resource1, "workload_domain_id", os.Getenv("VRA_VCF_WORKLOAD_DOMAIN_ID")),
					resource.TestCheckResourceAttr(
						resource1, "workload_domain_name", os.Getenv("VRA_VCF_WORKLOAD_DOMAIN_NAME")),
					resource.TestCheckResourceAttrPair(
						resource1, "dc_id", dataSource1, "id"),
					resource.TestCheckResourceAttr(
						resource1, "regions.#", "1"),
					resource.TestCheckResourceAttr(
						resource1, "tags.#", "2"),
				),
			},
			{
				Config: testAccCheckVRACloudAccountVCFUpdateConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRACloudAccountVCFExists(resource1),
					resource.TestMatchResourceAttr(
						resource1, "name", regexp.MustCompile("^test-vcf-account-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						resource1, "description", "your vcf cloud account"),
					resource.TestCheckResourceAttr(
						resource1, "nsx_hostname", os.Getenv("VRA_VCF_NSX_HOSTNAME")),
					resource.TestCheckResourceAttr(
						resource1, "nsx_username", os.Getenv("VRA_VCF_NSX_USERNAME")),
					resource.TestCheckResourceAttr(
						resource1, "vcenter_hostname", os.Getenv("VRA_VCF_VCENTER_HOSTNAME")),
					resource.TestCheckResourceAttr(
						resource1, "vcenter_username", os.Getenv("VRA_VCF_VCENTER_USERNAME")),
					resource.TestCheckResourceAttr(
						resource1, "workload_domain_id", os.Getenv("VRA_VCF_WORKLOAD_DOMAIN_ID")),
					resource.TestCheckResourceAttr(
						resource1, "workload_domain_name", os.Getenv("VRA_VCF_WORKLOAD_DOMAIN_NAME")),
					resource.TestCheckResourceAttrPair(
						resource1, "dc_id", dataSource1, "id"),
					resource.TestCheckResourceAttr(
						resource1, "regions.#", "1"),
					resource.TestCheckResourceAttr(
						resource1, "tags.#", "2"),
				),
			},
		},
	})
}

func TestAccVRACloudAccountVCF_Duplicate(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckVCF(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRACloudAccountVCFDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckVRACloudAccountVCFConfigDuplicateRegion(rInt),
				ExpectError: regexp.MustCompile("specified regions are not unique"),
			},
		},
	})
}

func testAccCheckVRACloudAccountVCFExists(n string) resource.TestCheckFunc {
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

func testAccCheckVRACloudAccountVCFDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vra_cloud_account_vcf" {
			continue
		}

		_, err := apiClient.CloudAccount.GetVcfCloudAccount(cloud_account.NewGetVcfCloudAccountParams().WithID(rs.Primary.ID))
		if err == nil {
			return err
		}
	}

	return nil
}

func testAccCheckVRACloudAccountVCFConfig(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	nsxHostName := os.Getenv("VRA_VCF_NSX_HOSTNAME")
	nsxUserName := os.Getenv("VRA_VCF_NSX_USERNAME")
	nsxPassword := os.Getenv("VRA_VCF_NSX_PASSWORD")
	vCenterHostName := os.Getenv("VRA_VCF_VCENTER_HOSTNAME")
	vCenterUserName := os.Getenv("VRA_VCF_VCENTER_USERNAME")
	vCenterPassword := os.Getenv("VRA_VCF_VCENTER_PASSWORD")
	workloadDomainID := os.Getenv("VRA_VCF_WORKLOAD_DOMAIN_ID")
	workloadDomainName := os.Getenv("VRA_VCF_WORKLOAD_DOMAIN_NAME")
	region := os.Getenv("VRA_VCF_REGION")
	dcName := os.Getenv("VRA_VCF_DATA_COLLECTOR_NAME")
	return fmt.Sprintf(`
	data "vra_data_collector" "dc" {
	  name = "%s"
	}

	resource "vra_cloud_account_vcf" "this" {
	  name                 = "test-vcf-account-%d"
	  description          = "my vcf cloud account"
	  workload_domain_id   = "%s"
	  workload_domain_name = "%s"

	  vcenter_username = "%s"
	  vcenter_password = "%s"
	  vcenter_hostname = "%s"

	  nsx_hostname = "%s"
	  nsx_password = "%s"
	  nsx_username = "%s"

	  dc_id                   = data.vra_data_collector.dc.id
	  regions                 = ["%s"]
	  accept_self_signed_cert = true

	  tags {
		key   = "foo"
		value = "bar"
	  }

	  tags {
		key   = "where"
		value = "waldo"
	  }
	}`, dcName, rInt, workloadDomainID, workloadDomainName, vCenterUserName, vCenterPassword, vCenterHostName, nsxHostName, nsxPassword, nsxUserName, region)
}

func testAccCheckVRACloudAccountVCFUpdateConfig(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	nsxHostName := os.Getenv("VRA_VCF_NSX_HOSTNAME")
	nsxUserName := os.Getenv("VRA_VCF_NSX_USERNAME")
	nsxPassword := os.Getenv("VRA_VCF_NSX_PASSWORD")
	vCenterHostName := os.Getenv("VRA_VCF_VCENTER_HOSTNAME")
	vCenterUserName := os.Getenv("VRA_VCF_VCENTER_USERNAME")
	vCenterPassword := os.Getenv("VRA_VCF_VCENTER_PASSWORD")
	workloadDomainID := os.Getenv("VRA_VCF_WORKLOAD_DOMAIN_ID")
	workloadDomainName := os.Getenv("VRA_VCF_WORKLOAD_DOMAIN_NAME")
	region := os.Getenv("VRA_VCF_REGION")
	dcName := os.Getenv("VRA_VCF_DATA_COLLECTOR_NAME")
	return fmt.Sprintf(`
	data "vra_data_collector" "dc" {
	  name = "%s"
	}

	resource "vra_cloud_account_vcf" "this" {
	  name                 = "test-vcf-account-%d"
	  description          = "your vcf cloud account"
	  workload_domain_id   = "%s"
	  workload_domain_name = "%s"

	  vcenter_username = "%s"
	  vcenter_password = "%s"
	  vcenter_hostname = "%s"

	  nsx_hostname = "%s"
	  nsx_password = "%s"
	  nsx_username = "%s"

	  dc_id                   = data.vra_data_collector.dc.id
	  regions                 = ["%s"]
	  accept_self_signed_cert = true

	  tags {
		key   = "foo"
		value = "bar"
	  }

	  tags {
		key   = "where"
		value = "waldo"
	  }
	}`, dcName, rInt, workloadDomainID, workloadDomainName, vCenterUserName, vCenterPassword, vCenterHostName, nsxHostName, nsxPassword, nsxUserName, region)
}

func testAccCheckVRACloudAccountVCFConfigDuplicateRegion(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	nsxHostName := os.Getenv("VRA_VCF_NSX_HOSTNAME")
	nsxUserName := os.Getenv("VRA_VCF_NSX_USERNAME")
	nsxPassword := os.Getenv("VRA_VCF_NSX_PASSWORD")
	vCenterHostName := os.Getenv("VRA_VCF_VCENTER_HOSTNAME")
	vCenterUserName := os.Getenv("VRA_VCF_VCENTER_USERNAME")
	vCenterPassword := os.Getenv("VRA_VCF_VCENTER_PASSWORD")
	workloadDomainID := os.Getenv("VRA_VCF_WORKLOAD_DOMAIN_ID")
	workloadDomainName := os.Getenv("VRA_VCF_WORKLOAD_DOMAIN_NAME")
	region := os.Getenv("VRA_VCF_REGION")
	dcName := os.Getenv("VRA_VCF_DATA_COLLECTOR_NAME")
	return fmt.Sprintf(`
	data "vra_data_collector" "dc" {
	  name = "%s"
	}

	resource "vra_cloud_account_vcf" "this" {
	  name                 = "test-vcf-account-%d"
	  description          = "my vcf cloud account"
	  workload_domain_id   = "%s"
	  workload_domain_name = "%s"

	  vcenter_username = "%s"
	  vcenter_password = "%s"
	  vcenter_hostname = "%s"

	  nsx_hostname = "%s"
	  nsx_password = "%s"
	  nsx_username = "%s"

	  dc_id                   = data.vra_data_collector.dc.id
	  regions                 = ["%s", "%s"]
	  accept_self_signed_cert = true

	  tags {
		key   = "foo"
		value = "bar"
	  }

	  tags {
		key   = "where"
		value = "waldo"
	  }
	}`, dcName, rInt, workloadDomainID, workloadDomainName, vCenterUserName, vCenterPassword, vCenterHostName, nsxHostName, nsxPassword, nsxUserName, region, region)
}
