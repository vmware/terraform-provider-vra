---
layout: "vra"
page_title: "VMware vRealize Automation: vra_fabric_compute"
description: |-
  Provides a data lookup for vRA fabric computes.
---

# Data Source: vra_fabric_compute

## Example Usages

This is an example of how to lookup fabric computes.

**Fabric compute data source by Id:**

```hcl
# Lookup fabric compute using its id
data "vra_fabric_compute" "this" {
  id = var.fabric_compute_id
}
```

**Fabric compute data source by filter query:**

```hcl
# Lookup fabric compute using its name
data "vra_fabric_compute" "this" {
  filter = "name eq '${var.fabric_compute_name}'"
}
```

A fabric compute data source supports the following arguments:

## Argument Reference

* `id` - (Optional) The id of the fabric compute resource instance. Only one of 'id' or 'filter' must be specified.

* `filter` - (Optional) Search criteria to narrow down the fabric compute resource instance. Only one of 'id' or 'filter' must be specified.

## Attribute Reference

* `created_at` - Date when the entity was created. The date is in ISO 8601 and UTC.

* `custom_properties` - A list of key value pair of custom properties for the fabric compute resource.

* `description` - A human-friendly description.

* `external_id` - The id of the external entity on the provider side.

* `external_region_id` - The external region id of the fabric compute.

* `external_zone_id` - The external zone id of the fabric compute.

* `lifecycle_state` - Lifecycle status of the compute instance.

* `links` - HATEOAS of the entity.

* `name` - A human-friendly name used as an identifier for the fabric compute resource instance.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `power_state` - Power state of fabric compute instance.

* `tags` -  A set of tag keys and optional values that were set on this resource:
  * `key` - Tag’s key.
  * `value` - Tag’s value.

* `type` - Type of the fabric compute instance.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
