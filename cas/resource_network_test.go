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

func TestAccCASNetwork_Basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCASNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCASNetworkConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCASNetworkExists("cas_network.my_network"),
					resource.TestMatchResourceAttr(
						"cas_network.my_network", "name", regexp.MustCompile("^terraform_cas_network-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"cas_network.my_network", "outbound_access", "false"),
					resource.TestCheckResourceAttr(
						"cas_network.my_network", "constraints.#", "1"),
					resource.TestCheckResourceAttr(
						"cas_network.my_network", "constraints.0.mandatory", "true"),
					resource.TestCheckResourceAttr(
						"cas_network.my_network", "constraints.0.expression", "pci"),

					// Currently not working as expected - possible bug on Provisioning service's side
					resource.TestCheckResourceAttr(
						"cas_network.my_network", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"cas_network.my_network", "tags.0.key", "stoyan"),
					resource.TestCheckResourceAttr(
						"cas_network.my_network", "tags.0.value", "genchev"),
				),
			},
		},
	})
}

func testAccCheckCASNetworkExists(n string) resource.TestCheckFunc {
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

func testAccCheckCASNetworkDestroy(s *terraform.State) error {
	/*
		apiClient := testAccProviderCAS.Meta().(*Client).apiClient

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "cas_network" {
				continue
			}

			_, err := apiClient.ReadResource("/iaas/networks/" + rs.Primary.ID)

			if err != nil && !strings.Contains(err.Error(), "404") {
				return fmt.Errorf(
					"Error waiting for network (%s) to be destroyed: %s",
					rs.Primary.ID, err)
			}
		}
	*/

	return nil
}

func testAccCheckCASNetworkConfig(rInt int) string {
	return fmt.Sprintf(`
resource "cas_network" "my_network" {
  name = "terraform_cas_network-%d"
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
