.PHONY: build-plugin build-sandbox run-plugin run-sandbox test race tf-init tf-apply tf-plan

vars:=$(shell grep -v '^\#' .env | xargs)

build-plugin:
	go build -o bin/terraform-provider-cphalo cmd/tf-plugin/plugin.go

build-sandbox:
	go build -o bin/sandbox cmd/sandbox/sandbox.go

run-plugin: build-plugin
	bin/terraform-provider-cphalo

run-sandbox: build-sandbox
	$(vars) bin/sandbox

test:
	go test -v ./api ./cphalo

race:
	go test -v -race ./api ./cphalo

testacc: build-plugin
	$(vars) TF_ACC=1 go test -v -timeout 1m ./cphalo

.env:
	cp .env.example .env

tf-init: build-plugin
	terraform init -plugin-dir=bin/

tf-apply: tf-init
	$(vars) terraform apply

tf-plan: tf-init
	$(vars) terraform plan

tf-destroy:
	$(vars) terraform destroy
