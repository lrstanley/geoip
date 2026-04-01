.DEFAULT_GOAL := build-all

export PROJECT := "geoip"
export PACKAGE := "github.com/lrstanley/geoip"

license:
	curl -sL https://liam.sh/-/gh/g/license-header.sh | bash -s

up: prepare go-up node-up
	@echo

prepare: clean node-prepare node-build go-prepare
	@echo

build-all: prepare go-build
	@echo

clean:
	/bin/rm -rfv "web/dist/*" "web/src/api" ${PROJECT}

docker-build:
	docker build \
		--pull \
		--tag ${PROJECT} \
		--force-rm .

node-up:
	cd web && bun update -i

node-fetch:
	cd web && bun install

node-prepare: node-fetch

node-debug: node-prepare
	cd web && bun run dev

node-build: node-prepare
	cd web && bun run build

node-preview: node-build
	cd web && bun run preview

node-test: node-prepare
	if [ -n "${CI}" ];then echo "output=${PWD}/web/tests/results/" >> "${GITHUB_OUTPUT}";fi
	cd web; if [ -n "${CI}" ];then bunx playwright install;fi
	cd web; bun run test:e2e

# backend
go-fetch:
	go mod tidy

go-up:
	go get -u ./...
	go mod tidy

go-prepare: go-fetch
	go generate -x ./...
	go run ${PACKAGE} generate-markdown > USAGE.md

go-dlv: go-prepare
	dlv debug \
		--headless --listen=:2345 \
		--api-version=2 --log \
		--allow-non-terminal-interactive \
		${PACKAGE} -- --debug

go-debug: go-prepare
	go run ${PACKAGE} \
		--http.limit 1000000 \
		--http.max-concurrent 0 \
		--dns.resolver "8.8.8.8" \
		--dns.resolver "1.1.1.1" \
		--debug

go-build: go-prepare
	CGO_ENABLED=0 \
	go build \
		-ldflags '-d -s -w -extldflags=-static' \
		-tags=netgo,osusergo,static_build \
		-installsuffix netgo \
		-trimpath \
		-o ${PROJECT} \
		${PACKAGE}
