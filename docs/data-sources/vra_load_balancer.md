---
page_title: "VMware vRealize Automation: vra_load_balancer"
description: |-
  Provides a data source for a load balancer resource.
---

# Data Source: vra_load_balancer

Provides a data source to retrieve a `vra_load_balancer`.

## Example Usage

This is an example of how to get a load balancer resource by its id.

```hcl
data "vra_load_balancer" "this" {
  id = "load-balancer-id"
}

output "load_balancer_name" {
  value = data.vra_load_balancer.this.name
}

## Argument Reference

* `id` - (Required) The id of the load balancer.

## Attribute Reference

* `address` - Primary address allocated or in use by this load balancer. The address could be an in the form of a publicly resolvable DNS name or an IP address.

* `created_at` - Date when the entity was created. The date is in ISO 6801 and UTC.

* `custom_properties` - A set of custom properties that were set on this resource instance. This is a key-value pair where the key is a string and the value can be any string.

* `description` - Description of the load balancer.

* `deployment_id` - The id of the deployment that is associated with this resource.

* `external_id` - External entity Id on the provider side.

* `external_region_id` - The external regionId of the resource.

* `external_zone_id` - The external regionId of the resource.

* `links` - HATEOAS of the entity.

* `name` - Name of the load balancer.

* `organization_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `project_id` - The id of the project that is associated with this resource.

* `routes` - The load balancer route configuration regarding ports and protocols.

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

* `targets` - A list of links to target load balancer pool members. Links can be to either a machine or a machine's network interface.

* `tags` - A set of tag keys and optional values that were set on this resource instance.
  * `key` - Tag’s key.
  * `value` - Tag’s value.

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
