---
layout: "vra"
page_title: "VMware vRealize Automation: vra_cloud_account_aws"
description: |-
  Creates a vra_cloud_account_aws resource.
---

# Resource: vra\_cloud\_account\_aws

Creates a VMware vRealize Automation AWS cloud account resource.

## Example Usages

The following example shows how to create an AWS cloud account resource.

```hcl
resource "vra_cloud_account_aws" "this" {
  name        = "tf-vra-cloud-account-aws"
  description = "terraform test cloud account aws"
  access_key  = var.access_key
  secret_key  = var.secret_key
  regions     = ["us-east-1", "us-west-1"]  // Regions to be enabled for this cloud account

  tags {
    key   = "foo"
    value = "bar"
  }
}
```

## Argument Reference

Create your AWS cloud account resource with the following arguments:

* `access_key` - (Required) Access key ID for AWS.

* `description` - (Optional) Human-friendly description.

* `name` - (Required) Name of AWS cloud account.

* `regions` - (Optional) Set of region names enabled for the cloud account.

* `secret_key` - (Required) AWS Secret Access Key

* `tags` - (Optional) Set of tag keys and values to apply to the cloud account.
Example: `[ { "key" : "vmware", "value": "provider" } ]`



## Attribute Reference

* `created_at` - Date when entity was created. Date and time format is ISO 8601 and UTC.

* `id` - ID of AWS cloud account.

* `links` - HATEOAS of entity.

* `org_id` - ID of organization that entity belongs to.

* `owner` - Email of entity owner.

* `updated_at` - Date when entity was last updated. Date and time format is ISO 8601 and UTC.


## Import

To import the AWS cloud account, use the ID as in the following example:

`$ terraform import vra_cloud_account_aws.new_aws 05956583-6488-4e7d-84c9-92a7b7219a15`
