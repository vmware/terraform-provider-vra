# cloud\_account\_vcf example

This is an example on how to setup a VCF cloud account in VMware Aria Automation (vRA).

## Getting Started

There are variables which need to be added to terraform.tfvars. The first are for connecting to the VMware Aria Automation (vRA) endpoint, and the remaining are credentials and identifiers for connecting to the VCF workload domain.

* `url` - The URL for the VMware Aria Automation (vRA) endpoint
* `refresh_token` - The refresh token for the vRA account
* `insecure` - Whether to skip TLS verification (`false` for vRA Cloud and `true` for vRA on-prem)
* `data_collector_name` - The name of the data collector (cloud proxy) already setup for the endpoint
* `nsx_hostname` - Hostname / IP address of the NSX Manager in the workload domain
* `nsx_password` - Password used to authenticate with NSX Manager
* `nsx_username` - Username used to authenticate with NSX Manager
* `regions` - List of region IDs to enable on the cloud account
* `sddc_manager_id` - SDDC manager integration id
* `vcenter_hostname` - Hostname / IP address of the vCenter in the workload domain
* `vcenter_password` - Password used to authenticate with vCenter
* `vcenter_username` - Username used to authenticate with vCenter
* `workload_domain_id` - ID of the workload domain to add as VCF cloud account
* `workload_domain_name` - Name of the workload domain to add as VCF cloud account

To facilitate adding these variables, a sample tfvars file can be copied first:
```shell
cp terraform.tfvars.sample terraform.tfvars
```

Once the information is added to `terraform.tfvars`, the cloud account can be created via:

```shell
terraform init
terraform apply
```

## Data source usage

Lookup by name:

```hcl
data "vra_cloud_account_vcf" "by_name" {
	name = vra_cloud_account_vcf.this.name
}
```

Lookup by id:

```hcl
data "vra_cloud_account_vcf" "by_id" {
	id = vra_cloud_account_vcf.this.id
}
```
