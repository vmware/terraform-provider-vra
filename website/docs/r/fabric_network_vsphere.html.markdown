---
layout: "vra"
page_title: "VMware vRealize Automation: fabric_network_vsphere"
description: |-
  Provides a VMware vRA fabric_network_vsphere resource.
---

# Resource: fabric_network_vsphere

Provides a VMware vRA fabric_network_vsphere resource.

## Example Usages

You cannot create a vSphere fabric network resource, however you can import using the command specified in the import section below.
Once a resource is imported, you can update it as shown below:

```hcl

resource "fabric_network_vsphere" "simple" {
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
* `filter` - Filter query string that is supported by vRA multi-cloud IaaS API. Only one of 'filter' or 'id' must be specified.

* `id` - The id of the vRA fabric network.  Only one of 'filter' or 'id'  must be specified.

## Attribute Reference
* `cloud_account_ids` - Set of ids of the cloud accounts this entity belongs to.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `cidr` - Network CIDR to be used.

* `is_public` - Indicates whether the sub-network supports public IP assignment.

* `is_default` - Indicates whether this is the default subnet for the zone.

* `description` - State object representing a network on a external cloud provider.

* `external_id` - External entity Id on the provider side.

* `external_region_id` - The id of the region for which this network is defined.

* `links` - HATEOAS of the entity

* `name` - A human-friendly name used as an identifier in APIs that support this option. 

* `organization_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `tags` -  A set of tag keys and optional values that were set on this resource.
                       example:[ { "key" : "ownedBy", "value": "Rainpole" } ]

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.


## Import

vSphere fabric network resource can be imported using the id, e.g.

`$ terraform import fabric_network_vsphere.new_fabric_network_vsphere 05956583-6488-4e7d-84c9-92a7b7219a15`