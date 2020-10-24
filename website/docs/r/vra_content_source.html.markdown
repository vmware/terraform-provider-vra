---
layout: "vra"
page_title: "VMware vRealize Automation: Resource vra_content_source"
description: A resource that can be used to create a content source in vRealize Automation(vRA).
---

# Resource: vra\_content_source

This resource provides a way to create a content source vRealize Automation(vRA).

## Example Usages

```hcl
resource "vra_content_source" "this" {
  name       = var.content_source_name
  project_id = var.project_id

  // type_id needs to be one of com.gitlab, com.github or com.vmware.marketplace
  type_id     = "com.gitlab"
  description = "Some content Source"

  //whether automatically sync content or not
  sync_enabled = "false"

  config {
    path           = "blueprint01"
    branch         = "master"
    repository     = "vracontent/vra8_content_source_test"
    content_type   = "BLUEPRINT"
    project_name   = var.project_name
    integration_id = var.integration_id
  }
}
```


## Argument Reference

* `config` - (Required) Content source custom configuration.

    * `branch` - Content source branch name.
    
    * `content_type` - Content source type. Supported values are `BLUEPRINT`, `IMAGE`, `ABX_SCRIPTS`, `TERRAFORM_CONFIGURATION`.
    
    * `integration_id` - Content source integration id as seen in vRA integrations. 
    
    * `path` - Path to refer to in the content source repository and branch.
    
    * `project_name` - Name of the project.
    
    * `repository` - Content source repository.

* `description` - (Optional) A human-friendly description.

* `name` - (Required) A human-friendly name for content source used as an identifier in APIs that support this option.

* `project_id` - (Required) The id of the project this entity belongs to.

* `sync_enabled` - (Required) Flag indicating whether sync is enabled for this content source.

* `type_id` - (Required) Content Source type. Supported values are `com.gitlab`, `com.github`, `com.vmware.marketplace`, `org.bitbucket`.


## Attribute Reference

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `created_by` - The user the entity was created by.

* `id` - The id of this cloud template.

* `last_updated_at` - Date when the entity was last updated. The date is in ISO 6801 and UTC.

* `last_updated_by` - The user the entity was last updated by.

* `org_id` - The id of the organization this entity belongs to.


## Import

Content source can be imported using the id, e.g.

`$ terraform import vra_content_source.this 05956583-6488-4e7d-84c9-92a7b7219a15`
