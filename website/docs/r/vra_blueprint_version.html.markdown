---
layout: "vra"
page_title: "VMware vRealize Automation: Resource vra_blueprint_version"
description: A resource that can be used to create a vRealize Automation cloud template version.
---

# Resource: vra\_blueprint\_version

This resource provides a way to create a vRealize Automation(vRA) cloud template (blueprint) version.

## Example Usages

```hcl
resource "vra_blueprint_version" "this" {
  blueprint_id = var.vra_blueprint_id
  change_log   = "First version"
  description  = "Released from vRA terraform provider"
  release      = true
  version      = (random_integer.suffix.result / random_integer.suffix.result)
}
```

## Argument Reference

* `blueprint_id` - (Required) The id of the cloud template.

* `change_log` - (Optional) Cloud template version log.

* `description` - (Optional) A human-friendly description for the cloud template version. 
 
* `releae` - (Optional) Flag to indicate whether to release the version.

* `version` - (Required) Cloud template version.


## Attribute Reference

* `blueprint_description` - Description of the cloud template.

* `content` - Blueprint YAML content.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `created_by` - The user the entity was created by.

* `id` - The id of this cloud template version.

* `name` - Name of the cloud template version.

* `org_id` - The id of the organization this entity belongs to.

* `project_id` - The id of the project this entity belongs to.

* `project_name` - The name of the project the entity belongs to.

* `status` - Status of the cloud template. Supported values: `DRAFT`, `VERSIONED`, `RELEASED`.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `updated_by` - The user the entity was last updated by.

* `valid` - Flag to indicate if the current content of the cloud template is valid.

## Import

Cloud template version can be imported using the id, e.g.

`$ terraform import vra_blueprint_version.this 05956583-6488-4e7d-84c9-92a7b7219a15`
