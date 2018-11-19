FROM golang:alpine as builder
RUN apk update && apk add git gcc musl-dev
WORKDIR /go/src/github.com/Toggly/core

COPY . ./
RUN version=$(git describe --always --tags) && \
    revision=${version}-$(date +%Y%m%d-%H:%M:%S) && \
    go build -buildmode=plugin -o in-memory-cache.so ./internal/plugin/in-memory-cache/cache.go && \
    go build -o toggly-server -ldflags "-X main.revision=${revision}" ./cmd/toggly-server

FROM alpine:latest
COPY --from=builder /go/src/github.com/Toggly/core/in-memory-cache.so in-memory-cache.so
COPY --from=builder /go/src/github.com/Toggly/core/toggly-server toggly-server
EXPOSE 8080
ENTRYPOINT [ "./toggly-server" ]
