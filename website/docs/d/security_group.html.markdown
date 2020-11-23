---
layout: "vra"
page_title: "VMware vRealize Automation: security_group"
description: |-
  Provides a data lookup for security groups.
---

# Data Source: security_group
## Example Usages
This is an example of how to lookup security groups.

**Security groups by Id:**

```hcl
# Lookup Security groups using its Id
data "security_group" "this" {
  id = var.security_group_id
}
```

**Security groups by filter query:**

```hcl
# Lookup Security groups using its name
data "security_group" "this" {
  filter = "name eq '${var.name}'"
}
```

A Security group supports the following arguments:

## Argument Reference
* `filter` - Search criteria to narrow down the Security groups. Only one of 'filter' or 'id' must be specified.

* `id` - The id of the security group. Only one of 'filter' or 'id' must be specified.

## Attribute Reference

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `description` - A human-friendly description of the security groups.

* `external_id` - External entity Id on the provider side.

* `external_region_id` - The id of the region for which this entity is defined.

* `links` - HATEOAS of the entity

* `name` - A human-friendly name used as an identifier in APIs that support this option.  

* `organization_id` - The id of the organization this entity belongs to.

* `rules` - List of security rules.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.