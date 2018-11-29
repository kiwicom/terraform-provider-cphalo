.PHONY: build build-plugin build-sandbox build-client run-plugin run-sandbox run-client test race testacc tf-init tf-apply tf-plan tf-destroy clean

vars:=$(shell test -f .env && grep -v '^\#' .env | xargs)

build: build-plugin build-sandbox build-client

build-plugin:
	@go build -o bin/terraform-provider-cphalo cmd/tf-plugin/plugin.go

build-sandbox:
	@go build -o bin/sandbox cmd/sandbox/sandbox.go

build-client:
	@go build -o bin/client cmd/client/client.go

run-plugin: build-plugin
	bin/terraform-provider-cphalo

run-sandbox: build-sandbox
	$(vars) bin/sandbox

run-client: endpoint=server_groups
run-client: build-client
	@$(vars) bin/client $(endpoint)

test:
	go test -v -cover ./api ./cphalo

race:
	go test -v -race ./api ./cphalo

testacc: build-plugin
	$(vars) TF_ACC=1 go test -cover -v -timeout 15m ./cphalo

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

clean:
	rm -f terraform.tfstate*
