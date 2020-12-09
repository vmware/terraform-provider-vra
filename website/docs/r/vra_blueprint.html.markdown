---
layout: "vra"
page_title: "VMware vRealize Automation: Resource vra_blueprint"
description: A resource that can be used to create a vRealize Automation cloud template, formerly know as blueprint.
---

# Resource: vra\_blueprint

Creates a VMware vRealize Automation (vRA) cloud template resource, formerly known as a blueprint.

## Example Usage

The following example shows how to create a blueprint resource.

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

Create your blueprint resource with the following arguments:

* `content` - (Optional) Blueprint YAML content.

* `description` - (Optional) Human-friendly description.

* `name` - (Required) Human-friendly name used as an identifier in APIs that support this option.

* `project_id` - (Required) ID of project that entity belongs to. 

* `request_scope_org` - (Optional) Flag to indicate whether blueprint can be requested from any project in the organization that entity belongs to.


## Attribute Reference

* `content_source_id` - ID of content source. 

* `content_source_path` - Content source path.

* `content_source_sync_at` - Date when content source was last synced. The date is in ISO 8601 and UTC.

* `content_source_sync_messages` - Content source last sync messages.

* `content_source_sync_status` - Content source last sync status. Supported values: `SUCCESSFUL`, `FAILED`.

* `content_source_type` - Content source type.

* `created_at` - Date when entity was created. The date is in ISO 8601 and UTC.

* `created_by` - The user who created entity.

* `id` - ID of cloud template.

* `org_id` - ID of organization that entity belongs to.

* `project_name` - Name of project that entity belongs to.

* `self_link` - HATEOAS of entity.

* `status` - Status of cloud template. Supported values: `DRAFT`, `VERSIONED`, `RELEASED`.

* `total_released_versions` - Total number of released versions. 

* `total_versions` - Total number of versions.

* `updated_at` - Date when entity was last updated. Date and time format is ISO 8601 and UTC.

* `updated_by` - The user who last updated the entity.

* `valid` - Flag to indicate if the current content of the cloud template/blueprint is valid.

* `validation_messages` - List of validations messages.
    * message - Validation message.
    
    * metadata - Validation metadata.
    
    * path - Validation path.
    
    * resource_name - Name of resource.
    
    * type - Message type. Supported values: `INFO`, `WARNING`, `ERROR`.


## Import

To import the cloud template, use the ID as in the following example:

`$ terraform import vra_blueprint.this 05956583-6488-4e7d-84c9-92a7b7219a15`
