---
layout: "vra"
page_title: "VMware vRealize Automation: vra_network_ip_range"
description: |-
  Creates a network_ip_range resource.
---

# Resource: vra_network_ip_range

Creates a VMware vRealize Automation network_ip_range resource.

## Example Usages

**Create vRA Network IP range resource:**

This is an example of how to create a vRA Network IP range resource.

```hcl

resource "vra_network_ip_range" "this" {
  name               = "example-ip-range"
  description        = "Internal Network IP Range Example"
  start_ip_address   = var.start_ip
  end_ip_address     = var.end_ip
  ip_version         = var.ip_version
  fabric_network_ids = [data.fabric_network.subnet.id]

  tags {
    key   = "foo"
    value = "bar"
  }
}

```


## Argument Reference

The following arguments are supported for a vRA Network IP range resource:

* `description` - (Optional) A human-friendly description.

* `end_ip_address` - (Required) End IP address of the range.

* `fabric_network_ids` - (Optional) The Ids of the fabric networks.

* `ip_version` - (Required) IP address version: IPv4 or IPv6.

* `name` - (Required) The name of the network IP range.

* `start_ip_address` - (Required) Start IP address of the range.

* `tags` -  (Optional) Set of tag keys and values to apply to the resource.
            Example: `[ { "key" : "vmware", "value": "provider" } ]`

## Attribute Reference

* `created_at` - Date when the entity was created. The date is in ISO 8601 and UTC.

* `external_id` - External entity Id on the provider side.

* `id` - ID of the network IP range

* `links` - HATEOAS of the entity

* `org_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.

## Import

To import the vRA Network IP range, use the ID as in the following example:

`$ terraform import network_ip_range.new_ip_range 05956583-6488-4e7d-84c9-92a7b7219a15`
