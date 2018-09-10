# nwei-server

A place to keep data nwei cares about.

## Authentication

Provided via mutual TLS.

## Example usage

### Run server
```
go run main.go \
  --tls-cert certs/nwei-server.crt \
  --tls-key certs/nwei-server.key \
  --client-ca certs/nwei-ca.crt \
  --port 8443
```

### Run testclient
```
go run testclient/client.go \
  --tls-cert certs/nwei-client.crt \
  --tls-key certs/nwei-client.key \
  --server-ca certs/nwei-ca.crt \
  --url https://localhost:8443/healthcheck
```
