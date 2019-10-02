package vra

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/vra-sdk-go/pkg/client/image_profile"
)

func TestAccVRAImageProfileBasic(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckMachine(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAImageProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRAImageProfileConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAImageProfileExists("vra_image_profile.my-image-profile"),
					resource.TestCheckResourceAttr(
						"vra_image_profile.my-image-profile", "name", "my-image-profile-"+strconv.Itoa(rInt)),
					resource.TestCheckResourceAttr(
						"vra_image_profile.my-image-profile", "description", "my image profile"),
					resource.TestCheckResourceAttr(
						"vra_image_profile.my-image-profile", "image_mapping.#", "1"),
				),
			},
			{
				Config: testAccCheckVRAImageProfileUpdateConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAImageProfileExists("vra_image_profile.my-image-profile"),
					resource.TestCheckResourceAttr(
						"vra_image_profile.my-image-profile", "name", "my-image-profile-"+strconv.Itoa(rInt)),
					resource.TestCheckResourceAttr(
						"vra_image_profile.my-image-profile", "description", "my image profile update"),
					resource.TestCheckResourceAttr(
						"vra_image_profile.my-image-profile", "image_mapping.#", "1"),
				),
			},
		},
	})
}

func testAccCheckVRAImageProfileExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no image profile ID is set")
		}

		return nil
	}
}

func testAccCheckVRAImageProfileDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "vra_image_profile" {
			_, err := apiClient.ImageProfile.GetImageProfile(image_profile.NewGetImageProfileParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_image_profile' still exists with id %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckVRAImageProfileConfig(rInt int) string {
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

resource "vra_image_profile" "my-image-profile" {
	name = "my-image-profile-%d"
    description = "my image profile"
    region_id = data.vra_region.us-east-1-region.id
    image_mapping {
        name = "ubuntu"
        image_name = "%s"
    }
}`, name, rInt, image)
}

func testAccCheckVRAImageProfileUpdateConfig(rInt int) string {
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
resource "vra_image_profile" "my-image-profile" {
	name = "my-image-profile-%d"
    description = "my image profile update"
    region_id = data.vra_region.us-east-1-region.id
    image_mapping {
        name = "ubuntu"
        image_name = "%s"
    }
}`, name, rInt, image)
}
