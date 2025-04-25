# Policy examples

This is an example on how to create policies in VMware Aria Automation (vRA):

* Approval Policy: Request approval for catalog deployments from specified users
* Day2 Action Policy: Manage what actions are available for deployed resources
* IaaS Resource Policy: Manage IaaS resource lifecycle at namespace level
* Lease Policy: Automate the expiration and destruction of deployed catalog items

## Getting Started

There are variables which need to be added to terraform.tfvars for connecting to the VMware Aria Automation endpoint:

* `url` - The URL for the VMware Aria Automation (vRA) endpoint
* `refresh_token` - The refresh token (API token) for the VMware Aria Automation (vRA) user account

To facilitate adding these variables, a sample tfvars file can be copied first:
```shell
cp terraform.tfvars.sample terraform.tfvars
```

Once the information is added to `terraform.tfvars`, the policies can be created via:

```shell
terraform init
terraform apply
```
