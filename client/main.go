package main

import (
	"flag"
	"log"
	"time"

	cl "chat_go/client/chatclient"
)

var (
	addr  = flag.String("addr", "localhost:50051", "the address to connect to")
	from  = flag.String("from", "A", "Login from")
	to    = flag.String("to", "B", "Login to")
	body  = flag.String("body", "Hello", "Message body")
	login = flag.String("login", "userA", "Login to subscribe")
	a     = flag.String("a", "", "Action")
)

// submitRequest calls the function and prints result
// depends on chosen action
func submitRequest(client *cl.Client) {
	flag.Parse()
	switch *a {
	case "users":
		users, err := client.GetUsers()
		if err != nil {
			log.Fatalf("didn't get users: %v", err)
		}
		log.Printf("Users: %s", users)
	case "message":
		created_at := int32(time.Now().Unix())
		m := cl.Message{LoginFrom: *from, LoginTo: *to, CreatedAt: created_at, Body: *body}
		resp, err := client.SendMessage(&m)
		if err != nil {
			log.Fatalf("didn't send message: %v", err)
		}
		log.Printf("Status: %s", resp)
	case "subscribe":
		messages, err := client.Subscribe(*login)
		if err != nil {
			log.Fatalf("didn't get messages: %v", err)
		}
		log.Printf("Messages: %s", messages)
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
