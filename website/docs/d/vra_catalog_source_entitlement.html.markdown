---
layout: "vra"
page_title: "VMware vRealize Automation: vra_catalog_source_entitlement"
description: A data source for catalog source entitlement.
---

# Data Source: vra\_catalog\_source\_entitlement

> **Note**:  Deprecated - please use `vra_content_sharing_policy` instead.

This data source provides information about a catalog source entitlement in vRA.

## Example Usages

This is an example of how to get a vRA catalog source entitlement by its id:

```hcl
data "vra_catalog_source_entitlement" "this" {
  id         = var.catalog_source_entitlement_id
  project_id = var.project_id
}
```

This is an example of how to get a vRA catalog source entitlement by its catalog source id:

```hcl
data "vra_catalog_source_entitlement" "this" {
  catalog_source_id = var.catalog_source_id
  project_id        = var.project_id
}
```

## Argument Reference

* `catalog_source_id` - (Optional) The id of the catalog source to find the entitlement. One of `catalog_source_id` or `id` must be provided.

* `id` - (Optional) The id of entitlement. One of `catalog_source_id` or `id` must be provided.

* `project_id` - (Required) The id of the project that this entitlement belongs to.

## Attribute Reference

* `definition` - Represents a catalog source that is linked to a project via an entitlement.

    * `description` - Description of the catalog source.

    * `icon_id` - Icon id of associated catalog source.

    * `id` - Id of the catalog source.

    * `name` - Name of the catalog source.

    * `number_of_items` - Number of items in the associated catalog source.

    * `source_name` - Catalog source name.

    * `source_type` - Catalog source type.

    * `type` - Content definition type.
