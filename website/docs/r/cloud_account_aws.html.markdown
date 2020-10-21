---
layout: "vra"
page_title: "VMware vRealize Automation: vra_cloud_account_aws"
description: |-
  Provides a VMware vRA vra_cloud_account_aws resource.
---

# Resource: vra\_cloud\_account\_aws

Provides a VMware vRA vra_cloud_account_aws resource.

## Example Usages

**Create AWS cloud account:**

This is an example of how to create an AWS cloud account resource.

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

The following arguments are supported for an AWS cloud account resource:

* `access_key` - (Required) Access key id for AWS.

* `description` - (Optional) A human-friendly description.

* `name` - (Required) The name of this AWS cloud account.

* `regions` - (Optional) A set of region names that are enabled for this account.

* `secret_key` - (Required) Aws Secret Access Key

* `tags` - (Optional) A set of tag keys and optional values that to set on this cloud account.
example:[ { "key" : "vmware", "value": "provider" } ]
  * `key` - Tag’s key.
  * `value` - Tag’s value.


## Attribute Reference

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `id` - The id of this AWS cloud account.

* `links` - HATEOAS of the entity.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.
  
* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.


## Import

AWS cloud account can be imported using the id, e.g.

`$ terraform import vra_cloud_account_aws.new_aws 05956583-6488-4e7d-84c9-92a7b7219a15`
