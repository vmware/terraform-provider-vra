---
layout: "vra"
page_title: "VMware vRealize Automation: vra_cloud_account_azure"
description: |-
  Provides a VMware vRA vra_cloud_account_azure resource.
---

# Resource: vra\_cloud\_account\_azure

Provides a VMware vRA vra_cloud_account_azure resource.

## Example Usages

**Create Azure cloud account:**

This is an example of how to create an Azure cloud account resource.

```hcl

resource "vra_cloud_account_azure" "this" {
  name = "my-cloud-account-%d"
  description = "test cloud account"
  subscription_id = "sample-subscription-id"
  tenant_id = "sample-tenant-id"
  application_id = "sample-application-id"
  application_key = "sample-application=key"
  regions = ["centralus"]
}

```



## Argument Reference

The following arguments are supported for an Azure cloud account resource:

* `application_id` - (Required) Azure Client Application ID.

* `application_key` - (Required) Azure Client Application Secret Key.

* `description` - (Optional) A human-friendly description.

* `name` - (Optional) The name of this Azure cloud account.

* `regions` - (Optional) A set of region names that are enabled for this account.

* `subscription_id` - (Required) Azure Subscription ID.

* `tags` - (Optional) A set of tag keys and optional values that to set on this resource.
example:[ { "key" : "vmware", "value": "provider" } ]
  * `key` - Tag’s key.
  * `value` - Tag’s value.

* `tenant_id` - (Required) Azure Tenant ID.

## Attribute Reference

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `links` - HATEOAS of the entity.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

## Import

Azure cloud account can be imported using the id, e.g.

`$ terraform import vra_cloud_account_azure.new_azure 05956583-6488-4e7d-84c9-92a7b7219a15`