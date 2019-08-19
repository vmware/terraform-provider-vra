package vra

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccVRABlockDevice_Basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRABlockDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRABlockDeviceConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRABlockDeviceExists("vra_block_device.my_block_device"),
					resource.TestMatchResourceAttr(
						"vra_block_device.my_block_device", "name", regexp.MustCompile("^terraform_vra_block_device-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"vra_block_device.my_block_device", "capacity_in_gb", "4"),
					resource.TestCheckResourceAttr(
						"vra_block_device.my_block_device", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"vra_block_device.my_block_device", "tags.0.key", "stoyan"),
					resource.TestCheckResourceAttr(
						"vra_block_device.my_block_device", "tags.0.value", "genchev"),
				),
			},
		},
	})
}

func testAccCheckVRABlockDeviceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Block Device ID is set")
		}

		return nil
	}
}

func testAccCheckVRABlockDeviceDestroy(s *terraform.State) error {
	/*
		apiClient := testAccProviderVRA.Meta().(*Client).apiClient

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "vra_block_device" {
				continue
			}

			_, err := client.ReadResource("/iaas/block-devices/" + rs.Primary.ID)

			if err != nil && !strings.Contains(err.Error(), "404") {
				return fmt.Errorf(
					"Error waiting for block device (%s) to be destroyed: %s",
					rs.Primary.ID, err)
			}
		}
	*/

	return nil
}

func testAccCheckVRABlockDeviceConfig(rInt int) string {
	return fmt.Sprintf(`
resource "vra_block_device" "my_block_device" {
  name = "terraform_vra_block_device-%d"
  capacity_in_gb = 4

  tags {
	key = "stoyan"
    value = "genchev"
  }
}`, rInt)
}
