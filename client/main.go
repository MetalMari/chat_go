package main


import (
	"context"
	"flag"
	"io"
	"log"
	"time"

    "google.golang.org/grpc"
	pb "chat_go/chat_protos"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	from = flag.String("from", "A", "Login from")
	to = flag.String("to", "B", "Login to")
	body = flag.String("body", "Hello", "Message body")
	login = flag.String("login", "userA", "Login to subscribe")
	a = flag.String("a", "", "Action")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewChatClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var p string = *a
	switch p {
	case "users":
		r, err := c.GetUsers(ctx, &pb.GetUsersRequest{})
	    if err != nil {
		    log.Fatalf("could not get users: %v", err)
	    }
	    log.Printf("Users: %s", r.GetUsers())
	case "message":
	    m := &pb.Message{LoginFrom: *from, LoginTo: *to, CreatedAt: 1234, Body: *body}

	    k, err := c.SendMessage(ctx, &pb.SendMessageRequest{Message: m})
        if err != nil {
            log.Fatalf("could not greet: %v", err)
        }
        log.Printf("Status: %s", k.Status)
    case "subscribe":
	    stream, err := c.Subscribe(ctx, &pb.SubscribeRequest{Login: *login})
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
	case "":
		log.Printf("Choose action")
	}
}
