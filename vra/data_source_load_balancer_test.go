// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVRALoadBalancerDataSource(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckLoadBalancerDataSource(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRALoadBalancerDataSourceConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.vra_load_balancer.test", "name", fmt.Sprintf("my-lb-%d", rInt)),
					resource.TestMatchResourceAttr("data.vra_load_balancer.test", "id", regexp.MustCompile("^lb-")),
					resource.TestCheckResourceAttr("data.vra_load_balancer.test", "routes.#", "1"),
					resource.TestCheckResourceAttr("data.vra_load_balancer.test", "nics.#", "1"),
					resource.TestCheckResourceAttr("data.vra_load_balancer.test", "targets.#", "1"),
				),
			},
		},
	})
}

func testAccCheckVRALoadBalancerDataSourceConfig(rInt int) string {
	return testAccCheckVRALoadBalancerDataSource(rInt) + fmt.Sprintf(`
resource "vra_load_balancer" "test" {
	name        = "my-lb-%d"
	project_id  = vra_project.my-project.id
	description = "load balancer description"

	targets {
		machine_id = vra_machine.my_machine.id
	}

	nics {
		network_id = data.vra_network.my-network.id
	}

	routes {
		protocol        = "TCP"
		port            = "80"
		member_protocol = "TCP"
		member_port     = "80"
		health_check_configuration = {
			protocol            = "TCP"
			port                = "80"
			interval_seconds    = 30
			timeout_seconds     = 10
			unhealthy_threshold = 2
			healthy_threshold   = 10
		}
	}
}

data "vra_load_balancer" "test" {
	id = vra_load_balancer.test.id
}
`, rInt)
}

func testAccPreCheckLoadBalancerDataSource(t *testing.T) {
	if os.Getenv("VRA_AWS_CLOUD_ACCOUNT_NAME") == "" {
		t.Fatal("VRA_AWS_CLOUD_ACCOUNT_NAME must be set for acceptance tests")
	}
	if os.Getenv("VRA_FABRIC_NETWORK") == "" {
		t.Fatal("VRA_FABRIC_NETWORK must be set for acceptance tests")
	}
	if os.Getenv("VRA_REGION") == "" {
		t.Fatal("VRA_REGION must be set for acceptance tests")
	}
	if os.Getenv("VRA_IMAGE") == "" {
		t.Fatal("VRA_IMAGE must be set for acceptance tests")
	}
	if os.Getenv("VRA_FLAVOR") == "" {
		t.Fatal("VRA_FLAVOR must be set for acceptance tests")
	}
}

func testAccCheckVRALoadBalancerDataSource(rInt int) string {
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
	region           = "%s"
}

resource "vra_zone" "my-zone" {
	name        = "my-zone-%d"
	description = "description my-vra-zone"
	region_id   = data.vra_region.my-region.id
}

resource "vra_project" "my-project" {
	name        = "my-project-%d"
	description = "test project"
	zone_assignments {
		zone_id       = vra_zone.my-zone.id
		priority      = 1
		max_instances = 2
	}
}

data "vra_fabric_network" "subnet" {
	filter = "name eq '%s'"
}

resource "vra_network_profile" "my-network-profile" {
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

data "vra_network" "my-network" {
	name       = data.vra_fabric_network.subnet.name
	depends_on = [vra_network_profile.my-network-profile]
}

resource "vra_image_profile" "my-image-profile" {
	name        = "my-image-profile-%d"
	description = "test image profile"
	region_id   = data.vra_region.my-region.id

	image_mapping {
		name       = "image"
		image_name = "%s"
	}
}

resource "vra_flavor_profile" "my-flavor-profile" {
	name        = "my-flavor-profile-%d"
	description = "my flavor"
	region_id   = data.vra_region.my-region.id
	flavor_mapping {
		name          = "flavor"
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
}
`, name, region, rInt, rInt, fabricNetwork, rInt, rInt, image, rInt, flavor, rInt)
}
