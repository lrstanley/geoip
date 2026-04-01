# build-web image
FROM oven/bun:1 AS build-web

COPY . /build/
WORKDIR /build/web
ENV NODE_ENV=production
RUN bun install && bun run build

# build-go image
FROM golang:alpine AS build-go

RUN apk add --no-cache make
COPY . /build/
COPY --from=build-web /build/web/dist/ /build/web/dist/
WORKDIR /build
RUN make go-build

# runtime image
FROM alpine:3.23
RUN apk add --no-cache ca-certificates
COPY --from=build-go /build/geoip /usr/local/bin/geoip

# runtime params
VOLUME /data
EXPOSE 8080
WORKDIR /
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
CMD ["geoip", "--http.bind-addr", "0.0.0.0:8080", "--db.geoip-path", "/data/geoip.db", "--db.asn-path", "/data/asn.db"]
