---
page_title: "VMware Aria Automation: vra_content_source"
description: A data source for content source.
---

# Data Source: vra_content_source

This is an example of how to lookup a content source data source:

**Content source data source by id:**

```hcl
data "vra_content_source" "this" {
  id = var.vra_content_source_id
}
```

**Content source data source by name:**

```hcl
data "vra_content_source" "this" {
  name = var.vra_content_source_name
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The id of the content source instance.

* `name` - (Optional) The name of the content source instance.

-> **Note:** One of `id` or `name` must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `config` - The content source custom configuration.

  * `branch` - The content source branch name.

  * `content_type` - The content source type.

  * `integration_id` - The content source integration id as seen integrations.

  * `path` - Path to refer to in the content source repository and branch.

  * `project_name` - The name of the project.

  * `repository` - The content source repository.

* `created_at` - Date when the entity was created. The date is in ISO 8601 and UTC.

* `created_by` - The user the entity was created by.

* `description` - A human-friendly description for the content source instance.

* `last_updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `last_updated_by` - The user the entity was last updated by.

* `org_id` - The id of the organization this entity belongs to.

* `project_id` - The id of the project this entity belongs to.

* `sync_enabled` - Wether or not sync is enabled for this content source.

* `type_id` - The type of this content source.
