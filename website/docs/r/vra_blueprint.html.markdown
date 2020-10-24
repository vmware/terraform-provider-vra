---
layout: "vra"
page_title: "VMware vRealize Automation: Resource vra_blueprint"
description: A resource that can be used to create a vRealize Automation cloud template, formerly know as blueprint.
---

# Resource: vra\_blueprint

This resource provides a way to create a vRealize Automation(vRA) cloud template, formerly known as blueprint.

## Example Usages

```hcl
resource "vra_blueprint" "this" {
  name        = var.blueprint_name
  description = "Created by vRA terraform provider"

  project_id = vra_project.this.id

  content = <<-EOT
    formatVersion: 1
    inputs:
      image:
        type: string
        description: "Image"
      flavor:
        type: string
        description: "Flavor"
    resources:
      Machine:
        type: Cloud.Machine
        properties:
          image: $${input.image}
          flavor: $${input.flavor}
  EOT
}
```


## Argument Reference

* `content` - (Optional) Blueprint YAML content.

* `description` - (Optional) A human-friendly description.

* `name` - (Required) A human-friendly name used as an identifier in APIs that support this option.

* `project_id` - (Required) The id of the project this entity belongs to. 

* `request_scope_org` - (Optional) Flag to indicate whether this blueprint can be requested from any project in the organization this entity belongs to.


## Attribute Reference

* `content_source_id` - The id of the content source. 

* `content_source_path` - Content source path.

* `content_source_sync_at` - Content source last sync at.

* `content_source_sync_messages` - Content source last sync messages.

* `content_source_sync_status` - Content source last sync status. Supported values: `SUCCESSFUL`, `FAILED`.

* `content_source_type` - Content source type.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `created_by` - The user the entity was created by.

* `id` - The id of this cloud template.

* `org_id` - The id of the organization this entity belongs to.

* `project_name` - The name of the project the entity belongs to.

* `self_link` - HATEOAS of the entity.

* `status` - Status of the cloud template. Supported values: `DRAFT`, `VERSIONED`, `RELEASED`.

* `total_released_versions` - Total number of released versions. 

* `total_versions` - Total number of versions.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `updated_by` - The user the entity was last updated by.

* `valid` - Flag to indicate if the current content of the cloud template/blueprint is valid.

* `validation_messages` - List of validations messages.
    * message - Validation message.
    
    * metadata - Validation metadata.
    
    * path - Validation path.
    
    * resource_name - Name of the resource.
    
    * type - Message type. Supported values: `INFO`, `WARNING`, `ERROR`.


## Import

Cloud template can be imported using the id, e.g.

`$ terraform import vra_blueprint.this 05956583-6488-4e7d-84c9-92a7b7219a15`
