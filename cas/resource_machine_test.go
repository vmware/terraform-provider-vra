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
				Config: testAccCheckTangoMachineConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTangoMachineExists("cas_machine.my_machine"),
					resource.TestMatchResourceAttr(
						"cas_machine.my_machine", "name", regexp.MustCompile("^terraform_cas_machine-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"cas_machine.my_machine", "image", "ubuntu"),
					resource.TestCheckResourceAttr(
						"cas_machine.my_machine", "flavor", "small"),
					resource.TestCheckResourceAttr(
						"cas_machine.my_machine", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"cas_machine.my_machine", "tags.0.key", "stoyan"),
					resource.TestCheckResourceAttr(
						"cas_machine.my_machine", "tags.0.value", "genchev"),
				),
			},
		},
	})
}

func testAccCheckTangoMachineExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no machine ID is set")
		}

		return nil
	}
}

func testAccCheckTangoMachineDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*tango.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cas_machine" {
			continue
		}

		_, err := client.ReadResource("/iaas/machines/" + rs.Primary.ID)

		if err != nil && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf(
				"error waiting for machine (%s) to be destroyed: %s",
				rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckTangoMachineConfig(rInt int) string {
	return fmt.Sprintf(`
resource "cas_machine" "my_machine" {
  name = "terraform_cas_machine-%d"
  image = "ubuntu"
  flavor = "small"

  tags {
	key = "stoyan"
    value = "genchev"
  }
}`, rInt)
}
