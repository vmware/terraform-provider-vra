---
layout: "vra"
page_title: "VMware vRealize Automation: Resource vra_blueprint_version"
description: A resource that can be used to create a vRealize Automation cloud template version.
---

# Resource: vra\_blueprint\_version

Creates a VMware vRealize Automation cloud template (blueprint) version resource.

## Example Usages

The following example shows how to create a cloud template (blueprint) version resource.

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

Create your cloud template (blueprint) version resource with the following arguments:

* `blueprint_id` - (Required) ID of the cloud template  (blueprint).

* `change_log` - (Optional) Cloud template  (blueprint) version log.

* `description` - (Optional) Human-friendly description for the cloud template  (blueprint) version. 
 
* `release` - (Optional) Flag to indicate whether to release the version.

* `version` - (Required) Cloud template  (blueprint) version.


## Attribute Reference

* `blueprint_description` - Description of cloud template (blueprint).

* `content` - Blueprint YAML content.

* `created_at` - Date when the entity was created. Date and time format is ISO 8601 and UTC.

* `created_by` - User who created the entity.

* `id` - ID of cloud template (blueprint) version.

* `name` - Name of cloud template (blueprint) version.

* `org_id` - ID of organization that entity belongs to.

* `project_id` - ID of project that entity belongs to.

* `project_name` - Name of project that entity belongs to.

* `status` - Status of the cloud template (blueprint). Supported values: `DRAFT`, `VERSIONED`, `RELEASED`.

* `updated_at` - Date when the entity was last updated. Date and time format is ISO 8601 and UTC.

* `updated_by` - User who last updated the entity.

* `valid` - Flag to indicate if the current content of the cloud template (blueprint) is valid.

## Import

To import the cloud template (blueprint) version, use the ID as in the following example:

`$ terraform import vra_blueprint_version.this 05956583-6488-4e7d-84c9-92a7b7219a15`
