.DEFAULT_GOAL := build

GOPATH := $(shell go env | grep GOPATH | sed 's/GOPATH="\(.*\)"/\1/')
PATH := $(GOPATH)/bin:$(PATH)
export $(PATH)

BINARY=geoip
COMPRESS_CONC ?= $(shell nproc)
VERSION=$(shell git describe --tags --abbrev=0 2>/dev/null | sed -r "s:^v::g")
RSRC=README_TPL.md
ROUT=README.md

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

readme-gen: ## Generates readme from template file
	cp -av "${RSRC}" "${ROUT}"
	sed -ri -e "s:\[\[tag\]\]:${VERSION}:g" -e "s:\[\[os\]\]:linux:g" -e "s:\[\[arch\]\]:amd64:g" "${ROUT}"

release: clean fetch generate ## Generate a release, but don't publish to GitHub.
	$(GOPATH)/bin/goreleaser --skip-validate --skip-publish

publish: clean fetch generate ## Generate a release, and publish to GitHub.
	$(GOPATH)/bin/goreleaser

snapshot: clean fetch generate ## Generate a snapshot release.
	$(GOPATH)/bin/goreleaser --snapshot --skip-validate --skip-publish

update-deps: fetch ## Updates all dependencies to the latest available versions.
	$(GOPATH)/bin/govendor add +external
	$(GOPATH)/bin/govendor remove +unused
	$(GOPATH)/bin/govendor update +external

fetch: ## Fetches the necessary dependencies to build.
	test -f $(GOPATH)/bin/govendor || go get -u -v github.com/kardianos/govendor
	test -f $(GOPATH)/bin/goreleaser || go get -u -v github.com/goreleaser/goreleaser
	test -f $(GOPATH)/bin/rice || go get -u -v github.com/GeertJohan/go.rice/rice
	$(GOPATH)/bin/govendor sync
	test -d public/node_modules || (cd public && npm install)

clean: ## Cleans up generated files/folders from the build.
	/bin/rm -rfv "dist" "public/dist" "rice-box.go" "${BINARY}"

compress: ## Uses upx to compress release binaries (if installed, uses all cores/parallel comp.)
	(which upx > /dev/null && find dist/*/* | xargs -I{} -n1 -P ${COMPRESS_CONC} upx --best "{}") || echo "not using upx for binary compression"

generate-dev: ## Generate public html/css/js for use when developing (faster, but larger files.)
	cd public && npm run dev

generate-watch: ## Generate public html/css/js when files change (faster, but larger files.)
	cd public && npm run watch

generate: ## Generate public html/css/js files for use in production (slower, smaller/minified files.)
	cd public && npm run build
	$(GOPATH)/bin/rice -v embed-go

debug: fetch clean generate-dev ## Runs the application in debug mode (with generate-dev.)
	go run *.go -d --http.limit 200000 --update-url "http://hq.hq.liam.sh/tmp/GeoLite2-City.mmdb.gz"

build: fetch clean generate ## Builds the application (with generate.)
	go build -ldflags '-d -s -w' -tags netgo -installsuffix netgo -v -x -o "${BINARY}"
