package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	chitchat "github.com/AGJ2000/chitchat/gRPC"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)
func max64(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}
func main() {
	var Lc uint64 = 0
	var cancelSub context.CancelFunc

	fmt.Println("ChitChat client is starting...")

	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	client := chitchat.NewChitChatClient(conn)

	nameFlag := flag.String("name", "", "Desired username")
	flag.Parse()

	envName, _ := os.LookupEnv("NAME")
	chosenName := strings.TrimSpace(envName)
	if chosenName == "" {
		chosenName = strings.TrimSpace(*nameFlag)
	}
	if chosenName == "" {
		chosenName = "Anonymous"
	}
	fmt.Println("Using name:", chosenName)

	joinCtx, joinCancel := context.WithTimeout(context.Background(), 3*time.Second)
	joinResp, err := client.Join(joinCtx, &chitchat.JoinRequest{DesiredName: chosenName})
	joinCancel()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to join chat:", err)
		os.Exit(1)
	}

	fmt.Printf("Joined chat with ID: %s and Effective Name: %s\n", joinResp.ParticipantId, joinResp.EffectiveName)
	fmt.Printf("Logical Time: %d\n", joinResp.LogicalTime)
	Lc = max64(Lc, joinResp.LogicalTime) + 1

	id := joinResp.ParticipantId

	// set up per-participant log file
	_ = os.MkdirAll("logs", 0755)
	logPath := filepath.Join("logs", id+".log")
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open log file %s: %v", logPath, err)
	}
	defer logFile.Close()

	subReq := &chitchat.SubscribeRequest{ParticipantId: id}
	
	var ctxSub context.Context
	ctxSub, cancelSub = context.WithCancel(context.Background())
	stream, err := client.Subscribe(ctxSub, subReq)
	if err != nil {
		log.Fatalf("failed to subscribe: %v", err)
	}

	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					return
				}
				if s, ok := status.FromError(err); ok {
					if s.Code() == codes.Canceled || s.Code() == codes.Unavailable {
						return
					}
				}
				fmt.Printf("stream closed: %v\n", err)
				return
			}
			switch msg.Type {
			case chitchat.Broadcast_JOIN:
				fmt.Printf("Participant %s joined Chit Chat at logical time %d\n", msg.Name, msg.LogicalTime)
				fmt.Fprintf(logFile, "JOIN [%d] %s\n", msg.LogicalTime, msg.Content)
			case chitchat.Broadcast_LEAVE:
				fmt.Printf("Participant %s left Chit Chat at logical time %d\n", msg.Name, msg.LogicalTime)
				fmt.Fprintf(logFile, "LEAVE [%d] %s\n", msg.LogicalTime, msg.Content)
			case chitchat.Broadcast_MESSAGE:
				fmt.Printf("[%d] %s: %s\n", msg.LogicalTime, msg.Name, msg.Content)
				fmt.Fprintf(logFile, "[%d] %s: %s\n", msg.LogicalTime, msg.Name, msg.Content)
			default:
				fmt.Printf("Unknown message type received: %v\n", msg.Type)
			}
			Lc = max64(Lc, msg.LogicalTime) + 1
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("read error: %v\n", err)
			return
		}
		line = strings.TrimSpace(line)

		if !utf8.ValidString(line) {
			fmt.Println("invalid UTF-8 input; not sent")
			continue
		}
		if utf8.RuneCountInString(line) > 128 {
			fmt.Println("message too long (max 128 characters); not sent")
			continue
		}

		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "/") {
			switch line {
			case "/quit":
				if cancelSub != nil {
					cancelSub()
				}
				leaveCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				_, _ = client.Leave(leaveCtx, &chitchat.LeaveRequest{ParticipantId: id})
				cancel()
				fmt.Println("Goodbye!")
				return
			default:
				fmt.Printf("unknown command %q\n", line)
			}
			continue
		}
		Lc++
		pubCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		pubAck, err := client.Publish(pubCtx, &chitchat.PublishRequest{
			ParticipantId: id,
			Content:       line,
			ClientTime:    Lc,
		})
		cancel()
		if err != nil {
			fmt.Printf("Publish failed: %v\n", err)
		} else {
			Lc = max64(Lc, pubAck.LogicalTime) + 1
		}
	}
}
