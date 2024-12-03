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

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vra-sdk-go/pkg/client/content_source"
)

func TestAccVRAContentSource_Valid(t *testing.T) {
	rInt := acctest.RandInt()
	resource1 := "vra_content_source.this"
	project := "vra_project.this"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckContentSource(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAContentSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRAContentSourceValidConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAContentSourceExists(resource1),
					resource.TestMatchResourceAttr(resource1, "name", regexp.MustCompile("^test-cs-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttrPair(resource1, "project_id", project, "id"),
					resource.TestCheckResourceAttr(resource1, "description", "terraform test content_source"),
					resource.TestCheckResourceAttr(resource1, "type_id", "com.gitlab"),
				),
			},
		},
	})
}

func TestAccVRAContentSource_Invalid(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckContentSource(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAContentSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckVRAContentSourceInvalidConfigContentType(rInt),
				ExpectError: regexp.MustCompile("content_type to be one of"),
			},
			{
				Config:      testAccCheckVRAContentSourceInvalidTypeID(rInt),
				ExpectError: regexp.MustCompile("expected type_id to be one of"),
			},
		},
	})
}

func testAccCheckVRAContentSourceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no content_source ID is set")
		}

		return nil
	}
}

func testAccCheckVRAContentSourceDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vra_content_source" {
			continue
		}

		_, err := apiClient.ContentSource.GetContentSourceUsingGET(content_source.NewGetContentSourceUsingGETParams().WithID(strfmt.UUID(rs.Primary.ID)))

		if err == nil {
			return fmt.Errorf("resource 'vra_content_source' still exists with id %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckVRAContentSourceValidConfig(rInt int) string {
	integrationID := os.Getenv("VRA_INTEGRATION_ID")
	repoFolder := os.Getenv("VRA_CONTENT_SOURCE_PATH")
	repoBranch := os.Getenv("VRA_CONTENT_SOURCE_BRANCH")
	repo := os.Getenv("VRA_CONTENT_SOURCE_REPO")

	return fmt.Sprintf(`
	resource "vra_project" "this" {
		name = "terraform-test-project %d"
	  }

	resource "vra_content_source" "this" {
	  name        = "test-cs-%d"
	  description = "terraform test content_source"
	  project_id = vra_project.this.id
	  type_id = "com.gitlab"
	  sync_enabled = "false"
	  config  {
			path = "%s"
			branch = "%s"
			repository = "%s"
			content_type = "BLUEPRINT"
			project_name = vra_project.this.name
			integration_id = "%s"
			}
	  }`, rInt, rInt, repoFolder, repoBranch, repo, integrationID)
}

func testAccCheckVRAContentSourceInvalidConfigContentType(rInt int) string {
	integrationID := os.Getenv("VRA_INTEGRATION_ID")
	repoFolder := os.Getenv("VRA_CONTENT_SOURCE_PATH")
	repoBranch := os.Getenv("VRA_CONTENT_SOURCE_BRANCH")
	repo := os.Getenv("VRA_CONTENT_SOURCE_REPO")

	return fmt.Sprintf(`
	resource "vra_content_source" "this" {
		name        = "test-cs-%d"
		description = "terraform test content_source"
		project_id = "9704b10f-ffff-aaaa-bbbb-7799029197d3"
		type_id = "com.gitlab"
		sync_enabled = "true"
		config  {
			  path = "%s"
			  branch=  "%s"
			  repository=  "%s"
			  content_type= "PANCAKE"
			  project_name= "some random name"
			  integration_id= "%s"
			  }
		}`, rInt, repoFolder, repoBranch, repo, integrationID)
}

func testAccCheckVRAContentSourceInvalidTypeID(rInt int) string {
	integrationID := os.Getenv("VRA_INTEGRATION_ID")
	repoFolder := os.Getenv("VRA_CONTENT_SOURCE_PATH")
	repoBranch := os.Getenv("VRA_CONTENT_SOURCE_BRANCH")
	repo := os.Getenv("VRA_CONTENT_SOURCE_REPO")

	return fmt.Sprintf(`
	resource "vra_content_source" "this" {
	  name        = "test-cs-%d"
	  description = "terraform test content_source"
	  project_id = "9704b10f-ffff-aaaa-bbbb-7799029197d3"
	  type_id = "com.subversion"
	  sync_enabled = "true"
	  config  {
			path = "%s"
			branch = "%s"
			repository = "%s"
			content_type = "BLUEPRINT"
			project_name = "some random name"
			integration_id = "%s"
			}
	  }`, rInt, repoFolder, repoBranch, repo, integrationID)
}
