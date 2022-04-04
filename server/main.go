package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "chat_go/chat_protos"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement chat.ChatServer.
type server struct {
	pb.UnimplementedChatServer
}

// GetUsers returns list of users.
func (s *server) GetUsers(ctx context.Context, in *pb.GetUsersRequest) (*pb.GetUsersReply, error) {
	return &pb.GetUsersReply{Users: []*pb.User{{Login: "userA", FullName: "aa aa"},
		{Login: "userB", FullName: "bb bb"}}}, nil
}

// SendMessage gets message and returns simple string if the message from client is received.
func (s *server) SendMessage(ctx context.Context, in *pb.SendMessageRequest) (*pb.SendMessageReply, error) {
	return &pb.SendMessageReply{Status: in.Message.LoginTo +
		" received message from " + in.Message.LoginFrom}, nil
}

// Subscribe returns stream of messages by subscription.
func (s *server) Subscribe(resp *pb.SubscribeRequest, stream pb.Chat_SubscribeServer) error {
	messages := [2]*pb.Message{
		{LoginFrom: resp.Login, LoginTo: "B", CreatedAt: 1234, Body: "Hello, B!"},
		{LoginFrom: resp.Login, LoginTo: "D", CreatedAt: 1234, Body: "Hello, D!"},
	}
	for _, message := range messages {
		if err := stream.Send(message); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterChatServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
