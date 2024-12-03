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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vra-sdk-go/pkg/client/storage_profile"
)

func TestAccVRAGenericStorageProfileAWS(t *testing.T) {
	rInt := acctest.RandInt()
	const resourceName = "vra_storage_profile.this"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckStorageProfile(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAStorageProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRAStorageProfileAWSConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAStorageProfileExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile("^my-aws-storage-profile-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(resourceName, "description", "my aws storage profile"),
					resource.TestCheckResourceAttr(resourceName, "default_item", "true"),
					resource.TestCheckResourceAttr(resourceName, "external_region_id", os.Getenv("VRA_AWS_REGION_NAME")),
					resource.TestCheckResourceAttr(resourceName, "disk_properties.deviceType", "instance-store"),
				),
			},
		},
	})
}

func TestAccVRAGenericStorageProfileAWSWithEBS(t *testing.T) {
	rInt := acctest.RandInt()
	const resourceName = "vra_storage_profile.this"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckStorageProfile(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAStorageProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRAStorageProfileAWSConfigWithEBS(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAStorageProfileExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile("^my-aws-storage-profile-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(resourceName, "description", "my aws storage profile"),
					resource.TestCheckResourceAttr(resourceName, "default_item", "true"),
					resource.TestCheckResourceAttr(resourceName, "external_region_id", os.Getenv("VRA_AWS_REGION_NAME")),
					resource.TestCheckResourceAttr(resourceName, "disk_properties.deviceType", "ebs"),
					resource.TestCheckResourceAttr(resourceName, "disk_properties.volumeType", "io1"),
					resource.TestCheckResourceAttr(resourceName, "disk_properties.iops", "1000"),
				),
			},
		},
	})
}

func TestAccVRAGenericStorageProfileAzureWithManagedDisks(t *testing.T) {
	rInt := acctest.RandInt()
	const resourceName = "vra_storage_profile.this"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckStorageProfile(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAStorageProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRAStorageProfileAzureConfigWithManagedDisks(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAStorageProfileExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile("^azure-with-managed-disks-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(resourceName, "description", "Azure Storage Profile with managed disks."),
					resource.TestCheckResourceAttr(resourceName, "default_item", "true"),
					resource.TestCheckResourceAttr(resourceName, "supports_encryption", "true"),
					resource.TestCheckResourceAttr(resourceName, "external_region_id", os.Getenv("VRA_ARM_REGION_NAME")),
					resource.TestCheckResourceAttr(resourceName, "disk_properties.azureDataDiskCaching", "None"),
					resource.TestCheckResourceAttr(resourceName, "disk_properties.azureManagedDiskType", "Standard_LRS"),
					resource.TestCheckResourceAttr(resourceName, "disk_properties.azureOsDiskCaching", "None"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "links.#", "2"), // Links for self, and region
				),
			},
		},
	})
}

func TestAccVRAGenericStorageProfileAzureWithUnmanagedDisks(t *testing.T) {
	rInt := acctest.RandInt()
	const resourceName = "vra_storage_profile.this"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckStorageProfile(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAStorageProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRAStorageProfileAzureConfigWithUnmanagedDisks(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAStorageProfileExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile("^azure-with-unmanaged-disks-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(resourceName, "description", "Azure Storage Profile with unmanaged disks."),
					resource.TestCheckResourceAttr(resourceName, "default_item", "false"),
					resource.TestCheckResourceAttr(resourceName, "supports_encryption", "false"),
					resource.TestCheckResourceAttr(resourceName, "external_region_id", os.Getenv("VRA_ARM_REGION_NAME")),
					resource.TestCheckResourceAttr(resourceName, "disk_properties.azureDataDiskCaching", "None"),
					resource.TestCheckResourceAttr(resourceName, "disk_properties.azureOsDiskCaching", "None"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "links.#", "3"), // Links for self, region, and storage-account
				),
			},
		},
	})
}

