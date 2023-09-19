.PHONY: build
build: dep
	ENV=prod go build -o be main.go

.PHONY: test
test: dep
	ENV=test go test -v ./... $(ARGS)

.PHONY: run
run: dep
	swag fmt && \
		swag init --parseDependency --parseInternal && \
		ENV=dev go run main.go

.PHONY: fmt
fmt:
	gofmt -w .

.PHONY: dep
dep:
	go get -u ./... && \
		go mod tidy
