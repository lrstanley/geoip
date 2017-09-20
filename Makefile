.DEFAULT_GOAL := build

GOPATH := $(shell go env | grep GOPATH | sed 's/GOPATH="\(.*\)"/\1/')
PATH := $(GOPATH)/bin:$(PATH)
export $(PATH)

BINARY=geoip
LD_FLAGS += -s -w

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

release: clean fetch ## Generate a release, but don't publish to GitHub.
	$(GOPATH)/bin/goreleaser --skip-validate --skip-publish

publish: clean fetch ## Generate a release, and publish to GitHub.
	$(GOPATH)/bin/goreleaser

snapshot: clean fetch ## Generate a snapshot release.
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

clean: ## Cleans up generated files/folders from the build.
	/bin/rm -rfv "public/dist" "rice-box.go" "${BINARY}"

generate-dev: ## Generate public html/css/js for use when developing (faster, but larger files).
	cd public && npm run dev

generate-watch: ## Generate public html/css/js when files change (faster, but larger files).
	cd public && npm run watch

generate: ## Generate public html/css/js files for use in production (slower, smaller/minified files).
	cd public && npm run build
	$(GOPATH)/bin/rice -v embed-go

debug: fetch clean generate-dev ## Runs the application in debug mode (with generate-dev).
	go run *.go -d --http.limit 200000 --update-url "http://hq.hq.liam.sh/tmp/GeoLite2-City.mmdb.gz"

build: fetch clean generate ## Builds the application (with generate).
	go build -ldflags "${LD_FLAGS}" -i -v -o ${BINARY}
