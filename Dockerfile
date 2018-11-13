FROM alpine:latest

COPY ./toggly-server toggly-server

ENTRYPOINT [ "./toggly-server" ]