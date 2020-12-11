---
layout: "vra"
page_title: "VMware vRealize Automation: network_ip_range"
description: |-
  Creates a network_ip_range resource.
---

# Resource: network_ip_range

Creates a VMware vRealize Automation network_ip_range resource.

## Example Usages

**Create vRA Network IP range resource:**

This is an example of how to create a vRA Network IP range resource.

```hcl

resource "network_ip_range" "this" {
  name              = "example-ip-range"
  description       = "Internal Network IP Range Example"
  start_ip_address  = var.start_ip
  end_ip_address    = var.end_ip
  ip_version        = var.ip_version
  fabric_network_id = data.fabric_network.subnet.id

  tags {
    key   = "foo"
    value = "bar"
  }
}

```


## Argument Reference

The following arguments are supported for a vRA Network IP range resource:

* `description` - State object representing a network on a external cloud provider.

* `end_ip_address` - End IP address of the range.

* `fabric_network_id` - Fabric network Id.

* `id` - ID of the network IP range. 

* `ip_version` - IP address version: IPv4 or IPv6. Default: IPv4.

* `name` - A human-friendly name used as an identifier in APIs that support this option.

* `start_ip_address` - Start IP address of the range.

## Attribute Reference

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `external_id` - External entity Id on the provider side.

* `links` - HATEOAS of the entity

* `org_id` - ID of organization that entity belongs to.

* `owner` - Email of the user that owns the entity.

* `tags` -  Set of tag keys and values to apply to the resource.
            Example:[ { "key" : "vmware", "value": "provider" } ]

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

## Import

To import the vRA Network IP range, use the ID as in the following example:

`$ terraform import network_ip_range.new_ip_range 05956583-6488-4e7d-84c9-92a7b7219a15`