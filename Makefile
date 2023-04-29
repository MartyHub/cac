.PHONY: clean lint mock_start mock_stop test

default: all

all: tidy lint test build

build:
	go build -ldflags="-X 'github.com/MartyHub/cac/cmd.Version=development'" -race

clean: mock_stop
	rm -f cac cac.exe coverage.out

lint:
	$(CURDIR)/scripts/lint.sh

mock_start:
	$(CURDIR)/scripts/mock_start.sh

mock_stop:
	$(CURDIR)/scripts/mock_stop.sh

test:
	$(CURDIR)/scripts/test.sh $(test)

tidy:
	go mod tidy
