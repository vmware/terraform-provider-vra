---
layout: "vra"
page_title: "VMware vRealize Automation: fabric_datastore_vsphere"
description: |-
  Provides a data lookup for vSphere fabric datastores.
---

# Data Source: fabric_datastore_vsphere
## Example Usages
This is an example of how to lookup vSphere fabric datastores.

**Fabric vSphere data source by Id:**

```hcl
# Lookup vSphere fabric datastore using its name
data "vra_fabric_datastore_vsphere" "this" {
  id = var.fabric_datastore_vsphere_id
}
```

**Fabric vSphere data store by filter query:**

```hcl
# Lookup vSphere fabric datastore using its name
data "fabric_datastore_vsphere" "this" {
  filter = "name eq '${var.datastore_name}'"
}
```

A storage profile data source supports the following arguments:

## Argument Reference
* `filter` - Filter query string that is supported by vRA multi-cloud IaaS API.  Only one of 'filter' or 'id' must be specified.

* `id` - The id of the vSphere data source.  Only one of 'filter' or 'id' must be specified.

## Attribute Reference
* `cloud_account_ids` - Set of ids of the cloud accounts this entity belongs to.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `external_id` - External entity Id on the provider side.

* `external_region_id` - Id of datacenter in which the datastore is present.

* `free_size_gb` - Indicates free size available in datastore.

* `links` - HATEOAS of the entity

* `name` - A human-friendly name used as an identifier in APIs that support this option. 

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `type` - Type of datastore.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.