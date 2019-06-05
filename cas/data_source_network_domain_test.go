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
)

func TestAccDataSourceCASNetworkDomainBasic(t *testing.T) {
	rInt := acctest.RandInt()
	dataSourceName := "data.cas_network_domain.my-network-domain"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCASRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceCASNetworkDomainNoneConfig(rInt),
				ExpectError: regexp.MustCompile("cas_network_domain filter did not match any network domain"),
			},
			{
				Config: testAccDataSourceCASNetworkDomainOneConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "external_region_id", "us-east-1"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "rainpole-dev"),
				),
			},
		},
	})
}

func testAccDataSourceCASNetworkDomainBaseConfig(rInt int) string {
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

func testAccDataSourceCASNetworkDomainNoneConfig(rInt int) string {
	return testAccDataSourceCASNetworkDomainBaseConfig(rInt) + `
		data "cas_network_domain" "my-network-domain" {
			filter = "name eq 'foobar rainpole-dev'"
		}`
}

func testAccDataSourceCASNetworkDomainOneConfig(rInt int) string {
	return testAccDataSourceCASNetworkDomainBaseConfig(rInt) + `
	data "cas_network_domain" "my-network-domain" {
			filter = "name eq 'rainpole-dev'"
		}`
}

func testAccCheckCASNetworkDomainDestroy(s *terraform.State) error {
	apiClient := testAccProviderCAS.Meta().(*Client).apiClient

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
