---
page_title: "VMware Aria Automation: vra_policy_day2_action"
description: A data source for Day2 Action policies.
---

# Data Source: vra_policy_day2_action

The following examples shows how to lookup for a day2 action policy:

**Day2 Action policy data source by its id:**

```hcl
data "vra_policy_day2_action" "this" {
  id = var.vra_day2_action_policy_id
}
```

**Day2 Action policy data source by name search:**

```hcl
data "vra_policy_day2_action" "this" {
  search = var.vra_day2_action_policy_search_name
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The id of the policy instance.

* `search` - (Optional) Search criteria to narrow down the policy instance.

-> **Note:** One of `id` or `search` must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `actions` - List of allowed actions for authority/authorities.

* `authorities` - List of authorities that will be allowed to perform certain actions.

* `created_at` - Date when the entity was created. The date is in ISO 8601 and UTC.

* `created_by` - The user the entity was created by.

* `criteria` - The policy criteria.

* `description` - A human-friendly description for the policy instance.

* `enforcement_type` - The type of enforcement for the policy.

* `last_updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `last_updated_by` - The user the entity was last updated by.

* `name` - A human-friendly name used as an identifier for the policy instance.

* `org_id` - The id of the organization this entity belongs to.

* `project_criteria` - The project based criteria.

* `project_id` - The id of the project this entity belongs to.
