SOURCES := $(shell find . -name *.go)
BINARY:=glooctl
VERSION:=$(shell cat version)
UNAME := $(shell uname)

build: $(BINARY)

$(BINARY): $(SOURCES)
	CGO_ENABLED=0 go build -ldflags "-X main.Version=$(VERSION)" -v -o $@ *.go

$(BINARY)-darwin: $(SOURCES)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -v -o $(BINARY) *.go

$(BINARY)-win: $(SOURCES)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -v -o $(BINARY).exe *.go

$(BINARY)-linux: $(SOURCES)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -v -o $(BINARY) *.go

release-darwin: $(BINARY)-darwin
	mkdir -p release
	tar cvzf release/$(BINARY)-$(VERSION)-macOS-64.tar.gz $(BINARY)

release-win: $(BINARY)-win
	mkdir -p release
	zip release/$(BINARY)-$(VERSION)-Windows-64.zip $(BINARY).exe

release-linux: $(BINARY)-linux
	mkdir -p release
	tar cvzf release/$(BINARY)-$(VERSION)-Linux-64.tar.gz $(BINARY)

release: release-darwin release-win release-linux
ifeq ($(UNAME),Darwin)
	cd release && shasum -a 256 * > $(BINARY)-$(VERSION)-checksums.txt
else
	cd release && sha256sum * > $(BINARY)-$(VERSION)-checksums.txt
endif

test:
	ginkgo -r -v 

clean:
	rm -f $(BINARY)
	rm -f $(BINARY).exe
	rm -rf release
