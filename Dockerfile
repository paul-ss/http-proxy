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
RUN service postgresql start && \
    psql "user=proxy password=jw8s0F4 host=localhost port=5432 dbname=proxy_db" -f /opt/proxy/configs/sql/proxy.sql && \
    service postgresql stop

EXPOSE 8000
EXPOSE 8080

RUN /opt/proxy/scripts/gen-ca.sh && /opt/proxy/scripts/install-ca.sh

WORKDIR /opt/proxy
CMD service postgresql start && ./build/main
