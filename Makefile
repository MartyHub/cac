.PHONY: lint, test

default: lint test

build:
	go build -race

lint:
	$(CURDIR)/bin/lint.sh

test:
	$(CURDIR)/bin/test.sh $(test)
