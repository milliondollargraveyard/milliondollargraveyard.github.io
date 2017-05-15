BINARY = $(notdir $(PWD))
VERSION := $(shell git describe --tags --dirty --always 2> /dev/null || echo "dev")
SOURCES = $(wildcard *.go **/*.go)

all: $(BINARY)

$(BINARY): $(SOURCES)
	go build -ldflags "-X main.version=`git describe --long --tags --dirty --always`" -o "$@"

deps:
	go get ./...

crawl: $(BINARY) map.json
	cat map.json | ./$(BINARY) -c=25 | tee out.jsons

build: $(BINARY)

clean:
	rm $(BINARY)

run: $(BINARY)
	./$(BINARY) --help

debug: $(BINARY)
	./$(BINARY) --pprof :6060 -vv

test:
	go test ./...
	golint ./...
