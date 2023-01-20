---
layout: "vra"
page_title: "VMware vRealize Automation: vra_storage_profile_aws"
description: |-
  Provides a data lookup for vra_storage_profile_aws.
---

# Data Source: vra_storage_profile_aws
## Example Usages
This is an example of how to create a storage profile aws resource.

**Storage profile aws data source by its id:**

```hcl
data "vra_storage_profile_aws" "this" {
  id = vra_storage_profile_aws.this.id
}
```

**Vra storage profile data source filter by external region id:**

```hcl
data "vra_storage_profile_aws" "this" {
  filter = "externalRegionId eq 'foobar'"
}
```

A storage profile aws data source supports the following arguments:

## Argument Reference

* `filter` - (Optional) Filter query string that is supported by vRA multi-cloud IaaS API. Example: `regionId eq '<regionId>' and cloudAccountId eq '<cloudAccountId>'`.

* `id` - (Optional) The id of the image profile instance.

* `description` - (Optional) A human-friendly description.

* `device_type` - (Optional) Indicates the type of storage device.

* `iops` - (Optional) Indicates maximum I/O operations per second in range(1-20,000).

* `supports_encryption` - (Optional) Indicates whether this storage profile supports encryption or not.

* `tags` - (Optional) A set of tag keys and optional values that were set on this Network Profile.
           example: `[ { "key" : "ownedBy", "value": "Rainpole" } ]`

* `volume_type` - (Optional) Indicates the type of volume associated with type of storage.

## Attributes Reference
* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `default_item` - Indicates if this storage profile is a default profile.

* `external_region_id` - The id of the region as seen in the cloud provider for which this profile is defined.

* `links` - HATEOAS of the entity

* `name` - A human-friendly name used as an identifier in APIs that support this option.

* `organization_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `region_id` - A link to the region that is associated with the storage profile.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
