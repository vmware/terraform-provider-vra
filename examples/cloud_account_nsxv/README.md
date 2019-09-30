# cloud\_account\_nsxv example

This is an example on how to setup a cloud account for NSX-V.

## Getting Started

There are variables which need to be added to terraform.tfvars. The first are for connecting to the vRealize Automation (vRA) endpoint, there are credentials for connecting to the NSX-V instance, and the name of the data collector already setup in vRA.

* `url` - The URL for the vRealize Automation (vRA) endpoint
* `refresh_token` - The refresh token for the vRA account
* `username` - username for NSX-V
* `password` - password for NSX-V
* `hostname` - hostname for NSX-V
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
