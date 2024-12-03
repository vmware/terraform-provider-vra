// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vra-sdk-go/pkg/client/storage_profile"
)

func TestAccVRAStorageProfileVsphereBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckStorageProfileVsphere(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAStorageProfileVsphereDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRAStorageProfileVsphereConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAStorageProfileVsphereExists("vra_storage_profile_vsphere.my-storage-profile-vsphere"),
					resource.TestCheckResourceAttr(
						"vra_storage_profile_vsphere.my-storage-profile-vsphere", "name", "vra-storage-profile-vsphere-test-0"),
					resource.TestCheckResourceAttr(
						"vra_storage_profile_vsphere.my-storage-profile-vsphere", "description", "my storage profile vsphere"),
					resource.TestCheckResourceAttr(
						"vra_storage_profile_vsphere.my-storage-profile-vsphere", "default_item", "true"),
					// TODO: Enable after https://jira.eng.vmware.com/browse/VCOM-20943 is resolved
					//resource.TestCheckResourceAttr(
					//	"vra_storage_profile_vsphere.my-storage-profile-vsphere", "disk_mode", "independent-persistent"),
					resource.TestCheckResourceAttr(
						"vra_storage_profile_vsphere.my-storage-profile-vsphere", "limit_iops", "1000"),
				),
			},
		},
	})
}

func testAccCheckVRAStorageProfileVsphereExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no storage profile vsphere ID is set")
		}

		return nil
	}
}

func testAccCheckVRAStorageProfileVsphereDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "vra_storage_profile_vsphere" {
			_, err := apiClient.StorageProfile.GetVSphereStorageProfile(storage_profile.NewGetVSphereStorageProfileParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_storage_profile_vsphere' still exists with id %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckVRAStorageProfileVsphereConfig() string {
	// Need valid credentials since this is creating a real cloud account
	cloudAccountName := os.Getenv("VRA_VSPHERE_CLOUD_ACCOUNT_NAME")

	region := os.Getenv("VRA_VSPHERE_REGION")
	return fmt.Sprintf(`
	data "vra_cloud_account_vsphere" "my-cloud-account" {
	name = "%s"
}

	data "vra_region" "my-region" {
	cloud_account_id = "${data.vra_cloud_account_vsphere.my-cloud-account.id}"
	region = "%s"
}

resource "vra_storage_profile_vsphere" "my-storage-profile-vsphere" {
	name = "vra-storage-profile-vsphere-test-0"
	description = "my storage profile vsphere"
	region_id = "${data.vra_region.my-region.id}"
	default_item = true
	disk_mode = "independent-persistent"
	limit_iops = "1000"
}`, cloudAccountName, region)
}
