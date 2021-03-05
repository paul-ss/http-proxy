#!/bin/sh

openssl genrsa -out rootCA.key 2048;
openssl req -x509 -new -key rootCA.key -days 10000  -subj '/CN=paul-s-http-proxy-ca' -out rootCA.crt;
openssl genrsa -out cert.key 2048;
mkdir serts/
