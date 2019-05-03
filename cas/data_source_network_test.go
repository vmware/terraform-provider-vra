package cas

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceCASNetwork(t *testing.T) {
	rInt := acctest.RandInt()
	dataSourceName1 := "data.cas_network.test-network"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCas(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceCASNetworkNoneConfig(rInt),
				ExpectError: regexp.MustCompile("network invalid-name not found"),
			},
			{
				Config: testAccDataSourceCASNetworkOneConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName1, "id", "6d25dcb5d510875582822c89a1d4"),
				),
			},
		},
	})
}

func testAccDataSourceCASNetworkNoneConfig(rInt int) string {
	return `
	    data "cas_network" "test-network" {
			name = "invalid-name"
		}`
}

func testAccDataSourceCASNetworkOneConfig(rInt int) string {
	return `
		data "cas_network" "test-network" {
			name = "foo1-mcm653-56201379059"
		}`
}
