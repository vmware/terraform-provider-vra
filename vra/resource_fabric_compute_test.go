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
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vra-sdk-go/pkg/client/fabric_compute"
)

func TestAccVRAFabricCompute_importBasic(t *testing.T) {
	resourceName := "vra_fabric_compute.this"
	fabricComputeID := os.Getenv("VRA_FABRIC_COMPUTE_ID")

	checkFn := func(s []*terraform.InstanceState) error {
		if len(s) != 1 {
			return fmt.Errorf("expected 1 state %#v", s)
		}

		fabricComputeState := s[0]
		if fabricComputeID != fabricComputeState.ID {
			return fmt.Errorf("expected fabric compute ID of %s,%s received instead", fabricComputeID, fabricComputeState.ID)
		}
		return nil
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFabricCompute(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAFabricComputeDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckVRAFabricComputeConfig(),
				ExpectError: regexp.MustCompile("vra_fabric_compute resources are only importable"),
			},
			{
				Config:           testAccCheckVRAFabricComputeConfig(),
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateId:    fabricComputeID,
				ImportStateCheck: checkFn,
			},
		},
	})
}

func testAccCheckVRAFabricComputeDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "vra_fabric_compute" {
			if _, err := apiClient.FabricCompute.GetFabricCompute(fabric_compute.NewGetFabricComputeParams().WithID(rs.Primary.ID)); err == nil {
				return fmt.Errorf("Resource 'vra_fabric_compute' with id `%s` does not exist", rs.Primary.ID)
			}
		}

	}

	return nil
}

func testAccCheckVRAFabricComputeConfig() string {
	regionName := os.Getenv("VRA_FABRIC_COMPUTE_NAME")
	return fmt.Sprintf(`
		resource "vra_fabric_compute" "this" {
			name = "%s"
		}`, regionName)
}
