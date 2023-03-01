.PHONY: clean lint mock_start mock_stop test

default: all

all: lint test build

build:
	go build -ldflags="-X 'github.com/MartyHub/cac/cmd.Version=development'" -race

clean: mock_stop
	rm -f cac cac.exe coverage.out

lint:
	$(CURDIR)/bin/lint.sh

mock_start:
	$(CURDIR)/bin/mock_start.sh

mock_stop:
	$(CURDIR)/bin/mock_stop.sh

test:
	$(CURDIR)/bin/test.sh $(test)
