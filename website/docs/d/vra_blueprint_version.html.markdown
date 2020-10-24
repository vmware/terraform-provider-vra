---
layout: "vra"
page_title: "VMware vRealize Automation: Data source vra_blueprint_version"
description: A blueprint version data source.
---

# Data Source: vra\_blueprint\_version

This data source provides information about a cloud template (blueprint) version in vRA.

## Example Usages

```hcl
data "vra_blueprint_version" "this" {
  blueprint_id = var.blueprint_id
  id           = var.blueprint_version_id
}
```

## Argument Reference

* `blueprint_id` - (Required) Name of the cloud template. One of `id` or `name` must be provided.

* `id` - (Required) The id of the cloud template version.


## Attribute Reference

* `blueprint_description` - Description of the cloud template.

* `content` - Blueprint YAML content.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `created_by` - The user the entity was created by.

* `description` - (Optional) Cloud template version description.

* `name` - Name of the cloud template version.

* `org_id` - The id of the organization this entity belongs to.

* `project_id` - The id of the project this entity belongs to.

* `project_name` - The name of the project the entity belongs to.

* `status` - Status of the cloud template. Supported values: `DRAFT`, `VERSIONED`, `RELEASED`.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `updated_by` - The user the entity was last updated by.

* `valid` - Flag to indicate if the current content of the cloud template is valid.

* `version` - Cloud template version.

* `version_change_log` - Cloud template version change log.
