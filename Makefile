.PHONY: build
build:
	go build -o multidiff main.go

.PHONY: ci
ci:
	go vet
	go test
