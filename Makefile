build:
	ENV=prod go build -o be main.go

test:
	ENV=test go test -v ./... $(ARGS)

run:
	swag init --parseDependency --parseInternal && ENV=dev go run main.go

fmt:
	gofmt -w .

dep:
	go get -u ./...
	go mod tidy
