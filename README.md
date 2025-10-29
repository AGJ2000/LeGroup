# ChitChat gRPC
README was AI generated 

Distributed chat service with Go and gRPC.

go build -o chitchat-server ./server
./chitchat-server

NAME=Bob go run ./client 


🗨️ ChitChat – Distributed gRPC Chat System

Course: Distributed Systems – IT University of Copenhagen
Author: Anders Grangaard, Math
Language: Go (v1.21+)
Framework: gRPC, Protocol Buffers

⸻

📦 Project Overview

ChitChat is a distributed chat system implemented using gRPC.
It allows multiple concurrent clients to join, exchange messages, and leave the chat,
all coordinated through a single server instance.

System Features

Spec	Requirement	Implemented
S1	gRPC for all communication	✅
S2	One service, multiple clients	✅
S3	Valid UTF-8 messages ≤ 128 characters	✅
S4	Broadcast each message with logical timestamp	✅
S5	Broadcast JOIN to all (incl. new participant)	✅
S6	Broadcast LEAVE to all remaining participants	✅
S7	Each client displays + logs all messages	✅


⸻

🧱 Project Structure

ChitChat/
├── client/          # Client process code
│   └── main.go
├── server/          # Server process code
│   └── main.go
├── gRPC/            # .proto file and generated Go stubs
│   └── chitchat.proto
├── logs/            # Client-side logs (auto-created)
│   ├── 0.log
│   ├── 1.log
│   └── 2.log
└── README.md


⸻

⚙️ Build & Run Instructions 

1️⃣  Build the server

go build -o chitchat-server ./server

2️⃣ Start the server

./chitchat-server | tee demo-logs/server.log

(Using tee will store a server.log file with all runtime logs)

3️⃣ Start multiple clients

In separate terminals:
NAME=Anna go run ./client | tee demo-logs/client-anna.log
NAME=Bo go run ./client   | tee demo-logs/client-bo.log
NAME=Cara go run ./client | tee demo-logs/client-cara.log

Each client will also create its own log file under logs/:

logs/
 ├── 0.log
 ├── 1.log
 └── 2.log


⸻

💬 Example Interaction

Server output

ChitChat server is starting...
Server is listening on 127.0.0.1:50051
Participant Anna joined with ID 0 at 1
Participant Bo joined with ID 1 at 3
Participant Cara joined with ID 2 at 5
Participant 2 disconnected (epoch 1)

Client Anna

Using name: Anna
Joined chat with ID: 0 and Effective Name: Anna
Logical Time: 1
Participant Anna joined Chit Chat at logical time 1
Participant Bo joined Chit Chat at logical time 3
[6] Bo: hej fra Bo
Participant Cara left Chit Chat at logical time 7
Goodbye!

Client Bo

Using name: Bo
Joined chat with ID: 1 and Effective Name: Bo
Logical Time: 3
Participant Bo joined Chit Chat at logical time 3
> hej fra Bo
[6] Bo: hej fra Bo
Goodbye!
⸻
