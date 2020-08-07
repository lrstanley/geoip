.DEFAULT_GOAL := build

DIRS=bin
BINARY=geoip

VERSION=$(shell git describe --tags --always --abbrev=0 --match=v* 2> /dev/null | sed -r "s:^v::g" || echo 0)

$(info $(shell mkdir -p $(DIRS)))
BIN=$(CURDIR)/bin
export GOBIN=$(CURDIR)/bin


help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

fetch: ## Fetches the necessary dependencies to build.
	which $(BIN)/rice 2>&1 > /dev/null || go get -v github.com/GeertJohan/go.rice/rice
	go mod download
	go mod tidy
	test -d public/node_modules || (cd public && npm install)

upgrade-deps: ## Upgrade all dependencies to the latest version.
	go get -u ./...

upgrade-deps-patch: ## Upgrade all dependencies to the latest patch release.
	go get -u=patch ./...

clean: ## Cleans up generated files/folders from the build.
	/bin/rm -rfv "public/dist" "rice-box.go" "${BINARY}"

clean-cache: ## Cleans up generated cache (speeds up during dev time).
	/bin/rm -rfv "public/.cache"

generate-watch: ## Generate public html/css/js when files change (faster, but larger files.)
	cd public && npm run watch

generate-frontend: ## Generate public html/css/js files for use in production (slower, smaller/minified files.)
	cd public && npm run build

generate-go: ## Generate go bundled files from frontend
	$(BIN)/rice -v embed-go

compile:
	go build -ldflags '-s -w' -tags netgo -installsuffix netgo -v -o "${BINARY}"

build: fetch clean clean-cache generate-frontend generate-go compile ## Builds the application (with generate.)
	echo

debug: fetch clean generate-frontend ## Runs the application in debug mode (with generate-dev.)
	go run *.go -d --http.limit 200000 --http.proxy
