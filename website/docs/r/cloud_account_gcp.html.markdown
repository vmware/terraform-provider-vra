---
layout: "vra"
page_title: "VMware vRealize Automation: vra_cloud_account_gcp"
description: |-
    Provides a VMware vRA vra_cloud_account_gcp resource.
---

# Resource: vra\_cloud\_account\_gcp

Provides a VMware vRA vra_cloud_account_gcp resource.

## Example Usages

**Create GCP cloud account:**

This is an example of how to create a GCP cloud account resource.

```hcl

resource "vra_cloud_account_gcp" "this" {
  name           = "tf-vra-cloud-account-gcp"
  description    = "terraform test cloud account gcp"
  client_email   = "client_email"
  private_key_id = "private_key_id"
  private_key    = "private_key"
  project_id     = "project_id"
  regions        = ["us-west1", "us-west2"]

  tags {
    key   = "foo"
    value = "bar"
  }
}

```



## Argument Reference


The following arguments are supported for an GCP cloud account resource:

* `client_email` - (Required) GCP Client email.

* `description` - (Optional) A human-friendly description.

* `name` - (Required) The name of this GCP cloud account.

* `private_key` - (Required) GCP Private key.

* `private_key_id` - (Required)  GCP Private key ID.

* `project_id` - (Required) GCP Project ID.

* `tags` - (Optional) A set of tag keys and optional values that to set on this resource.
example:[ { "key" : "vmware", "value": "provider" } ]
  * `key` - Tag’s key.
  * `value` - Tag’s value.


## Attribute Reference

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `id` - The id of this GCP cloud account.

* `links` - HATEOAS of the entity.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `regions` - A set of region names that are enabled for this account.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.


## Import

GCP cloud account can be imported using the id, e.g.

`$ terraform import vra_cloud_account_gcp.new_gcp 05956583-6488-4e7d-84c9-92a7b7219a15`