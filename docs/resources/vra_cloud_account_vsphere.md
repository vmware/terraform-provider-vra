---
page_title: "VMware Aria Automation: vra_cloud_account_vsphere"
description: |-
    Creates a vra_cloud_account_vsphere resource.
---

# Resource: vra_cloud_account_vsphere

Creates a VMware Aria Automation vSphere cloud account resource.

## Example Usages

The following example shows how to create a vSphere cloud account resource.

```hcl
resource "vra_cloud_account_vsphere" "this" {
  name                         = "tf-vSphere-account"
  description                  = "foobar"
  username                     = var.username
  password                     = var.password
  hostname                     = var.hostname
  dc_id                        = var.vra_data_collector_id // Required for VMware Aria Automation Cloud, Optional for VMware Aria Automation 8.X
  associated_cloud_account_ids = [var.vra_cloud_account_nsxt_id]
  accept_self_signed_cert      = true

  enabled_regions  {
    external_region_id = var.region_external_id
    name               = var.region_name
  }

  tags {
    key   = "foo"
    value = "bar"
  }
}
```

## Argument Reference

Create your vSphere cloud account resource with the following arguments:

* `accept_self_signed_cert` - (Optional) Accept self-signed certificate when connecting to the cloud account.

* `dc_id` - (Optional) Identifier of a data collector VM deployed in the on premise infrastructure.

* `description` - (Optional) Human-friendly description.

* `enabled_regions` - (Required) A set of region names that are enabled for the cloud account.

  * `external_region_id` - Unique identifier of the region on the provider side.

  * `name` - Name of the region on the provider side.

* `hostname` - (Required) IP address or FQDN of the vCenter Server. The cloud proxy belongs on this vCenter.

* `name` - (Optional) Name of the vSphere cloud account.

* `password` - (Required) Password used to authenticate to the cloud account.

* `regions` - (Required) A set of region names that are enabled for the cloud account.

  > **Note**:  Deprecated - please use `enabled_regions` instead.

* `tags` - (Optional) A set of tag keys and optional values to apply to the cloud account. Example: `[ { "key" : "vmware", "value": "provider" } ]`

* `username` - (Required) vSphere username used to authenticate to the cloud account.

## Attribute Reference

* `associated_cloud_account_ids` - Cloud accounts associated with the cloud account.

* `created_at` - Date when  entity was created. Date and time format is ISO 8601 and UTC.

* `id` - (Optional) ID of the vSphere cloud account.

* `links` - Hypermedia as the Engine of Application State (HATEOAS) of entity.

* `org_id` - ID of organization that entity belongs to.

* `owner` - Email of entity owner.

* `updated_at` - Date when the entity was last updated. Date and time format is ISO 8601 and UTC.

## Import

To import the vSphere cloud account, use the ID as in the following example:

`$ terraform import vra_cloud_account_vsphere.new_vsphere 05956583-6488-4e7d-84c9-92a7b7219a15`
