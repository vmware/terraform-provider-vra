---
page_title: "VMware Aria Automation: vra_cloud_account_vsphere"
description: |-
    Provides a data lookup for vra_cloud_account_vsphere.
---

# Data Source: vra_cloud_account_vsphere

Provides a vra_cloud_account_vsphere data source.

## Example Usages

**vSphere cloud account data source by its id:**

This is an example of how to read the cloud account data source using its id.

```hcl

data "vra_cloud_account_vsphere" "this" {
  id = var.vra_cloud_account_vsphere_id
}
```

**vSphere cloud account data source by its name:**

This is an example of how to read the cloud account data source using its name.

```hcl

data "vra_cloud_account_vsphere" "this" {
  name = var.vra_cloud_account_vsphere_name
}
```

## Argument Reference

The following arguments are supported for an vSphere cloud account data source:

* `id` - (Optional) The id of this vSphere cloud account.

* `name` - (Optional) The name of this vSphere cloud account.

## Attribute Reference

* `associated_cloud_account_ids` - Cloud accounts associated with this cloud account.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `custom_properties` - Additional properties that may be used to extend the base info.

* `dc_id` - Identifier of a data collector vm deployed in the on premise infrastructure.

* `description` - A human-friendly description.

* `enabled_regions` - A set of regions that are enabled for this cloud account.

  * `external_region_id` - Unique identifier of the region on the provider side. Example: `Datacenter:datacenter-2`

  * `id` - Unique identifier of the region. Example: `6196725b-6f5d-497c-8777-abe7814ee5ec`

  * `name` - Name of the region on the provider side. Example: `vcfcons-mgmt-vc01`

* `hostname` - The IP address or FQDN of the vCenter Server. The cloud proxy belongs on this vCenter.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `regions` - A set of region IDs that are enabled for this account.

* `tags` - A set of tag keys and optional values that were set on this resource. Example: `[ { "key" : "vmware", "value": "provider" } ]`

  * `key` - Tag’s key.

  * `value` - Tag’s value.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `username` - The vSphere username to authenticate the vsphere account.
