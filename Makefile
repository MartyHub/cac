.PHONY: lint mock test

default: all

all: lint test build

build:
	go build -ldflags="-X 'github.com/MartyHub/cac/cmd.Version=development'" -race

lint:
	$(CURDIR)/bin/lint.sh

mock:
	podman run -it --rm \
	--name    wiremock \
	--publish 8443:8443 \
	--volume  $(PWD)/wiremock:/home/wiremock:ro \
	wiremock/wiremock:2.35.0 \
	--disable-banner \
	--disable-http \
	--global-response-templating \
	--https-port 8443 \
	--verbose

test:
	$(CURDIR)/bin/test.sh $(test)
