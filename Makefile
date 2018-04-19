SOURCES := $(shell find . -name *.go)
BINARY:=glooctl
VERSION:=$(shell cat version)

build: $(BINARY)

$(BINARY): $(SOURCES)
	CGO_ENABLED=0 go build -ldflags "-X main.Version=$(VERSION)" -v -o $@ *.go

$(BINARY)-darwin: $(SOURCES)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -v -o $(BINARY)-macOS-64 *.go

$(BINARY)-win: $(SOURCES)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -v -o $(BINARY)-Windows-64 *.go

$(BINARY)-linux: $(SOURCES)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -v -o $(BINARY)-Linux-64 *.go

release: $(BINARY)-darwin $(BINARY)-win $(BINARY)-linux
	mkdir release
	tar cvzf release/$(BINARY)-$(VERSION)-macOS-64.tar.gz $(BINARY)-macOS-64
	tar cvzf release/$(BINARY)-$(VERSION)-Linux-64.tar.gz $(BINARY)-Linux-64
	zip release/$(BINARY)-$(VERSION)-Windows-64.zip $(BINARY)-Windows-64

test:
	go test -cover ./...

clean:
	rm -f $(BINARY)
	rm -f *-64
	rm -rf release
