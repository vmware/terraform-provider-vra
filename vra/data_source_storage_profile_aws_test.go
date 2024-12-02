// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"fmt"
	"os"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"testing"
)

func TestAccDataSourceStorageProfileAws(t *testing.T) {
	resourceName1 := "vra_storage_profile_aws.this"
	dataSourceName1 := "data.vra_storage_profile_aws.this"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckStorageProfileAws(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAStorageProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceVRAStorageProfileAwsNotFound(),
				ExpectError: regexp.MustCompile("vra_storage_profile_aws filter did not match any storage profile"),
			},
			{
				Config: testAccDataSourceVRAStorageProfileAwsByID(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName1, "id", dataSourceName1, "id"),
					resource.TestCheckResourceAttrPair(resourceName1, "description", dataSourceName1, "description"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName1, "name"),
					resource.TestCheckResourceAttrPair(resourceName1, "default_item", dataSourceName1, "default_item"),
					resource.TestCheckResourceAttrPair(resourceName1, "external_region_id", dataSourceName1, "external_region_id"),
				),
			},
		},
	})
}

func testAccDataSourceVRAStorageProfileAwsNotFound() string {
	return `
	data "vra_storage_profile_aws" "this" {
		filter = "externalRegionId eq 'foobar'"
	}`
}

func testAccDataSourceVRAStorageProfileAwsByID() string {
	return testAccDataSourceVRAStorageProfileAWSConfig() + `
	data "vra_storage_profile_aws" "this" {
		id = vra_storage_profile_aws.this.id
	}`
}

func testAccDataSourceVRAStorageProfileAWSConfig() string {
	awsCloudAccountName := os.Getenv("VRA_AWS_CLOUD_ACCOUNT_NAME")
	awsRegionName := os.Getenv("VRA_AWS_REGION_NAME")
	return fmt.Sprintf(`
	data "vra_cloud_account_aws" "this" {
  		name = "%s"
	}

	data "vra_region" "this" {
  		cloud_account_id = data.vra_cloud_account_aws.this.id
  		region           = "%s"
	}

	resource "vra_storage_profile_aws" "this" {
  		name                = "aws-with-instance-store-1"
  		description         = "AWS Storage Profile with instance store device type."
  		region_id           = data.vra_region.this.id
  		default_item        = false
  		supports_encryption = false

  		device_type = "ebs"

  		// Volume Types: gp2 - General Purpose SSD, io1 - Provisioned IOPS SSD, sc1 - Cold HDD, ST1 - Throughput Optimized HDD, standard - Magnetic
  		volume_type = "io1"  // Supported values: gp2, io1, sc1, st1, standard.
  		iops       = "1000" // Required only when volumeType is io1.

		tags {
			key   = "foo"
			value = "bar"
		}
	}`, awsCloudAccountName, awsRegionName)
}
