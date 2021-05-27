---
layout: "vra"
page_title: "VMware vRealize Automation: Data source vra_catalog_item_entitlement"
description: A data source for catalog item entitlement.
---

# Data Source: vra\_catalog\_item\_entitlement

This data source provides information about a catalog item entitlement in vRA.

## Example Usages

This is an example of how to get a vRA cloud template catalog item entitlement its id.

```hcl
data "vra_catalog_item_entitlement" "this" {
  id         = var.catalog_item_entitlement_id
  project_id = var.project_id
}
```

This is an example of how to get a vRA cloud template catalog item entitlement by catalog item id.

```hcl
data "vra_catalog_item_entitlement" "this" {
  catalog_item_id = var_catalog_item_id
  project_id        = var.project_id
}
```


## Argument Reference

* `catalog_item_id` - (Optional) The id of the catalog item to find the entitlement. One of `catalog_item_id` or `id` must be provided.

* `id` - (Optional) The id of entitlement. One of `catalog_item_id` or `id` must be provided.

* `project_id` - (Required) The id of the project that this entitlement belongs to.


## Attribute Reference

* `definition` - Represents a catalog item or content item that is linked to a project via an entitlement.

    * `description` - Description of either the catalog item, or the catalog item.

    * `id` - ID of either the catalog item, or the catalog item.

    * `name` - Name of the catalog item.

    * `source_type` - Type of the catalog item.

    * `type` - Content definition type.
