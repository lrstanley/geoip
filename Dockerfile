FROM --platform=$BUILDPLATFORM oven/bun:1-alpine AS build-web
# required by tree-sitter
RUN apk add --no-cache python3 make g++
WORKDIR /build/web
COPY ./web/bun.lock ./web/package.json /build/web/
# see https://github.com/thedotmack/claude-mem/issues/2280
RUN CXXFLAGS="-std=c++20" CXX_STANDARD=20 bun install --frozen-lockfile
COPY ./internal/handlers/apihandler/openapi_v2.yaml /build/internal/handlers/apihandler/openapi_v2.yaml
COPY ./web /build/web
RUN bun run build

FROM golang:alpine AS build-go
WORKDIR /build
COPY go.* /build/
RUN go mod download
COPY . /build/
COPY --from=build-web /build/web/dist/ /build/web/dist/
RUN apk add --no-cache make
RUN make go-build

FROM alpine:3.24
RUN apk add --no-cache ca-certificates
COPY --from=build-go /build/geoip /usr/local/bin/geoip

# runtime params
VOLUME /data
EXPOSE 8080
WORKDIR /
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
CMD ["geoip", "--http.bind-addr", "0.0.0.0:8080", "--db.geoip-path", "/data/geoip.db", "--db.asn-path", "/data/asn.db"]
