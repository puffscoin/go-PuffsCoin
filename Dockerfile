# Build Geth in a stock Go builder container
FROM golang:1.12-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers

ADD . /go-puffscoin
RUN cd /go-puffscoin && make gpuffs

# Pull Geth into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go-puffscoin/build/bin/gpuffs /usr/local/bin/

EXPOSE 11363 11364 31313 31313/udp
ENTRYPOINT ["gpuffs"]
