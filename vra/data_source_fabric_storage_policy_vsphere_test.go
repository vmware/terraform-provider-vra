// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFabricStoragePolicyVsphere_Basic(t *testing.T) {
	dsName := os.Getenv("VRA_VSPHERE_STORAGE_POLICY_NAME")
	datasourceName := "data.vra_fabric_storage_policy_vsphere.sp"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVsphereForStoragePolicy(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccFabricStoragePolicyVsphereConfig(dsName + "foo"),
				ExpectError: regexp.MustCompile("fabric vSphere storage policies filter doesn't match to any storage policy"),
			},
			{
				Config: testAccFabricStoragePolicyVsphereConfig(dsName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", dsName),
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
				),
			},
		},
	})
}

func testAccFabricStoragePolicyVsphereConfig(dsName string) string {
	return fmt.Sprintf(`
		data "vra_fabric_storage_policy_vsphere" "sp" {
		  filter = "name eq '%s'"
		}`, dsName)
}
