.DEFAULT_GOAL := build-all

export PROJECT := "geoip"
export PACKAGE := "github.com/lrstanley/geoip"

license:
	curl -sL https://liam.sh/-/gh/g/license-header.sh | bash -s

prepare: clean node-prepare node-build go-prepare
	@echo

build-all: prepare go-build
	@echo

clean:
	/bin/rm -rfv "public/dist/*" ${PROJECT}

docker-build:
	docker build \
		--pull \
		--tag ${PROJECT} \
		--force-rm .

# frontend
node-fetch:
	command -v pnpm >/dev/null >&2 || npm install \
		--no-audit \
		--no-fund \
		--quiet \
		--global pnpm
	cd public; pnpm install

node-upgrade-deps:
	cd public && \
		pnpm up -i

node-prepare: node-fetch
	cd public; pnpm exec openapi \
		--input ../internal/handlers/apihandler/openapi_v2.yaml \
		--output src/lib/api/openapi/ \
		--client fetch \
		--useOptions \
		--indent 2 \
		--name HTTPClient

node-lint: node-build # needed to generate eslint auto-import ignores.
	cd public; pnpm exec eslint \
		--ignore-path ../.gitignore \
		--ext .js,.ts,.vue .

node-test: node-prepare
	@echo "output=${PWD}/public/tests/results/" >> "${GITHUB_OUTPUT}"
	cd public; if [ -n "${CI}" ];then pnpm exec playwright install-deps;fi
	cd public; pnpm exec playwright test

node-debug: node-prepare
	cd public; pnpm exec vite

node-build: node-prepare
	cd public; pnpm exec vite build

node-preview: node-build
	cd public; pnpm exec vite preview

# backend
go-fetch:
	go mod download
	go mod tidy

go-upgrade-deps:
	go get -u ./...
	go mod tidy

go-upgrade-deps-patch:
	go get -u=patch ./...
	go mod tidy

go-prepare: go-fetch
	go generate -x ./...
	{ echo '## :gear: Usage'; go run ${PACKAGE} --generate-markdown --db.license-key ""; } > USAGE.md

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
		--http.metrics \
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
