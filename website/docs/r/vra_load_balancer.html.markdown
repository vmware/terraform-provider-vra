---
layout: "vra"
page_title: "VMware vRealize Automation: vra_load_balancer"
description: |-
  Creates a vra_load_balancer resource.
---

# Resource: vra_load_balancer

Creates a VMware vRealize Automation load balancer resource.

## Example Usages

The following example shows how to create a load balancer resource.


```hcl
resource "vra_load_balancer" "this" {
  name        = "my-load-balancer"
  project_id  = vra_project.my-project.id
  description = "My Load Balancer"
  custom_properties = {
    "edgeClusterRouterStateLink"  = "/resources/routers/<uuid>"
    "tier0LogicalRouterStateLink" = "/resources/routers/<uuid>"
  }

  targets {
    machine_id = vra_machine.my_machine.id
  }

  nics {
    network_id = data.vra_network.my-network.id
  }

  routes {
    protocol        = "TCP"
    port            = "80"
    member_protocol = "TCP"
    member_port     = "80"
    health_check_configuration {
      protocol            = "TCP"
      port                = "80"
      interval_seconds    = 30
      timeout_seconds     = 10
      unhealthy_threshold = 2
      healthy_threshold   = 10
    }
  }
}
```

A block device resource supports the following arguments:

## Argument Reference

* `custom_properties` - (Optional) Additional custom properties that may be used to extend the machine.

* `deployment_id` - (Optional) The id of the deployment that is associated with this resource.

* `description` - (Optional) Describes machine within the scope of your organization and is not propagated to the cloud.

* `internet_facing` - (Optional) An Internet-facing load balancer has a publicly resolvable DNS name, so it can route requests from clients over the Internet to the instances that are registered with the load balancer.

* `name` - (Required) A human-friendly name used as an identifier in APIs that support this option.

* `nics` - (Required) A set of network interface specifications for this load balancer.

* `project_id` - (Required) The id of the project the current user belongs to.

* `routes` - (Required) The load balancer route configuration regarding ports and protocols.

    * `algorithm` - Algorithm employed for load balancing.

    * `algorithm_parameters` - Parameters need for load balancing algorithm.Use newline to separate multiple parameters.

    * `health_check_configuration` - Load balancer health check configuration.

        * `healthy_threshold` - Number of consecutive successful checks before considering a particular back-end instance as healthy.

        * `http_method` - HTTP or HTTPS method to use when sending a health check request.

        * `interval_seconds` - Interval (in seconds) at which the health checks will be performed.

        * `passive_monitor` - Enable passive monitor mode. This setting only applies to NSX-T.

        * `port` - Port on the back-end instance machine to use for the health check.

        * `protocol` - The protocol used for the health check.

        * `request_body` - Request body. Used by HTTP, HTTPS, TCP, UDP.

        * `response_body` - Expected response body. Used by HTTP, HTTPS, TCP, UDP.

        * `timeout_seconds` - Timeout (in seconds) to wait for a response from the back-end instance.

        * `unhealthy_threashold` - Number of consecutive check failures before considering a particular back-end instance as unhealthy.

        * `urlPath` - URL path on the back-end instance against which a request will be performed for the health check. Useful when the health check protocol is HTTP/HTTPS.

    * `member_port` - Member port where the traffic is routed to.

    * `member_protocol` - The protocol of the member traffic.

    * `port` - Port which the load balancer is listening to.

    * `protocol` - The protocol of the incoming load balancer requests.

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
example: `[ { "key" : "vmware.enumeration.type", "value": "nebs_block" } ]`
  * `key` - Tag’s key.
  * `value` - Tag’s value.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
