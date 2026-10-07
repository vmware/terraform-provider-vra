#!/bin/bash

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
# WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
# COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
# OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

# Generates and returns an `access_token` and `refresh_token` from VCF Automation 9 for use by the vRA Terraform provider.
#
# Run with `source` so they persist in your shell.
# Set VRA_INSECURE=true to skip TLS verification (not recommended).
#
# Sets environment variables:
# * `VRA_URL`: ThE URL to the VCF Automation 9 endpoint.
# * `VCFA_ORGANIZATION`: The name of the organization.
# * `VRA_CLIENT_NAME`: The client name of the API Token.
# * `VRA_CLIENT_ID`: The client ID of the API Token.
# * `VRA_ACCESS_TOKEN`: The access token.
# * `VRA_REFRESH_TOKEN`: The refresh token.

### Check for an installation of jq. ###

if ! [ -x "$(command -v jq)" ]; then
	echo -e "\nThe jq utility is missing. See https://stedolan.github.io/jq/ for installation instructions.\n"
	return 1 2>/dev/null || exit 1
fi

### Check for an existing endpoint value. ###

if [[ -n "${fqdn:-}" ]]; then
	echo -e "\nFQDN variable found: $fqdn. Skipping...\n"
else
	echo -e "\nEnter the VCF Automation FQDN:"
	read -r fqdn
fi
export VRA_URL="https://${fqdn}"

### Check for an existing username value. ###

if [[ -n "${username:-}" ]]; then
	echo -e "\nUsername variable found: $username. Skipping...\n"
else
	echo -e "\nEnter the username to authenticate with VCF Automation:"
	read -r username
fi

### Check for an existing password value. ###

if [[ -n "${password:-}" ]]; then
	echo -e "\nPassword variable found. Skipping...\n"
else
	echo -e "\nEnter the password to authenticate with VCF Automation:"
	read -r -s password
fi

### Check for an existing organization value. ###

if [[ -n "${organization:-}" ]]; then
	echo -e "\nOrganization variable found: $organization. Skipping...\n"
else
	echo -e "\nEnter the organization:"
	read -r organization
fi
export VCFA_ORGANIZATION="${organization:-}"

### Check for an existing organization type value. ###

if [[ -z "${app_type:-}" ]]; then
	echo -e "\nSelect the organization type:"
	echo "  1) VM Apps"
	echo "  2) All Apps"
	printf "Enter 1 or 2: "
	read -r app_type
fi
while true; do
	case "${app_type}" in
		1|"VM Apps") app_type="VM Apps"; break ;;
		2|"All Apps") app_type="All Apps"; break ;;
	esac
	echo -e "\nInvalid selection: ${app_type}."
	printf "Enter 1 (VM Apps) or 2 (All Apps): "
	read -r app_type
done
export VCFA_APP_TYPE="${app_type}"

tls_opts=()
if [[ "${VRA_INSECURE:-false}" == "true" ]]; then
	tls_opts+=(-k)
fi

if [[ "${app_type}" == "All Apps" ]]; then

	### All Apps: Generate the user access token. ###

	echo -e "\nAuthenticating User..."

	curl_opts=(-s -I -o /dev/null)
	if [[ "${VRA_INSECURE:-false}" == "true" ]]; then
		curl_opts+=(-k)
	fi

	# Pass credentials through a config file descriptor so the password is not visible in the process list.
	cred="${username}@${organization}:${password}"
	cred="${cred//\\/\\\\}"
	cred="${cred//\"/\\\"}"

	access_token=$(curl "${curl_opts[@]}" -X POST \
		"${VRA_URL}/tm/cloudapi/1.0.0/sessions" \
		--header "Accept: application/json;version=40.0" \
		-w "%header{x-vmware-vcloud-access-token}" \
		-K <(printf 'user = "%s"\n' "$cred"))

	unset cred

	if [[ -z "${access_token}" ]]; then
		echo -e "\nFailed to obtain an access token. Check the FQDN, credentials and organization.\n"
		unset password
		return 1 2>/dev/null || exit 1
	fi

	export VRA_ACCESS_TOKEN="${access_token}"

