---
page_title: "VMware Aria Automation: vra_policy_approval"
description: A data source for Approval policies.
---

# Data Source: vra_policy_approval

The following examples shows how to lookup for an approval policy:

**Approval policy data source by its id:**

```hcl
data "vra_policy_approval" "this" {
  id = var.vra_approval_policy_id
}
```

**Approval policy data source by name search:**

```hcl
data "vra_policy_approval" "this" {
  search = var.vra_approval_policy_search_name
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The id of the policy instance.

* `search` - (Optional) Search criteria to narrow down the policy instance.

-> **Note:** One of `id` or `search` must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `actions` - List of actions to trigger approval.

* `approval_level` - The level defines the order in which the policy is enforced. Level 1 approvals are applied first, followed by level 2 approvals, and so on.

* `approval_mode` - Who must approve the request.

* `approval_type` - Approval Type.

* `approvers` - List of approvers of the policy.

* `auto_approval_decision` - Automatically approve or reject a request after the number of days specified in the Auto expiry trigger field.

* `auto_approval_expiry` - The number of days the approvers have to respond before the Auto action is triggered.

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
