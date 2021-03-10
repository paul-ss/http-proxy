#!/bin/sh

cp certs/common/rootCA.crt /usr/local/share/ca-certificates/
update-ca-certificates