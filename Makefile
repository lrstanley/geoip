.DEFAULT_GOAL := build

DIRS=bin
BINARY=geoip

VERSION=$(shell git describe --tags --always --abbrev=0 --match=v* 2> /dev/null | sed -r "s:^v::g" || echo 0)

$(info $(shell mkdir -p $(DIRS)))
BIN=$(CURDIR)/bin
export GOBIN=$(CURDIR)/bin


help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

fetch-go: ## Fetches the necessary Go dependencies to build.
	go mod download
	go mod tidy

fetch-node: ## Fetches the necessary NodeJS dependencies to build.
	test -d public/node_modules || (cd public && npm run run-install)

upgrade-deps: ## Upgrade all dependencies to the latest version.
	go get -u ./...

upgrade-deps-patch: ## Upgrade all dependencies to the latest patch release.
	go get -u=patch ./...

clean: ## Cleans up generated files/folders from the build.
	/bin/rm -rfv "public/dist" "${BINARY}"
	mkdir -p public/dist
	touch public/dist/.gitkeep

generate-node: ## Generate public html/css/js files for use in production (slower, smaller/minified files.)
	cd public && npm run build
	ls -lah public/dist/

frontend-watch: ## Use this to spin up vite, and proxy calls to the backend.
	cd public && npm run server

debug: fetch-go fetch-node clean ## Runs the application in debug mode (with generate-dev.)
	go run *.go -d --http.limit 200000 --http.proxy

prepare: fetch-go fetch-node clean generate-node ## Runs preparation steps for build.
	@echo

build: prepare ## Builds the application (with generate.)
	go generate ./...
	CGO_ENABLED=0 go build -ldflags '-d -s -w -extldflags=-static' -tags=netgo,osusergo,static_build -installsuffix netgo -buildvcs=false -trimpath -o "${BINARY}"
