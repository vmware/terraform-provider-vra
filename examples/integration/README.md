# integration example

This is an example on how to create integrations in VMware vRealize Automation (vRA).

Integrations enable you to add external systems like VMware vRealize Orchestrator, configuration management and other external systems such as GitHub, Ansible, Puppet, and external IPAM providers such as Infoblox.

## Getting Started

There are variables which need to be added to terraform.tfvars. The first are for connecting to the VMware vRealize Automation endpoint:

* `url` - The URL for the vRealize Automation (vRA) endpoint
* `refresh_token` - The refresh token (API token) for the vRA user account

To create an Active Directory integration, you will need the following variables:

* `ad_server` - The LDAP host / IP.
* `ad_endpoint_id` - The id of the runtime environment.
* `ad_user` - The LDAP user.
* `ad_password` - The LDAP password.
* `ad_default_ou` - The base DN.

To create a Github integration, you will need the following variables:

* `github_token` - The GitHub token.

To create a SaltStack integration, you will need the following variables:

* `saltstack_hostname` - The hostname of the SaltStack Config server.
* `saltstack_endpoint_id` - The id of the runtime environment.
* `saltstack_username` - The username for the SaltStack Config server.
* `saltstack_password` - The password for the SaltStack Config server.

To facilitate adding these variables, a sample tfvars file can be copied first:
```shell
cp terraform.tfvars.sample terraform.tfvars
```

Once the information is added to `terraform.tfvars`, the integrations can be created via:

```shell
terraform init
terraform apply
```
