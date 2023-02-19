.PHONY: lint, test

default: all

all: lint test build

build:
	go build -ldflags="-X 'github.com/MartyHub/cac/cmd.Version=development'" -race

lint:
	$(CURDIR)/bin/lint.sh

test:
	$(CURDIR)/bin/test.sh $(test)
