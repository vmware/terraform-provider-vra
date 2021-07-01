package vra

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFabricDatastoreVsphere_Basic(t *testing.T) {
	dsName := os.Getenv("VRA_VSPHERE_DATASTORE_NAME")
	datasourceName := "data.vra_fabric_datastore_vsphere.datastore_vsphere"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVsphereForDataStore(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckFabricDatastoreVsphereConfig(dsName + "foo"),
				ExpectError: regexp.MustCompile("fabric vSphere datastore filter doesn't match to any datastore"),
			},
			{
				Config: testAccCheckFabricDatastoreVsphereConfig(dsName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", dsName),
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
				),
			},
		},
	})
}

func testAccCheckFabricDatastoreVsphereConfig(dsName string) string {
	return fmt.Sprintf(`
		data "vra_fabric_datastore_vsphere" "datastore_vsphere" {
		  filter = "name eq '%s'"
		}`, dsName)
}
