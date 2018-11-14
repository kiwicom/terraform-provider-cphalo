.PHONY: build run

build:
	go build -o bin/terraform-provider-cloudpassage

run: build
	bin/terraform-provider-cloudpassage
