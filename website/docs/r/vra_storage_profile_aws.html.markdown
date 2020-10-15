---
layout: "vra"
page_title: "VMware vRealize Automation: vra_storage_profile_aws"
sidebar_current: "docs-vra-resource-vra-storage-profile_aws"
description: |-
  Provides a data lookup for vra_storage_profile_aws.
---

# Resource: vra_storage_profile_aws
## Example Usages
This is an example of how to create a storage profile aws resource.

**Vra storage profile aws:**

```hcl
# AWS storage profile using generic vra_storage_profile resource. Use 'vra_storage_profile_aws' resource as an alternative.
resource "vra_storage_profile" "this" {
  name         = "aws-with-instance-store"
  description  = "AWS Storage Profile with instance store device type."
  region_id    = data.vra_region.this.id
  default_item = false

  disk_properties = {
    deviceType = "instance-store"
  }

  tags {
    key   = "foo"
    value = "bar"
  }
}
```

A storage profile aws resource supports the following arguments:

## Required arguments

* `default_item` - Indicates if this storage profile is a default profile.

* `name` - A human-friendly name used as an identifier in APIs that support this option.  Only one of 'filter', 'id', 'name' or 'region_id' must be specified.

* `region_id` - A link to the region that is associated with the storage profile.

## Optional arguments

* `description` - A human-friendly description.

* `device_type` - Indicates the type of storage device.

* `iops` -  Indicates maximum I/O operations per second in range(1-20,000).

* `supports_encryption` - Indicates whether this storage profile supports encryption or not.

* `tags` - A set of tag keys and optional values that were set on this Network Profile.
           example:[ { "key" : "ownedBy", "value": "Rainpole" } ]

* `volume_type` - Indicates the type of volume associated with type of storage.

## Imported attributes
* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `external_region_id` - The id of the region as seen in the cloud provider for which this profile is defined.

* `links` - HATEOAS of the entity

* `organization_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
