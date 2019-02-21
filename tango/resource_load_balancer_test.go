package tango

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"tango-terraform-provider/tango/client"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTangoLoadBalancer_Basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTangoLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTangoLoadBalancerConfig_basic(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTangoLoadBalancerExists("tango_load_balancer.my_load_balancer"),
					resource.TestMatchResourceAttr(
						"tango_load_balancer.my_load_balancer", "name", regexp.MustCompile("^terraformtangoloadbalancer-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"tango_load_balancer.my_load_balancer", "nics.#", "1"),
					resource.TestCheckResourceAttr(
						"tango_load_balancer.my_load_balancer", "target_links.#", "1"),
					resource.TestCheckResourceAttr(
						"tango_load_balancer.my_load_balancer", "routes.#", "1"),
					resource.TestCheckResourceAttr(
						"tango_load_balancer.my_load_balancer", "routes.0.port", "80"),
					resource.TestCheckResourceAttr(
						"tango_load_balancer.my_load_balancer", "routes.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(
						"tango_load_balancer.my_load_balancer", "routes.0.member_port", "80"),
					resource.TestCheckResourceAttr(
						"tango_load_balancer.my_load_balancer", "routes.0.member_protocol", "TCP"),
					resource.TestCheckResourceAttr(
						"tango_load_balancer.my_load_balancer", "routes.0.health_check_configuration.#", "1"),
					resource.TestCheckResourceAttr(
						"tango_load_balancer.my_load_balancer", "routes.0.health_check_configuration.0.port", "80"),
					resource.TestCheckResourceAttr(
						"tango_load_balancer.my_load_balancer", "routes.0.health_check_configuration.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(
						"tango_load_balancer.my_load_balancer", "routes.0.health_check_configuration.0.unhealthy_threshold", "2"),
					resource.TestCheckResourceAttr(
						"tango_load_balancer.my_load_balancer", "routes.0.health_check_configuration.0.timeout_seconds", "5"),
					resource.TestCheckResourceAttr(
						"tango_load_balancer.my_load_balancer", "routes.0.health_check_configuration.0.interval_seconds", "30"),
					resource.TestCheckResourceAttr(
						"tango_load_balancer.my_load_balancer", "routes.0.health_check_configuration.0.healthy_threshold", "10"),
				),
			},
		},
	})
}

func testAccCheckTangoLoadBalancerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Load Balancer ID is set")
		}

		return nil
	}
}

func testAccCheckTangoLoadBalancerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*tango.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tango_load_balancer" {
			continue
		}

		_, err := client.ReadResource("/iaas/load-balancers/" + rs.Primary.ID)

		if err != nil && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf(
				"Error waiting for load balancer (%s) to be destroyed: %s",
				rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckTangoLoadBalancerConfig_basic(rInt int) string {
	return fmt.Sprintf(`
resource "tango_network" "my_network" {
	name = "terraform_tango_network"

	constraints {
		mandatory = true
		expression = "pci"
	}
}

resource "tango_machine" "my_machine" {
	name = "terraform_tango_machine"
	
	image = "ubuntu"
	flavor = "small"	

	nics {
        network_id = "${tango_network.my_network.id}"
    }
}	

resource "tango_load_balancer" "my_load_balancer" {
	name = "terraformtangoloadbalancer-%d"

    nics {
        network_id = "${tango_network.my_network.id}"
    }

    routes {
        protocol = "TCP"
        port = "80"
        member_protocol = "TCP"
		member_port = "80"
		health_check_configuration {
			protocol = "TCP"
            port = "80"
            interval_seconds = 30
            timeout_seconds = 5
            unhealthy_threshold = 2
            healthy_threshold = 10
		}
    }

    target_links = ["${tango_machine.my_machine.self_link}"]
}`, rInt)
}
