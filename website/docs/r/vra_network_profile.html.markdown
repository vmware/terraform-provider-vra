---
layout: "vra"
page_title: "VMware vRealize Automation: vra_network_profile"
description: |-
  Provides a data lookup for vra_network_profile.
---

# Resource: vra_network_profile
## Example Usages
This is an example of how to create a network profile resource.

**Network profile:**

```hcl
resource "vra_network_profile" "simple" {
  name        = "no-isolation"
  description = "Simple Network Profile with no isolation."
  region_id   = data.vra_region.this.id

  fabric_network_ids = [
    data.vra_fabric_network.subnet_1.id,
    data.vra_fabric_network.subnet_2.id
  ]

  isolation_type = "NONE"

  tags {
    key   = "foo"
    value = "bar"
  }
}
```

A network profile resource supports the following arguments:

## Argument Reference

* `custom_properties` - (Optional) Additional properties that may be used to extend the Network Profile object that is produced from this specification. For isolationType security group, datastoreId identifies the Compute Resource Edge datastore. computeCluster and resourcePoolId identify the Compute Resource Edge cluster. For isolationType subnet, distributedLogicalRouterStateLink identifies the on-demand network distributed local router. onDemandNetworkIPAssignmentType identifies the on-demand network IP range assignment type static, dynamic, or mixed.

* `description` - (Optional) A human-friendly description.

* `fabric_network_ids` - (Optional) A list of fabric network Ids which are assigned to the network profile.
                         example: `[ "6543" ]`

* `isolated_external_fabric_network_id` - (Optional) The id of the fabric network used for outbound access.

* `name` - (Required) A human-friendly name used as an identifier in APIs that support this option.

* `region_id` - (Required) The id of the region for which this profile is defined as in vRealize Automation(vRA).

## Attributes Reference

* * `cloud_account_id` - The ID of the cloud account this flavor profile belongs to.

* `created_at` - Date when  entity was created. Date and time format is ISO 8601 and UTC.

* `external_region_id` - The external regionId of the resource.

* `isolated_network_cidr_prefix` - The CIDR prefix length to be used for the isolated networks that are created with the network profile.

* `isolated_network_domain_cidr` - CIDR of the isolation network domain.

* `isolated_network_domain_id` - The id of the network domain used for creating isolated networks.

* `isolation_type` - Specifies the isolation type e.g. none, subnet or security group

* `links` - HATEOAS of the entity

* `org_id` - ID of organization that entity belongs to.

* `organization_id` - The id of the organization this entity belongs to. Deprecated, refer to org_id instead.

* `owner` - Email of the user that owns the entity.

* `security_group_ids` - A list of security group Ids which are assigned to the network profile.
                         example: `[ "6545" ]`

* `tags` - A set of tag keys and optional values that were set on this Network Profile.
           example: `[ { "key" : "ownedBy", "value": "Rainpole" } ]`

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
