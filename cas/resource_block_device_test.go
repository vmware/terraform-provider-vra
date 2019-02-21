package cas

import (
	"fmt"
	"github.com/vmware/terraform-provider-cas/sdk"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTangoBlockDevice_Basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTangoBlockDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTangoBlockDeviceConfig_basic(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTangoBlockDeviceExists("tango_block_device.my_block_device"),
					resource.TestMatchResourceAttr(
						"tango_block_device.my_block_device", "name", regexp.MustCompile("^terraform_tango_block_device-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"tango_block_device.my_block_device", "capacity_in_gb", "4"),
					resource.TestCheckResourceAttr(
						"tango_block_device.my_block_device", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"tango_block_device.my_block_device", "tags.0.key", "stoyan"),
					resource.TestCheckResourceAttr(
						"tango_block_device.my_block_device", "tags.0.value", "genchev"),
				),
			},
		},
	})
}

func testAccCheckTangoBlockDeviceExists(n string) resource.TestCheckFunc {
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

func testAccCheckTangoBlockDeviceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*tango.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tango_block_device" {
			continue
		}

		_, err := client.ReadResource("/iaas/block-devices/" + rs.Primary.ID)

		if err != nil && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf(
				"Error waiting for block device (%s) to be destroyed: %s",
				rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckTangoBlockDeviceConfig_basic(rInt int) string {
	return fmt.Sprintf(`
resource "tango_block_device" "my_block_device" {
  name = "terraform_tango_block_device-%d"
  capacity_in_gb = 4

  tags {
	key = "stoyan"
    value = "genchev"
  }
}`, rInt)
}
