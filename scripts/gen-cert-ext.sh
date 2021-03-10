#!/bin/sh

echo "authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names
[alt_names]
DNS.1 = $1" > certs/ext/domains.ext

openssl req -new -nodes -key certs/common/cert.key -subj "/C=RU/L=Moscow/O=Example-Certificates/CN=$1" |
openssl x509 -req -sha256 -days 1024 -CA certs/common/rootCA.crt -CAkey certs/common/rootCA.key -CAcreateserial -extfile certs/ext/domains.ext
