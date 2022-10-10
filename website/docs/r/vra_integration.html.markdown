---
layout: "vra"
page_title: "VMware vRealize Automation: vra_integration"
description: |-
    Creates a vra_integration resource.
---

# Resource: vra\_integration

Creates a VMware vRealize Automation Integration resource.

## Example Usages

The following example shows how to create an Integration resource:

```hcl
resource "vra_integration" "this" {
  name                   = "saltstack"
  description            = "SaltStack Integration"
  integration_properties = {
    hostName: var.hostname
  }
  integration_type       = "saltstack"
  private_key_id         = var.username
  private_key            = var.password

  tags {
    key   = "created_by"
    value = "vra-terraform-provider"
  }
}
```

## Argument Reference

Create your integration resource with the following arguments:

* `associated_cloud_account_ids` - (Optional) Ids of the cloud accounts to associate with this integration.

* `certificate` - (Optional) Certificate to be used to connect to the integration.

* `custom_properties` - (Optional) Additional custom properties that may be used to extend the Integration.

* `description` - (Optional) A human-friendly description.

* `integration_properties` - (Required) Integration specific properties supplied in as name value pairs.

* `integration_type` - (Required) Integration type.

* `name` - (Required) The name of the integration.

* `private_key` - (Optional) Secret access key or password to be used to authenticate with the integration.

* `private_key_id` - (Optional) Access key id or username to be used to authenticate with the integration.

* `tags` - (Optional) A set of tag keys and optional values to apply to the integration. Example: `[ { "key" : "provider", "value": "vmware" } ]`.

## Attribute Reference

* `created_at` - Date when the entity was created. The date is in ISO 8601 and UTC.

* `id` - (Optional) The id of the integration.

* `links` - HATEOAS of entity.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.


## Import

To import an existing integration, use the id as in the following example:

`$ terraform import vra_integration 90b9a230-bd61-4d39-a082-b12a17cd03c8`
