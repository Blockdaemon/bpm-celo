VERSION:=$(CI_COMMIT_REF_NAME)

ifeq ($(VERSION),)
	# Looks like we are not running in the CI so default to current branch
	VERSION:=$(shell $ git symbolic-ref HEAD | sed -e "s/^refs\/heads\///")
endif

# Need to wrap in "bash -c" so env vars work in the compiler as well as on the cli to specify the output
BUILD_CMD:=bash -c 'go build -ldflags "-X main.version=$(VERSION)" -o bin/bpm-$(VERSION)-$$GOOS-$$GOARCH cmd/*'

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 $(BUILD_CMD)
	GOOS=darwin GOARCH=amd64 $(BUILD_CMD)
	GOOS=windows GOARCH=amd64 $(BUILD_CMD)

.PHONY: check
check: test lint

.PHONY: test
test:
	go test -v ./...

.PHONY: test-run-all
test-run-all:
	@bash -c 'go build -ldflags "-X main.version=$(VERSION)" -o build/bpm-$(VERSION)-TEST-LOCAL cmd/*'
	@chmod +x build/bpm-$(VERSION)-TEST-LOCAL
	./scripts/runTests.sh bpm-$(VERSION)-TEST-LOCAL proxy attestation-node fullnode validator attestation-service

.PHONY: lint
lint:
	golangci-lint run --enable gofmt ./...
.PHONY: build
