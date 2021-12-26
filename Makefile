# The root package
PKGROOT=github.com/Ladicle/hack

# VERSION is the git commit hash.
VERSION ?= $(shell git rev-parse --short HEAD)

# OUTDIR is directory where artifacts are stored.
# This must be a relative to PKGROOT.
OUTDIR := build/out

# GOLDFLAGS is a flag used for a build command.
GOLDFLAGS := -w -X $(PKGROOT)/cmd.version=$(VERSION)


.PHONY: build build-linux build-darwin build-windows install test vet lint fmt check clean

build:
	CGO_ENABLED=0 go build -ldflags "$(GOLDFLAGS)" -o $(OUTDIR)/hack

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(GOLDFLAGS)" -o $(OUTDIR)/linux-amd64/hack

build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "$(GOLDFLAGS)" -o $(OUTDIR)/darwin-amd64/hack

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "$(GOLDFLAGS)" -o $(OUTDIR)/windows-amd64/hack.exe

install:
	CGO_ENABLED=0 go install -ldflags "$(GOLDFLAGS)"

test:
	go test $(PKGROOT)/cmd/... $(PKGROOT)/pkg/...

vet:
	go vet -printfuncs Infof,Warningf,Errorf,Fatalf,Exitf,Logf $(PKGROOT)/cmd/... $(PKGROOT)/pkg/...

lint:
	hack/golangci-lint.sh

fmt:
	go fmt $(PKGROOT)/cmd/... $(PKGROOT)/pkg/...

check: fmt vet lint test

clean:
	-rm -rf $(OUTDIR)/*
