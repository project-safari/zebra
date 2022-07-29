#!/bin/bash -xe

OUTPUT=$1
DOMAIN=$2
IPSAN=$3

cleanup() {
  echo "Removing ${OUTPUT}/san"
  rm -f ${OUTPUT}/san.cnf
}

echo "[req]
default_bits  = 4096
distinguished_name = req_distinguished_name
req_extensions = req_ext
x509_extensions = v3_req
prompt = no
[req_distinguished_name]
countryName = US
stateOrProvinceName = California
localityName = San Jose
organizationName = zebra-server
commonName = ${DOMAIN}
[req_ext]
subjectAltName = @alt_names
[v3_req]
subjectAltName = @alt_names
[alt_names]
IP.1 = ${IPSAN}
IP.2 = 127.0.0.1
DNS.1 = ${DOMAIN}
" > ${OUTPUT}/san.cnf

# CA Certificate
if [[ ! -f "${OUTPUT}/zebra-ca.crt" || ! -f "${OUTPUT}/zebra-ca.key" ]]; then
openssl req \
    -newkey rsa:4096 \
    -nodes \
    -days 3650 \
    -x509 \
    -keyout ${OUTPUT}/zebra-ca.key \
    -out ${OUTPUT}/zebra-ca.crt \
    -subj "/CN=*"
fi

# Server Key
if [[ ! -f "${OUTPUT}/zebra-server.crt" || ! -f "${OUTPUT}/zebra-server.key" ]]; then
openssl req \
    -newkey rsa:4096 \
    -nodes \
    -keyout ${OUTPUT}/zebra-server.key \
    -out ${OUTPUT}/zebra-server.csr \
    -config ${OUTPUT}/san.cnf

# Server Certificate
openssl x509 \
    -req \
    -days 3650 \
    -sha256 \
    -in ${OUTPUT}/zebra-server.csr \
    -CA ${OUTPUT}/zebra-ca.crt \
    -CAkey ${OUTPUT}/zebra-ca.key \
    -CAcreateserial \
    -out ${OUTPUT}/zebra-server.crt \
    -extensions req_ext \
    -extfile ${OUTPUT}/san.cnf
fi

# Client Key
if [[ ! -f "${OUTPUT}/zebra-client.crt" || ! -f "${OUTPUT}/zebra-client.key" ]]; then
openssl req \
    -newkey rsa:4096 \
    -nodes \
    -keyout ${OUTPUT}/zebra-client.key \
    -out ${OUTPUT}/zebra-client.csr \
    -subj "/C=US/ST=California/L=San Jose/O=safari/OU=zebra-client/CN=*"

openssl x509 \
    -req \
    -days 3650 \
    -sha256 \
    -in ${OUTPUT}/zebra-client.csr \
    -CA ${OUTPUT}/zebra-ca.crt \
    -CAkey ${OUTPUT}/zebra-ca.key \
    -CAcreateserial \
    -out ${OUTPUT}/zebra-client.crt
fi

rm -f *.csr

trap cleanup EXIT