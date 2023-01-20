---layout: "vra"
page_title: "VMware vRealize Automation: vra_machine"
description: |-
  Provides a data lookup for vra_machine.
---

# Data Source: vra_machine
## Example Usages

This is an example of how to read a machine data source.

```hcl

data "vra_machine" "this" {
  id = var.my_machine_id
}

```

**Machine data source filter by name:**
```hcl

data "vra_machine" "this" {
  filter = "name eq '${var.machine_name}'"
}

```
## Argument Reference
* `description` - (Optional) A human-friendly description.

* `filter` - (Optional) Filter query string that is supported by vRA multi-cloud IaaS API. Example: `regionId eq '<regionId>' and cloudAccountId eq '<cloudAccountId>'`.

* `id` - (Optional) The id of the image profile instance.

## Attribute Reference

* `address` - Primary address allocated or in use by this machine. The actual type of the address depends on the adapter type. Typically it is either the public or the external IP address.

* `cloud_account_ids` - Set of ids of the cloud accounts this resource belongs to.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `custom_properties` - Additional properties that may be used to extend the base resource.

* `deployment_id` - Deployment id that is associated with this resource.

* `external_id` - External entity Id on the provider side.

* `external_region_id` - The external regionId of the resource.

* `external_zone_id` - The external zoneId of the resource.

* `links` - HATEOAS of the entity

* `name` - A human-friendly name used as an identifier in APIs that support this option.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `power_state` - Power state of machine.

* `project_id` - The id of the project this resource belongs to.

* `tags` - A set of tag keys and optional values that were set on this resource.
           example: `[ { "key" : "ownedBy", "value": "Rainpole" } ]`

* `update_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