func TestAccVRAGenericStorageProfileVSphere(t *testing.T) {
	rInt := acctest.RandInt()
	const resourceName = "vra_storage_profile.this"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckStorageProfile(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAStorageProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRAStorageProfileVSphereConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAStorageProfileExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestMatchResourceAttr(
						resourceName, "name", regexp.MustCompile("^vSphere-standard-independent-non-persistent-disk-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(resourceName, "description", "vSphere Storage Profile with standard independent non-persistent disk."),
					resource.TestCheckResourceAttr(resourceName, "default_item", "false"),
					resource.TestCheckResourceAttr(resourceName, "external_region_id", os.Getenv("VRA_VSPHERE_REGION_NAME")),
					resource.TestCheckResourceAttr(resourceName, "disk_properties.independent", "true"),
					resource.TestCheckResourceAttr(resourceName, "disk_properties.persistent", "false"),
					resource.TestCheckResourceAttr(resourceName, "disk_properties.limitIops", "2000"),
					resource.TestCheckResourceAttr(resourceName, "disk_properties.provisioningType", "thin"),
					resource.TestCheckResourceAttr(resourceName, "disk_properties.sharesLevel", "custom"),
					resource.TestCheckResourceAttr(resourceName, "disk_properties.shares", "1500"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "links.#", "4"), // For self, region, datastore and storage policy.
				),
			},
		},
	})
}

func TestAccVRAGenericStorageProfileVSphereWithFCD(t *testing.T) {
	rInt := acctest.RandInt()
	const resourceName = "vra_storage_profile.this"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckStorageProfile(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAStorageProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRAStorageProfileVSphereFCDConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAStorageProfileExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile("^vSphere-first-class-disk-"+strconv.Itoa(rInt))),
					resource.TestCheckResourceAttr(resourceName, "description", "vSphere Storage Profile with first class disk."),
					resource.TestCheckResourceAttr(resourceName, "default_item", "false"),
					resource.TestCheckResourceAttr(resourceName, "external_region_id", os.Getenv("VRA_VSPHERE_REGION_NAME")),
					resource.TestCheckResourceAttr(resourceName, "disk_properties.diskType", "firstClass"),
					resource.TestCheckResourceAttr(resourceName, "disk_properties.provisioningType", "thick"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "links.#", "4"), // For self, region, datastore and storage policy.
				),
			},
		},
	})
}

func testAccCheckVRAStorageProfileExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no storage profile ID is set")
		}

		return nil
	}
}

func testAccCheckVRAStorageProfileDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "vra_storage_profile" {
			_, err := apiClient.StorageProfile.GetStorageProfile(storage_profile.NewGetStorageProfileParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_storage_profile' still exists with id %s", rs.Primary.ID)
			}
		}
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "vra_storage_profile_aws" {
			_, err := apiClient.StorageProfile.GetStorageProfile(storage_profile.NewGetStorageProfileParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_storage_profile_aws' still exists with id %s", rs.Primary.ID)
			}
		}
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "vra_storage_profile_azure" {
			_, err := apiClient.StorageProfile.GetStorageProfile(storage_profile.NewGetStorageProfileParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_storage_profile_azure' still exists with id %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckVRAStorageProfileAWSBaseConfig() string {
	// Need valid AWS cloud account name and region name from vRA as the test uses them to create storage profile.
	cloudAccountName := os.Getenv("VRA_AWS_CLOUD_ACCOUNT_NAME")
	regionName := os.Getenv("VRA_AWS_REGION_NAME")
	return fmt.Sprintf(`
	data "vra_cloud_account_aws" "this" {
		name = "%s"
	}

	data "vra_region" "this" {
		cloud_account_id = "${data.vra_cloud_account_aws.this.id}"
  		region           = "%s"
	}`, cloudAccountName, regionName)
}

func testAccCheckVRAStorageProfileAWSConfig(rInt int) string {
	return testAccCheckVRAStorageProfileAWSBaseConfig() + fmt.Sprintf(`
	resource "vra_storage_profile" "this" {
		name = "my-aws-storage-profile-%d"
		description = "my aws storage profile"
		region_id = "${data.vra_region.this.id}"
		default_item = true
		disk_properties = {
			deviceType = "instance-store"
		}
	}`, rInt)
}

func testAccCheckVRAStorageProfileAWSConfigWithEBS(rInt int) string {
	return testAccCheckVRAStorageProfileAWSBaseConfig() + fmt.Sprintf(`
	resource "vra_storage_profile" "this" {
		name = "my-aws-storage-profile-%d"
		description = "my aws storage profile"
		region_id = "${data.vra_region.this.id}"
		default_item = true
		disk_properties = {
			deviceType = "ebs"

			// Volume Types: gp2 - General Purpose SSD, io1 - Provisioned IOPS SSD, sc1 - Cold HDD, ST1 - Throughput Optimized HDD, standard - Magnetic
    		volumeType = "io1" // Supported values: gp2, io1, sc1, st1, standard.
			iops       = "1000" // Required only when volumeType is io1.
		}
	}`, rInt)
}

func testAccCheckVRAStorageProfileAzureBaseConfig() string {
	// Need valid Azure cloud account name and region name from vRA as the test uses them to create storage profile.
	cloudAccountName := os.Getenv("VRA_ARM_CLOUD_ACCOUNT_NAME")
	regionName := os.Getenv("VRA_ARM_REGION_NAME")
	return fmt.Sprintf(`
	data "vra_cloud_account_azure" "this" {
  		name = "%s"
	}

	data "vra_region" "this" {
  		cloud_account_id = data.vra_cloud_account_azure.this.id
  		region           = "%s"
	}`, cloudAccountName, regionName)
}

func testAccCheckVRAStorageProfileAzureConfigWithManagedDisks(rInt int) string {
	return testAccCheckVRAStorageProfileAzureBaseConfig() + fmt.Sprintf(`
	resource "vra_storage_profile" "this" {
  		name                = "azure-with-managed-disks-%d"
  		description         = "Azure Storage Profile with managed disks."
  		region_id           = data.vra_region.this.id
  		default_item        = true
  		supports_encryption = true

  		disk_properties = {
    		azureDataDiskCaching = "None"         // Supported Values: None, ReadOnly, ReadWrite
    		azureManagedDiskType = "Standard_LRS" // Supported Values: Standard_LRS, StandardSSD_LRS, Premium_LRS
    		azureOsDiskCaching   = "None"         // Supported Values: None, ReadOnly, ReadWrite
  		}

  		tags {
    		key   = "foo"
    		value = "bar"
  		}
	}`, rInt)
}

func testAccCheckVRAStorageProfileAzureConfigWithUnmanagedDisks(rInt int) string {
	storageAccountID := os.Getenv("VRA_ARM_FABRIC_STORAGE_ACCOUNT_NAME")
	regionName := os.Getenv("VRA_ARM_REGION_NAME")
	return testAccCheckVRAStorageProfileAzureBaseConfig() + fmt.Sprintf(`
	data "vra_fabric_storage_account_azure" "this" {
		filter = "name eq '%s' and externalRegionId eq '%s' and cloudAccountId eq '${data.vra_cloud_account_azure.this.id}'"
	}

	resource "vra_storage_profile" "this" {
  		name                = "azure-with-unmanaged-disks-%d"
  		description         = "Azure Storage Profile with unmanaged disks."
  		region_id           = data.vra_region.this.id
  		default_item        = false
  		supports_encryption = false

  		disk_properties = {
    		azureDataDiskCaching = "None" // Supported Values: None, ReadOnly, ReadWrite
    		azureOsDiskCaching   = "None" // Supported Values: None, ReadOnly, ReadWrite
  		}

  		disk_target_properties = {
    		storageAccountId = data.vra_fabric_storage_account_azure.this.id
  		}

  		tags {
    		key   = "foo"
    		value = "bar"
  		}
	}`, storageAccountID, regionName, rInt)
}

func testAccCheckVRAStorageProfileVSphereBaseConfig() string {
	// Need valid vSphere cloud account name and region name from vRA as the test uses them to create storage profile.
	cloudAccountName := os.Getenv("VRA_VSPHERE_CLOUD_ACCOUNT_NAME")
	regionName := os.Getenv("VRA_VSPHERE_REGION_NAME")
	return fmt.Sprintf(`
	data "vra_cloud_account_vsphere" "this" {
		name = "%s"
	}

	data "vra_region" "this" {
		cloud_account_id = "${data.vra_cloud_account_vsphere.this.id}"
  		region           = "%s"
	}`, cloudAccountName, regionName)
}

func testAccCheckVRAStorageProfileVSphereConfig(rInt int) string {
	regionName := os.Getenv("VRA_VSPHERE_REGION_NAME")
	datastoreName := os.Getenv("VRA_VSPHERE_FABRIC_DATASTORE_NAME")
	storagePolicyName := os.Getenv("VRA_VSPHERE_FABRIC_STORAGE_POLICY_NAME")
	return testAccCheckVRAStorageProfileVSphereBaseConfig() + fmt.Sprintf(`
	data "vra_fabric_datastore_vsphere" "this" {
		filter = "name eq '%s' and externalRegionId eq '%s' and cloudAccountId eq '${data.vra_cloud_account_vsphere.this.id}'"
	}

	data "vra_fabric_storage_policy_vsphere" "this" {
		filter = "name eq '%s' and cloudAccountId eq '${data.vra_cloud_account_vsphere.this.id}'"
	}

	resource "vra_storage_profile" "this" {
  		name         = "vSphere-standard-independent-non-persistent-disk-%d"
  		description  = "vSphere Storage Profile with standard independent non-persistent disk."
  		region_id    = data.vra_region.this.id
  		default_item = false

		disk_properties = {
    		independent      = "true"
    		persistent       = "false"
    		limitIops        = "2000"
    		provisioningType = "thin" 			  // Supported values: "thin", "thick", "eagerZeroedThick"
    		sharesLevel      = "custom"           // Supported values: "low", "normal", "high", "custom"
    		shares           = "1500"             // Required only when sharesLevel is "custom".
  		}

		disk_target_properties = {
    		datastoreId     = data.vra_fabric_datastore_vsphere.this.id
    		storagePolicyId = data.vra_fabric_storage_policy_vsphere.this.id // Remove to select the datastore default storage policy.
  		}

  		tags {
    		key   = "foo"
    		value = "bar"
		}
	}`, datastoreName, regionName, storagePolicyName, rInt)
}

func testAccCheckVRAStorageProfileVSphereFCDConfig(rInt int) string {
	regionName := os.Getenv("VRA_VSPHERE_REGION_NAME")
	datastoreName := os.Getenv("VRA_VSPHERE_FABRIC_DATASTORE_NAME")
	storagePolicyName := os.Getenv("VRA_VSPHERE_FABRIC_STORAGE_POLICY_NAME")
	return testAccCheckVRAStorageProfileVSphereBaseConfig() + fmt.Sprintf(`
	data "vra_fabric_datastore_vsphere" "this" {
		filter = "name eq '%s' and externalRegionId eq '%s' and cloudAccountId eq '${data.vra_cloud_account_vsphere.this.id}'"
	}

	data "vra_fabric_storage_policy_vsphere" "this" {
		filter = "name eq '%s' and cloudAccountId eq '${data.vra_cloud_account_vsphere.this.id}'"
	}

	resource "vra_storage_profile" "this" {
  		name         = "vSphere-first-class-disk-%d"
  		description  = "vSphere Storage Profile with first class disk."
  		region_id    = data.vra_region.this.id
  		default_item = false

  		disk_properties = {
    		diskType         = "firstClass"
    		provisioningType = "thick" // Supported values: "thin", "thick", "eagerZeroedThick"
  		}

		disk_target_properties = {
    		datastoreId     = data.vra_fabric_datastore_vsphere.this.id
    		storagePolicyId = data.vra_fabric_storage_policy_vsphere.this.id // Remove to select the datastore default storage policy.
  		}

  		tags {
    		key   = "foo"
    		value = "bar"
  		}

		tags {
			key   = "environment"
			value = "Test"
		}
	}`, datastoreName, regionName, storagePolicyName, rInt)
}
