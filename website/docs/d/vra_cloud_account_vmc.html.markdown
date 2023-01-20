---
layout: "vra"
page_title: "VMware vRealize Automation: vra_cloud_account_vmc"
description: |-
    Provides a data lookup for vra_cloud_account_vmc.
---

# Data Source: vra\_cloud\_account\_vmc

Provides a VMware vRA vra_cloud_account_vmc data source.

## Example Usages

**VMC cloud account data source by its id:**

This is an example of how to read the cloud account data source using its id.

```hcl

data "vra_cloud_account_vmc" "this" {
  id = var.vra_cloud_account_vmc_id
}

```

**vmc cloud account data source by its name:**

This is an example of how to read the cloud account data source using its name.

```hcl

data "vra_cloud_account_vmc" "this" {
  name = var.vra_cloud_account_vmc_name
}

```



## Argument Reference

The following arguments are supported for an vmc cloud account data source:

* `id` - (Optional) The id of this vmc cloud account.

* `name` - (Optional) The name of this vmc cloud account.

## Attribute Reference

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `dc_id` - Identifier of a data collector vm deployed in the on premise infrastructure. Refer to the data-collector API to create or list data collector.

* `description` - A human-friendly description.

* `links` - HATEOAS of the entity.

* `nsx_hostname` - The IP address of the NSX Manager server in the specified SDDC / FQDN.

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `regions` - A set of region names that are enabled for this account.

* `sddc_name` - Identifier of the on-premise SDDC to be used by this cloud account. Note that NSX-V SDDCs are not supported.

* `tags` - A set of tag keys and optional values that were set on this resource.
example: `[ { "key" : "vmware", "value": "provider" } ]`
  * `key` - Tag’s key.
  * `value` - Tag’s value.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

* `vcenter_hostname` - The IP address or FQDN of the vCenter Server in the specified SDDC. The cloud proxy belongs on this vCenter.

* `vcenter_username` - vCenter user name for the specified SDDC. The specified user requires CloudAdmin credentials. The user does not require CloudGlobalAdmin credentials.
