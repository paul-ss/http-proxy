#!/bin/sh
mkdir -p certs/common certs/ext
openssl genrsa -out certs/common/rootCA.key 2048;
openssl req -x509 -new -key certs/common/rootCA.key -days 10000  -subj '/CN=paul-s-http-proxy-ca' -out certs/common/rootCA.crt;
openssl genrsa -out certs/common/cert.key 2048;


#openssl req -x509 -nodes -new -sha256 -days 1024 -newkey rsa:2048 -keyout RootCA.key -out RootCA.pem -subj "/C=RU/CN=Example"
#openssl x509 -outform pem -in RootCA.pem -out RootCA.crt