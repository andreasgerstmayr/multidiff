.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o zgit main.go

.PHONY: ci
ci:
	go vet
	go test
