// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vra-sdk-go/pkg/client/compute"
	"github.com/vmware/vra-sdk-go/pkg/client/flavor_profile"
	"github.com/vmware/vra-sdk-go/pkg/client/image_profile"
	"github.com/vmware/vra-sdk-go/pkg/client/load_balancer"
	"github.com/vmware/vra-sdk-go/pkg/client/location"
	"github.com/vmware/vra-sdk-go/pkg/client/network_profile"
	"github.com/vmware/vra-sdk-go/pkg/client/project"
)

func TestAccVRALoadBalancer_Basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckLoadBalancer(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRALoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRALoadBalancerConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRALoadBalancerExists("vra_load_balancer.my_load_balancer"),
					resource.TestMatchResourceAttr(
						"vra_load_balancer.my_load_balancer", "name", regexp.MustCompile("^my-lb-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"vra_load_balancer.my_load_balancer", "routes.#", "1"),
					resource.TestCheckResourceAttr(
						"vra_load_balancer.my_load_balancer", "nics.#", "1"),
					resource.TestCheckResourceAttr(
						"vra_load_balancer.my_load_balancer", "targets.#", "1"),
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
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "vra_project" {
			_, err := apiClient.Project.GetProject(project.NewGetProjectParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_project' still exists with id %s", rs.Primary.ID)
			}
		}
		if rs.Type == "vra_zone" {
			_, err := apiClient.Location.GetZone(location.NewGetZoneParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_zone' still exists with id %s", rs.Primary.ID)
			}
		}
		if rs.Type == "vra_machine" {
			_, err := apiClient.Compute.GetMachine(compute.NewGetMachineParams().WithID(rs.Primary.ID))

			_, ok := err.(*compute.GetMachineNotFound)
			if err != nil && !ok {
				return fmt.Errorf("error waiting for machine (%s) to be destroyed: %s", rs.Primary.ID, err)
			}
		}
		if rs.Type == "vra_image_profile" {
			_, err := apiClient.ImageProfile.GetImageProfile(image_profile.NewGetImageProfileParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_image_profile' still exists with id %s", rs.Primary.ID)
			}
		}
		if rs.Type == "vra_flavor_profile" {
			_, err := apiClient.FlavorProfile.GetFlavorProfile(flavor_profile.NewGetFlavorProfileParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_flavor_profile' still exists with id %s", rs.Primary.ID)
			}
		}
		if rs.Type == "vra_network_profile" {
			_, err := apiClient.NetworkProfile.GetNetworkProfile(network_profile.NewGetNetworkProfileParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_network_profile' still exists with id %s", rs.Primary.ID)
			}
		}
		if rs.Type == "vra_load_balancer" {
			_, err := apiClient.LoadBalancer.GetLoadBalancer(load_balancer.NewGetLoadBalancerParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_load_balancer' still exists with id %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckVRALoadBalancer(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	name := os.Getenv("VRA_AWS_CLOUD_ACCOUNT_NAME")
	fabricNetwork := os.Getenv("VRA_FABRIC_NETWORK")
	region := os.Getenv("VRA_REGION")
	image := os.Getenv("VRA_IMAGE")
	flavor := os.Getenv("VRA_FLAVOR")
	return fmt.Sprintf(`

	data "vra_cloud_account_aws" "my-cloud-account" {
		name = "%s"
	  }

data "vra_region" "my-region" {
    cloud_account_id = data.vra_cloud_account_aws.my-cloud-account.id
    region = "%s"
}

resource "vra_zone" "my-zone" {
    name = "my-zone-%d"
    description = "description my-vra-zone"
	region_id = data.vra_region.my-region.id
}

resource "vra_project" "my-project" {
	name = "my-project-%d"
	description = "test project"
	zone_assignments {
		zone_id       = vra_zone.my-zone.id
		priority      = 1
		max_instances = 2
	  }
 }
 data vra_fabric_network "subnet" {
    filter = "name eq '%s'"
}

resource vra_network_profile "my-network-profile" {
  name        = "my-network-profile-%d"
  description = "test network profile"
  region_id   = data.vra_region.my-region.id

  fabric_network_ids = [
    data.vra_fabric_network.subnet.id,
  ]

  isolation_type = "NONE"

  tags {
    key   = "foo"
    value = "bar"
  }
}

data vra_network "my-network" {
  name = data.vra_fabric_network.subnet.name
  depends_on = [vra_network_profile.my-network-profile]
}

resource "vra_image_profile" "my-image-profile" {
	name        = "my-image-profile-%d"
	description = "test image profile"
	region_id = data.vra_region.my-region.id

	image_mapping {
	  name       = "image"
	  image_name = "%s"
	}
  }

resource "vra_flavor_profile" "my-flavor-profile" {
	name = "my-flavor-profile-%d"
	description = "my flavor"
	region_id = data.vra_region.my-region.id
	flavor_mapping {
		name = "flavor"
		instance_type = "%s"
	}
}

resource "vra_machine" "my-machine" {
	name        = "my-machine-%d"
	description = "test machine updated"
	project_id  = vra_project.my-project.id
	image       = "image"
	flavor      = "flavor"

	tags {
	  key   = "foo"
	  value = "bar"
	}
}`, name, region, rInt, rInt, fabricNetwork, rInt, rInt, image, rInt, flavor, rInt)
}

// TODO: Enable this test when the issue https://jira.eng.vmware.com/browse/VCOM-13736 is fixed
/*
func testAccCheckVRALoadBalancerNoTargetLinkConfig(rInt int) string {

	return testAccCheckVRALoadBalancer(rInt) + fmt.Sprintf(`
	resource "vra_load_balancer" "my-load-balancer" {
		name = "my-load-balancer"
		project_id = vra_project.my-project.id
		description = "My Load Balancer"
		custom_properties = {
			"edgeClusterRouterStateLink" = "/resources/routers/<uuid>"
			"tier0LogicalRouterStateLink" = "/resources/routers/<uuid>"
		}
		targets {
			machine_id = data.vra_machine.my-machine.id
		}

		nics {
			network_id = data.vra_network.my-network.id
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
				timeout_seconds = 10
				unhealthy_threshold = 2
				healthy_threshold = 10
			}
		}
		}`, rInt)
} */

func testAccCheckVRALoadBalancerConfig(rInt int) string {

	return testAccCheckVRALoadBalancer(rInt) + fmt.Sprintf(`
	resource "vra_load_balancer" "my_load_balancer" {
		name = "my-lb-%d"
		project_id = vra_project.my-project.id
		description = "load balancer description"

		targets {
			machine_id = vra_machine.my_machine.id
		}

		nics {
			network_id = data.vra_network.my-network.id
		}

		routes {
			protocol = "TCP"
			port = "80"
			member_protocol = "TCP"
			member_port = "80"
			health_check_configuration = {
				protocol = "TCP"
				port = "80"
				interval_seconds = 30
				timeout_seconds = 10
				unhealthy_threshold = 2
				healthy_threshold = 10
			}
		}
		}`, rInt)
}
