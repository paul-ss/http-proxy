#!/bin/sh

openssl req -new -key cert.key -subj "/CN=$1" |
openssl x509 -req -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -out "certs/$1.crt" -days 5000;
