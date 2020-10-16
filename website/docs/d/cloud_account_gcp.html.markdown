---
layout: "vra"
page_title: "VMware vRealize Automation: vra_cloud_account_gcp"
description: |-
    Provides a data lookup for vra_cloud_account_gcp.
---

# Data Source: vra\_cloud\_account\_gcp

Provides a VMware vRA vra_cloud_account_gcp data source.

## Example Usages

**GCP cloud account data source by its id:**

This is an example of how to create an GCP cloud account resource and read it as a data source using its id.
NOTE: The GCP cloud account resource need not be created through terraform.
To create an GCP cloud account, follow the resource GCP cloud account documentation:

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

data "vra_cloud_account_gcp" "this" {
  id = "vra_cloud_account_gcp.this.id"
}

```

**GCP cloud account data source by its name:**

This is an example of how to create an GCP cloud account resource and read it as a data source using its name.
NOTE: The GCP cloud account resource need not be created through terraform.
To create an GCP cloud account , follow the resource GCP cloud account documentation:

```hcl

resource "vra_cloud_account_gcp" "this" {
  name = "my-cloud-account-gcp"
  description = "test cloud account"
  subscription_id = "sample-subscription-id"
  tenant_id = "sample-tenant-id"
  application_id = "sample-application-id"
  application_key = "sample-application=key"
  regions = ["centralus"]
}

data "vra_cloud_account_gcp" "this" {
  name = "vra_cloud_account_gcp.this.name"
}

```



## Argument Reference

The following arguments are supported for an GCP cloud account data source:

* `id` - (Optional) The id of this GCP cloud account.

* `name` - (Optional) The name of this GCP cloud account.

## Attribute Reference

* `client_email` - GCP Client email.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `description` - A human-friendly description.

* `links` - HATEOAS of the entity.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `private_key_id` - GCP Private key ID.

* `project_id` - GCP Project ID.

* `regions` - A set of region names that are enabled for this account.

* `tags` - A set of tag keys and optional values that were set on this resource.
example:[ { "key" : "vmware", "value": "provider" } ]
  * `key` - Tag’s key.
  * `value` - Tag’s value.
  
* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.