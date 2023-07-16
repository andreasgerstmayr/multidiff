.PHONY: build
build:
	go build -o zgit main.go

.PHONY: ci
ci:
	go vet
	go test
