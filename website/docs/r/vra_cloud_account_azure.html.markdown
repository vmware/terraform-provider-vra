---
layout: "vra"
page_title: "VMware vRealize Automation: vra_cloud_account_azure"
description: |-
  Creates a vra_cloud_account_azure resource.
---

# Resource: vra\_cloud\_account\_azure

Creates a VMware vRealize Automation Azure cloud account resource.

## Example Usages

The following example shows how to create an Azure cloud account resource.

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

Create your Azure cloud account resource with the following arguments:

* `application_id` - (Required) Azure Client Application ID.

* `application_key` - (Required) Azure Client Application Secret Key.

* `description` - (Optional) Human-friendly description.

* `name` - (Optional) Name of Azure cloud account.

* `regions` - (Optional) Set of region names enabled for the cloud account.

* `subscription_id` - (Required) Azure Subscription ID.

* `tags` - (Optional) Set of tag keys and values to apply to the cloud account.
Example: `[ { "key" : "vmware", "value": "provider" } ]`

* `tenant_id` - (Required) Azure Tenant ID.

## Attribute Reference

* `created_at` - Date when entity was created. Date and time format is ISO 8601 and UTC.

* `links` - HATEOAS of entity.

* `org_id` - ID of organization that entity belongs to.

* `owner` - Email of entity owner.

* `updated_at` - Date when entity was last updated. Date and time format is ISO 8601 and UTC.

## Import

To import the Azure cloud account, use the ID as in the following example:

`$ terraform import vra_cloud_account_azure.new_azure 05956583-6488-4e7d-84c9-92a7b7219a15`
