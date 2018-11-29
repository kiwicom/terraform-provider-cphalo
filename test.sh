#!/usr/bin/env bash

set -eu

token=$( echo -n $1:$2 | openssl base64)
echo "SHA auth: ${token}"

response=$(curl -s -X POST \
     -H 'Content-type: application/json' \
     -H "Authorization: Basic ${token}" \
     "https://api.cloudpassage.com/oauth/access_token?grant_type=client_credentials")

echo "Response: ${response}"

accessToken=$(echo ${response} | jq -r '.access_token')
echo "Access token: ${accessToken}"

#echo "All the groups:"
#curl -X GET \
#     -H 'Content-type: application/json' \
#     -H "Authorization: Bearer ${accessToken}" \
#     "https://api.cloudpassage.com/v1/groups" | jq .

echo "All the CSP Accounts:"
curl -X GET \
     -H 'Content-type: application/json' \
     -H "Authorization: Bearer ${accessToken}" \
     "https://api.cloudpassage.com/v1/csp_accounts" | jq .

#echo "One CSP Accounts:"
#curl -X GET \
#     -H 'Content-type: application/json' \
#     -H "Authorization: Bearer ${accessToken}" \
#     "https://api.cloudpassage.com/v1/csp_accounts/716d070a-c23d-489c-bf15-663a11669b56" | jq .

#curl -X POST -v \
#     --data '{"csp_account_type":"AWS","external_id":"this-is-external-id-1","role_arn":"arn:aws:iam::782106534067:role/CloudPassage-Service-Role","group_id":"fff04606e97b11e896d9252f8ed31fc8"}' \
#     -H 'Content-type: application/json' \
#     -H "Authorization: Bearer ${accessToken}" \
#     "https://api.cloudpassage.com/v1/csp_accounts"
#
#echo "Delete the CSP Accounts:"
#curl -X DELETE \
#     -H 'Content-type: application/json' \
#     -H "Authorization: Bearer ${accessToken}" \
#     "https://api.cloudpassage.com/v1/csp_accounts/716d070a-c23d-489c-bf15-663a11669b56" | jq .

#echo "All firewall policies:"
#curl -X GET \
#     -H 'Content-type: application/json' \
#     -H "Authorization: Bearer ${accessToken}" \
#     "https://api.cloudpassage.com/v1/firewall_policies" | jq .
#
#echo "Firewall policy details:"
#curl -X GET \
#     -H 'Content-type: application/json' \
#     -H "Authorization: Bearer ${accessToken}" \
#     "https://api.cloudpassage.com/v1/firewall_policies/be28b106ee5b11e8b7f1017da54e9117" | jq .
#
#echo "Firewall policy rules:"
#curl -X GET \
#     -H 'Content-type: application/json' \
#     -H "Authorization: Bearer ${accessToken}" \
#     "https://api.cloudpassage.com/v1/firewall_policies/be28b106ee5b11e8b7f1017da54e9117/firewall_rules" | jq .
#
#
#echo "Create new group"
#
#createResponse=$(curl -X POST \
#     --data '{"group":{"name":"from_bash"}}' \
#     -H 'Content-type: application/json' \
#     -H "Authorization: Bearer ${accessToken}" \
#     "https://api.cloudpassage.com/v1/groups")
#
#echo ${createResponse} | jq .
#
#createdToken=$(echo ${createResponse} | jq -r .group.id)
#echo "Created token: ${createdToken}"
#
#echo "Create new child group"
#createResponseChild=$(curl -X POST \
#     --data "{\"group\":{\"name\":\"child_group\",\"parent_id\":\"${createdToken}\"}}" \
#     -H "Content-type: application/json" \
#     -H "Authorization: Bearer ${accessToken}" \
#     "https://api.cloudpassage.com/v1/groups")
#
#createdTokenChild=$(echo ${createResponseChild} | jq -r .group.id)
#echo "Created child token: ${createdTokenChild}"
#
#echo "Get the group again:"
#curl -X GET \
#     -H 'Content-type: application/json' \
#     -H "Authorization: Bearer ${accessToken}" \
#     "https://api.cloudpassage.com/v1/groups/${createdToken}" | jq .
#
#echo "Get the child group again:"
#curl -X GET \
#     -H 'Content-type: application/json' \
#     -H "Authorization: Bearer ${accessToken}" \
#     "https://api.cloudpassage.com/v1/groups/${createdTokenChild}" | jq .
#
#
#echo "Delete the child group:"
#curl -X DELETE \
#     -H 'Content-type: application/json' \
#     -H "Authorization: Bearer ${accessToken}" \
#     "https://api.cloudpassage.com/v1/groups/${createdTokenChild}" | jq .
#
#echo "Get the deleted child group again:"
#curl -X GET \
#     -H 'Content-type: application/json' \
#     -H "Authorization: Bearer ${accessToken}" \
#     "https://api.cloudpassage.com/v1/groups/${createdTokenChild}" | jq .
#
#echo "Get non-existing group:"
#curl -X GET -v \
#     -H 'Content-type: application/json' \
#     -H "Authorization: Bearer ${accessToken}" \
#     "https://api.cloudpassage.com/v1/groups/fd63ce2aee4011e8948e2f3994aefbx6" | jq .
#
#echo "Delete the group:"
#curl -X DELETE \
#     -H 'Content-type: application/json' \
#     -H "Authorization: Bearer ${accessToken}" \
#     "https://api.cloudpassage.com/v1/groups/${createdToken}" | jq .
