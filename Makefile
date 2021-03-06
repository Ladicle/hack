REPO_NAME=hack
PKGROOT=github.com/Ladicle/hack

# VERSION is the git commit hash.
VERSION ?= $(shell git rev-parse --short HEAD)

UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)
ifeq (${UNAME_S},Linux)
	OS=linux
endif
ifeq (${UNAME_S},Darwin)
	OS=darwin
endif
ifeq (${UNAME_M},x86_64)
	ARCH=amd64
endif

OUTDIR=_output

build:
	GO111MODULE=on CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -ldflags "-w -X $(PKGROOT)/cmd.version=$(VERSION) -X $(PKGROOT)/cmd.gitRepo=$(REPO_NAME)" -o $(OUTDIR)/hack

build_darwin64:
	GO111MODULE=on CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-w -X $(PKGROOT)/cmd.version=$(VERSION) -X $(PKGROOT)/cmd.gitRepo=$(REPO_NAME)" -o $(OUTDIR)/hack_darwin64

build_linux64:
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -X $(PKGROOT)/cmd.version=$(VERSION) -X $(PKGROOT)/cmd.gitRepo=$(REPO_NAME)" -o $(OUTDIR)/hack_linux64

install:
	GO111MODULE=on CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go install -ldflags "-w -X $(PKGROOT)/cmd.version=$(VERSION) -X $(PKGROOT)/cmd.gitRepo=$(REPO_NAME)"

check:
	GO111MODULE=on go vet $(PKGROOT)/...
	GO111MODULE=on go test $(PKGROOT)/...

.PHONY: clean
clean:
	-rm -r $(OUTDIR)
