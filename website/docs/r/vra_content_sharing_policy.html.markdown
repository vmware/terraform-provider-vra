---
layout: "vra"
page_title: "VMware vRealize Automation: vra_content_sharing_policy"
description: Creates a vra_content_sharing_policy resource.
---

# Resource: vra\_content\_sharing\_policy

Creates a VMware vRealize Automation Content Sharing Policy resource.

## Example Usages

The following example shows how to create a Content Sharing Policy resource:

```hcl
resource "vra_content_sharing_policy" "this" {
  name               = "content-sharing-policy"
  description        = "My Content Sharing Policy"
  project_id         = var.project_id

  catalog_item_ids   = [
    var.catalog_item_id
  ]

  catalog_source_ids = [
    var.catalog_source_id
  ]
}
```

## Argument Reference

Create your resource with the following arguments:

* `catalog_item_ids` - (Optional) List of catalog item ids to share.

* `catalog_source_ids` - (Optional) List of catalog source ids to share.

* `description` - (Optional) The policy description.

* `name` - (Required) The policy name.

* `project_id` - (Required) The ID of the project to which the policy belongs.

-> **Note:** One of `catalog_item_ids` or `catalog_source_ids` must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `created_at` - Policy creation timestamp.

* `created_by` - Policy author.

* `last_updated_at` - Most recent policy update timestamp.

* `last_updated_by` - Most recent policy editor.

* `org_id` - The ID of the organization to which the policy belongs.

## Import

To import an existing content sharing policy, use the `id` as in the following example:

`$ terraform import vra_content_sharing_policy 87c17193-39ee-4921-9a11-7e03e3df6029`
