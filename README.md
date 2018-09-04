# nwei-server

Things nwei cares about

### Run server
```
go run main.go --tls-cert certs/nwei-server.crt --tls-key certs/nwei-server.key --client-ca certs/nwei-ca.crt
```

### Run test client
```
go run testclient/client.go --tls-cert certs/nwei-client.crt --tls-key certs/nwei-client.key --server-ca certs/nwei-ca.crt --url https://localhost:8443
```
