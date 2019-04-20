deps:
	go get github.com/rakyll/statik

static: deps
	statik -f -src=./public -dest=server/ -p public

release: static
	go build -tags release

.PHONY: static deps