---
layout: "vra"
page_title: "VMware vRealize Automation: vra_storage_profile_aws"
description: |-
  Provides a data lookup for vra_storage_profile_aws.
---

# Resource: vra_storage_profile_aws
## Example Usages
This is an example of how to create a storage profile aws resource.

**Vra storage profile aws:**

```hcl
# AWS storage profile using vra_storage_profile_aws resource and EBS disk type.
resource "vra_storage_profile_aws" "this" {
  name                = "aws-with-instance-store-1"
  description         = "AWS Storage Profile with instance store device type."
  region_id           = data.vra_region.this.id
  default_item        = false
  supports_encryption = false

  device_type         = "ebs"

  // Volume Types: gp2 - General Purpose SSD, io1 - Provisioned IOPS SSD, sc1 - Cold HDD, ST1 - Throughput Optimized HDD, standard - Magnetic
  volume_type         = "io1"  // Supported values: gp2, io1, sc1, st1, standard.
  iops                = "1000" // Required only when volumeType is io1.

  tags {
    key   = "foo"
    value = "bar"
  }
}

# AWS storage profile using vra_storage_profile_aws resource and instance store disk type.
resource "vra_storage_profile_aws" "this" {
  name         = "aws-with-instance-store-1"
  description  = "AWS Storage Profile with instance store device type."
  region_id    = data.vra_region.this.id
  default_item = false

  device_type  = "instance-store"

  tags {
    key   = "foo"
    value = "bar"
  }
}

```

A storage profile aws resource supports the following arguments:

## Argument Reference

* `default_item` - (Required) Indicates if this storage profile is a default profile.

* `description` - (Optional) A human-friendly description.

* `device_type` - (Optional) Indicates the type of storage device.

* `iops` - (Optional) Indicates maximum I/O operations per second in range(1-20,000).

* `name` - (Required) A human-friendly name used as an identifier in APIs that support this option.

* `region_id` - (Required) A link to the region that is associated with the storage profile.

* `supports_encryption` - (Optional) Indicates whether this storage profile supports encryption or not.

* `tags` - (Optional) A set of tag keys and optional values that were set on this Network Profile.
           example: `[ { "key" : "ownedBy", "value": "Rainpole" } ]`

* `volume_type` - (Optional) Indicates the type of volume associated with type of storage.

## Attributes Reference

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `external_region_id` - The id of the region as seen in the cloud provider for which this profile is defined.

* `links` - HATEOAS of the entity

* `organization_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
