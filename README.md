# Be

## Stack

- Golang
- Postgresql

```bash
# create database
# createdb be-master-test
createdb be-master-dev
```

## Interacting with API

```bash
# start dev server
make run

# interact directly from swagger or with curl
# API endpoints available at http://localhost:8080/swagger/index.html

# register user
./curl/register.sh <username> <password>

# login user
./curl/login.sh <username> <password>

# add entry
./curl/login.sh <username> <password> | jq -r '.jwt' | ./curl/add-entry.sh <content>

# get all entries for user
./curl/login.sh <username> <password> | jq -r '.jwt' | ./curl/get-all-entries.sh
```

## Testing

```bash
# run all tests
make test

# run single test
make test ARGS="-testify.m <function pattern>"
```

## What I don't like currently

- `SetupTest` and `TearDownTest` run before and after each tests
- I Think there are too many `panic` methods in tests
- I don't like that test API endpoints may not be aligned with the actual endpoints
