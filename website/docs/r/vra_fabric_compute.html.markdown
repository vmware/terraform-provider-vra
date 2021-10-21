---
layout: "vra"
page_title: "VMware vRealize Automation: vra_fabric_compute"
description: |-
  Updates a fabric_compute resource.
---

# Resource: vra_fabric_compute

Updates a VMware vRealize Automation fabric_compute resource.

## Example Usages

You cannot create a fabric compute resource, however you can import it using the command specified in the import section below.

Once a resource is imported, you can update it as shown below:

```hcl
resource "vra_fabric_compute" "this" {
  tags {
    key   = "foo"
    value = "bar"
  }
}
```
## Argument Reference

* `tags` -  A set of tag keys and optional values that were set on this resource:
  * `key` - Tag’s key.
  * `value` - Tag’s value.

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

* `type` - Type of the fabric compute instance.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

## Import

To import the fabric compute resource, use the ID as in the following example:

`$ terraform import vra_fabric_compute.this 88fdea8b-92ed-4aa9-b6ee-4670412961b0
