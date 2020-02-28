.DEFAULT_GOAL := build

DIRS=bin dist
BINARY=geoip

VERSION=$(shell git describe --tags --always --abbrev=0 --match=v* 2> /dev/null | sed -r "s:^v::g" || echo 0)
VERSION_FULL=$(shell git describe --tags --always --dirty --match=v* 2> /dev/null | sed -r "s:^v::g" || echo 0)

RSRC=README_TPL.md
ROUT=README.md

$(info $(shell mkdir -p $(DIRS)))
BIN=$(CURDIR)/bin
export GOBIN=$(CURDIR)/bin


help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

readme-gen: ## Generates readme from template file
	cp -av "${RSRC}" "${ROUT}"
	sed -ri -e "s:\[\[tag\]\]:${VERSION}:g" -e "s:\[\[os\]\]:linux:g" -e "s:\[\[arch\]\]:amd64:g" "${ROUT}"

release: clean clean-cache fetch generate ## Generate a release, but don't publish to GitHub.
	$(BIN)/goreleaser --skip-validate --skip-publish

publish: clean clean-cache fetch generate ## Generate a release, and publish to GitHub.
	$(BIN)/goreleaser

snapshot: clean clean-cache fetch generate ## Generate a snapshot release.
	$(BIN)/goreleaser --snapshot --skip-validate --skip-publish

fetch: ## Fetches the necessary dependencies to build.
	which $(BIN)/rice 2>&1 > /dev/null || go get -v github.com/GeertJohan/go.rice/rice
	which $(BIN)/goreleaser 2>&1 > /dev/null || wget -qO- "https://github.com/goreleaser/goreleaser/releases/download/v0.122.0/goreleaser_Linux_x86_64.tar.gz" | tar -xz -C $(BIN) goreleaser
	go mod download
	go mod tidy
	go mod vendor
	test -d public/node_modules || (cd public && npm install)

upgrade-deps: ## Upgrade all dependencies to the latest version.
	go get -u ./...

upgrade-deps-patch: ## Upgrade all dependencies to the latest patch release.
	go get -u=patch ./...

clean: ## Cleans up generated files/folders from the build.
	/bin/rm -rfv "dist" "public/dist" "rice-box.go" "${BINARY}-${VERSION_FULL}"

clean-cache: ## Cleans up generated cache (speeds up during dev time).
	/bin/rm -rfv "public/.cache"

generate-watch: ## Generate public html/css/js when files change (faster, but larger files.)
	cd public && npm run watch

generate: ## Generate public html/css/js files for use in production (slower, smaller/minified files.)
	cd public && npm run build
	$(BIN)/rice -v embed-go

build: fetch clean clean-cache generate ## Builds the application (with generate.)
	go build -ldflags '-s -w' -tags netgo -installsuffix netgo -v -o "${BINARY}-${VERSION_FULL}"

debug: fetch clean generate ## Runs the application in debug mode (with generate-dev.)
	go run *.go -d --http.limit 200000 --http.proxy
