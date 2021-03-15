FROM golang:buster AS build
WORKDIR /go/src/proxy
COPY . .
RUN go build -o build/main ./cmd/main.go

FROM ubuntu:20.04
ENV TZ=Europe/Moscow
ENV DEBIAN_FRONTEND=noninteractive

#ENV PGVER 12
RUN apt-get update && apt-get -y install postgresql postgresql-contrib ca-certificates
COPY --from=build /go/src/proxy /opt/proxy

USER postgres
RUN service postgresql start && psql -f /opt/proxy/configs/sql/init.sql && service postgresql stop
USER root

EXPOSE 8000

RUN /opt/proxy/scripts/gen-ca.sh && /opt/proxy/scripts/install-ca.sh

WORKDIR /opt/proxy
CMD service postgresql start && ./build/main
