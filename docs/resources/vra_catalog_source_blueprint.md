---
page_title: "VMware Aria Automation: Resource vra_catalog_source_blueprint"
description: A resource that can be used to create a catalog source of type cloud template (blueprint).
---

# Resource: vra_catalog_source_blueprint

Creates a catalog source resource of type cloud template (blueprint).

## Example Usages

The following example shows how to create a catalog source resource.

```hcl
resource "vra_catalog_source_blueprint" "this" {
  name       = var.catalog_source_name
  project_id = var.vra_project_id
}
```

## Argument Reference

Create your catalog resource with the following arguments:

* `description` - (Optional) A human-friendly description for the blueprint content source instance.

* `name` - (Required) The name of the blueprint content source instance.

* `project_id` - (Required) The id of the project the blueprint content source instance belongs to.

## Attribute Reference

* `config` - The content source custom configuration.

* `created_at` - Date when the entity was created. The date is in ISO 8601 and UTC.

* `created_by` - The user the entity was created by.

* `description` - A human-friendly description for the blueprint content source instance.

* `global` - Global flag indicating that all the items can be requested across all projects.

* `icon_id` - Default Icon Identifier.

* `items_found` - Number of items found.

* `items_imported` - Number of items imported.

* `last_import_completed_at` - Date when the last import was completed. The date is in ISO 8601 and UTC.

* `last_import_errors` - A list of errors seen at last time the content source is imported.

* `last_import_started_at` - Date when the last import was started. The date is in ISO 8601 and UTC.

* `last_updated_by` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `type_id` - The type of this content source. Example: `blueprint`, `CFT`, etc.

## Import

To import the cloud template catalog source, use the ID as in the following example:

`$ terraform import vra_catalog_source_blueprint.this 05956583-6488-4e7d-84c9-92a7b7219a15`
