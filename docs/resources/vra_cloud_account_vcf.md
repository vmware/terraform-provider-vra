---
page_title: "VMware Aria Automation: vra_cloud_account_vcf"
description: |-
    Creates a vra_cloud_account_vcf resource.
---

# Resource: vra_cloud_account_vcf

Creates a VMware Aria Automation VCF cloud account resource.

## Example Usages

The following example shows how to create a VCF cloud account resource.

```hcl
resource "vra_cloud_account_vcf" "this" {
  name                 = "tf-vra-cloud-account-vcf"
  description          = "tf test vcf cloud account"
  workload_domain_id   = var.workload_domain_id
  workload_domain_name = var.workload_domain_name

  vcenter_hostname = var.vcenter_hostname
  vcenter_password = var.vcenter_password
  vcenter_username = var.vcenter_username

  nsx_hostname = var.nsx_hostname
  nsx_password = var.nsx_password
  nsx_username = var.nsx_username

  dc_id                   = var.data_collector_id // Required for VMware Aria Automation Cloud, Optional for VMware Aria Automation on-prem
  sddc_manager_id         = var.sddc_manager_id
  accept_self_signed_cert = true

  enabled_regions {
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

Create your VCF cloud account resource with the following arguments:

* `accept_self_signed_cert` - (Optional) Accept self-signed certificate when connecting to the cloud account.

* `dc_id` - (Optional) Identifier of a data collector VM deployed in the on premise infrastructure.

* `description` - (Optional) Human-friendly description.

* `enabled_regions` - (Required) A set of region names that are enabled for the cloud account.

  * `external_region_id` - Unique identifier of the region on the provider side.

  * `name` - Name of the region on the provider side.

* `name` - (Required) Human-friendly name used as an identifier in APIs that support this option.

* `nsx_hostname` - (Required) IP address or FQDN of the NSX Manager server in the specified workload domain.

* `nsx_password` - (Required) Password used to authenticate to the NSX Manager in the specified workload domain.

* `nsx_username` - (Required) Username used to authenticate to the NSX Manager in the specified workload domain.

* `regions` - (Required) A set of region names that are enabled for the cloud account.

  > **Note**: Deprecated - please use `enabled_regions` instead.

* `sddc_manager_id` - (Optional) SDDC manager integration id.

* `tags` - (Optional) Set of tag keys and optional values to apply to the cloud account. Example: `[ { "key" : "vmware", "value": "provider" } ]`

* `vcenter_hostname` - (Required) IP address or FQDN of the vCenter Server in the specified workload domain.

* `vcenter_password` - (Required) Password used to authenticate to the vCenter Server in the specified workload domain.

* `vcenter_username` - (Required) Username used to authenticate to the vCenter Server in the specified workload domain.

* `workload_domain_id` - (Required) ID of the workload domain to add as VCF cloud account.

* `workload_domain_name` - (Required) Name of the workload domain to add as VCF cloud account.

## Attribute Reference

* `created_at` - Date when entity was created. Date and time format is ISO 8601 and UTC.

* `id` - ID of the VCF cloud account.

* `links` - Hypermedia as the Engine of Application State (HATEOAS) of entity.

* `org_id` - ID of organization that entity belongs to.

* `owner` - Email of entity owner.

* `updated_at` - Date when the entity was last updated. Date and time format is ISO 8601 and UTC.

## Import

To import the VCF cloud account, use the ID as in the following example:

`$ terraform import vra_cloud_account_vcf.new_vcf 05956583-6488-4e7d-84c9-92a7b7219a15`