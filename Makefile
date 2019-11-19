
BINARY := thesaurus

LDFLAGS_DEV = -ldflags "-X main.version=${VERSION}"

build:
	@go build ${LDFLAGS_DEV} -o bin/${BINARY} cmd/thesaurus/main.go

test:
	@go test -v -race ./...

# update the golden files used for the integration tests
update-tests:
	@go test integration/cli_test.go -update