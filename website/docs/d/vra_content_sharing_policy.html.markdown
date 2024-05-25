---
layout: "vra"
page_title: "VMware vRealize Automation: vra_content_sharing_policy"
description: A data source for content sharing policy.
---

# Data Source: vra\_content\_sharing\_policy

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

* `id` - (Optional) The policy ID.

* `name` - (Optional) The policy name.

-> **Note:** One of `id` or `name` must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `catalog_item_ids` - List of catalog item ids to share.

* `catalog_source_ids` - List of catalog source ids to share.

* `created_at` - Policy creation timestamp.

* `created_by` - Policy author.

* `description` - The policy description.

* `last_updated_at` - Most recent policy update timestamp.

* `last_updated_by` - Most recent policy editor.

* `org_id` - The ID of the organization to which the policy belongs.

* `project_id` - The ID of the project to which the policy belongs.
