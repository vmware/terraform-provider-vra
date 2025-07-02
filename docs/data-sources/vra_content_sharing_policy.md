---
page_title: "VMware Aria Automation: vra_content_sharing_policy"
description: A data source for content sharing policy.
---

# Data Source: vra_content_sharing_policy

This is an example of how to lookup a content sharing policy data source:

**Content sharing policy data source by id:**

```hcl
data "vra_content_sharing_policy" "this" {
  id = var.vra_content_sharing_policy_id
}
```

**Content sharing policy data source by name:**

```hcl
data "vra_content_sharing_policy" "this" {
  name = var.vra_content_sharing_policy_name
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The id of the policy instance.

* `name` - (Optional) The name of the policy instance.

-> **Note:** One of `id` or `name` must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `catalog_item_ids` - List of catalog item ids to share.

* `catalog_source_ids` - List of catalog source ids to share.

* `created_at` - Date when the entity was created. The date is in ISO 8601 and UTC.

* `created_by` - The user the entity was created by.

* `description` - A human-friendly description for the policy instance.

* `enforcement_type` - The type of enforcement for the policy.

* `entitlement_type` - Entitlement type.

* `last_updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `last_updated_by` - The user the entity was last updated by.

* `org_id` - The id of the organization this entity belongs to.

* `principals` - List of users or roles that can share content:

  * `reference_id` - The reference ID of the principal.

  * `type` - The type of the principal.

* `project_criteria` - The project based criteria.

* `project_id` - The id of the project this entity belongs to.
