#!/bin/bash

echo "Generating..."

openssl genpkey -algorithm RSA -out ./hmac_key.pem -pkeyopt rsa_keygen_bits:4096

openssl req -new -newkey rsa:4096 -x509 -sha256 -days 365 -nodes -keyout api_key.pem -out api_cert.pem \
    -subj "/C=NL/ST=Noord-Holland/L=Amsterdam/O=ExOrg/OU=IT/CN=localhost"