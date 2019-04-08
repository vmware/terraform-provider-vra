package cas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceCASCloudAccountAWS(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName1 := "cas_cloud_account_aws.my-cloud-account"
	dataSourceName1 := "data.cas_cloud_account_aws.test-cloud-account"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckAWS(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCASCloudAccountAWS(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "access_key", dataSourceName1, "access_key"),
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
				),
			},
		},
	})
}

func testAccDataSourceCASCloudAccountAWS(rInt int) string {
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
	 
	 data "cas_cloud_account_aws" "test-cloud-account" {
     name = "${cas_cloud_account_aws.my-cloud-account.name}"
	 }`, rInt, id, secret)
}
