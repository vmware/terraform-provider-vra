// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceVRACloudAccountVCF(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName1 := "vra_cloud_account_vcf.this"
	dataSourceName1 := "data.vra_cloud_account_vcf.this"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVCF(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRACloudAccountVCFNotFound(rInt),
				ExpectError: regexp.MustCompile("cloud account with name 'foobar' not found"),
			},
			{
				Config: testAccDataSourceVRACloudAccountVCF(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "dc_id", dataSourceName1, "dc_id"),
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
					resource.TestCheckResourceAttrPair(resourceName1, "nsx_hostname", dataSourceName1, "nsx_hostname"),
					resource.TestCheckResourceAttrPair(resourceName1, "regions", dataSourceName1, "regions"),
					resource.TestCheckResourceAttrPair(resourceName1, "tags", dataSourceName1, "tags"),
					resource.TestCheckResourceAttrPair(resourceName1, "sddc_manager_id", dataSourceName1, "sddc_manager_id"),
					resource.TestCheckResourceAttrPair(resourceName1, "vcenter_hostname", dataSourceName1, "vcenter_hostname"),
					resource.TestCheckResourceAttrPair(resourceName1, "vcenter_username", dataSourceName1, "vcenter_username"),
					resource.TestCheckResourceAttrPair(resourceName1, "workload_domain_id", dataSourceName1, "workload_domain_id"),
					resource.TestCheckResourceAttrPair(resourceName1, "workload_domain_name", dataSourceName1, "workload_domain_name"),
				),
			},
		},
	})
}

func testAccDataSourceVRACloudAccountVCFBase(rInt int) string {
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
	  description          = "tf test vcf cloud account"
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

func testAccDataSourceVRACloudAccountVCFNotFound(rInt int) string {
	return testAccDataSourceVRACloudAccountVCFBase(rInt) + `
	data "vra_cloud_account_vcf" "this" {
		name = "foobar"
	}`
}

func testAccDataSourceVRACloudAccountVCF(rInt int) string {
	return testAccDataSourceVRACloudAccountVCFBase(rInt) + `
	data "vra_cloud_account_vcf" "this" {
		name = "${vra_cloud_account_vcf.this.name}"
	}`
}
