# nwei-server

Things nwei cares about

### Run server
```
go run main.go --tls-cert certs/nwei-server.crt --tls-key certs/nwei-server.key
```

### Run test client
```
go run testclient/client.go --tls-ca certs/nwei-ca.crt --url https://localhost:8443
```
