SOURCES := $(shell find . -name *.go)
BINARY:=glooctl
OUTPUT_DIR ?= _output
OUTPUT_BINARY ?= $(OUTPUT_DIR)/$(BINARY)
VERSION:=$(shell cat version)
UNAME := $(shell uname)

$(OUTPUT_DIR):
	mkdir -p $@

build: $(OUTPUT_BINARY)

$(OUTPUT_BINARY): $(SOURCES)
	CGO_ENABLED=0 go build -ldflags "-X main.Version=$(VERSION)" -v -o $@ *.go

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


#----------------------------------------------------------------------------------
# Release
#----------------------------------------------------------------------------------

RELEASE_BINARIES := $(OUTPUT_DIR)/glooctl-linux-amd64 $(OUTPUT_DIR)/glooctl-darwin-amd64 $(OUTPUT_DIR)/glooctl-win-amd64

$(OUTPUT_DIR)/qlooctl-linux-amd64: $(SOURCES)
	GOOS=linux CGO_ENABLED=0 go build -ldflags "-X main.Version=$(VERSION)" -v -o $@ *.go

$(OUTPUT_DIR)/qlooctl-darwin-amd64: $(SOURCES)
	GOOS=darwin CGO_ENABLED=0 go build -ldflags "-X main.Version=$(VERSION)" -v -o $@ *.go

$(OUTPUT_DIR)/qlooctl-darwin-amd64: $(SOURCES)
	GOOS=windows CGO_ENABLED=0 go build -ldflags "-X main.Version=$(VERSION)" -v -o $@ *.go

.PHONY: release-binaries
release-binaries: $(RELEASE_BINARIES)

.PHONY: release
release: release-binaries
	hack/create-release.sh github_api_token=$(GITHUB_TOKEN) owner=solo-io repo=glooctl tag=v$(VERSION)
	@$(foreach BINARY,$(RELEASE_BINARIES),hack/upload-github-release-asset.sh github_api_token=$(GITHUB_TOKEN) owner=solo-io repo=qloo tag=v$(VERSION) filename=$(BINARY);)
