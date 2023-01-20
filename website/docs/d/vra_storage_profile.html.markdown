---
layout: "vra"
page_title: "VMware vRealize Automation: vra_storage_profile"
description: |-
  Provides a data lookup for vra_storage_profile.
---

# Data Source: vra_storage_profile
## Example Usages
This is an example of how to create a storage profile data source.

**Storage profile data source by its id:**

```hcl
data "vra_storage_profile" "this" {
  id = vra_storage_profile.this.id
}
```

**Vra storage profile data source filter by external region id:**

```hcl
data "vra_storage_profile" "this" {
  filter = "externalRegionId eq 'foobar'"
}
```

A storage profile data source supports the following arguments:

## Argument Reference
* `description` - (Optional) A human-friendly description.

* `disk_properties` - (Optional) Map of storage properties that are to be applied on disk while provisioning.

* `filter` - (Optional) Filter query string that is supported by vRA multi-cloud IaaS API. Example: `regionId eq '<regionId>' and cloudAccountId eq '<cloudAccountId>'`.

* `id` - (Optional) The id of the image profile instance.

## Attributes Reference
* `cloud_account_id` - Id of the cloud account this storage profile belongs to.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `external_region_id` - The id of the region as seen in the cloud provider for which this profile is defined.

* `default_item` - Indicates if this storage profile is a default profile.

* `links` - HATEOAS of the entity

* `name` - A human-friendly name used as an identifier in APIs that support this option.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `supports_encryption` - Indicates whether this storage profile supports encryption or not.

* `tags` - A set of tag keys and optional values that were set on this Network Profile.
           example: `[ { "key" : "ownedBy", "value": "Rainpole" } ]`

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
