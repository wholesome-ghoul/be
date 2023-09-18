#!/bin/bash

# Usage: ./login.sh <username> <password>

username=$1
password=$2

curl -s -X POST \
	-H "Content-Type: application/json" \
	-d "{\"username\": \"$username\", \"password\": \"$password\"}" \
	localhost:8080/api/v1/auth/login
