.DEFAULT_GOAL := build-all

export PROJECT := "geoip"
export PACKAGE := "github.com/lrstanley/geoip"

license:
	curl -sL https://liam.sh/-/gh/g/license-header.sh | bash -s

build-all: clean node-fetch go-fetch node-build go-build
	@echo

clean:
	/bin/rm -rfv "public/dist/*" ${PROJECT}

docker-build:
	docker build \
		--tag ${PROJECT} \
		--force-rm .

# frontend
node-fetch:
	cd public; npm install --no-fund --no-audit

node-debug:
	cd public; npm run server

node-build: node-fetch
	cd public; npm run build

# backend
go-prepare:
	go generate -x ./...

go-fetch:
	go mod download
	go mod tidy

go-upgrade-deps:
	go get -u ./...
	go mod tidy

go-upgrade-deps-patch:
	go get -u=patch ./...
	go mod tidy

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
		--log.quiet
		# --debug \

go-build: go-prepare go-fetch
	CGO_ENABLED=0 \
	go build \
		-ldflags '-d -s -w -extldflags=-static' \
		-tags=netgo,osusergo,static_build \
		-installsuffix netgo \
		-trimpath \
		-o ${PROJECT} \
		${PACKAGE}
