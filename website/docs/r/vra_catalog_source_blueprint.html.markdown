---
layout: "vra"
page_title: "VMware vRealize Automation: Resource vra_catalog_source_blueprint"
description: A resource that can be used to create a vRealize Automation catalog source of type cloud template (blueprint).
---

# Resource: vra\_catalog\_source\_blueprint

This resource provides a way to create a vRealize Automation(vRA) catalog source of type catalog item (blueprint).

## Example Usages

```hcl
resource "vra_catalog_source_blueprint" "this" {
  name       = var.catalog_source_name
  project_id = var.vra_project_id
}
```


## Argument Reference

* `config` - (Optional) Custom configuration of the catalog source as a map of key values.

* `description` - (Optional) A human-friendly description.

* `name` - (Required) A human-friendly name used as an identifier in APIs that support this option.

* `project_id` - (Required) The id of the project this entity belongs to. 


## Attribute Reference 

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `created_by` - The user the entity was created by.

* `global` - A flag indicating that all the items can be requested across all projects.

* `id` - The id of this catalog source.

* `items_found` - Number of items found in the catalog source.

* `items_imported` - Number of items imported from the catalog source.

* `last_import_completed_at` - Time at which the last import was completed at.

* `last_import_errors` - A list of errors seen at last time the catalog source is imported.

* `last_import_started_at` - Time at which the last import was started at.

* `last_updated_by` - The user that last updated the catalog source. 

* `type_id` - Type of catalog source. Example: `blueprint`, `CFT`, etc.


## Import

Cloud template Catalog source can be imported using the id, e.g.

`$ terraform import vra_catalog_source_blueprint.this 05956583-6488-4e7d-84c9-92a7b7219a15`
