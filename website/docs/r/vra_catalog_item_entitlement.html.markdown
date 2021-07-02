---
layout: "vra"
page_title: "VMware vRealize Automation: Resource vra_catalog_item_entitlement"
description: A resource that can be used to create a vRealize Automation catalog item entitlement.
---

# Resource: vra\_catalog\_item\_entitlement

This resource provides a way to create a vRealize Automation(vRA) catalog item entitlement.

## Example Usages

```hcl
resource "vra_catalog_item_entitlement" "this" {
  catalog_item_id = var.catalog_item_blueprint_id
  project_id        = var.project_id
}
```


## Argument Reference

* `catalog_item_id` - (Required) The id of the catalog item to create the entitlement.

* `project_id` - (Required) The id of the project this entity belongs to.


## Attribute Reference

* `definition` - Represents a catalog item that is linked to a project via an entitlement.

    * `description` - Description of the catalog item.

    * `id` - ID of the catalog item.

    * `name` - Name of the catalog item.

    * `item_type` - Type of the catalog item.

    * `type` - Content definition type.

* `id` - The id of this catalog item entitlement.


## Import

Catalog item entitlement can be imported using the id, e.g.

`$ terraform import vra_catalog_item_entitlement.this 05956583-6488-4e7d-84c9-92a7b7219a15`
