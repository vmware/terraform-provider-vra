# Get Your Refresh Token for the VMware Aria Automation API

Before making a call to VMware Aria Automation, you request an API token that authenticates you for authorized API connections. The API token is also known as a "refresh token".

The Terraform provider for VMware Aria Automation accepts either a `refresh_token` or an `access_token`, but not both at the same time.

* Refresh token are valid for **90 days**, when using the API.
* Access tokens are valid for **8 hours**, but times out after **25 minutes** of inactivity.

## Procedures

### UI Procedure

The following procedure is only applicable to VMware Aria Automation Cloud.

1. Login to the VMware Cloud Services console using your credentials at <https://console.cloud.vmware.com>.

2. Once logged in, click the drop-down arrow by your name and select **My Account**.

3. On the `My Account` page, click the **API Tokens** tab.

4. In the `API Tokens` section, click **Generate Token**.

    a. On the `Generate a New API Token` page, enter a **Token Name**.

    b. Select a **Token TTL** (Time to Line). The default is 6 months in the UI.

    c. Under the `Define Scopes` section select the **Organization Owner** for the **Organization Roles**.

    d. Under the `Define Scopes` section expand **VMware Cloud Assembly** and select the **Cloud Assembly Administrator** in for the **Service Roles**.

    e. Under the `Define Scopes` section expand **VMware Service Broker** and select the **Service Broker Administrator** in for the **Service Roles**.

    f. (Optional) Under the `Email Preferences` section check the option to send expiration reminders.

    g. Click **Generate**.

    The **Token Generated** window displays a token with the name that you specified and the name of your organization.

5. Click **COPY**.

    Clicking **COPY** ensures that you capture the exact string.

    > Note: Once you click **Continue**, you will not be able to retrieve this token again.

6. Use the `refresh_token` in the Terraform provider configuration. For example:

    ```hcl
    provider "vra" {
      url           = "https://api.mgmt.cloud.vmware.com"
      refresh_token = "mx7w9**********************zB3UC"
      insecure      = false
    }
    ```

### API Procedure

To request a `refresh_token` using the API, you will need your user credentials:

* `username`
* `password`
* `domain` (optional)

In addition, you will need the fully qualified domain name (FQDN) of the endpoint associated with the identity access service.

* For VMware Aria Automation, this will be the fully qualified domain name of the VMware Aria Automation cluster VIP or appliance. For example, `cloud.example.com`.

* For VMware Aria Automation (SaaS) is available in multiple global regions. When making a API request to the service hosted in the United States, use `api.mgmt.cloud.vmware.com`. For organizations located outside of the United States, prefix the URL with the country abbreviation for your API endpoint as shown in the following:

  * Australia: `au.api.mgmt.cloud.vmware.com`
  * Brazil: `br.api.mgmt.cloud.vmware.com`
  * Canada: `ca.api.mgmt.cloud.vmware.com`
  * Germany: `de.api.mgmt.cloud.vmware.com`
  * Japan: `jp.api.mgmt.cloud.vmware.com`
  * Singapore: `sg.api.mgmt.cloud.vmware.com`
  * United Kingdom: `uk.api.mgmt.cloud.vmware.com`

You then pass a JSON body containing the credentials to the API.

  **Example**: JSON body with domain.

  ```json
  {
    "username":"john.doe",
    "password":"VMw@re1!",
    "domain":"example.com"
  }
  ```

  **Example**: JSON body without domain.

  ```json
  {
    "username":"john.doe",
    "password":"VMw@re1!"
  }
  ```

If successful, a JSON response will be returned with the value for the `refresh_token`.

### PowerShell Example

1. Set the variables:

    ```powershell

    $vraFqdn="cloud.example.com"

    $vraUsername="john.doe"

    $vraPassword="VMw@re1!"

    $vraDomain="example.com"

    $vraUrl="https://$vraFqdn/csp/gateway/am/api/login?access_token"

    $vraBody="{""username"":""$vraUsername"",""password"":""$vraPassword"",""domain"":""$vraDomain""}"
    ```

2. `POST` request to the API:

    ```powershell
    $vraResponse = Invoke-RestMethod -Method POST -ContentType "application/json" -URI $vraUrl -Body $vraBody
    ```

3. Get the `refresh_token`:

    ```powershell
    $vraResponse.refresh_token
    ```

    The `refresh_token` is returned.

    ```powershell
    mx7w9**********************zB3UC
    ```

4. Use the `refresh_token` in the Terraform provider configuration. For example:

    ```hcl
    provider "vra" {
      url           = "https://api.mgmt.cloud.vmware.com"
      refresh_token = "mx7w9**********************zB3UC"
      insecure      = false
    }
    ```

### Bash Example

1. Set the variables:

    ```shell
    vraFqdn=cloud@example.com

    vraUsername=john.doe

    vraPassword=VMw@re1!

    vraDomain=example.com

    vraUrl="https://"$vraFqdn"/csp/gateway/am/api/login?access_token"

    vraBody="{\"username\":\"$vraUsername\",\"password\":\"$vraPassword\",\"domain\":\"$vraDomain\"}"
    ```

2. `POST` request to the API:

    ```shell
    curl -k -X POST $vraUrl -H "Accept: application/json" -H "Content-Type: application/json" -s -d $vraBody
    ```

    The `refresh_token` is returned.

    ```shell
    {"refresh_token":"mx7w9**********************zB3UC"}
    ```

3. Use the `refresh_token` in the Terraform provider configuration. For example:

    ```hcl
    provider "vra" {
      url           = "https://api.mgmt.cloud.vmware.com"
      refresh_token = "mx7w9**********************zB3UC"
      insecure      = false
    }
    ```

## Scripts

Scripts for both PowerShell and Bash are included in the project repository in the `docs` directory. These scripts will prompt you for the values and return the `refresh_token`.

* PowerShell Script: [`get_token.ps1`](./get_token.ps1)
* Bash Script: [`get_token.sh`](./get_token.sh)

### PowerShell Script on Windows: `get_token.ps1`

```powershell
> ./get_token.ps1

Enter the FQDN for the VMware Aria Automation: cloud.example.com

Enter the username to authenticate with VMware Aria Automation: john.doe

Enter the password to authenticate with VMware Aria Automation: ********

Enter the domain or press enter to skip: example.com

Successfully connected to the endpoint for VMware Aria Automation services: cloud.example.com

Generating Refresh Token...

----------Refresh Token---------
mx7w9**********************zB3UC
--------------------------------

Saving environmental variables...

VRA_URL = https://cloud.example.com
VRA_REFRESH_TOKEN = mx7w9**********************zB3UC
```

### Bash Script on Linux or macOS: `get_token.sh`

```shell
$ ./get_token.sh

Enter the FQDN for the VMware Aria Automation:
cloud.example.com

Enter the username to authenticate with VMware Aria Automation:
john.doe

Enter the password to authenticate with VMware Aria Automation:
********

Enter the domain or press enter to skip:
example.com

Generating Refresh Token...

----------Refresh Token---------
mx7w9**********************zB3UC
--------------------------------

Environmental variables...

VRA_URL = https://cloud.example.com
VRA_REFRESH_TOKEN = mx7w9**********************zB3UC
```