else

	### VM Apps: generate the user access token. ###

	echo -e "\nAuthenticating User..."

	# Credentials are passed through a config file descriptor so they are not visible in the process list.
	vm_username="${username//\\/\\\\}"
	vm_username="${vm_username//\"/\\\"}"
	vm_password="${password//\\/\\\\}"
	vm_password="${vm_password//\"/\\\"}"

	token=$(curl -f -sS "${tls_opts[@]}" -X POST \
		"${VRA_URL}/oauth/tenant/${organization}/token" \
		-H "Accept: application/json" \
		-H "Content-Type: application/x-www-form-urlencoded" \
		--data-urlencode "grant_type=password" \
		-K <(printf 'data-urlencode = "username=%s"\ndata-urlencode = "password=%s"\n' "${vm_username}" "${vm_password}") | jq -r '.access_token // empty')

	unset vm_username vm_password

	if [[ -z "${token}" ]]; then
		echo -e "\nFailed to obtain an access token. Check the FQDN, credentials and organization.\n"
		unset password
		return 1 2>/dev/null || exit 1
	fi

	### VM Apps: exchange the token at the OIDC endpoint. ###

	echo -e "\nExchanging access token..."

	access_token=$(curl -f -sS "${tls_opts[@]}" -X POST \
		"${VRA_URL}/oidc/oauth2/token" \
		-H "Accept: application/json" \
		-H "Content-Type: application/x-www-form-urlencoded" \
		--data-urlencode "grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer" \
		--data-urlencode "scope=openid+profile+email+phone+groups+vcd_idp" \
		-K <(printf 'data-urlencode = "assertion=%s"\n' "${token}") | jq -r '.access_token // empty')

	unset token

	if [[ -z "${access_token}" ]]; then
		echo -e "\nFailed to exchange the access token at the OIDC endpoint.\n"
		unset password
		return 1 2>/dev/null || exit 1
	fi

	export VRA_ACCESS_TOKEN="${access_token}"
fi

### Register an OAuth client to generate a client_id. ###

echo -e "\nRegistering OAuth client..."

client_name="terraform-$(LC_ALL=C tr -dc 'a-z0-9' </dev/urandom | head -c 8)"

# The bearer token is passed through a config file descriptor so it is not visible in the process list.
client_id=$(curl -s "${tls_opts[@]}" -X POST \
	"${VRA_URL}/tm/oauth/tenant/${organization}/register" \
	--header "Accept: application/json;version=40.0" \
	--header "Content-Type: application/json" \
	-K <(printf 'header = "Authorization: Bearer %s"\n' "${access_token}") \
	-d "$(jq -n --arg name "${client_name}" '{client_name: $name}')" | jq -r '.client_id // empty')

if [[ -z "${client_id}" ]]; then
	echo -e "\nFailed to register an OAuth client. Check the organization and access token.\n"
	unset password
	return 1 2>/dev/null || exit 1
fi

export VRA_CLIENT_NAME="${client_name}"
export VRA_CLIENT_ID="${client_id}"

### Exchange the access token for a refresh token. ###

echo -e "\nGenerating Refresh Token..."

# The assertion is passed through a config file descriptor so it is not visible in the process list.
refresh_token=$(curl -s "${tls_opts[@]}" -X POST \
	"${VRA_URL}/tm/oauth/tenant/${organization}/token" \
	--header "Accept: application/json;version=40.0" \
	--header "Content-Type: application/x-www-form-urlencoded" \
	--data-urlencode "grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer" \
	--data-urlencode "client_id=${client_id}" \
	-K <(printf 'data-urlencode = "assertion=%s"\n' "${access_token}") | jq -r '.refresh_token // empty')

if [[ -z "${refresh_token}" ]]; then
	echo -e "\nFailed to obtain a refresh token. Check the client_id and access token.\n"
	unset password
	return 1 2>/dev/null || exit 1
fi

export VRA_REFRESH_TOKEN="${refresh_token}"

echo ""
echo "Environmental variables..."
echo ""
echo "VRA_URL = ${VRA_URL}"
echo "VCFA_ORGANIZATION = ${VCFA_ORGANIZATION}"
echo "VCFA_APP_TYPE = ${VCFA_APP_TYPE}"
echo "VRA_CLIENT_NAME = ${VRA_CLIENT_NAME}"
echo "VRA_CLIENT_ID = ${VRA_CLIENT_ID}"
echo "VRA_ACCESS_TOKEN is set"
echo "VRA_REFRESH_TOKEN is set"

### Clear the password value. ###
unset password
