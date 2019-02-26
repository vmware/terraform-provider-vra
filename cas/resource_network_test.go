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

func TestAccTangoNetwork_Basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTangoNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTangoNetworkConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTangoNetworkExists("tango_network.my_network"),
					resource.TestMatchResourceAttr(
						"tango_network.my_network", "name", regexp.MustCompile("^terraform_tango_network-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"tango_network.my_network", "outbound_access", "false"),
					resource.TestCheckResourceAttr(
						"tango_network.my_network", "constraints.#", "1"),
					resource.TestCheckResourceAttr(
						"tango_network.my_network", "constraints.0.mandatory", "true"),
					resource.TestCheckResourceAttr(
						"tango_network.my_network", "constraints.0.expression", "pci"),

					// Currently not working as expected - possible bug on Provisioning service's side
					resource.TestCheckResourceAttr(
						"tango_network.my_network", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"tango_network.my_network", "tags.0.key", "stoyan"),
					resource.TestCheckResourceAttr(
						"tango_network.my_network", "tags.0.value", "genchev"),
				),
			},
		},
	})
}

func testAccCheckTangoNetworkExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Network ID is set")
		}

		return nil
	}
}

func testAccCheckTangoNetworkDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*tango.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tango_network" {
			continue
		}

		_, err := client.ReadResource("/iaas/networks/" + rs.Primary.ID)

		if err != nil && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf(
				"Error waiting for network (%s) to be destroyed: %s",
				rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckTangoNetworkConfig(rInt int) string {
	return fmt.Sprintf(`
resource "tango_network" "my_network" {
  name = "terraform_tango_network-%d"
  outbound_access = false

  tags {
	key = "stoyan"
    value = "genchev"
  }

  constraints {
	  mandatory = true
	  expression = "pci"
  }
}`, rInt)
}
