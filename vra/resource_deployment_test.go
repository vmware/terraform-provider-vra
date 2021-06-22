package vra

import (
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/vmware/vra-sdk-go/pkg/client/deployments"

	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccVRADeployment_CatalogItem(t *testing.T) {
	rInt := acctest.RandInt()
	resource1 := "vra_deployment.this"
	project := "data.vra_project.this"
	catalogItem := "data.vra_catalog_item.this"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeployment(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRADeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRADeploymentCatalogItemConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRADeploymentExists(resource1),
					resource.TestMatchResourceAttr(resource1, "name", regexp.MustCompile("^test-deployment-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttrPair(resource1, "project_id", project, "id"),
					resource.TestCheckResourceAttrPair(resource1, "catalog_item_id", catalogItem, "id"),
				),
			},
		},
	})
}

func TestAccVRADeployment_Empty(t *testing.T) {
	rInt := acctest.RandInt()
	resource1 := "vra_deployment.this"
	project := "data.vra_project.this"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeployment(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRADeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRADeploymentEmptyConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRADeploymentExists(resource1),
					resource.TestMatchResourceAttr(resource1, "name", regexp.MustCompile("^test-deployment-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttrPair(resource1, "project_id", project, "id"),
				),
			},
			{
				Config:      testAccCheckVRADeploymentDuplicateConfig(rInt),
				ExpectError: regexp.MustCompile(fmt.Sprintf("a deployment with name '%v' exists already. Try with a differnet name", "test-deployment-"+strconv.Itoa(rInt))),
			},
		},
	})
}

func TestAccVRADeployment_Blueprint(t *testing.T) {
	rInt := acctest.RandInt()
	resource1 := "vra_deployment.this"
	project := "data.vra_project.this"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeployment(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRADeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRADeploymentBlueprintConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRADeploymentExists(resource1),
					resource.TestMatchResourceAttr(resource1, "name", regexp.MustCompile("^test-deployment-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttrPair(resource1, "project_id", project, "id"),
				),
			},
		},
	})
}

func testAccCheckVRADeploymentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no deployment ID is set")
		}

		return nil
	}
}

func testAccCheckVRADeploymentDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vra_deployment" {
			continue
		}

		_, err := apiClient.Deployments.GetDeploymentByIDUsingGET(deployments.NewGetDeploymentByIDUsingGETParams().WithDeploymentID(strfmt.UUID(rs.Primary.ID)))
		if err == nil {
			return fmt.Errorf("resource 'vra_deployment' still exists with id %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckVRADeploymentCatalogItemConfig(rInt int) string {
	// Need valid details since this is creating a real deployment
	catalogItemName := os.Getenv("VRA_CATALOG_ITEM_NAME")
	projectName := os.Getenv("VRA_PROJECT_NAME")

	return fmt.Sprintf(`
	data "vra_project" "this" {
	  name = "%s"
	}

	data "vra_catalog_item" "this" {
	  name = "%s"
	}

	resource "vra_deployment" "this" {
	  name        = "test-deployment-%d"
	  description = "terraform test deployment"

	  catalog_item_id = data.vra_catalog_item.this.id
	  project_id = data.vra_project.this.id
	}`, projectName, catalogItemName, rInt)
}

func testAccCheckVRADeploymentEmptyConfig(rInt int) string {
	// Need valid details since this is creating a real deployment without a catalog item and blueprint
	projectName := os.Getenv("VRA_PROJECT_NAME")

	return fmt.Sprintf(`
	data "vra_project" "this" {
	  name = "%s"
	}

	resource "vra_deployment" "this" {
	  name        = "test-deployment-%d"
	  description = "terraform test deployment"

	  project_id = data.vra_project.this.id
	}`, projectName, rInt)
}

func testAccCheckVRADeploymentDuplicateConfig(rInt int) string {
	// Need valid details since this is creating a real deployment without a catalog item and blueprint
	return testAccCheckVRADeploymentEmptyConfig(rInt) + fmt.Sprintf(`
	resource "vra_deployment" "duplicate" {
	  name        = "test-deployment-%d"
	  description = "terraform test deployment"

	  project_id = data.vra_project.this.id
	}`, rInt)
}

func testAccCheckVRADeploymentBlueprintConfig(rInt int) string {
	// Need valid details since this is creating a real deployment
	blueprintID := os.Getenv("VRA_BLUEPRINT_ID")
	blueprintVersion := os.Getenv("VRA_BLUEPRINT_VERSION")
	projectName := os.Getenv("VRA_PROJECT_NAME")

	return fmt.Sprintf(`
	data "vra_project" "this" {
	  name = "%s"
	}

	resource "vra_deployment" "this" {
	  name        = "test-deployment-%d"
	  description = "terraform test deployment"

	  blueprint_id = "%s"
	  blueprint_version = "%s"
	  project_id = data.vra_project.this.id
	}`, projectName, rInt, blueprintID, blueprintVersion)
}

func TestAccVRADeployment_test(t *testing.T) {
	rInt := acctest.RandInt()
	resource1 := "vra_deployment.this"
	project := "data.vra_project.this"
	catalogItem := "data.vra_catalog_item.this"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDeployment(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRADeploymentCreateConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resource1, "name", regexp.MustCompile("^test-deployment-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttrPair(resource1, "project_id", project, "id"),
					resource.TestCheckResourceAttrPair(resource1, "catalog_item_id", catalogItem, "id"),
					resource.TestCheckResourceAttrPair(resource1, "inputs.0.cpuCores", catalogItem, "1"),
				),
			},
			{
				Config: testAccCheckVRADeploymentUpdateConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resource1, "name", regexp.MustCompile("^test-deployment-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttrPair(resource1, "project_id", project, "id"),
					resource.TestCheckResourceAttrPair(resource1, "catalog_item_id", catalogItem, "id"),
					resource.TestCheckResourceAttrPair(resource1, "inputs.0.cpuCores", catalogItem, "2"),
				),
			},
		},
	})
}

func testAccCheckVRADeploymentCreateConfig(rInt int) string {
	// Need valid details since this is creating a real deployment
	catalogItemName := os.Getenv("VRA_CATALOG_ITEM_NAME")
	projectName := os.Getenv("VRA_PROJECT_NAME")

	return fmt.Sprintf(`
	data "vra_project" "this" {
	  name = "%s"
	}

	data "vra_catalog_item" "this" {
	  name = "%s"
	}

	resource "vra_deployment" "this" {
	  name        = "test-deployment-%d"
	  description = "terraform test deployment"

	  catalog_item_id 		= data.vra_catalog_item.this.id
	  catalog_item_version 	= 2
	  project_id = data.vra_project.this.id

	  inputs = {
        cpuCores = 1
	  }
	}`, projectName, catalogItemName, rInt)
}

func testAccCheckVRADeploymentUpdateConfig(rInt int) string {
	// Need valid details since this is creating a real deployment
	catalogItemName := os.Getenv("VRA_CATALOG_ITEM_NAME")
	projectName := os.Getenv("VRA_PROJECT_NAME")

	return fmt.Sprintf(`
	data "vra_project" "this" {
	  name = "%s"
	}

	data "vra_catalog_item" "this" {
	  name = "%s"
	}

	resource "vra_deployment" "this" {
	  name        = "test-deployment-%d"
	  description = "terraform test deployment"

	  catalog_item_id 		= data.vra_catalog_item.this.id
	  catalog_item_version 	= 2
	  project_id = data.vra_project.this.id

	  inputs = {
        cpuCores = 2
	  }
	}`, projectName, catalogItemName, rInt)
}
