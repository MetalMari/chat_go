package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "chat_go/chat_protos"
	st "chat_go/storage"

	"google.golang.org/grpc"
)

const (
	dialTimeout = 5 * time.Second
)

var (
	servHost = flag.String("servHost", "localhost", "The server host")
	servPort = flag.Int("servPort", 50051, "The server port")

	storHost = flag.String("storHost", "localhost", "The storage host")
	storPort = flag.Int("storPort", 2379, "The storage port")
)

// server is used to implement chat.ChatServer.
type server struct {
	pb.UnimplementedChatServer
	storage st.EtcdStorage
}

// GetUsers returns list of users.
func (s *server) GetUsers(ctx context.Context, in *pb.GetUsersRequest) (*pb.GetUsersReply, error) {
	users, err := s.storage.GetUsers()
	if err != nil {
		log.Fatal(err)
	}
	var uu []*pb.User
	for _, u := range users {
		uu = append(uu, &pb.User{Login: u.Login, FullName: u.FullName})
	}
	return &pb.GetUsersReply{Users: uu}, nil
}

// SendMessage gets message and returns simple string if the message from client is received.
func (s *server) SendMessage(ctx context.Context, in *pb.SendMessageRequest) (*pb.SendMessageReply, error) {
	created_at := int32(time.Now().Unix())
	m := st.Message{LoginFrom: in.Message.LoginFrom,
		LoginTo:   in.Message.LoginTo,
		CreatedAt: created_at,
		Body:      in.Message.Body}
	s.storage.CreateMessage(m)
	statusMessage := in.Message.LoginTo + " received message from " + in.Message.LoginFrom
	return &pb.SendMessageReply{Status: statusMessage}, nil
}

// Subscribe returns stream of messages by subscription.
func (s *server) Subscribe(resp *pb.SubscribeRequest, stream pb.Chat_SubscribeServer) error {
	defer log.Printf("Finish subscription for user %v", resp.Login)
	for {
		messages, err := s.storage.GetMessages(resp.Login)
		if err != nil {
			log.Fatal(err)
		}
		for _, mes := range messages {
			time.Sleep(2 * time.Second)
			message := &pb.Message{
				LoginFrom: mes.LoginFrom,
				LoginTo:   mes.LoginTo,
				CreatedAt: mes.CreatedAt,
				Body:      mes.Body}
			log.Println("Send message..")
			if err := stream.Send(message); err != nil {
				return err
			}
			s.storage.DeleteMessage(mes)
		}
	}
}

// Creates storage on defined address and port.
func createStorage(stor_host string, stor_port int) (storage *st.EtcdStorage) {
	endpoint := fmt.Sprintf("%v:%v", stor_host, stor_port)
	endpoints := []string{endpoint}
	stor, err := st.NewEtcdStorage(endpoints, dialTimeout)
	if err != nil {
		log.Fatalf("failed to create storage: %v", err)
	}
	return stor
}

// Creates users and saves them to storage.
func fillUsers(stor *st.EtcdStorage) {
	users, err := stor.GetUsers()
	if err != nil {
		log.Fatal(err)
	}
	if len(users) == 0 {
		for i := 0; i < 4; i++ {
			login := fmt.Sprintf("user%d", i)
			full_name := fmt.Sprintf("user%d_user%d", i, i)
			user := st.User{Login: login, FullName: full_name}
			stor.CreateUser(user)
		}
	}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *servHost, *servPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	storage := createStorage(*storHost, *storPort)
	fillUsers(storage)
	pb.RegisterChatServer(s, &server{storage: *storage})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
