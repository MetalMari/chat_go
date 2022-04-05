# Chat app using grpc

## Description
Chat is  client-server app uses gRPC protocol for communication. It supports next operations:
 - Get list of users
 - Send a message to user
 - Get messages in queue and subscribe to new ones.

As data storage is used `etcd` storage.

## Checkout
Create new directory and go to it (optionally):
```bash
mkdir my_directory
cd my_directory
```
Clone repository:
```bash
git clone https://github.com/MetalMari/chat_go.git
```
Go to `chat_protos` directory and update submodules:
```bash
cd chat_grpc/chat_protos
git submodule update --init --recursive
```

## gRPC Code Generation
In `chat_go` directory:
- generate gRPC files for go using Makefile:
```bash
make proto
```
- generate gRPC code without make:
```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative chat_protos/chat.proto
```

## Usage Instructions
In terminal in `chat_go`directory run server:
```bash
go run server/main.go
```
Open new terminal, in `chat_go` directory run client:
```bash
go run client/main.go
```
Choose action:
1. to get list of users:
```bash
go run client/main.go --a=users
```
2. to send message:
```bash
go run client/main.go --a=message --from=login_from --to==login_to --body="Message content"
```
3. to subscribe for getting messages:
```bash
go run client/main.go --a=subscribe --login=login_to_subscribe
```
