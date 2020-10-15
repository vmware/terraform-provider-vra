---
layout: "vra"
page_title: "VMware vRealize Automation: vra_storage_profile"
sidebar_current: "docs-vra-resource-vra-storage-profile"
description: |-
  Provides a data lookup for vra_storage_profile.
---

# Resource: vra_storage_profile
## Example Usages
This is an example of how to create a storage profile resource.

**Vra storage profile:**

```hcl
# vSphere storage profile using generic vra_storage_profile resource.
resource "vra_storage_profile" "this" {
  name         = "vSphere-standard-independent-non-persistent-disk"
  description  = "vSphere Storage Profile with standard independent non-persistent disk."
  region_id    = data.vra_region.this.id
  default_item = false

  disk_properties = {
    independent      = "true"
    persistent       = "false"
    limitIops        = "2000"
    provisioningType = "eagerZeroedThick" // Supported values: "thin", "thick", "eagerZeroedThick"
    sharesLevel      = "custom"           // Supported values: "low", "normal", "high", "custom"
    shares           = "1500"             // Required only when sharesLevel is "custom".
  }

  disk_target_properties = {
    datastoreId     = data.vra_fabric_datastore_vsphere.this.id
    storagePolicyId = data.vra_fabric_storage_policy_vsphere.this.id // Remove it if datastore default storage policy needs to be selected.
  }

  tags {
    key   = "foo"
    value = "bar"
  }
}
```

A storage profile resource supports the following arguments:

## Optional arguments
* `description` - A human-friendly description.

* `disk_properties` -  Map of storage properties that are to be applied on disk while provisioning.

* `disk_target_properties` - Map of storage placements to know where the disk is provisioned.

* `supports_encryption` - Indicates whether this storage profile supports encryption or not.

## Imported attributes
* `cloud_account_id` - Id of the cloud account this storage profile belongs to.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `external_region_id` - The id of the region as seen in the cloud provider for which this profile is defined.

* `default_item` - Indicates if this storage profile is a default profile.

* `links` - HATEOAS of the entity

* `name` - A human-friendly name used as an identifier in APIs that support this option.  Only one of 'filter', 'id', 'name' or 'region_id' must be specified.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `tags` - A set of tag keys and optional values that were set on this Network Profile.
           example:[ { "key" : "ownedBy", "value": "Rainpole" } ]

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
