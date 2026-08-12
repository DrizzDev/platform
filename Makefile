.PHONY: build crossbuild release hygiene format fix module vet staticcheck lint architecture test race vulnerability license secret smoke hook prepush verify

build:
	go build -trimpath ./...

crossbuild:
	scripts/crossbuild

# release packages the archives, checksums, Homebrew cask, and GitHub release. It requires the compiled device helpers
# (device-bridge: npm run compile), a git tag, and a GitHub token. The token is read from GITHUB_TOKEN, falling back to
# DRIZZ_GITHUB_PAT_TOKEN loaded from a gitignored .env. Sequential (-p 1) because each target injects its own helper
# into the one embedded asset path.
release:
	@set -a; [ -f .env ] && . ./.env; set +a; \
	GITHUB_TOKEN="$${GITHUB_TOKEN:-$${DRIZZ_GITHUB_PAT_TOKEN:-}}" goreleaser release --clean --parallelism 1

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
