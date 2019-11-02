#Script to generate an API refresh token for accessing vRA8/CAS. This is needed for
# the terraform provider to connect successfully
param(
    [Parameter(HelpMessage= "Username to connect to vRA with")][string]$vRAUser,
    [Parameter(HelpMessage= "Password to connect to vRA with")][string]$vRApassword,
    [Parameter(HelpMessage= "User's Domain connect to vRA with")][string]$vRAdomain,
    [Parameter(HelpMessage= "vRA/identity server hostname/fqdn")][string]$vRAServer
)

if ($PSBoundParameters.Keys.Contains("vRAUser")) { 
    Write-Host "Found value for vRAUser param: $vRAUser"
} else {
    $vRAUser = Read-Host -Prompt "Enter a username to connect to vRA with"
}

if ($PSBoundParameters.Keys.Contains("vRAdomain")) { 
    Write-Host "Found value for vRAdomain param: $vRADomain"
} else {
    $vRAdomain = Read-Host -Prompt "Enter a domain to connect to vRA with(AD/LDap) or  press enter to leave empty"
}

if ($PSBoundParameters.Keys.Contains("vRAPassword")) { 
    Write-Host "Found value for vRAPassword param"
} else {
    $vrapassword = Read-Host -Prompt "Enter a password to connect to vRA with"
}

if ($PSBoundParameters.Keys.Contains("vRAServer")) { 
    Write-Host "Found value for vRAServer param: $vRAServer"
} else {
    $vRAServer = Read-Host -Prompt "Enter a hostname/fqdn to connect to vRA with"
}

$loginurl="https://$vraserver/csp/gateway/am/api/login?access_token"
if ($vradomain.length -gt 1) {
    $body = "{ ""username"":""$vRAUser"",""password"":""$vRAPassword"",""domain"":""$vRADomain""}"    
} else {
    $body = "{ ""username"":""$vRAUser"",""password"":""$vRAPassword""}"
}

$resp = Invoke-RestMethod -Method POST -ContentType "application/json" -URI $loginurl -Body $body
Write-Host "`n---------Refresh Token---------"
$resp.refresh_token
Write-Host "-------------------------------`n"

#Set ENV Variables for those wanting to use them for the Terraform Provider
$ENV:VRA_URL="https://$vRAServer"
$ENV:VRA_REFRESH_TOKEN=$resp.refresh_token

