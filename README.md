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

class server {
  main.go
  chitchat-server (binary)
}

class client {
  main.go
}

class GRPC {
  chitchat.proto
  chitchat.pb.go
  chitchat_grpc.pb.go
}

class demo_logs {
  client-anna.term.log
  client-bo.term.log
  client-cara.term.log
  server.log
}

class logs {
  0.log
  1.log
  2.log
}

class grpc_runtime

server --> GRPC : uses stubs
client --> GRPC : uses stubs

GRPC ..> grpc_runtime : generated against gRPC
server ..> grpc_runtime : server runtime
client ..> grpc_runtime : client runtime
