.PHONY: dep build

# REPO_INFO is the URL of git repository.
REPO_INFO ?= $(shell git config --get remote.origin.url)

# VERSION is the git commit hash prefixed with "git-".
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

dep:
	dep ensure -update

build:
	CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -ldflags "-w -X main/cmd.version=$(VERSION) -X main.gitRepo=$(REPO_INFO)" -o $(OUTDIR)/hack

build_darwin64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-w -X main/cmd.version=$(VERSION) -X main.gitRepo=$(REPO_INFO)" -o $(OUTDIR)/hack_darwin64

build_linux64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -X main.version=$(VERSION) -X main.gitRepo=$(REPO_INFO)" -o $(OUTDIR)/hack_linux64

install:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go install -ldflags "-w -X main.version=$(VERSION) -X main.gitRepo=$(REPO_INFO)"
