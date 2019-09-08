# cloud\_account\_vmc example

This is an example on how to setup a VMware Cloud on AWS (VMC) cloud account in vRealize Automation(vRA) along with regions.

## Getting Started

There are variables which need to be added to terraform.tfvars. The first are for connecting to the vRealize Automation (vRA) endpoint, and the remaining are credentials for connecting to the VMC cloud account.

* `url` - The URL for the vRealize Automation (vRA) endpoint
* `refresh_token` - The refresh token for the vRA account
* `api_token` - The VMware Cloud on AWS API token used to access your SDDC environment
* `data_collector_name` - The name for the data collector (aka Cloud proxy) already setup for the vCenter account
* `nsx_hostname` - The hostname / IP address of the NSX Manager server in the specified SDDC
* `sddc_name` - SDDC name
* `vcenter_hostname` - The hostname / IP address of the vCenter
* `vcenter_password` - The password to use in combination with username to connect to vCenter
* `vcenter_username` - The username to use in combination with password to connect to vCenter

To facilitate adding these variables, a sample tfvars file can be copied first:
```shell
cp terraform.tfvars.sample terraform.tfvars
```

Once the information is added to `terraform.tfvars`, the cloud account can be brought up via:

```shell
terraform init
terraform apply
```
