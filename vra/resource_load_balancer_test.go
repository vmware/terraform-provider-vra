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

func TestAccVRALoadBalancer_Basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRALoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRALoadBalancerConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRALoadBalancerExists("vra_load_balancer.my_load_balancer"),
					resource.TestMatchResourceAttr(
						"vra_load_balancer.my_load_balancer", "name", regexp.MustCompile("^terraformcasloadbalancer-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"vra_load_balancer.my_load_balancer", "nics.#", "1"),
					resource.TestCheckResourceAttr(
						"vra_load_balancer.my_load_balancer", "target_links.#", "1"),
					resource.TestCheckResourceAttr(
						"vra_load_balancer.my_load_balancer", "routes.#", "1"),
					resource.TestCheckResourceAttr(
						"vra_load_balancer.my_load_balancer", "routes.0.port", "80"),
					resource.TestCheckResourceAttr(
						"vra_load_balancer.my_load_balancer", "routes.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(
						"vra_load_balancer.my_load_balancer", "routes.0.member_port", "80"),
					resource.TestCheckResourceAttr(
						"vra_load_balancer.my_load_balancer", "routes.0.member_protocol", "TCP"),
					resource.TestCheckResourceAttr(
						"vra_load_balancer.my_load_balancer", "routes.0.health_check_configuration.#", "1"),
					resource.TestCheckResourceAttr(
						"vra_load_balancer.my_load_balancer", "routes.0.health_check_configuration.0.port", "80"),
					resource.TestCheckResourceAttr(
						"vra_load_balancer.my_load_balancer", "routes.0.health_check_configuration.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(
						"vra_load_balancer.my_load_balancer", "routes.0.health_check_configuration.0.unhealthy_threshold", "2"),
					resource.TestCheckResourceAttr(
						"vra_load_balancer.my_load_balancer", "routes.0.health_check_configuration.0.timeout_seconds", "5"),
					resource.TestCheckResourceAttr(
						"vra_load_balancer.my_load_balancer", "routes.0.health_check_configuration.0.interval_seconds", "30"),
					resource.TestCheckResourceAttr(
						"vra_load_balancer.my_load_balancer", "routes.0.health_check_configuration.0.healthy_threshold", "10"),
				),
			},
		},
	})
}

func testAccCheckVRALoadBalancerExists(n string) resource.TestCheckFunc {
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

func testAccCheckVRALoadBalancerDestroy(s *terraform.State) error {
	/*
		apiClient := testAccProviderVRA.Meta().(*Client).apiClient


		for _, rs := range s.RootModule().Resources {
			if rs.Type != "vra_load_balancer" {
				continue
			}

			_, err := client.ReadResource("/iaas/load-balancers/" + rs.Primary.ID)

			if err != nil && !strings.Contains(err.Error(), "404") {
				return fmt.Errorf(
					"Error waiting for load balancer (%s) to be destroyed: %s",
					rs.Primary.ID, err)
			}
		}
	*/

	return nil
}

func testAccCheckVRALoadBalancerConfig(rInt int) string {
	return fmt.Sprintf(`
resource "vra_network" "my_network" {
	name = "terraform_vra_network"

	constraints {
		mandatory = true
		expression = "pci"
	}
}

resource "vra_machine" "my_machine" {
	name = "terraform_vra_machine"
	
	image = "ubuntu"
	flavor = "small"	

	nics {
        network_id = "${vra_network.my_network.id}"
    }
}	

resource "vra_load_balancer" "my_load_balancer" {
	name = "terraformcasloadbalancer-%d"

    nics {
        network_id = "${vra_network.my_network.id}"
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

    target_links = ["${vra_machine.my_machine.self_link}"]
}`, rInt)
}
