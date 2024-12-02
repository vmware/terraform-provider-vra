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
	"github.com/vmware/vra-sdk-go/pkg/client/fabric_vsphere_datastore"
)

func TestAccVRAFabricDatastoreVsphere_importBasic(t *testing.T) {
	resourceName := "vra_fabric_datastore_vsphere"
	fabricDatastoreVsphereID := os.Getenv("VRA_FABRIC_DATASTORE_VSPHERE_ID")

	checkFn := func(s []*terraform.InstanceState) error {
		if len(s) != 1 {
			return fmt.Errorf("expected 1 state %#v", s)
		}

		fabricDatastoreVsphereState := s[0]
		if fabricDatastoreVsphereID != fabricDatastoreVsphereState.ID {
			return fmt.Errorf("expected fabric datastore vSphere ID of %s,%s received instead", fabricDatastoreVsphereID, fabricDatastoreVsphereState.ID)
		}
		return nil
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFabricDatastoreVsphere(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAFabricDatastoreVsphereDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckVRAFabricDatastoreVsphereConfig(),
				ExpectError: regexp.MustCompile("vra_fabric_datastore_vsphere resources are only importable"),
			},
			{
				Config:           testAccCheckVRAFabricDatastoreVsphereConfig(),
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateId:    fabricDatastoreVsphereID,
				ImportStateCheck: checkFn,
			},
		},
	})
}

func testAccCheckVRAFabricDatastoreVsphereDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "vra_fabric_datastore_vsphere" {
			if _, err := apiClient.FabricvSphereDatastore.GetFabricVSphereDatastore(fabric_vsphere_datastore.NewGetFabricVSphereDatastoreParams().WithID(rs.Primary.ID)); err == nil {
				return fmt.Errorf("Resource 'vra_fabric_datastore_vsphere' with id `%s` does not exist", rs.Primary.ID)
			}
		}

	}

	return nil
}

func testAccCheckVRAFabricDatastoreVsphereConfig() string {
	fabricDatastoreVsphereName := os.Getenv("VRA_FABRIC_DATASTORE_VSPHERE_NAME")
	return fmt.Sprintf(`
		resource "vra_fabric_datastore_vsphere" "this" {
			name = "%s"
		}`, fabricDatastoreVsphereName)
}
