package vra

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/vmware/vra-sdk-go/pkg/client/compute"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
				Config: testAccCheckVRAMachineConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAMachineExists("vra_machine.my-machine"),
					resource.TestMatchResourceAttr(
						"vra_machine.my-machine", "name", regexp.MustCompile("^my-machine-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(
						"vra_machine.my-machine", "description", "test machine"),
					resource.TestCheckResourceAttr(
						"vra_machine.my-machine", "image", "ubuntu"),
					resource.TestCheckResourceAttr(
						"vra_machine.my-machine", "flavor", "small"),
					resource.TestCheckResourceAttr(
						"vra_machine.my-machine", "tags.#", "1"),
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
		switch rs.Type {
		case "vra_machine":
			{
				_, err := apiClient.Compute.GetMachine(compute.NewGetMachineParams().WithID(rs.Primary.ID))

				_, ok := err.(*compute.GetMachineNotFound)
				if err != nil && !ok {
					return fmt.Errorf("error waiting for machine (%s) to be destroyed: %s", rs.Primary.ID, err)
				}
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
	flavor      = "small"
  
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
	image       = "ubuntu"
	flavor      = "small"
  
	tags {
	  key   = "foo"
	  value = "bar"
	}
}`, rInt)
}

func testAccCheckVRAMachine(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	name := os.Getenv("VRA_AWS_CLOUD_ACCOUNT_NAME")
	image := os.Getenv("VRA_IMAGE")
	return fmt.Sprintf(`

	data "vra_cloud_account_aws" "my-cloud-account" {
		name = "%s"
	  }

data "vra_region" "us-east-1-region" {
    cloud_account_id = data.vra_cloud_account_aws.my-cloud-account.id
    region = "us-east-1"
}

resource "vra_zone" "my-zone" {
    name = "my-zone-%d"
    description = "description my-vra-zone"
	region_id = data.vra_region.us-east-1-region.id
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
	region_id = data.vra_region.us-east-1-region.id
  
	image_mapping {
	  name       = "ubuntu"
	  image_name = "%s"
	}
  }

resource "vra_flavor_profile" "my-flavor-profile" {
	name = "my-flavor-profile-%d"
	description = "my flavor"
	region_id = data.vra_region.us-east-1-region.id
	flavor_mapping {
		name = "small"
		instance_type = "t2.small"
	}
	flavor_mapping {
		name = "medium"
		instance_type = "t2.medium"
	}
}`, name, rInt, rInt, rInt, image, rInt)
}
