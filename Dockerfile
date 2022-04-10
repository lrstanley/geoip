# build-node image
FROM node:16 as build-node

# for cache reasons, copy these separately.
COPY public/package.json /build/public/package.json
COPY public/package-lock.json /build/public/package-lock.json
COPY Makefile /build/

COPY public/ /build/public/
WORKDIR /build
RUN make fetch-node generate-node

# build-go image
FROM golang:alpine as build-go

# for cache reasons, copy these separately.
COPY go.mod /build/go.mod
COPY go.sum /build/go.sum

RUN apk add --no-cache g++ make
COPY . /build/
COPY --from=build-node /build/public/dist/ /build/public/dist/
WORKDIR /build
RUN make fetch-go build

# runtime image
FROM alpine:latest
RUN apk add --no-cache ca-certificates
# set up nsswitch.conf for Go's "netgo" implementation
# - https://github.com/docker-library/golang/blob/1eb096131592bcbc90aa3b97471811c798a93573/1.14/alpine3.12/Dockerfile#L9
RUN [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf
COPY --from=build-go /build/geoip /usr/local/bin/geoip

# runtime params
VOLUME /data
EXPOSE 80
WORKDIR /
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
CMD ["geoip", "--http.bind", "0.0.0.0:80", "--http.proxy", "--db", "/data/geoip.db"]
