---
page_title: "VMware Aria Automation: Resource vra_content_source"
description: A resource that can be used to create a content source in VMware Aria Automation.
---

# Resource: vra_content_source

This resource provides a way to create a content source VMware Aria Automation.

## Example Usages

```hcl
resource "vra_content_source" "this" {
  name        = var.content_source_name
  description = "Some content Source"
  project_id  = var.project_id
  sync_enabled = "false"
  type_id      = "com.gitlab"

  config {
    branch         = "main"
    content_type   = "BLUEPRINT"
    integration_id = var.integration_id
    repository     = "vracontent/vra8_content_source_test"
    path           = "blueprint01"
  }
}
```

## Argument Reference

* `config` - (Required) The content source custom configuration.

  * `branch` - (Required) The content source branch name.

  * `content_type` - (Required) The content source type. Supported values are `BLUEPRINT`, `IMAGE`, `ABX_SCRIPTS`, `TERRAFORM_CONFIGURATION`.

  * `integration_id` - (Required) The content source integration id as seen integrations.

  * `path` - (Optional) Path to refer to in the content source repository and branch.

  * `repository` - (Required) The content source repository.

* `description` - (Optional) A human-friendly description for the content source instance.

* `name` - (Required) The name of the content source instance.

* `project_id` - (Required) The id of the project this entity belongs to.

* `sync_enabled` - (Required) Wether or not sync is enabled for this content source.

* `type_id` - (Required) The type of this content source. Supported values are `com.gitlab`, `com.github`, `org.bitbucket`.

## Attribute Reference

* `created_at` - Date when the entity was created. The date is in ISO 8601 and UTC.

* `created_by` - The user the entity was created by.

* `last_updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `last_updated_by` - The user the entity was last updated by.

* `org_id` - The id of the organization this entity belongs to.

## Import

Content source can be imported using the id, e.g.

`$ terraform import vra_content_source.this 05956583-6488-4e7d-84c9-92a7b7219a15`
