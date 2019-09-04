package vra

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/client/location"
	"github.com/vmware/vra-sdk-go/pkg/client/project"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

func TestAccVRAProjectBasic(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRAProjectConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAProjectExists("vra_project.my-project"),
					resource.TestCheckResourceAttr(
						"vra_project.my-project", "name", "my-project-"+strconv.Itoa(rInt)),
					resource.TestCheckResourceAttr(
						"vra_project.my-project", "description", "test project"),
				),
			},
			{
				Config: testAccCheckVRAProjectUpdateConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAProjectExists("vra_project.my-project"),
					resource.TestCheckResourceAttr(
						"vra_project.my-project", "description", "update test project"),
					resource.TestCheckResourceAttr(
						"vra_project.my-project", "zone_assignments.priority", "1"),
					resource.TestCheckResourceAttr(
						"vra_project.my-project", "zone_assignments.max_instances", "2"),
				),
			},
		},
	})
}

func testAccCheckVRAProjectExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no project ID is set")
		}

		return nil
	}
}

func testAccCheckVRAProjectDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "vra_cloud_account_aws" {
			_, err := apiClient.CloudAccount.GetAwsCloudAccount(cloud_account.NewGetAwsCloudAccountParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_cloud_account_aws' still exists with id %s", rs.Primary.ID)
			}
		}
		if rs.Type == "vra_project" {
			_, err := apiClient.Project.UpdateProject(project.NewUpdateProjectParams().WithID(rs.Primary.ID).WithBody(&models.ProjectSpecification{
				ZoneAssignmentConfigurations: []*models.ZoneAssignmentConfig{},
			}))
			if err != nil {
				return err
			}

			_, err = apiClient.Project.DeleteProject(project.NewDeleteProjectParams().WithID(rs.Primary.ID))
			if err != nil {
				return fmt.Errorf("Resource 'vra_project' still exists with id %s", rs.Primary.ID)
			}
		}
		if rs.Type == "vra_zone" {
			_, err := apiClient.Location.GetZone(location.NewGetZoneParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_zone' still exists with id %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckVRAProjectConfig(rInt int) string {
	id := os.Getenv("VRA_AWS_ACCESS_KEY_ID")
	secret := os.Getenv("VRA_AWS_SECRET_ACCESS_KEY")
	return fmt.Sprintf(`
resource "vra_cloud_account_aws" "my-cloud-account" {
	name = "my-cloud-account"
	description = "test cloud account"
	access_key = "%s"
	secret_key = "%s"
	regions = ["us-east-1"]
 }

 resource "vra_zone" "my-zone" {
    name = "my-vra-zone-%d"
    description = "description my-vra-zone"
    region_id = "${element(vra_cloud_account_aws.my-cloud-account.0.region_ids, 0)}"
    tags {
        key = "mykey"
        value = "myvalue"
    }
    tags {
        key = "foo"
        value = "bar"
    }
    tags {
        key = "faz"
        value = "baz"
    }
}
resource "vra_project" "my-project" {
	name = "my-project-%d"
	description = "test project"
 }`, id, secret, rInt, rInt)
}

func testAccCheckVRAProjectUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "vra_project" "my-project" {
		name = "my-project-%d"
		description = "update test project"
		zone_assignments {
			zone_id       = vra_zone.my-zone.id
			priority      = 1
			max_instances = 2
		  }
	 }`, rInt)
}
