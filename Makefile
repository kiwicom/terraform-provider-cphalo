.PHONY: build run lint testacc test clean release

vars:=$(shell test -f .env && grep -v '^\#' .env | xargs)
VERSION:=v0.0.0-local

#? build: build binary for current system
build: bin/current_system/terraform-provider-cphalo_$(VERSION)

#? run: run plugin
run: build
	bin/current_system/terraform-provider-cphalo

#? lint: run a meta linter
lint:
	@hash golangci-lint || (echo "Download golangci-lint from https://github.com/golangci/golangci-lint#install" && exit 1)
	golangci-lint run

#? testacc: run acceptance tests
testacc:
	$(vars) TF_ACC=1 go test -cover -v -timeout 15m -failfast ./cphalo

#? test: run unit tests
test:
	$(vars) go test -v ./...

.env:
	cp .env.example .env

#? clean: removes all artificats
clean:
	rm -fr bin/

bin/current_system/terraform-provider-cphalo_%:  GOARGS =
bin/darwin_amd64/terraform-provider-cphalo_%:  GOARGS = GOOS=darwin GOARCH=amd64
bin/linux_amd64/terraform-provider-cphalo_%:  GOARGS = GOOS=linux GOARCH=amd64
bin/linux_386/terraform-provider-cphalo_%:  GOARGS = GOOS=linux GOARCH=386
bin/linux_arm/terraform-provider-cphalo_%:  GOARGS = GOOS=linux GOARCH=arm
bin/windows_amd64/terraform-provider-cphalo_%:  GOARGS = GOOS=windows GOARCH=amd64
bin/windows_386/terraform-provider-cphalo_%:  GOARGS = GOOS=windows GOARCH=386

bin/%/terraform-provider-cphalo_$(VERSION): clean
	$(GOARGS) CGO_ENABLED=0 go build -o $@ -ldflags="-s -w" .

#? release: make a release for all systems
release: \
	bin/release/terraform-provider-cphalo_darwin_amd64.zip \
	bin/release/terraform-provider-cphalo_linux_amd64.zip \
	bin/release/terraform-provider-cphalo_linux_386.zip \
	bin/release/terraform-provider-cphalo_linux_arm.zip \
	bin/release/terraform-provider-cphalo_windows_amd64.zip \
	bin/release/terraform-provider-cphalo_windows_386.zip

bin/release/terraform-provider-cphalo_%.zip: NAME=terraform-provider-cphalo_$(VERSION)_$*
bin/release/terraform-provider-cphalo_%.zip: DEST=bin/release/$(VERSION)/$(NAME)
bin/release/terraform-provider-cphalo_%.zip: bin/%/terraform-provider-cphalo_$(VERSION)
	mkdir -p $(DEST)
	cp bin/$*/terraform-provider-cphalo_$(VERSION) readme.md $(DEST)
	cd $(DEST) && zip -r ../$(NAME).zip . && cd .. && sha256sum $(NAME).zip > $(NAME).sha256 && rm -rf $(NAME)

#? help: display help
help: Makefile
	@printf "Available make targets:\n\n"
	@sed -n 's/^#?//p' $< | column -t -s ':' |  sed -e 's/^/ /'
