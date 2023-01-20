---
layout: "vra"
page_title: "VMware vRealize Automation: vra_network_profile"
description: |-
  Provides a data lookup for vra_network_profile.
---

# Data Source: vra_network_profile
## Example Usages
This is an example of how to create a network profile resource.

**Network profile data source by its id:**

```hcl
data "vra_network_profile" "this" {
  filter = "name eq '${vra_network_profile.this.name}'"
}
```

**Vra network profile data source filter by region id:**

```hcl
data "vra_network_profile" "this" {
  filter = "regionId eq '${data.vra_region.this.id}'"
}
```

A network profile data source supports the following arguments:

## Argument Reference

* `filter` - (Optional) Filter query string that is supported by vRA multi-cloud IaaS API. Example: `regionId eq '<regionId>' and cloudAccountId eq '<cloudAccountId>'`.

* `id` - (Optional) The id of the image profile instance.

* `isolated_external_fabric_network_id` - (Optional) The Id of the fabric network used for outbound access.

* `isolated_network_domain_id` - (Optional) The Id of the network domain used for creating isolated networks.

## Attributes Reference

* `custom_properties` - Additional properties that may be used to extend the Network Profile object that is produced from this specification. For isolationType security group, datastoreId identifies the Compute Resource Edge datastore. computeCluster and resourcePoolId identify the Compute Resource Edge cluster. For isolationType subnet, distributedLogicalRouterStateLink identifies the on-demand network distributed local router. onDemandNetworkIPAssignmentType identifies the on-demand network IP range assignment type static, dynamic, or mixed.

* `description` - A human-friendly description.

* `external_region_id` - The external regionId of the resource.

* `fabric_network_ids` - A list of fabric network Ids which are assigned to the network profile.
                         example: `[ "6543" ]`
* `isolated_network_cidr_prefix` - The CIDR prefix length to be used for the isolated networks that are created with the network profile.

* `isolated_network_domain_cidr` - CIDR of the isolation network domain.

* `isolation_type` - Specifies the isolation type e.g. none, subnet or security group

* `links` - HATEOAS of the entity

* `name` - A human-friendly name used as an identifier in APIs that support this option.

* `organization_id` - The id of the organization this entity belongs to.

* `owner` - Email of the user that owns the entity.

* `region_id` - The id of the region for which this profile is defined as in vRealize Automation(vRA).

* `security_group_ids` - A list of security group Ids which are assigned to the network profile.
                         example: `[ "6545" ]`

* `tags` - A set of tag keys and optional values that were set on this Network Profile.
           example: `[ { "key" : "ownedBy", "value": "Rainpole" } ]`

* `updated_at` - Date when the entity was last updated. The date is ISO 8601 and UTC.
