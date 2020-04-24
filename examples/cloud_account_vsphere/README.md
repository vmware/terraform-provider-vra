# cloud\_account\_vsphere example

This is an example on how to setup a cloud account for vSphere along with zones and a project.

## Getting Started

There are variables which need to be added to terraform.tfvars. The first are for connecting to the vRealize Automation (vRA) endpoint, there are credentials for connecting to the vCenter instance, and the name of the data collector already setup in vRA.

* `url` - The URL for the vRealize Automation (vRA) endpoint
* `refresh_token` - The refresh token for the vRA account
* `username` - username for vCenter server
* `password` - password for vCenter server
* `hostname` - hostname for vCenter server
* `nsxt_username` - username for NSX-T
* `nsxt_password` - password for NSX-T
* `nsxt_hostname` - hostname for NSX-T
* `datacollector` - the name for the data collector already setup for the vCenter account

To facilitate adding these variables, a sample tfvars file can be copied first:
```shell
cp terraform.tfvars.sample terraform.tfvars
```

Once the information is added to `terraform.tfvars`, the cloud account can be brought up via:

```shell
terraform init
terraform apply
```
