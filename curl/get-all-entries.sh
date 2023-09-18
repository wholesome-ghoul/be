#!/bin/bash

# Usage: ./login.sh <username> <password> | jq -r '.jwt' | ./get-all-entries.sh

read token

curl \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer $token" \
	localhost:8080/api/v1/entry
