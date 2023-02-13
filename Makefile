.PHONY: lint, test

default: lint test

build:
	go build -ldflags="-X 'github.com/MartyHub/cac/internal.Version=development'" -race

lint:
	$(CURDIR)/bin/lint.sh

test:
	$(CURDIR)/bin/test.sh $(test)
