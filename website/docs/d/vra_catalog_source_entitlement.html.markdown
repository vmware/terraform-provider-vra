---
layout: "vra"
page_title: "VMware vRealize Automation: Data source vra_catalog_source_entitlement"
description: A data source for catalog source entitlement.
---

# Data Source: vra\_catalog\_source\_entitlement

This data source provides information about a catalog source entitlement in vRA.

## Example Usages

This is an example of how to get a vRA cloud template catalog source entitlement its id.

```hcl
data "vra_catalog_source_entitlement" "this" {
  id         = var.catalog_source_entitlement_id
  project_id = var.project_id
}
```

This is an example of how to get a vRA cloud template catalog source entitlement by catalog source id.

```hcl
data "vra_catalog_source_entitlement" "this" {
  catalog_source_id = var_catalog_source_id
  project_id        = var.project_id
}
```


## Argument Reference

* `catalog_source_id` - (Optional) The id of the catalog source to find the entitlement. One of `catalog_source_id` or `id` must be provided.

* `id` - (Optional) The id of entitlement. One of `catalog_source_id` or `id` must be provided.

* `project_id` - (Required) The id of the project that this entitlement belongs to.


## Attribute Reference

* `definition` - Represents a catalog item or content source that is linked to a project via an entitlement.

    * `description` - Description of either the catalog item, or the catalog source.
    
    * `id` - ID of either the catalog source, or the catalog item.
    
    * `name` - Name of either the catalog item, or the catalog source.

    * `number_of_items` - Number of items in the associated catalog source.
    
    * `source_type` - Type of the catalog source.
    
    * `type` - Content definition type.