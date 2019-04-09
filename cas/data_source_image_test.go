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

func TestAccDataSourceCASImageBasic(t *testing.T) {
	rInt := acctest.RandInt()
	dataSourceName := "data.cas_image.ubuntu"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCASRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceCASImageNoneConfig(rInt),
				ExpectError: regexp.MustCompile("cas_image filter did not match any images"),
			},
			{
				Config:      testAccDataSourceCASImageManyConfig(rInt),
				ExpectError: regexp.MustCompile("cas_image must filter to a single image"),
			},
			{
				Config: testAccDataSourceCASImageOneConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "description", "Canonical, Ubuntu, 16.04 LTS, amd64 xenial image build on 2017-10-26"),
					resource.TestCheckResourceAttr(dataSourceName, "external_id", "ami-da05a4a0"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "ubuntu/images/hvm-ssd/ubuntu-xenial-16.04-amd64-server-20171026.1"),
					resource.TestCheckResourceAttr(dataSourceName, "private", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "region", "us-east-1"),
				),
			},
		},
	})
}

func testAccDataSourceCASImageBaseConfig(rInt int) string {
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
	`, rInt, id, secret)
}

func testAccDataSourceCASImageNoneConfig(rInt int) string {
	return testAccDataSourceCASImageBaseConfig(rInt) + `
		data "cas_image" "ubuntu" {
			filter = "name eq 'foobar ubuntu/images/hvm-ssd/ubuntu-xenial-16.04-amd64-server-20171026.1'"
		}`
}

func testAccDataSourceCASImageManyConfig(rInt int) string {
	return testAccDataSourceCASImageBaseConfig(rInt) + `
		data "cas_image" "ubuntu" {
			filter = "name eq 'ubuntu/images/hvm-ssd/ubuntu-xenial-16.04-amd64-server-20171026.1'"
		}`
}

func testAccDataSourceCASImageOneConfig(rInt int) string {
	return testAccDataSourceCASImageBaseConfig(rInt) + `
		data "cas_image" "ubuntu" {
			filter = "name eq 'ubuntu/images/hvm-ssd/ubuntu-xenial-16.04-amd64-server-20171026.1' and externalRegionId eq 'us-east-1'"
		}`
}

func testAccCheckCASImageDestroy(s *terraform.State) error {
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
