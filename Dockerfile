FROM golang:latest as builder
WORKDIR /go/src/github.com/Toggly/core
COPY . ./
RUN version=$(git describe --always --tags) && \
    revision=${version}-$(date +%Y%m%d-%H:%M:%S) && \
    CGO_ENABLED=0 GOOS=linux go build -o toggly-server -ldflags "-X main.revision=${revision}" ./cmd/toggly-server

FROM alpine:latest
COPY --from=builder /go/src/github.com/Toggly/core/toggly-server toggly-server
EXPOSE 8080
ENTRYPOINT [ "./toggly-server" ]