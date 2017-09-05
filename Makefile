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
	$(GOPATH)/bin/govendor sync

clean: ## Cleans up generated files/folders from the build.
	/bin/rm -rfv "dist/" "static/" "${BINARY}"

generate-dev:
	cd public && npm run dev
	mkdir -vp "static"
	cp -av public/index.html static/
	cp -av public/dist static/

generate:
	cd public && npm run build
	mkdir -vp "static"
	cp -av public/index.html static/
	cp -av public/dist static/
	$(GOPATH)/bin/rice -v embed-go

debug: fetch clean generate-dev ## Runs the application in debug mode.
	go run *.go -d --http.limit 2000 --update-url "http://hq.liam.sh/tmp/GeoLite2-City.mmdb.gz"

build: fetch clean generate ## Builds the application.
	go build -ldflags "${LD_FLAGS}" -i -v -o ${BINARY}
