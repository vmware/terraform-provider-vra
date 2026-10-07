# Get Your Refresh Token for the VMware Cloud Foundation Automation API

Before making a call to VMware Cloud Foundation (VCF) Automation, you request an API token that authenticates you for authorized API connections. The API token is also known as a "refresh token".

The Terraform provider accepts either a `refresh_token` or an `access_token`, but not both at the same time. When using VCF Automation, the `organization` is also required.

> Note: For VMware Aria Automation, see [Get Your Refresh Token for the VMware Aria Automation API](./refresh_token_vra.md).

## Procedures

### UI Procedure

The following procedure is applicable to both **VM Apps** and **All Apps** organizations.

1. Open the VCF Automation UI (`https://<vcfa-fqdn>`) and change the organization.

2. Log in with the organization administrator credentials.

3. Click the organization name (top-right) and select **My Account** under **User Settings**.

4. Click the **API Tokens** tab and click **New**.

5. Enter a token name and click **Create**.

6. Click **Copy**.

    > Note: Save the token now. You will not be able to retrieve it again.

7. Use the token as the `refresh_token` in the Terraform provider configuration. For example:

    ```hcl
    provider "vra" {
      url           = "https://vcfa.example.com"
      organization  = "my-org"
      refresh_token = "mx7w9**********************zB3UC"
      insecure      = false
    }
    ```

Reference: [Getting a Refresh Token for the VM Apps Tenant](https://techdocs.broadcom.com/us/en/vmware-cis/vcf/vcf-9-0-and-later/9-1/administration-sdks-cli-and-tools/about-the-vcf-automation-api/what-are-the-automation-apis-and-how-do-i-use-them_1/getting-your-authentication-token/getting-a-refresh-token-for-the-vm-apps-tenant.html).

### Script Procedure

The Bash script [`get_vcfa_token.sh`](./get_vcfa_token.sh) is included in the project repository in the `docs` directory. It prompts you for the values, requests the tokens from the VCF Automation API, and exports the results as environment variables. It works for both **VM Apps** and **All Apps** organizations.

#### Prerequisites

* The fully qualified domain name (FQDN) of your VCF Automation endpoint. For example, `vcfa.example.com`.
* The `username` and `password` of a user in the organization.
* The name of the `organization`.
* The organization type: **VM Apps** or **All Apps**.
* The [`jq`](https://jqlang.github.io/jq/) utility.

#### Usage

Run the script with `source` so the environment variables persist in your shell:

```shell
$ source ./get_vcfa_token.sh

Enter the VCF Automation FQDN:
vcfa.example.com

Enter the username to authenticate with VCF Automation:
john.doe

Enter the password to authenticate with VCF Automation:
********

Enter the organization:
my-org

Select the organization type:
  1) VM Apps
  2) All Apps
Enter 1 or 2: 2

Authenticating User...

Registering OAuth client...

Generating Refresh Token...

Environmental variables...

VRA_URL = https://vcfa.example.com
VCFA_ORGANIZATION = my-org
VCFA_APP_TYPE = All Apps
VRA_CLIENT_NAME = terraform-a1b2c3d4
VRA_CLIENT_ID = <client id>
VRA_ACCESS_TOKEN is set
VRA_REFRESH_TOKEN is set
```

The script sets the following environment variables:

| Variable | Description |
| --- | --- |
| `VRA_URL` | The URL of the VCF Automation endpoint. |
| `VCFA_ORGANIZATION` | The name of the organization. |
| `VRA_CLIENT_NAME` | The name of the registered OAuth client. |
| `VRA_CLIENT_ID` | The ID of the registered OAuth client. |
| `VRA_ACCESS_TOKEN` | The access token. |
| `VRA_REFRESH_TOKEN` | The refresh token. |

You can skip prompts by setting the `fqdn`, `username`, `password`, `organization` and `app_type` shell variables before running the script. To skip TLS verification (not recommended), set `VRA_INSECURE=true`.

Because the provider reads `VRA_URL`, `VRA_REFRESH_TOKEN` and `VCFA_ORGANIZATION` from the environment, the provider block can be left minimal:

```hcl
provider "vra" {}
```

Alternatively, set the values explicitly:

```hcl
provider "vra" {
  url           = "https://vcfa.example.com"
  organization  = "my-org"
  refresh_token = "mx7w9**********************zB3UC"
  insecure      = false
}
```
