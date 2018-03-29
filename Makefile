SOURCES := $(shell find . -name *.go)
BINARY:=glooctl
VERSION:=$(shell cat version)

build: $(BINARY)

$(BINARY): $(SOURCES)
	CGO_ENABLED=0 go build -ldflags "-X main.Version=$(VERSION)" -v -o $@ *.go

test:
	go test -cover ./...

clean:
	rm -f $(BINARY)
