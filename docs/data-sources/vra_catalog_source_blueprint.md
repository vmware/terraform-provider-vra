---
page_title: "VMware Aria Automation: Data source vra_catalog_source_blueprint"
description: A data source for catalog source of type cloud template (blueprint).
---

# Data Source: vra_catalog_source_blueprint

This data source provides information about a catalog source of type cloud template (blueprint).

## Example Usages

This is an example of how to get a cloud template catalog source by its name.

```hcl
data "vra_catalog_source_blueprint" "this" {
  name = var.catalog_source_name
}
```

This is an example of how to get a cloud template catalog source by its id.

```hcl
data "vra_catalog_source_blueprint" "this" {
  id = var.catalog_source_id
}
```

This is an example of how to get a cloud template catalog source by the project id it is associated with.

```hcl
data "vra_catalog_source_blueprint" "this" {
  project_id = var.project_id
}
```

## Argument Reference

* `id` - (Optional) The id of the blueprint content source instance.

* `name` - (Optional) The name of the blueprint content source instance.

* `project_id` - (Optional) The id of the project the blueprint content source instance belongs to.

-> **Note:** One of `id`, `name` or `project_id`  must be provided.

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
