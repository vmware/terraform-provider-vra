<#
    .NOTES
    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
    WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
    COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
    OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

    .SYNOPSIS
        Generates and returns a `refresh_token` from VMware Aria Automation for use by the Terraform provider.

    .DESCRIPTION
        The Request-vRARefreshToken function connects to the specified VMware Aria Automation endpoint and obtains a refresh token that is needed the Terraform provider.
            terraform {
                required_providers {
                    vra = {
                        source  = "vmware/vra"
                        version = ">= x.y.z"
                }
            }
                required_version = ">= 0.13"
            }
            provider "vra" {
                url           = "https://cloud.example.com"
                refresh_token = "mx7w9**********************zB3UC"
                insecure      = false
            }

    .PARAMETER FQDN
        (string) FQDN of the vVMware Aria Automation instance.

    .PARAMETER username
        (string) Username used to connect to the VMware Aria Automation instance.

    .PARAMETER password
        (string) Password used to authenticate to the VMware Aria Automation instance.

    .PARAMETER domain
        (string) Authentication domain configured in the VMware Aria Automation instance.

    .PARAMETER skipCertValidation
        (switch) Skip the certificate validation when connecting to the VMware Aria Automation instance.

    .EXAMPLE
        get_token.ps1
        Enter the FQDN for the VMware Aria Automation: cloud.example.com
        Enter the username to authenticate with VMware Aria Automation: john.doe
        Enter the password to authenticate with VMware Aria Automation: ********
        Enter the domain or press enter to skip: example.com
        Successfully connected to endpoint for VMware Aria Automation: cloud.example.com
        Generating token...

        ---------Refresh Token---------
        N5WD************************kfgm
        -------------------------------

        Saving environmental variables...

        VRA_URL = https://cloud.example.com
        VRA_REFRESH_TOKEN = N5WD************************kfgm

    .EXAMPLE
        get_token.ps1 -fqdn cloud.example.com -username john.doe -password ****** -domain example.com

        Successfully connected to endpoint for VMware Aria Automation: cloud.example.com
        Generating token...

        ---------Refresh Token---------
        N5WD************************kfgm
        -------------------------------

        Saving environmental variables...

        VRA_URL = https://cloud.example.com
        VRA_REFRESH_TOKEN = N5WD************************kfgm

    .EXAMPLE
        get_token.ps1 -FQDN cloud.example.com -Username john.doe -Domain example.com

        Enter the password to authenticate with VMware Aria Automation: *****************
        Successfully connected to endpoint for VMware Aria Automation: cloud.example.com
        Generating token...

        ---------Refresh Token---------
        N5WD************************kfgm
        -------------------------------

        Saving environmental variables...

        VRA_URL = https://cloud.example.com
        VRA_REFRESH_TOKEN = N5WD************************kfgm

    .INPUTS
    None

    .OUTPUTS
    Text

    .LINK
    https://github.com/vmware/terraform-provider-vra

#>
[CmdletBinding()]
Param (
    [Parameter ()]
    [ValidateNotNullOrEmpty()]
    [Alias('Server')]
    [String]$FQDN = (Read-Host "Enter the FQDN for the VMware Aria Automation"),
    [Parameter ()]
    [ValidateNotNullOrEmpty()]
    [Alias('User')]
    [String]$Username = (Read-Host "Enter the username to authenticate with VMware Aria Automation"),
    [Parameter ()]
    [ValidateNotNullOrEmpty()]
    [Alias('Pass')]
    [string]$Password = (Read-Host -MaskInput -Prompt "Enter the password to authenticate with VMware Aria Automation"),
    [Parameter ()] [String]$domain = (Read-Host -Prompt "Enter the domain or press enter to skip"),
    [Parameter ()] [switch]$skipCertValidation = $false
)


Function Set-BasicAuthHeader {
    $base64AuthInfo = [Convert]::ToBase64String([Text.Encoding]::ASCII.GetBytes(("{0}:{1}" -f $username, $password))) # Create Basic Authentication Encoded Credentials
    $headers = @{"Accept" = "application/json" }
    $headers.Add("Authorization", "Basic $base64AuthInfo")
    $headers.Add("Content-Type", "application/json")
    $headers
}
Function Request-vRARefreshToken {
    Param (
        [String]$fqdn,
        [String]$username,
        [String]$password,
        [String]$domain,
        [switch]$skipCertValidation
    )

    if ($skipCertValidation) {
        add-type @"
        using System.Net;
        using System.Security.Cryptography.X509Certificates;
        public class TrustAllCertsPolicy : ICertificatePolicy {
            public bool CheckValidationResult(
            ServicePoint srvPoint, X509Certificate certificate,
            WebRequest request, int certificateProblem) {
            return true;
            }
        }
"@
        $allProtocols = [System.Net.SecurityProtocolType]'Ssl3,Tls,Tls11,Tls12'
        [System.Net.ServicePointManager]::SecurityProtocol = $allProtocols
        [System.Net.ServicePointManager]::CertificatePolicy = New-Object TrustAllCertsPolicy
    }

    $ariaAutoBasicHeaders = Set-BasicAuthHeader $username $password
    $Global:ariaAutoFqdn = $fqdn

    Try {
        $uri = "https://$ariaAutoFqdn/csp/gateway/am/api/login?access_token"
        if ($PsBoundParameters.ContainsKey("domain")) {
            $body = "{ ""username"":""$username"",""password"":""$password"",""domain"":""$tenant""}"
        } else {
            $body = "{ ""username"":""$username"",""password"":""$password""}"
        }

        if ($PSEdition -eq 'Core') {
            $ariaAutoResponse = Invoke-WebRequest -Method POST -Uri $uri -Headers $ariaAutoBasicHeaders -Body $body -SkipCertificateCheck # PS Core has -SkipCertificateCheck implemented, PowerShell 5.x does not.
        } else {
            $ariaAutoResponse = Invoke-WebRequest -Method POST -Uri $uri -Headers $ariaAutoBasicHeaders -Body $body
        }

        if ($ariaAutoResponse.StatusCode -eq 200) {
            $Global:ariaAutoHeaders = New-Object "System.Collections.Generic.Dictionary[[String],[String]]"
            $ariaAutoHeaders.Add("Accept", "application/json")
            $ariaAutoHeaders.Add("Content-Type", "application/json")
            $ariaAutoHeaders.Add("Authorization", "Bearer " + $ariaAutoResponse.Headers.'Csp-Auth-Token')
            Write-Output "Successfully connected to endpoint for VMware Aria Automation: $ariaAutoFqdn"
            Write-Output "Generating token..."
            Write-Output "`n---------Refresh Token---------"
            ((Select-String -InputObject $ariaAutoResponse -Pattern '"refresh_token":') -Split ('"'))[3]
            Write-Output "-------------------------------`n"
            Write-Output "Saving environmental variables...`n"
            $ENV:VRA_URL = "https://$vRAFqdn"
            $ENV:VRA_REFRESH_TOKEN = ((Select-String -InputObject $ariaAutoResponse -Pattern '"refresh_token":') -Split ('"'))[3]
            Write-Output "VRA_URL = $ENV:VRA_URL"
            Write-Output "VRA_REFRESH_TOKEN = $ENV:VRA_REFRESH_TOKEN"
        }
    } Catch {
        Write-Error $_.Exception.Message
    }
}

# Execute Functions
Request-vRARefreshToken -fqdn $fqdn -username $username -password $password -domain $domain -skipCertValidation $skipCertValidation
