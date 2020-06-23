NAME:=celo

.PHONY: build
build:
	@echo "IMPORTANT: This is just a quick development build. Use 'goreleaser' for production ready builds!"
	go build -ldflags "-X main.version=development" -o bin/$(NAME)-development ./cmd

.PHONY: check
check: lint test

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	golangci-lint --enable gofmt run

.PHONY: pre-release
pre-release:
	@ test -n "$(VERSION)" || (echo 'ERROR: version is not set. Call like this: make version=1.14.0-rc1 release'; exit 1)

	@ test -n "$(GITLAB_TOKEN)" || (echo 'ERROR: GITLAB_TOKEN is not set. See: https://goreleaser.com/quick-start/'; exit 1)

	@ test -z "$$(git status --porcelain)" || (echo "ERROR: git is dirty - clean up first"; exit 1)

	@ echo "CHANGELOG.md starting here"
	@ echo "--------------------------"
	@ cat CHANGELOG.md
	@ read -p "Press enter to continue if the changelog looks ok. CTRL+C to abort."

.PHONY: test-release
test-release: check
	goreleaser --snapshot --skip-publish --rm-dist

.PHONY: release
release: pre-release check
	# tag it
	# git tag v$(VERSION)
	# git push origin v$(VERSION)

	# finally run the actually release
	goreleaser release --rm-dist

.PHONY: test-run-all
test-run-all:
	@bash -c 'go build -ldflags "-X main.version=$(VERSION)" -o build/bpm-$(VERSION)-TEST-LOCAL cmd/*'
	@chmod +x build/bpm-$(VERSION)-TEST-LOCAL
	./scripts/runTests.sh bpm-$(VERSION)-TEST-LOCAL proxy attestation-node fullnode validator attestation-service
	# ./scripts/runTests.sh bpm-$(VERSION)-TEST-LOCAL proxy validator
