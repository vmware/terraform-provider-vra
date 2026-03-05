---
page_title: "VMware Aria Automation: vra_cloud_account_vcf"
description: |-
    Provides a data lookup for vra_cloud_account_vcf.
---

# Data Source: vra_cloud_account_vcf

Provides a vra_cloud_account_vcf data source.

## Example Usages

**VCF cloud account data source by its id:**

This is an example of how to read the cloud account data source using its id.

```hcl

data "vra_cloud_account_vcf" "this" {
  id = var.vra_cloud_account_vcf_id
}
```

**VCF cloud account data source by its name:**

This is an example of how to read the cloud account data source using its name.

```hcl

data "vra_cloud_account_vcf" "this" {
  name = var.vra_cloud_account_vcf_name
}
```

## Argument Reference

The following arguments are supported for an VCF cloud account data source:

* `id` - (Optional) The id of this VCF cloud account.

* `name` - (Optional) The name of this VCF cloud account.

## Attribute Reference

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `dc_id` - Identifier of a data collector vm deployed in the on premise infrastructure.

* `description` - A human-friendly description.

* `enabled_regions` - A set of regions that are enabled for this cloud account.

  * `external_region_id` - Unique identifier of the region on the provider side.

  * `id` - Unique identifier of the region.

  * `name` - Name of the region on the provider side.

* `links` - Hypermedia as the Engine of Application State (HATEOAS) of the entity.

* `nsx_hostname` - IP address or FQDN of the NSX Manager Server in the specified workload domain.

* `nsx_username` - Username used to authenticate to the NSX Manager in the specified workload domain.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `regions` - A set of region IDs that are enabled for this account.

* `sddc_manager_id` - SDDC manager integration id.

* `tags` - A set of tag keys and optional values that were set on this resource. Example: `[ { "key" : "vmware", "value": "provider" } ]`

  * `key` - Tag’s key.

  * `value` - Tag’s value.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `vcenter_hostname` - IP address or FQDN of the vCenter Server in the specified workload domain.

* `vcenter_username` - Username used to authenticate to the vCenter Server in the specified workload domain.

* `workload_domain_id` - Id of the workload domain.

* `workload_domain_name` - Name of the workload domain.
