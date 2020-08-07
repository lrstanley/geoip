# build-node image
FROM node:latest as build-node
RUN mkdir /build
COPY public/ /build/
WORKDIR /build
RUN npm install && npm run build

# build-go image
FROM golang:latest as build-go
RUN mkdir -p /build/public/dist
COPY . /build/
COPY --from=build-node /build/dist/ /build/public/dist/
WORKDIR /build
RUN make fetch generate-go compile

FROM alpine:latest

RUN apk add --no-cache ca-certificates

# set up nsswitch.conf for Go's "netgo" implementation
# - https://github.com/docker-library/golang/blob/1eb096131592bcbc90aa3b97471811c798a93573/1.14/alpine3.12/Dockerfile#L9
RUN [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf

COPY --from=build-go /build/geoip /usr/local/bin/geoip

VOLUME /data
EXPOSE 80
WORKDIR /
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
CMD ["geoip", "--http", "0.0.0.0:80", "--behind-proxy", "--db", "/data/store.db"]