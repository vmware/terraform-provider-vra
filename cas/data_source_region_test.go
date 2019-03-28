package cas

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/cas-sdk-go/pkg/client/cloud_account"
	tango "github.com/vmware/terraform-provider-cas/sdk"
)

func TestAccDataSourceCASRegionBasic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "cas_cloud_account_aws.my-cloud-account"
	dataSourceName := "data.cas_region.east-region"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCASRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCASRegionConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "id", dataSourceName, "cloud_account_id"),
					resource.TestCheckResourceAttrPair(resourceName, "region_ids.0", dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "region", "us-east-1"),
				),
			},
		},
	})
}

func testAccDataSourceCASRegionConfig(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	id := os.Getenv("CAS_AWS_ACCESS_KEY_ID")
	secret := os.Getenv("CAS_AWS_SECRET_ACCESS_KEY")
	return fmt.Sprintf(`
	resource "cas_cloud_account_aws" "my-cloud-account" {
		name = "my-cloud-account-%d"
		description = "test cloud account"
		access_key = "%s"
		secret_key = "%s"
		regions = ["us-east-1"]
	}

	data "cas_region" "east-region" {
		cloud_account_id = "${cas_cloud_account_aws.my-cloud-account.id}"
		region = "us-east-1"
	}`, rInt, id, secret)
}

func TestAccDataSourceCASRegionInvalidRegion(t *testing.T) {
	rInt := acctest.RandInt()
	dataSourceName := "data.cas_region.east-region"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCASRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceCASRegionInvalidConfig(rInt),
				ExpectError: regexp.MustCompile("region us-west-1 not found in cloud account"),
			},
			{
				Config: testAccDataSourceCASRegionConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "region", "us-east-1"),
				),
			},
		},
	})
}

func testAccDataSourceCASRegionInvalidConfig(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	id := os.Getenv("CAS_AWS_ACCESS_KEY_ID")
	secret := os.Getenv("CAS_AWS_SECRET_ACCESS_KEY")
	return fmt.Sprintf(`
	resource "cas_cloud_account_aws" "my-cloud-account" {
		name = "my-cloud-account-%d"
		description = "test cloud account"
		access_key = "%s"
		secret_key = "%s"
		regions = ["us-east-1"]
	}

	data "cas_region" "east-region" {
		cloud_account_id = "${cas_cloud_account_aws.my-cloud-account.id}"
		region = "us-west-1"
	}`, rInt, id, secret)
}

func testAccCheckCASRegionDestroy(s *terraform.State) error {
	client := testAccProviderCAS.Meta().(*tango.Client)
	apiClient := client.GetAPIClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cas_cloud_account_aws" {
			continue
		}

		_, err := apiClient.CloudAccount.GetAwsCloudAccount(cloud_account.NewGetAwsCloudAccountParams().WithID(rs.Primary.ID))
		if err == nil {
			return err
		}
	}

	return nil
}
