---
layout: "vra"
page_title: "VMware vRealize Automation: vra_fabric_datastore_vsphere"
description: |-
  Updates a fabric_datastore_vsphere resource.
---

# Resource: vra_fabric_datastore_vsphere

Updates a VMware vRealize Automation fabric_datastore_vsphere resource.

## Example Usages

You cannot create a fabric datastore vSphere resource, however you can import it using the command specified in the import section below.

Once a resource is imported, you can update it as shown below:

```hcl
resource "vra_fabric_datastore_vsphere" "this" {
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

## Import

To import the fabric datastore vSphere resource, use the ID as in the following example:

`$ terraform import vra_fabric_datastore_vsphere.this 8e0c9a4c-3ab8-48e8-b9d5-0751c871e282
