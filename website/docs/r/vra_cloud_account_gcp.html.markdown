---
layout: "vra"
page_title: "VMware vRealize Automation: vra_cloud_account_gcp"
description: |-
    Creates a vra_cloud_account_gcp resource.
---

# Resource: vra\_cloud\_account\_gcp

Creates a VMware vRealize Automation GCP cloud account resource.

## Example Usages

The following example shows how to create a GCP cloud account resource.

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

Create your GCP cloud account resource with the following arguments:

* `client_email` - (Required) GCP Client email.

* `description` - (Optional) Human-friendly description.

* `name` - (Required) Name of GCP cloud account.

* `regions` - (Optional) Set of region names enabled for the cloud account.

* `private_key` - (Required) GCP Private key.

* `private_key_id` - (Required) GCP Private key ID.

* `project_id` - (Required) GCP Project ID.

* `tags` - (Optional) Set of tag keys and values to apply to the cloud account.
Example: `[ { "key" : "vmware", "value": "provider" } ]`

## Attribute Reference

* `created_at` - Date when entity was created. Date and time format is ISO 8601 and UTC.

* `id` - ID of GCP cloud account.

* `links` - HATEOAS of entity.

* `org_id` - ID of organization that entity belongs to.

* `owner` - Email of entity owner.

* `updated_at` - Date when entity was last updated. Date and time format is ISO 8601 and UTC.


## Import

To import the GCP cloud account, use the ID as in the following example:

`$ terraform import vra_cloud_account_gcp.new_gcp 05956583-6488-4e7d-84c9-92a7b7219a15`
