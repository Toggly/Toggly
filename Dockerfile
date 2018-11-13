FROM alpine:latest

COPY ./toggly-server toggly-server

EXPOSE 8080

ENTRYPOINT [ "./toggly-server" ]