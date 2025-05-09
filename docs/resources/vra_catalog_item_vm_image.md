---
page_title: "VMware Aria Automation: vra_catalog_item_vm_image"
description: A resource for VM Image Catalog Items.
---

# Resource: vra_catalog_item_vm_image

Creates a Catalog Item resource from a VM Image.

## Example Usages

The following example shows how to create a catalog item resource from a VM image:

```hcl
resource "vra_catalog_item_vm_image" "catalog_item_vm_image" {
  name        = "terraform-vm-image"
  description = "Catalog Item [terraform-vm-image] created by Terraform"
  project_id  = var.project_id

  image_name   = var.image_name
  cloud_config = "#cloud-config\n# Sample cloud-config creating an admin user. Items in {{}} will appear as user inputs when deploying the VM.\nusers:\n  - name: {{user_name}}\n    plain_text_passwd: {{password}}\n    lock_passwd: false\n    sudo:\n      - ALL=(ALL) NOPASSWD:ALL\n    groups: sudo"
}
```

## Argument Reference

Create your resource with the following arguments:

* `cloud_config` - (Optional) Cloud config script to be applied to VMs provisioned from this image.

* `description` - (Optional) A human-friendly description for the catalog item.

* `global` - (Optional) Whether to allow this catalog to be shared with multiple projects or to restrict it to the specified project.

* `icon_id` - (Optional) ID of the icon to associate with this catalog item.

* `image_name` - (Required) Name of the VM image to publish.

* `name` - (Required) The name of the catalog item.

* `project_id` - (Required) ID of the project to share this catalog item with.

* `select_zone` - (Optional) Whether to create a zone input for the published catalog item.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `created_at` - Date when the entity was created. The date is in ISO 8601 and UTC.

* `created_by` - The user the entity was created by.

* `last_updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `last_updated_by` - The user the entity was last updated by.

## Import

To import an existing Catalog Item, use the `id` as in the following example:

`$ terraform import vra_catalog_item_vm_image.catalog_item_vm_image "b580ef9f-f191-4de2-b6f4-96f29d8c8d3d"`
