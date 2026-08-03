.PHONY: build crossbuild hygiene format fix module vet staticcheck lint architecture test race vulnerability license secret smoke hook prepush verify

build:
	go build -trimpath ./...

crossbuild:
	scripts/crossbuild

hygiene:
	pre-commit run --all-files --show-diff-on-failure

format:
	scripts/format

fix:
	go tool goimports -w .
	go mod tidy

module:
	go mod tidy -diff

vet:
	go vet ./...

staticcheck:
	go tool staticcheck ./...

lint:
	go tool golangci-lint run

architecture:
	go test ./tests/architecture

test:
	go test ./...

race:
	go test -race ./...

vulnerability:
	go tool govulncheck ./...

secret:
	scripts/secret

license:
	scripts/license

smoke:
	scripts/smoke

hook:
	pre-commit install --install-hooks

prepush:
	pre-commit run --hook-stage pre-push --all-files --show-diff-on-failure

verify: hygiene module build crossbuild vet staticcheck lint architecture test race vulnerability license secret smoke
