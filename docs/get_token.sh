#!/bin/bash
#
# Script to generate a refresh token for vRA8 on prem or vRA Cloud. 
# This will prompt for the following values if they are not already set:
#    username
# Sets environment variables for VRA_REFRESH_TOKEN and VRA_URL which can be consumed by the 
# TF provider more securely than leaving the token in cleartext. 
#
#
if  ! [ -x "$(command -v jq)" ]
then
	echo -e "\n\nthe jq utility is missing. See https://stedolan.github.io/jq/ for instructions to get it\n\n"
	return 1
fi

#Check for an already existing username value
if [[ -v username ]] 
then
	echo -e "\nusername variable found: $username\n"
else 
	echo -e "\n\nPlease enter username to connect to vra with"
	read username
fi

#Check for an already existing password value
if [[ -v password ]]
then
    echo -e "\npassword variable found\n"
else
    echo -e "\n\nPlease enter password to connect to vra with\n"
    read password
fi

#Check for an already existing LDAP/AD domain value
if [[ -v domain ]]
then
    echo -e "\nExisting domain variable found: $domain\n"
else
	echo -e "\n\nPlease enter domain to connect to vra with (for AD/LDAP users) or press Enter"
	read domain
fi

if [[ -v VRA_URL || -v host ]]
then 
	echo -e "\nfound a value for the vra/cas server\n"
else
 	echo -e "\n\nPlease enter the hostname/fqdn of the VRA8 server/ or cloud identity server"
	read host
	export VRA_URL="https://$host"
fi

#use different json bodies with curl depending on whether or not a domain 
# was specified
echo -e "\nGetting Token"
if [[ $domain == "" ]]
then
	export VRA_REFRESH_TOKEN=`curl -k -X POST \
  		"$VRA_URL/csp/gateway/am/api/login?access_token" \
  		-H 'Content-Type: application/json' \
  		-s \
  		-d '{
  		"username": "'"$username"'",
  		"password": "'"$password"'"
		}' | jq -r .refresh_token`

else
	export VRA_REFRESH_TOKEN=`curl -k -X POST \
  		"$VRA_URL/csp/gateway/am/api/login?access_token" \
  		-H 'Content-Type: application/json' \
  		-s \
  		-d '{
  		"username": "'"$username"'",
  		"password": "'"$password"'",
  		"domain": "'"$domain"'"
		}' | jq -r .refresh_token`
fi


#clean up password 
unset password

echo -e "\n\nRefresh Token"
echo "----------------------------"
echo $VRA_REFRESH_TOKEN
echo "----------------------------"
