package main

import (
	"context"
	"fmt"
	"net"
	"sync"

	"strings"
	"unicode/utf8"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	chitchat "github.com/AGJ2000/chitchat/gRPC"
	"google.golang.org/grpc"
)

type Participant struct {
	name     string
	subEpoch uint64
	joinTime uint64
	stream   chitchat.ChitChat_SubscribeServer
}
type Server struct {
	chitchat.UnimplementedChitChatServer
	L            uint64
	nextID       uint64
	participants map[string]*Participant
	sync.Mutex
}

func (s *Server) Join(ctx context.Context, req *chitchat.JoinRequest) (*chitchat.JoinAck, error) {
	s.Lock()
	id := fmt.Sprintf("%d", s.nextID)
	s.nextID++

	participant := &Participant{
		name:   req.DesiredName,
		stream: nil,
	}
	s.participants[id] = participant
	s.L++
	participant.joinTime = s.L

	pid := id
	name := participant.name
	lt := s.L

	ack := &chitchat.JoinAck{
		ParticipantId: pid,
		EffectiveName: name,
		LogicalTime:   lt,
	}
	joinMsg := &chitchat.Broadcast{
		Type:          chitchat.Broadcast_JOIN,
		ParticipantId: pid,
		Name:          name,
		Content:       fmt.Sprintf("Participant %s joined Chit Chat at logical time %d.", name, lt),
		LogicalTime:   lt,
	}
	s.Unlock()

	_ = s.broadcastMessage(joinMsg)

	fmt.Printf("Participant %s joined with ID %s at %d\n", name, pid, lt)
	return ack, nil
}

func (s *Server) Subscribe(req *chitchat.SubscribeRequest, stream chitchat.ChitChat_SubscribeServer) error {
	s.Lock()
	participant, exists := s.participants[req.ParticipantId]
	if !exists {
		s.Unlock()
		return status.Errorf(codes.NotFound, "participant not found")
	}

	participant.subEpoch++
	myEpoch := participant.subEpoch
	participant.stream = stream

	name := participant.name
	jt := participant.joinTime

	fmt.Printf("Participant %s subscribed (epoch: %d) at logical time %d\n", req.ParticipantId, myEpoch, s.L)

	s.Unlock()

	if err := stream.Send(&chitchat.Broadcast{
		Type:          chitchat.Broadcast_JOIN,
		ParticipantId: req.ParticipantId,
		Name:          name,
		Content:       fmt.Sprintf("Participant %s joined Chit Chat at logical time %d.", name, jt),
		LogicalTime:   jt,
	}); err != nil {
		fmt.Printf("Failed to send self join message to %s: %v\n", name, err)
	}

	// Keep the stream open until the client disconnects or context is cancelled
	<-stream.Context().Done()

	// On stream close, only clean up if this stream is still the active one
	s.Lock()
	participant, exists = s.participants[req.ParticipantId]
	if !exists {
		s.Unlock()
		return nil
	}
	if participant.subEpoch == myEpoch {
		delete(s.participants, req.ParticipantId)
		s.L++
		s.Unlock()
		fmt.Printf("Participant %s disconnected (epoch %d)\n", req.ParticipantId, myEpoch)
		return nil
	}
	s.Unlock()
	fmt.Printf("Participant %s stale stream closed (epoch %d ignored)\n", req.ParticipantId, myEpoch)
	return nil
}
func (s *Server) broadcastMessage(b *chitchat.Broadcast) error {
	s.Lock()
	slice := make([]*Participant, 0, len(s.participants))
	for _, p := range s.participants {
		if p.stream != nil {
			slice = append(slice, p)
		}
	}
	s.Unlock()
	for _, p := range slice {
		err := p.stream.Send(b)
		if err != nil {
			fmt.Printf("Failed to send message to %s: %v\n", p.name, err)
		}
	}
	return nil
}

func (s *Server) Publish(ctx context.Context, req *chitchat.PublishRequest) (*chitchat.PublishAck, error) {
	if utf8.RuneCountInString(req.Content) > 128 {
		return nil, status.Error(codes.InvalidArgument, "message too long (max 128 characters)")
	}

	if strings.TrimSpace(req.Content) == "" {
		return nil, status.Error(codes.InvalidArgument, "empty message")
	}

	s.Lock()
	p, ok := s.participants[req.ParticipantId]
	if !ok {
		s.Unlock()
		return nil, status.Error(codes.NotFound, "participant not found")
	}

	s.L = max(s.L,req.ClientTime) + 1
	current := s.L
	snapshotName := p.name
	s.Unlock()

	s.broadcastMessage(&chitchat.Broadcast{
		Type:          chitchat.Broadcast_MESSAGE,
		ParticipantId: req.ParticipantId,
		Name:          snapshotName,
		Content:       req.Content,
		LogicalTime:   current,
	})

	return &chitchat.PublishAck{LogicalTime: current}, nil
}

func main() {
	fmt.Println("ChitChat server is starting...")

	lis, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	srv := &Server{
		participants: make(map[string]*Participant),
	}
	chitchat.RegisterChitChatServer(s, srv)

	fmt.Println("Server is listening on 127.0.0.1:50051")
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
