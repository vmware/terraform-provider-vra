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

	"github.com/vmware/vra-sdk-go/pkg/client/compute"
	"github.com/vmware/vra-sdk-go/pkg/client/flavor_profile"
	"github.com/vmware/vra-sdk-go/pkg/client/image_profile"
	"github.com/vmware/vra-sdk-go/pkg/client/location"
	"github.com/vmware/vra-sdk-go/pkg/client/project"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVRAMachine_Basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckMachine(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckVRAMachineNoImageConfig(rInt),
				ExpectError: regexp.MustCompile("image or image_ref required"),
			},
			{
				// Machine creation
				Config: testAccCheckVRAMachineConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAMachineExists("vra_machine.my_machine"),
					resource.TestMatchResourceAttr(
						"vra_machine.my_machine", "name", regexp.MustCompile("^my-machine-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"vra_machine.my_machine", "description", "test machine"),
					resource.TestCheckResourceAttr(
						"vra_machine.my_machine", "image", "image_name1"),
					resource.TestCheckResourceAttr(
						"vra_machine.my_machine", "flavor", "flavor1"),
					resource.TestCheckResourceAttr(
						"vra_machine.my_machine", "tags.#", "1"),
				),
			},
			{
				// Machine resize due to change in flavor
				Config: testAccCheckVRAMachineUpdateConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAMachineExists("vra_machine.my_machine"),
					resource.TestMatchResourceAttr(
						"vra_machine.my_machine", "name", regexp.MustCompile("^my-machine-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"vra_machine.my_machine", "description", "test machine updated"),
					resource.TestCheckResourceAttr(
						"vra_machine.my_machine", "image", "image_name1"),
					resource.TestCheckResourceAttr(
						"vra_machine.my_machine", "flavor", "flavor2"),
					resource.TestCheckResourceAttr(
						"vra_machine.my_machine", "tags.#", "1"),
				),
			},
			{
				// Machine recreate (destroy and create) due to change in image
				Config: testAccCheckVRAMachineReCreateConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAMachineExists("vra_machine.my_machine"),
					resource.TestMatchResourceAttr(
						"vra_machine.my_machine", "name", regexp.MustCompile("^my-machine-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"vra_machine.my_machine", "description", "test machine updated"),
					resource.TestCheckResourceAttr(
						"vra_machine.my_machine", "image", "image_name2"),
					resource.TestCheckResourceAttr(
						"vra_machine.my_machine", "flavor", "flavor2"),
					resource.TestCheckResourceAttr(
						"vra_machine.my_machine", "tags.#", "1"),
				),
			},
		},
	})
}

func testAccCheckVRAMachineExists(n string) resource.TestCheckFunc {
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

func testAccCheckVRAMachineDestroy(s *terraform.State) error {
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
	}

	return nil
}

func testAccCheckVRAMachineNoImageConfig(rInt int) string {

	return testAccCheckVRAMachine(rInt) + fmt.Sprintf(`
resource "vra_machine" "my_machine" {
	name        = "my-machine-%d"
	description = "test machine"
	project_id  = vra_project.my-project.id
	flavor      = "flavor"

	tags {
	  key   = "foo"
	  value = "bar"
	}
}`, rInt)
}

func testAccCheckVRAMachineConfig(rInt int) string {

	return testAccCheckVRAMachine(rInt) + fmt.Sprintf(`
resource "vra_machine" "my_machine" {
	name        = "my-machine-%d"
	description = "test machine"
	project_id  = vra_project.my-project.id
	image       = "image_name1"
	flavor      = "flavor1"

	tags {
	  key   = "foo"
	  value = "bar"
	}
}`, rInt)
}

func testAccCheckVRAMachineUpdateConfig(rInt int) string {

	return testAccCheckVRAMachine(rInt) + fmt.Sprintf(`
resource "vra_machine" "my_machine" {
	name        = "my-machine-%d"
	description = "test machine updated"
	project_id  = vra_project.my-project.id
	image       = "image_name1"
	flavor      = "flavor2"

	tags {
	  key   = "foo"
	  value = "bar"
	}
}`, rInt)
}

func testAccCheckVRAMachineReCreateConfig(rInt int) string {

	return testAccCheckVRAMachine(rInt) + fmt.Sprintf(`
resource "vra_machine" "my_machine" {
	name        = "my-machine-%d"
	description = "test machine updated"
	project_id  = vra_project.my-project.id
	image       = "image_name2"
	flavor      = "flavor2"

	tags {
	  key   = "foo"
	  value = "bar"
	}
}`, rInt)
}

func testAccCheckVRAMachine(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	name := os.Getenv("VRA_AWS_CLOUD_ACCOUNT_NAME")
	image1 := os.Getenv("VRA_IMAGE_1")
	image2 := os.Getenv("VRA_IMAGE_2")
	flavor1 := os.Getenv("VRA_FLAVOR_1")
	flavor2 := os.Getenv("VRA_FLAVOR_2")
	region := os.Getenv("VRA_REGION")
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

resource "vra_image_profile" "this" {
	name        = "my-image-profile-%d"
	description = "test image profile"
	region_id = data.vra_region.my-region.id

	image_mapping {
	  name       = "image_name1"
	  image_name = "%s"
	}

    image_mapping {
	  name       = "image_name2"
	  image_name = "%s"
	}
  }

resource "vra_flavor_profile" "my-flavor-profile" {
	name = "my-flavor-profile-%d"
	description = "my flavor"
	region_id = data.vra_region.my-region.id
	flavor_mapping {
		name = "flavor1"
		instance_type = "%s"
	}
	flavor_mapping {
		name = "flavor2"
		instance_type = "%s"
	}
}`, name, region, rInt, rInt, rInt, image1, image2, rInt, flavor1, flavor2)
}
