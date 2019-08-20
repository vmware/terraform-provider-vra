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

func TestAccVRAMachineBasic(t *testing.T) {
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
					testAccCheckVRAMachineExists("vra_machine.my_machine"),
					resource.TestMatchResourceAttr(
						"vra_machine.my_machine", "name", regexp.MustCompile("^terraform_vra_machine-"+strconv.Itoa(rInt))),
					// TODO: Enable when https://jira.eng.vmware.com/browse/VCOM-10339 is resolved.
					//resource.TestCheckResourceAttr(
					//	"vra_machine.my_machine", "description", "Created by terraform provider test"),
					resource.TestCheckResourceAttr(
						"vra_machine.my_machine", "image", "ubuntu"),
					resource.TestCheckResourceAttr(
						"vra_machine.my_machine", "flavor", "small"),
					resource.TestCheckResourceAttr(
						"vra_machine.my_machine", "tags.#", "2"),
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

func testAccCheckVRAMachineConfig(rInt int) string {
	// Need valid details since this is using existing project
	image := os.Getenv("VRA_IMAGE")
	flavor := os.Getenv("VRA_FLAVOR")
	projectID := os.Getenv("VRA_PROJECT_ID")
	return fmt.Sprintf(`
resource "vra_machine" "my_machine" {
  name = "terraform_vra_machine-%d"
  project_id = "%s"
  image = "%s"
  flavor = "%s"

  tags {
	key = "description"
    value = "Testing Terraform Provider for VRA"
  }

  tags {
    key = "foo"
    value = "bar"
  }
}`, rInt, projectID, image, flavor)
}

func testAccCheckVRAMachineNoImageConfig(rInt int) string {
	flavor := os.Getenv("VRA_FLAVOR")
	projectID := os.Getenv("VRA_PROJECT_ID")
	return fmt.Sprintf(`
resource "vra_machine" "my_machine" {
  name = "terraform_vra_machine-%d"
  project_id = "%s"
  flavor = "%s"

  tags {
	key = "description"
    value = "Testing Terraform Provider for VRA"
  }
}`, rInt, projectID, flavor)
}
