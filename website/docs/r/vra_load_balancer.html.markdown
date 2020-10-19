---
layout: "vra"
page_title: "VMware vRealize Automation: vra_load_balancer"
description: |-
  Provides a VMware vRA vra_load_balancer resource.
---
# Resource: vra_load_balancer
## Example Usages

This is an example of how to read a load balancer resource.

```hcl
resource "vra_load_balancer" "my_load_balancer" {
    name = "my-lb-%d"
    project_id = vra_project.my-project.id
    description = "load balancer description"
    
    targets {
        machine_id = vra_machine.my_machine.id
    }

    nics {
        network_id = data.vra_network.my-network.id
    }

    routes {
        protocol = "TCP"
        port = "80"
        member_protocol = "TCP"
        member_port = "80"
        health_check_configuration = {
            protocol = "TCP"
            port = "80"
            interval_seconds = 30
            timeout_seconds = 10
            unhealthy_threshold = 2
            healthy_threshold = 10
        }
    }
}
```

A block device resource supports the following arguments:

## Argument Reference
* `custom_properties` - Additional custom properties that may be used to extend the machine.

* `deployment_id` - The id of the deployment that is associated with this resource.

* `description` - Describes machine within the scope of your organization and is not propagated to the cloud.

* `internet_facing` - An Internet-facing load balancer has a publicly resolvable DNS name, so it can route requests from clients over the Internet to the instances that are registered with the load balancer.

* `name` - A human-friendly name used as an identifier in APIs that support this option.

* `nics` - A set of network interface specifications for this load balancer.

* `project_id` - The id of the project the current user belongs to.

* `routes` - The load balancer route configuration regarding ports and protocols.

## Attribute Reference
* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `external_id` - External entity Id on the provider side.

* `external_region_id` - The external regionId of the resource.

* `external_zone_id` - The external regionId of the resource.

* `organization_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `links` - HATEOAS of the entity.

* `targets` - A list of links to target load balancer pool members. Links can be to either a machine or a machine's network interface.

* `address` - Primary address allocated or in use by this load balancer. The address could be an in the form of a publicly resolvable DNS name or an IP address.

* `tags` - A set of tag keys and optional values that were set on this resource instance.
example:[ { "key" : "vmware.enumeration.type", "value": "nebs_block" } ]
  * `key` - Tag’s key.
  * `value` - Tag’s value.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
