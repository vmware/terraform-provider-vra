package vra

import (
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/vmware/vra-sdk-go/pkg/client/deployments"

	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccVRABlueprint_Valid(t *testing.T) {
	rInt := acctest.RandInt()
	resource1 := "vra_blueprint.this"
	project := "vra_project.this"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckBlueprint(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRABlueprintDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRABlueprintValidContentConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resource1, "name", regexp.MustCompile("^test-blueprint-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttrPair(resource1, "project_id", project, "id"),
					resource.TestCheckResourceAttr(resource1, "description", "terraform test blueprint"),
					resource.TestCheckResourceAttr(resource1, "request_scope_org", "false"),
				),
			},
		},
	})
}

func TestAccVRABlueprint_Invalid(t *testing.T) {
	rInt := acctest.RandInt()
	resource1 := "vra_blueprint.this"
	project := "vra_project.this"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckBlueprint(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRABlueprintDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRABlueprintInvalidContentConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRADeploymentExists(resource1),
					resource.TestMatchResourceAttr(resource1, "name", regexp.MustCompile("^test-blueprint-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttrPair(resource1, "project_id", project, "id"),
					resource.TestCheckResourceAttr(resource1, "description", "terraform test blueprint"),
					resource.TestCheckResourceAttr(resource1, "request_scope_org", "false"),
				),
			},
		},
	})
}

func testAccCheckVRABlueprintDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vra_blueprint" {
			continue
		}

		_, err := apiClient.Deployments.GetDeploymentByIDUsingGET(deployments.NewGetDeploymentByIDUsingGETParams().WithDeploymentID(strfmt.UUID(rs.Primary.ID)))
		if err == nil {
			return fmt.Errorf("resource 'vra_blueprint' still exists with id %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckVRABlueprintValidContentConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "vra_project" "this" {
	  name = "terraform-test-project"
	}

	resource "vra_blueprint" "this" {
	  name        = "test-blueprint-%d"
	  description = "terraform test blueprint"

	  project_id = vra_project.this.id

	  content = <<-EOT
				 formatVersion: 1
				 inputs: {}
				 resources:  {}
				EOT
	}`, rInt)
}

func testAccCheckVRABlueprintInvalidContentConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "vra_project" "this" {
	  name = "terraform-test-project"
	}

	resource "vra_blueprint" "this" {
	  name        = "test-blueprint-%d"
	  description = "terraform test blueprint"

	  project_id = vra_project.this.id

	  content = <<-EOT
				 formatVersion: 1
				 inputs: {}
				 resources:
				   Cloud_Machine_1:
    				 type: Cloud.Machine
					 properties:
      				   image: 'ubuntu'
					   flavor: ''"
				EOT
	}`, rInt)
}
