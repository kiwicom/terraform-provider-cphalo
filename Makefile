.PHONY: build build-plugin build-sandbox build-client run-plugin run-sandbox run-client test race testacc tf-init tf-apply tf-plan tf-destroy clean release

vars:=$(shell test -f .env && grep -v '^\#' .env | xargs)
VERSION:=alpha

build: build-plugin build-sandbox build-client

build-plugin: bin/plugin/current_system/terraform-provider-cphalo

build-sandbox:
	@go build -o bin/sandbox cmd/sandbox/sandbox.go

build-client:
	@go build -o bin/client cmd/client/client.go

run-plugin: build-plugin
	bin/plugin/current_system/terraform-provider-cphalo

run-sandbox: build-sandbox
	$(vars) bin/sandbox

run-client: endpoint=server_groups
run-client: build-client
	@$(vars) bin/client $(endpoint)

test:
	go test -v -cover -timeout 1m ./api ./cphalo

race:
	go test -v -race -timeout 2m ./api ./cphalo

testacc:
	$(vars) TF_ACC=1 go test -cover -v -timeout 15m -failfast ./cphalo

.env:
	cp .env.example .env

tf-init: build-plugin
	terraform init -plugin-dir=bin/plugin/current_system/

tf-apply: tf-init
	$(vars) terraform apply

tf-plan: tf-init
	$(vars) terraform plan

tf-destroy:
	$(vars) terraform destroy

clean:
	rm -fr terraform.tfstate* crash.log bin/*

bin/plugin/current_system/terraform-provider-cphalo:  GOARGS =
bin/plugin/darwin_amd64/terraform-provider-cphalo:  GOARGS = GOOS=darwin GOARCH=amd64
bin/plugin/linux_amd64/terraform-provider-cphalo:  GOARGS = GOOS=linux GOARCH=amd64
bin/plugin/linux_386/terraform-provider-cphalo:  GOARGS = GOOS=linux GOARCH=386
bin/plugin/linux_arm/terraform-provider-cphalo:  GOARGS = GOOS=linux GOARCH=arm
bin/plugin/windows_amd64/terraform-provider-cphalo:  GOARGS = GOOS=windows GOARCH=amd64
bin/plugin/windows_386/terraform-provider-cphalo:  GOARGS = GOOS=windows GOARCH=386

bin/plugin/%/terraform-provider-cphalo: clean
	$(GOARGS) go build -o $@ -a cmd/tf-plugin/plugin.go

release: \
	bin/release/terraform-provider-cphalo_darwin_amd64.zip \
	bin/release/terraform-provider-cphalo_linux_amd64.zip \
	bin/release/terraform-provider-cphalo_linux_386.zip \
	bin/release/terraform-provider-cphalo_linux_arm.zip \
	bin/release/terraform-provider-cphalo_windows_amd64.zip \
	bin/release/terraform-provider-cphalo_windows_386.zip
	$(MAKE) checksum

bin/release/terraform-provider-cphalo_%.zip: NAME=terraform-provider-cphalo_$(VERSION)_$*
bin/release/terraform-provider-cphalo_%.zip: DEST=bin/release/$(VERSION)/$(NAME)
bin/release/terraform-provider-cphalo_%.zip: bin/plugin/%/terraform-provider-cphalo
	mkdir -p $(DEST)
	cp bin/plugin/$*/terraform-provider-cphalo readme.md $(DEST)
	cd $(DEST) && zip -r ../$(NAME).zip . && cd .. && rm -rf $(NAME)

checksum:
	cd bin/release/$(VERSION) && shasum -a 256 * > checksum.sha256
