---
page_title: "VMware Aria Automation: vra_policy_lease"
description: A data source for Lease policies.
---

# Data Source: vra_policy_lease

The following examples shows how to lookup for a lease policy:

**Lease policy data source by its id:**

```hcl
data "vra_policy_lease" "this" {
  id = var.vra_lease_policy_id
}
```

**Lease policy data source by name search:**

```hcl
data "vra_policy_lease" "this" {
  search = var.vra_lease_policy_search_name
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The id of the policy instance.

* `search` - (Optional) Search criteria to narrow down the policy instance.

-> **Note:** One of `id` or `search` must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `created_at` - Date when the entity was created. The date is in ISO 8601 and UTC.

* `created_by` - The user the entity was created by.

* `criteria` - The policy criteria.

* `description` - A human-friendly description for the policy instance.

* `enforcement_type` - The type of enforcement for the policy.

* `last_updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `last_updated_by` - The user the entity was last updated by.

* `lease_grace` - The duration in days that an expired object should be held before it is deleted.

* `lease_term_max` - The maximum duration in days between creation (or renewal) and expiration.

* `lease_total_term_max` - The maximum duration in days between creation and expiration. Unaffected by renewal.

* `name` - A human-friendly name used as an identifier for the policy instance.

* `org_id` - The id of the organization this entity belongs to.

* `project_criteria` - The project based criteria.

* `project_id` - The id of the project this entity belongs to.
