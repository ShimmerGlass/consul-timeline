GOOS ?= linux
GOARCH ?= amd64
OUT ?= consul-timeline

deps:
	go get github.com/rakyll/statik

static: deps
	statik -f -src=./public -dest=server/ -p public

release: static
	env GOOS=$(GOOS) GOARCH=$(GOARCH) go build -tags release -o $(OUT)

.PHONY: static deps