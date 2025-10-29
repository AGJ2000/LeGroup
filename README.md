# ChitChat gRPC

**Course:** Distributed Systems — IT University of Copenhagen  
**Authors:** Anders Grangaard Jensen · Mathias Vestergaard Djurhuus · Theodor Monberg

## How to run

> Prerequisites: Go **1.21+**

### 1) Build the server
```bash
go build -o chitchat-server ./server
```

### 2) Start the server
```bash
./chitchat-server
```
(Default bind: `127.0.0.1:50051`.)

### 3) Start one or more clients (each in its own terminal)
```bash
NAME=Anna go run ./client
NAME=Bo   go run ./client
NAME=Cara go run ./client
```

**Windows (PowerShell)**
```powershell
$env:NAME="Anna"; go run .\client
```
## 4) Repository Structure (Relationships-only UML)

```mermaid
classDiagram
direction LR

class ServerPkg {
  <<entrypoint>>
  main.go
}

class ClientPkg {
  <<entrypoint>>
  main.go
}

class ProtoPkg {
  <<proto>>
  chitchat.proto
  chitchat.pb.go
  chitchat_grpc.pb.go
}

class GrpcRuntime {
  <<library>>
}

ServerPkg --> ProtoPkg : uses stubs
ClientPkg --> ProtoPkg : uses stubs

ProtoPkg ..> GrpcRuntime : generated against gRPC
ServerPkg ..> GrpcRuntime : server runtime
ClientPkg ..> GrpcRuntime : client runtime
