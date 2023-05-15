# build-node image
FROM node:18 as build-node

COPY . /build/
WORKDIR /build
ENV NODE_ENV=production
RUN make node-build

# build-go image
FROM golang:alpine as build-go

RUN apk add --no-cache g++ make
COPY . /build/
COPY --from=build-node /build/public/dist/ /build/public/dist/
WORKDIR /build
RUN make go-build

# runtime image
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
# set up nsswitch.conf for Go's "netgo" implementation
# - https://github.com/docker-library/golang/blob/1eb096131592bcbc90aa3b97471811c798a93573/1.14/alpine3.12/Dockerfile#L9
RUN if [ ! -e /etc/nsswitch.conf ];then echo 'hosts: files dns' > /etc/nsswitch.conf;fi
COPY --from=build-go /build/geoip /usr/local/bin/geoip

# runtime params
VOLUME /data
EXPOSE 80
WORKDIR /
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
CMD ["geoip", "--http.bind-addr", "0.0.0.0:80", "--db.geoip-path", "/data/geoip.db", "--db.asn-path", "/data/asn.db"]
