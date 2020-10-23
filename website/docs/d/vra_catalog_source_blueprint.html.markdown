---
layout: "vra"
page_title: "VMware vRealize Automation: Data source vra_catalog_source_blueprint"
description: A data source for catalog source of type cloud template (blueprint).
---

# Data Source: vra\_catalog\_source\_blueprint

This data source provides information about a catalog source of type cloud template (blueprint) in vRA.

## Example Usages

This is an example of how to get a vRA cloud template catalog source by its name.

```hcl
data "vra_catalog_source_blueprint" "this" {
  name = var.catalog_source_name
}
```

This is an example of how to get a vRA cloud template catalog source by its id.

```hcl
data "vra_catalog_source_blueprint" "this" {
  id = var.catalog_source_id
}
```

This is an example of how to get a vRA cloud template catalog source by the project id it is associated with.

```hcl
data "vra_catalog_source_blueprint" "this" {
  project_id = var.project_id
}
```


## Argument Reference

* `id` - (Optional) The id of catalog source. One of `id`, `name` or `project_id`  must be provided.

* `name` - (Optional) Name of the catalog source. One of `id`, `name` or `project_id` must be provided.

* `project_id` - (Optional) The id of the project.  One of `id`, `name` or `project_id` must be provided.


## Attribute Reference

* `config` - Custom configuration of the catalog source as a map of key values.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `created_by` - The user the entity was created by.

* `description` - Catalog source description.

* `global` - A flag indicating that all the items can be requested across all projects.

* `items_found` - Number of items found in the catalog source.

* `items_imported` - Number of items imported from the catalog source.

* `last_import_completed_at` - Time at which the last import was completed at.

* `last_import_errors` - A list of errors seen at last time the catalog source is imported.

* `last_import_started_at` - Time at which the last import was started at.

* `last_updated_by` - The user that last updated the catalog source. 

* `type_id` - Type of catalog source. Example: `blueprint`, `CFT`, etc.
