---
layout: "vra"
page_title: "VMware vRealize Automation: vra_zone"
description: |-
  Provides a data lookup for vra_zone.
---

# Data Source: vra_zone

## Example Usages

This is an example of how to read a zone data source.

```hcl
data "vra_zone" "test-zone" {
  name = var.zone_name
}
```

A zone data source supports the following arguments:

## Argument Reference

* `id` - (Optional) The id of the zone resource instance.

* `name` - (Optional) A human-friendly name used as an identifier for the zone resource instance.

## Attributes Reference

* `cloud_account_id` - The ID of the cloud account this zone belongs to.

* `compute_ids` - The ids of the compute resources that has been explicitly assigned to this zone.

* `created_at` - Date when the entity was created. The date is in ISO 8601 and UTC.

* `custom_properties` - A list of key value pair of properties that will be used.

* `description` - A human-friendly description.

* `external_region_id` - The id of the region for which this zone is defined.

* `folder` - The folder relative path to the datacenter where resources are deployed to (only applicable for vSphere cloud zones).

* `links` - HATEOAS of the entity.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `placement_policy` - The placement policy for the zone. One of `DEFAULT`, `SPREAD` or `BINPACK`.

* `tags` - A set of tag keys and optional values that were set on this resource:
  * `key` - Tag’s key.
  * `value` - Tag’s value.

* `tags_to_match` - A set of tag keys and optional values for compute resource filtering:
  * `key` - Tag’s key.
  * `value` - Tag’s value.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
