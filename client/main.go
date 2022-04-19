package main

import (
	"flag"
	"log"

	pb "chat_go/chat_protos"
	cl "chat_go/client/chatclient"
	st "chat_go/storage"
)

var (
	addr  = flag.String("addr", "localhost:50051", "the address to connect to")
	from  = flag.String("from", "A", "Login from")
	to    = flag.String("to", "B", "Login to")
	body  = flag.String("body", "Hello", "Message body")
	login = flag.String("login", "userA", "Login to subscribe")
	a     = flag.String("action", "", "Action")
)

// submitRequest calls the function and prints result
// depends on chosen action
func submitRequest(client *cl.Client) {
	switch *a {
	case "users":
		users, err := client.GetUsers()
		if err != nil {
			log.Fatalf("didn't get users: %v", err)
		}
		log.Printf("Users: %s", users)
	case "message":
		m := st.Message{LoginFrom: *from, LoginTo: *to, Body: *body}

		resp, err := client.SendMessage(&m)

		if err != nil {
			log.Fatalf("didn't send message: %v", err)
		}
		log.Printf("Status: %s", resp)
	case "subscribe":
		channel := make(chan *pb.Message)

		go client.Subscribe(*login, channel)

		for message := range channel {
			log.Printf("Message: %v", message)
		}
	case "":
		log.Printf("Choose action")
	}
}

func main() {
	flag.Parse()
	client, err := cl.NewClient(*addr)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer client.Close()
	submitRequest(client)
}
