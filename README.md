# Golang Key-Value Database
A simple, no nonsense key-value database server in Go using gRPC.

## Setting up the environment
This project requires gRPC and the protobuf Go generator plugin:
```
go get google.golang.org/grpc
go get -u github.com/golang/protobuf/protoc-gen-go
```

Generating the Go gRPC API:
```
protoc -I api/ api.proto --go_out=plugins=grpc:api
```