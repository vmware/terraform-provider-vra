// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/client/location"
	"github.com/vmware/vra-sdk-go/pkg/client/project"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

func TestExpandProjectConstraints(t *testing.T) {
	c1 := map[string]interface{}{"mandatory": true, "expression": "Foo:Bar"}
	c2 := map[string]interface{}{"mandatory": false, "expression": "Env:Test"}

	projectConstraints := make([]interface{}, 0)
	expandedConstraints := expandProjectConstraints(projectConstraints)

	if len(expandedConstraints["extensibility"]) != 0 {
		t.Errorf("error while expanding when there are no project extensibility constraints")
	}
	if len(expandedConstraints["network"]) != 0 {
		t.Errorf("error while expanding when there are no project network constraints")
	}
	if len(expandedConstraints["storage"]) != 0 {
		t.Errorf("error while expanding when there are no project storage constraints")
	}

	constraints := make([]interface{}, 0)
	constraints = append(constraints, c1)
	constraints = append(constraints, c2)

	pc1 := make(map[string]interface{})
	pc1["extensibility"] = schema.NewSet(testSetFunc, constraints)
	pc1["storage"] = schema.NewSet(testSetFunc, constraints)
	pc1["network"] = schema.NewSet(testSetFunc, constraints)

	pc := make([]interface{}, 0)

	pc = append(pc, pc1)

	expandedConstraints = expandProjectConstraints(pc)

	if expandedConstraints == nil {
		t.Errorf("expanded constraints is nil")
	}

	if len(expandedConstraints) != 3 {
		t.Errorf("not all project constraints expanded correctly")
	}

	if len(expandedConstraints["extensibility"]) != 2 || len(expandedConstraints["network"]) != 2 || len(expandedConstraints["storage"]) != 2 {
		t.Errorf("not all extensibility / network / storage constraints expanded correctly")
	}
}

func testSetFunc(_ interface{}) int {
	return rand.Int()
}

func TestFlattenProjectConstraints(t *testing.T) {
	projectConstraints := make(map[string][]models.Constraint)
	flattenedConstraints := flattenProjectConstraints(projectConstraints)

	if len(flattenedConstraints) != 0 {
		t.Errorf("error while flattening when there are no project constraints")
	}

	constraint1 := models.Constraint{Expression: withString("Foo:Bar"), Mandatory: withBool(true)}
	constraint2 := models.Constraint{Expression: withString("Env:Test"), Mandatory: withBool(false)}

	constraints := make([]models.Constraint, 0)
	constraints = append(constraints, constraint1)
	constraints = append(constraints, constraint2)

	projectConstraints["extensibility"] = constraints
	projectConstraints["network"] = constraints
	projectConstraints["storage"] = constraints

	flattenedConstraints = flattenProjectConstraints(projectConstraints)

	if len(flattenedConstraints) != 1 {
		t.Errorf("not all project constraints are flattened correctly")
	}

	fc1 := flattenedConstraints[0]
	if len(fc1["extensibility"].([]interface{})) != 2 {
		t.Errorf("extensibility constraints are not flattened correctly")
	}

	if len(fc1["network"].([]interface{})) != 2 {
		t.Errorf("network constraints are not flattened correctly")
	}

	if len(fc1["storage"].([]interface{})) != 2 {
		t.Errorf("storage constraints are not flattened correctly")
	}
}

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
						"vra_project.my-project", "zone_assignments.#", "1"),
					resource.TestCheckResourceAttr("vra_project.my-project", "constraints.#", "1"),
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
	}

	return nil
}

func testAccCheckVRAProjectUpdateConfig(rInt int) string {
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

 data "vra_region" "us-east-1-region" {
    cloud_account_id = "${vra_cloud_account_aws.my-cloud-account.id}"
    region = "us-east-1"
}

 resource "vra_zone" "my-zone" {
    name 		= "my-vra-zone-%d"
	description = "description my-vra-zone"
	region_id 	= "${data.vra_region.us-east-1-region.id}"

    tags {
        key 	= "mykey"
        value 	= "myvalue"
    }
    tags {
        key 	= "foo"
        value	= "bar"
    }
    tags {
        key 	= "faz"
        value 	= "baz"
    }
}
resource "vra_project" "my-project" {
	name 		= "my-project-%d"
	description = "update test project"

	zone_assignments {
		zone_id       	 = vra_zone.my-zone.id
		priority      	 = 1
		max_instances 	 = 2
		cpu_limit	  	 = 1024
		memory_limit_mb  = 8192
		storage_limit_gb = 65536
	  }

	constraints {
    	extensibility {
      		expression = "foo:bar"
      		mandatory  = false
    	}
    	extensibility {
      		expression = "environment:Test"
      		mandatory  = true
		}

    	network {
      		expression = "foo:bar"
      		mandatory  = false
    	}
    	network {
      		expression = "environment:Test"
      		mandatory  = true
    	}

    	storage {
      		expression = "foo:bar"
      		mandatory  = false
		}
    	storage {
      		expression = "environment:Test"
      		mandatory  = true
    	}
  	}
 }`, id, secret, rInt, rInt)
}

func testAccCheckVRAProjectConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "vra_project" "my-project" {
		name = "my-project-%d"
		description = "test project"
	 }`, rInt)
}
