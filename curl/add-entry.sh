#!/bin/bash

# Usage: ./login.sh <username> <password> | jq -r '.jwt' | ./add-entry.sh <content>

read token
content=$1

curl -X POST \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer $token" \
	-d "{\"content\": \"$content\"}" \
	localhost:8080/api/v1/entry
