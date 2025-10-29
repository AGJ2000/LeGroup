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

class ServerFolder as "server/" {
  <<directory>>
  main.go
  chitchat-server (binary)
}

class ClientFolder as "client/" {
  <<directory>>
  main.go
}

class GRPCFolder as "gRPC/" {
  <<proto>>
  chitchat.proto
  chitchat.pb.go
  chitchat_grpc.pb.go
}

class DemoLogs as "demo-logs/" {
  <<logs>>
  client-anna.term.log
  client-bo.term.log
  client-cara.term.log
  server.log
}

class LogsFolder as "logs/" {
  <<logs>>
  0.log
  1.log
  2.log
}

class GrpcRuntime {
  <<library>>
  google.golang.org/grpc
  status/codes
}

%% Relationships
ServerFolder --> GRPCFolder : uses stubs
ClientFolder --> GRPCFolder : uses stubs

GRPCFolder ..> GrpcRuntime : generated against gRPC
ServerFolder ..> GrpcRuntime : server runtime
ClientFolder ..> GrpcRuntime : client runtime
