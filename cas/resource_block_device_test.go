package cas

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCASBlockDevice_Basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCASBlockDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCASBlockDeviceConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCASBlockDeviceExists("cas_block_device.my_block_device"),
					resource.TestMatchResourceAttr(
						"cas_block_device.my_block_device", "name", regexp.MustCompile("^terraform_cas_block_device-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"cas_block_device.my_block_device", "capacity_in_gb", "4"),
					resource.TestCheckResourceAttr(
						"cas_block_device.my_block_device", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"cas_block_device.my_block_device", "tags.0.key", "stoyan"),
					resource.TestCheckResourceAttr(
						"cas_block_device.my_block_device", "tags.0.value", "genchev"),
				),
			},
		},
	})
}

func testAccCheckCASBlockDeviceExists(n string) resource.TestCheckFunc {
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

func testAccCheckCASBlockDeviceDestroy(s *terraform.State) error {
	/*
		apiClient := testAccProviderCAS.Meta().(*Client).apiClient

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "cas_block_device" {
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

func testAccCheckCASBlockDeviceConfig(rInt int) string {
	return fmt.Sprintf(`
resource "cas_block_device" "my_block_device" {
  name = "terraform_cas_block_device-%d"
  capacity_in_gb = 4

  tags {
	key = "stoyan"
    value = "genchev"
  }
}`, rInt)
}
