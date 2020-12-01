---
layout: "vra"
page_title: "VMware vRealize Automation: Resource vra_catalog_source_blueprint"
description: A resource that can be used to create a vRealize Automation catalog source of type cloud template (blueprint).
---

# Resource: vra\_catalog\_source\_blueprint

Creates a VMware vRealize Automation catalog source resource of type cloud template, formerly known as a blueprint.

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

* `config` - (Optional) Custom configuration of the catalog source as a map of key values.

* `description` - (Optional) Human-friendly description.

* `name` - (Required) Human-friendly name used as an identifier in APIs that support this option.

* `project_id` - (Required) ID of the project this entity belongs to. 


## Attribute Reference 

* `created_at` - Date when entity was created. Date and time format is ISO 8601 and UTC.

* `created_by` - User who created the entity.

* `global` - Flag indicating that all items can be requested across all projects.

* `id` - ID of catalog source.

* `items_found` - Number of items found in the catalog source.

* `items_imported` - Number of items imported from the catalog source.

* `last_import_completed_at` - Time at which the last import completed.

* `last_import_errors` - List of errors seen when the catalog source was last imported.

* `last_import_started_at` - Time at which the last import started.

* `last_updated_by` - User who last updated the catalog source. 

* `type_id` - Type of catalog source. Example: `blueprint`, `CFT`, etc.


## Import

To import the cloud template catalog source, use the ID as in the following example:

`$ terraform import vra_catalog_source_blueprint.this 05956583-6488-4e7d-84c9-92a7b7219a15`
