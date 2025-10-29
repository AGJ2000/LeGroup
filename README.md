# ChitChat gRPC
README was AI generated 

🗨️ ChitChat – Distributed gRPC Chat System

Course: Distributed Systems – IT University of Copenhagen
Author: Anders Grangaard, Mathias Vestergaard Djurhuus, Theodor Monberg
Language: Go (v1.21+)
Framework: gRPC, Protocol Buffers

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
