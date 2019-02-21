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

func TestAccTangoMachine_Basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTangoMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTangoMachineConfig_basic(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTangoMachineExists("tango_machine.my_machine"),
					resource.TestMatchResourceAttr(
						"tango_machine.my_machine", "name", regexp.MustCompile("^terraform_tango_machine-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"tango_machine.my_machine", "image", "ubuntu"),
					resource.TestCheckResourceAttr(
						"tango_machine.my_machine", "flavor", "small"),
					resource.TestCheckResourceAttr(
						"tango_machine.my_machine", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"tango_machine.my_machine", "tags.0.key", "stoyan"),
					resource.TestCheckResourceAttr(
						"tango_machine.my_machine", "tags.0.value", "genchev"),
				),
			},
		},
	})
}

func testAccCheckTangoMachineExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Machine ID is set")
		}

		return nil
	}
}

func testAccCheckTangoMachineDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*tango.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tango_machine" {
			continue
		}

		_, err := client.ReadResource("/iaas/machines/" + rs.Primary.ID)

		if err != nil && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf(
				"Error waiting for machine (%s) to be destroyed: %s",
				rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckTangoMachineConfig_basic(rInt int) string {
	return fmt.Sprintf(`
resource "tango_machine" "my_machine" {
  name = "terraform_tango_machine-%d"
  image = "ubuntu"
  flavor = "small"

  tags {
	key = "stoyan"
    value = "genchev"
  }
}`, rInt)
}
