---
layout: "vra"
page_title: "VMware vRealize Automation: Resource vra_catalog_source_entitlement"
description: A resource that can be used to create a vRealize Automation catalog source entitlement.
---

# Resource: vra\_catalog\_source\_entitlement

This resource provides a way to create a vRealize Automation(vRA) catalog source entitlement.

## Example Usages

```hcl
resource "vra_catalog_source_entitlement" "this" {
  catalog_source_id = var.catalog_source_blueprint_id
  project_id        = var.project_id
}
```


## Argument Reference

* `catalog_source_id` - (Required) The id of the catalog source to create the entitlement.

* `project_id` - (Required) The id of the project this entity belongs to. 


## Attribute Reference 

* `definition` - Represents a catalog item or content source that is linked to a project via an entitlement.

    * `description` - Description of either the catalog item, or the catalog source.
    
    * `id` - ID of either the catalog source, or the catalog item.
    
    * `name` - Name of either the catalog item, or the catalog source.

    * `number_of_items` - Number of items in the associated catalog source.
    
    * `source_type` - Type of the catalog source.
    
    * `type` - Content definition type.
    
* `id` - The id of this catalog source entitlement.


## Import

Catalog source entitlement can be imported using the id, e.g.

`$ terraform import vra_catalog_source_entitlement.this 05956583-6488-4e7d-84c9-92a7b7219a15`
