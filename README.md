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

%% === Groups (folders) ===
namespace server {
  class server_main as "main.go"
}

namespace client {
  class client_main as "main.go"
}

namespace gRPC {
  class proto_file as "chitchat.proto"
  class proto_pb as "chitchat.pb.go"
  class proto_grpc as "chitchat_grpc.pb.go"
}

class grpc_runtime

%% === Relationships ===
server_main --> proto_grpc : imports / uses stubs
client_main --> proto_grpc : imports / uses stubs

proto_grpc ..> grpc_runtime : gRPC runtime
proto_pb ..> grpc_runtime : protobuf runtime
