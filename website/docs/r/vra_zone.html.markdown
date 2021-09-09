---
layout: "vra"
page_title: "VMware vRealize Automation: vra_zone"
description: |-
  Provides a VMware vRA vra_zone resource.
---

# Resource: vra_zone

## Example Usages

This is an example of how to create a zone resource.

```hcl
resource "vra_zone" "this" {
  name        = "tf-vra-zone"
  description = "my terraform test cloud zone"
  region_id   = data.vra_region.this.id

  tags {
    key   = "my-tf-key"
    value = "my-tf-value"
  }

  tags {
    key   = "tf-foo"
    value = "tf-bar"
  }
}
```

A zone resource supports the following arguments:

## Argument Reference

* `compute_ids` - (Optional) The ids of the compute resources that will be explicitly assigned to this zone.

* `custom_properties` - (Optional) A list of key value pair of properties that will be used.

* `description` - (Optional) A human-friendly description.

* `folder` - (Optional) The folder relative path to the datacenter where resources are deployed to (only applicable for vSphere cloud zones).

* `name` - (Required) A human-friendly name used as an identifier for the zone resource instance.

* `placement_policy` - (Optional) The placement policy for the zone. One of `DEFAULT`, `SPREAD` or `BINPACK`. Default is `DEFAULT`.

* `region_id` - (Required) The id of the region for which this zone is created.

* `tags` - (Optional) A set of tag keys and optional values that were set on this resource:
  * `key` - Tag’s key.
  * `value` - Tag’s value.

* `tags_to_match` - (Optional) A set of tag keys and optional values for compute resource filtering:
  * `key` - Tag’s key.
  * `value` - Tag’s value.


## Attribute Reference

* `cloud_account_id` - The ID of the cloud account this zone belongs to.

* `created_at` - Date when the entity was created. The date is in ISO 8601 and UTC.

* `external_region_id` - The id of the region for which this zone is defined.

* `links` - HATEOAS of entity.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
