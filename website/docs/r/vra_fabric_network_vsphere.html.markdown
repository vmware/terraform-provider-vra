---
layout: "vra"
page_title: "VMware vRealize Automation: vra_fabric_network_vsphere"
description: |-
  Updates a fabric_network_vsphere resource.
---

# Resource: vra_fabric_network_vsphere

Updates a VMware vRealize Automation fabric_network_vsphere resource.

## Example Usages

You cannot create a vSphere fabric network resource, however you can import using the command specified in the import section below.
Once a resource is imported, you can update it as shown below:

```hcl

resource "vra_fabric_network_vsphere" "simple" {
  cidr            = var.cidr
  default_gateway = var.gateway
  domain          = var.domain
  tags {
    key   = "foo"
    value = "bar"
  }
}

```
## Argument Reference

* `cidr` - Network CIDR to be used.

* `default_gateway` - IPv4 default gateway to be used.

* `default_ipv6_gateway` - IPv6 default gateway to be used.

* `dns_search_domains` - List of dns search domains for the vSphere network.

* `dns_server_addresses` - A human-friendly name used as an identifier in APIs that support this option.

* `domain` - Domain for the vSphere network.

* `ipv6_cidr` -  Network IPv6 CIDR to be used.

* `is_default` - Indicates whether this is the default subnet for the zone.

* `is_public` - Indicates whether the sub-network supports public IP assignment.

## Attribute Reference
* `cloud_account_ids` - Set of ids of the cloud accounts this entity belongs to.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `external_id` - External entity Id on the provider side.

* `external_region_id` - The id of the region for which this network is defined.

* `id` - ID of the vRA fabric network.

* `links` - HATEOAS of the entity

* `name` - A human-friendly name used as an identifier in APIs that support this option.

* `organization_id` - ID of organization that entity belongs to.

* `tags` -  Set of tag keys and values to apply to the resource.
            Example: `[ { "key" : "vmware", "value": "provider" } ]`

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.


## Import

To import the vSphere fabric network resource, use the ID as in the following example:

`$ terraform import fabric_network_vsphere.new_fabric_network_vsphere 05956583-6488-4e7d-84c9-92a7b7219a15`
