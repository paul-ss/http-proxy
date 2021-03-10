#!/bin/sh

openssl req -new -key certs/common/cert.key -subj "/CN=$1" |
openssl x509 -req -CA certs/common/rootCA.crt -CAkey certs/common/rootCA.key -CAcreateserial  -days 5000;
