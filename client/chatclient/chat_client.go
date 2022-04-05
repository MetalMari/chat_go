package chatclient

import (
	"context"
	"io"
	"log"

	pb "chat_go/chat_protos"

	"google.golang.org/grpc"
)

type Client struct {
	Endpoint string

	conn   io.Closer
	client pb.ChatClient
}

// NewClient creates new chat client.
func NewClient(endpoint string) (*Client, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Client{Endpoint: endpoint, conn: conn, client: pb.NewChatClient(conn)}, nil
}

// Close closes gRPC connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

// GetUsers gets and prints list of users.
func (c *Client) GetUsers() {
	ctx := context.Background()

	r, err := c.client.GetUsers(ctx, &pb.GetUsersRequest{})
	if err != nil {
		log.Fatalf("could not get users: %v", err)
	}
	log.Printf("Users: %s", r.Users)
}

type Message struct {
	LoginFrom string
	LoginTo   string
	CreatedAt int32
	Body      string
}

// SendMessage sends message contained sender's login, recipient's login,
// creation timestamp, body-content and prints response from server.
func (c *Client) SendMessage(m *Message) {
	mes := &pb.Message{LoginFrom: m.LoginFrom, LoginTo: m.LoginTo, CreatedAt: m.CreatedAt, Body: m.Body}
	ctx := context.Background()
	r, err := c.client.SendMessage(ctx, &pb.SendMessageRequest{Message: mes})
	if err != nil {
		log.Fatalf("could not send message: %v", err)
	}
	log.Printf("Status: %s", r.Status)
}

// Gets and prints all messages, given in stream by subscription.
func (c *Client) Subscribe(login string) {
	ctx := context.Background()
	stream, err := c.client.Subscribe(ctx, &pb.SubscribeRequest{Login: login})
	if err != nil {
		log.Fatalf("Cannot receive: %v", err)
	}
	for {
		mes, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Cannot receive: %v", err)
		}
		log.Printf("Message: %s", mes)
	}
}
