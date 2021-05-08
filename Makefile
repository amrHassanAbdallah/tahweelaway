NAME = $(notdir $(PWD))
VERSION = $(shell printf "%s.%s" \
		$$(git rev-list --count HEAD) \
		$$(git rev-parse --short HEAD))



generate:
	@echo :: getting generator
	go get -v -d
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen
	@echo :: generating code

	@echo :: install sqlc
	go get github.com/kyleconroy/sqlc/cmd/sqlc
	@echo :: generate sql

	(cd ./persistence && sqlc generate)
	oapi-codegen -package api -generate chi-server,types api/api.yml > api/api.gen.go
	oapi-codegen -package client -generate client,types api/api.yml > client/client.gen.go

test: generate
	@echo :: run tests
	go test ./... -race

standards:
	@echo :: exec go vet to  reports suspicious constructs
	go vet .
	@echo :: end go vet

	@echo :: exec gosec to  check security issues
	go get github.com/securego/gosec/cmd/gosec
	gosec -quiet ./...
	@echo :: end gosec

build:  $(OUTPUT)
	CGO_ENABLED=0 GOOS=linux go build -o bin/app \
		-ldflags "-X main.version=$(VERSION)" \
		-gcflags "-trimpath $(GOPATH)/src"