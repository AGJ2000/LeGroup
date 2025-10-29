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
## ChitChat gRPC-arkitektur

```mermaid
flowchart LR

  %% Server
  subgraph S["server"]
    direction TB
    S1["main.go"]
    S2["chitchat-server (binary)"]
  end

  %% Client
  subgraph C["client"]
    direction TB
    C1["main.go"]
  end

  %% gRPC stubs / proto
  subgraph G["gRPC"]
    direction TB
    G1["chitchat.proto"]
    G2["chitchat.pb.go"]
    G3["chitchat_grpc.pb.go"]
  end

  %% Runtime (lib)
  subgraph R["gRPC runtime"]
    direction TB
    R1["grpc_runtime"]
  end

  %% Demo logs
  subgraph L["demo_logs"]
    direction TB
    L1["client-anna.term.log"]
    L2["client-bo.term.log"]
    L3["client-cara.term.log"]
    L4["server.log"]
  end

  %% Relations
  S -->|uses stubs| G
  C -->|uses stubs| G
  G -->|generated against gRPC| R
  S -.->|server runtime| R
  C -.->|client runtime| R
