---
layout: "vra"
page_title: "VMware vRealize Automation: vra_fabric_datastore_vsphere"
description: |-
  Provides a data lookup for vSphere fabric datastores.
---

# Data Source: vra_fabric_datastore_vsphere

## Example Usages
This is an example of how to lookup vSphere fabric datastores.

**vSphere fabric datastore data source by Id:**

```hcl
# Lookup vSphere fabric datastore using its id
data "vra_fabric_datastore_vsphere" "this" {
  id = var.fabric_datastore_vsphere_id
}
```

**vSphere fabric datastore data source by filter query:**

```hcl
# Lookup vSphere fabric datastore using its name
data "vra_fabric_datastore_vsphere" "this" {
  filter = "name eq '${var.datastore_name}'"
}
```

A vSphere fabric datastore data source supports the following arguments:

## Argument Reference

* `id` - (Optional) The id of the vSphere fabric datastore resource instance. Only one of 'id' or 'filter' must be specified.

* `filter` - (Optional) Search criteria to narrow down the vSphere fabric datastore resource instance. Only one of 'id' or 'filter' must be specified.

## Attribute Reference

* `cloud_account_ids` - Set of ids of the cloud accounts this entity belongs to.

* `created_at` - Date when the entity was created. The date is in ISO 8601 and UTC.

* `description` - A human-friendly description.

* `external_id` - External entity Id on the provider side.

* `external_region_id` - Id of datacenter in which the datastore is present.

* `free_size_gb` - Indicates free size available in datastore.

* `links` - HATEOAS of the entity

* `name` - A human-friendly name used as an identifier for the vSphere fabric datastore resource instance.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `tags` -  A set of tag keys and optional values that were set on this resource:
  * `key` - Tag’s key.
  * `value` - Tag’s value.

* `type` - Type of datastore.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
