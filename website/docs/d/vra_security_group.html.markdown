---
layout: "vra"
page_title: "VMware vRealize Automation: vra_security_group"
description: |-
  Provides a data lookup for security groups.
---

# Data Source: vra_security_group
## Example Usages
This is an example of how to lookup security groups.

**Security groups by filter query:**

```hcl
# Lookup Security groups using its name
data "vra_security_group" "this" {
  filter = "name eq '${var.name}'"
}
```

A Security group supports the following arguments:

## Argument Reference
* `filter` - (Required) Search criteria to narrow down the Security groups. 

## Attribute Reference

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `description` - A human-friendly description of the security groups.

* `external_id` - External entity Id on the provider side.

* `external_region_id` - The id of the region for which this entity is defined.

* `id` - ID of the security group.

* `links` - HATEOAS of the entity

* `name` - Name of the security group.

* `organization_id` - ID of organization that entity belongs to.

* `rules` - List of security rules.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.