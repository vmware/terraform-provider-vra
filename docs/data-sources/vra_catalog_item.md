---
page_title: "VMware Aria Automation: Data source vra_catalog_item"
description: A data source for a catalog item.
---

# Data Source: vra_catalog_item

This data source provides information about a catalog item.

## Example Usages

This is an example of how to get a catalog item by its name.

```hcl
data "vra_catalog_item" "this" {
  name            = var.catalog_item_name
  expand_versions = true
}
```

This is an example of how to get a catalog item by its id.

```hcl
data "vra_catalog_item" "this" {
  id              = var.catalog_item_id
  expand_versions = true
}
```

## Argument Reference

* `expand_projects` - (Optional) Flag to indicate whether to expand detailed project data for the catalog item.

* `expand_versions` - (Optional) Flag to indicate whether to expand detailed versions of the catalog item.

* `id` - (Optional) The id of catalog item. One of `id` or `name` must be provided.

* `name` - (Optional) Name of the catalog item. One of `id` or `name` must be provided.

* `project_id` - (Optional) The id of the project to narrow the search while looking for catalog items.

## Attribute Reference

* `created_at` - Date when the entity was created. The date is in ISO 8601 and UTC.

* `created_by` - The user the entity was created by.

* `description` - A human-friendly description for the catalog item.

* `form_id` - ID of the form associated with this catalog item.

* `global` - Whether to allow this catalog to be shared with multiple projects or to restrict it to the specified project.

* `icon_id` - ID of the icon associated with this catalog item.

* `last_updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `last_updated_by` - The user the entity was last updated by.

* `project_ids` - List of associated project IDs that can be used for requesting this catalog item.

* `projects` - List of associated projects that can be used for requesting this catalog item.

  * `description` - A human friendly description.

  * `id` - Id of the entity.

  * `name` - Name of the entity.

  * `version` - Version of the entity, if applicable.

* `schema` - JSON schema describing request parameters, a simplified version of <http://json-schema.org/latest/json-schema-validation.html#rfc.section.5>

* `source_id` - LibraryItem source ID.

* `source_name` - LibraryItem source name.

* `source_project_id` - Project ID with which the catalog item was associated when created.

* `type` -

  * `description` - A human friendly description.

  * `id` - Id of the entity.

  * `name` - Name of the entity.

  * `version` - Version of the entity, if applicable.

* `versions` - Catalog item versions.

  * `created_at` - Date-time when catalog item version was created at.

  * `description` - A human-friendly description.

  * `id` - Id of the catalog item version.
