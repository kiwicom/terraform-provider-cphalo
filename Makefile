.PHONY: build-plugin build-sandbox run-plugin run-sandbox test

vars:=$(shell cat .env | xargs)

build-plugin:
	go build -o bin/terraform-provider-cloudpassage cmd/tf-plugin/plugin.go

build-sandbox:
	go build -o bin/sandbox cmd/sandbox/sandbox.go

run-plugin: build-plugin
	bin/terraform-provider-cloudpassage

run-sandbox: build-sandbox
	$(vars) bin/sandbox

test:
	go test -v ./api

.env:
	cp .env.example .env
