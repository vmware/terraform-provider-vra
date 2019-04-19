package cas

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/vmware/cas-sdk-go/pkg/client/compute"
	tango "github.com/vmware/terraform-provider-cas/sdk"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTangoMachineBasic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckMachine(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTangoMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckTangoMachineNoImageConfig(rInt),
				ExpectError: regexp.MustCompile("image or image_ref required"),
			},
			{
				Config: testAccCheckTangoMachineConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTangoMachineExists("cas_machine.my_machine"),
					resource.TestMatchResourceAttr(
						"cas_machine.my_machine", "name", regexp.MustCompile("^terraform_cas_machine-"+strconv.Itoa(rInt))),
					// TODO: Enable when https://jira.eng.vmware.com/browse/VCOM-10339 is resolved.
					//resource.TestCheckResourceAttr(
					//	"cas_machine.my_machine", "description", "Created by terraform provider test"),
					resource.TestCheckResourceAttr(
						"cas_machine.my_machine", "image", "ubuntu"),
					resource.TestCheckResourceAttr(
						"cas_machine.my_machine", "flavor", "small"),
					resource.TestCheckResourceAttr(
						"cas_machine.my_machine", "tags.#", "2"),
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
	client := testAccProviderCAS.Meta().(*tango.Client)
	apiClient := client.GetAPIClient()

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "cas_machine":
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

func testAccCheckTangoMachineConfig(rInt int) string {
	// Need valid details since this is using existing project
	image := os.Getenv("CAS_IMAGE")
	flavor := os.Getenv("CAS_FLAVOR")
	projectID := os.Getenv("CAS_PROJECT_ID")
	return fmt.Sprintf(`
resource "cas_machine" "my_machine" {
  name = "terraform_cas_machine-%d"
  project_id = "%s"
  image = "%s"
  flavor = "%s"

  tags {
	key = "description"
    value = "Testing Terraform Provider for CAS"
  }

  tags {
    key = "foo"
    value = "bar"
  }
}`, rInt, projectID, image, flavor)
}

func testAccCheckTangoMachineNoImageConfig(rInt int) string {
	flavor := os.Getenv("CAS_FLAVOR")
	projectID := os.Getenv("CAS_PROJECT_ID")
	return fmt.Sprintf(`
resource "cas_machine" "my_machine" {
  name = "terraform_cas_machine-%d"
  project_id = "%s"
  flavor = "%s"

  tags {
	key = "description"
    value = "Testing Terraform Provider for CAS"
  }
}`, rInt, projectID, flavor)
}
