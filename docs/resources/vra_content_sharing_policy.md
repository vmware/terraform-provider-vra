---
page_title: "VMware Aria Automation: vra_content_sharing_policy"
description: Creates a vra_content_sharing_policy resource.
---

# Resource: vra_content_sharing_policy

Creates a VMware Aria Automation content sharing policy resource.

## Example Usages

The following example shows how to create a content sharing policy resource:

**To share the content with a specific set of users:**

```hcl
resource "vra_content_sharing_policy" "this" {
  name        = "content-sharing-policy"
  description = "My Content Sharing Policy"
  project_id  = var.project_id

  entitlement_type = "USER"
  principals {
    reference_id   = var.username
    type           = "USER"
  }

  catalog_item_ids = [
    var.catalog_item_id
  ]

  catalog_source_ids = [
    var.catalog_source_id
  ]
}
```

**To share the content with a specific set of roles:**

```hcl
resource "vra_content_sharing_policy" "this" {
  name        = "content-sharing-policy"
  description = "My Content Sharing Policy"
  project_id  = var.project_id

  entitlement_type = "ROLE"
  principals {
    reference_id   = var.role
    type           = "ROLE"
  }

  catalog_item_ids = [
    var.catalog_item_id
  ]

  catalog_source_ids = [
    var.catalog_source_id
  ]
}
```

**To share the content with all users and groups in the project or organization:**

```hcl
resource "vra_content_sharing_policy" "this" {
  name        = "content-sharing-policy"
  description = "My Content Sharing Policy"
  project_id  = var.project_id

  entitlement_type = "USER"
  principals {
    reference_id   = ""
    type           = "PROJECT"
  }

  catalog_item_ids = [
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

* `catalog_source_ids` - (Optional) List of catalog source ids to share

* `description` - (Optional) A human-friendly description for the policy instance.

* `entitlement_type` - (Optional) Entitlement type. Supported values: `USER`, `ROLE`.

* `name` - (Required) The name of the policy instance.

* `principals` - (Optional) List of users or roles that can share content:

  * `reference_id` - (Optional) The reference ID of the principal.

  * `type` - (Required) The type of the principal.

* `project_criteria` - (Optional) The project based criteria. Updating this argument triggers a recreation of the resource. It cannot be specified when `project_id` is set.

* `project_id` - (Optional) The id of the project this entity belongs to. Updating this argument triggers a recreation of the resource.

-> **Note:** One of `catalog_item_ids` or `catalog_source_ids` must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `created_at` - Date when the entity was created. The date is in ISO 8601 and UTC.

* `created_by` - he user the entity was created by.

* `enforcement_type` - The type of enforcement for the policy.

* `last_updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `last_updated_by` - The user the entity was last updated by.

* `org_id` - The id of the organization this entity belongs to.

## Import

To import an existing content sharing policy, use the `id` as in the following example:

`$ terraform import vra_content_sharing_policy 87c17193-39ee-4921-9a11-7e03e3df6029`
