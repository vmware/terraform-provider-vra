---
layout: "vra"
page_title: "VMware vRealize Automation: Data source vra_blueprint"
description: A blueprint data source.
---

# Data Source: vra\_blueprint

This data source provides information about a cloud template (blueprint) in vRA.

## Example Usages

This is an example of how to get a vRA cloud template by its name.

```hcl
data "vra_blueprint" "this" {
  name = vra_blueprint.this.name
}
```

This is an example of how to get a vRA cloud template by its id.

```hcl
data "vra_blueprint" "this" {
  id = vra_blueprint.this.id
}
```

## Argument Reference

* `id` - (Optional) The id of this cloud template. One of `id` or `name` must be provided.

* `name` - (Optional) Name of the cloud template. One of `id` or `name` must be provided.

* `project_id` - (Optional) The id of the project to narrow the search while looking for cloud templates. 


## Attribute Reference

* `content` - Blueprint YAML content.

* `content_source_id` - The id of the content source. 

* `content_source_path` - Content source path.

* `content_source_sync_at` - Content source last sync at.

* `content_source_sync_messages` - Content source last sync messages.

* `content_source_sync_status` - Content source last sync status. Supported values: `SUCCESSFUL`, `FAILED`.

* `content_source_type` - Content source type.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `created_by` - The user the entity was created by.

* `description` - A human-friendly description.

* `org_id` - The id of the organization this entity belongs to.

* `project_name` - The name of the project the entity belongs to.

* `self_link` - HATEOAS of the entity.

* `request_scope_org` - Flag to indicate whether this blueprint can be requested from any project in the organization this entity belongs to.

* `status` - Status of the cloud template. Supported values: `DRAFT`, `VERSIONED`, `RELEASED`.

* `total_released_versions` - Total number of released versions. 

* `total_versions` - Total number of versions.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `updated_by` - The user the entity was last updated by.

* `valid` - Flag to indicate if the current content of the cloud template is valid.

* `validation_messages` - List of validations messages.
    * message - Validation message.
    
    * metadata - Validation metadata.
    
    * path - Validation path.
    
    * resource_name - Name of the resource.
    
    * type - Message type. Supported values: `INFO`, `WARNING`, `ERROR`.
