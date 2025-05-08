---
page_title: "VMware Aria Automation: vra_catalog_item_vro_workflow"
description: A resource for vRO Workflow Catalog Items.
---

# Resource: vra_catalog_item_vro_workflow

Creates a Catalog Item resource from a vRO Workflow.

## Example Usages

The following example shows how to create a catalog item resource from a vRO Workflow:

```hcl
resource "vra_catalog_item_vro_workflow" "catalog_item_vro_workflow" {
  name        = "terraform-vro-workflow"
  description = "Catalog Item [terraform-vro-workflow] created by Terraform"
  project_id  = var.project_id

  workflow_id = var.workflow_id
}
```

## Argument Reference

Create your resource with the following arguments:

* `description` - (Optional) A human-friendly description for the catalog item.

* `global` - (Optional) Whether to allow this catalog to be shared with multiple projects or to restrict it to the specified project.

* `icon_id` - (Optional) ID of the icon to associate with this catalog item.

* `name` - (Required) The name of the catalog item.

* `project_id` - (Required) ID of the project to share this catalog item with.

* `workflow_id` - (Required) ID of the vRO workflow to publish.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `created_at` - Date when the entity was created. The date is in ISO 8601 and UTC.

* `created_by` - The user the entity was created by.

* `last_updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `last_updated_by` - The user the entity was last updated by.

## Import

To import an existing Catalog Item, use the `id` as in the following example:

`$ terraform import vra_catalog_item_vro_workflow.catalog_item_vro_workflow "a090b0c2-5b49-4fb5-9e69-c1b84b01c908"`
