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

func main() {
	flag.Parse()
	client, err := cl.NewClient(*addr)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer client.Close()

	switch *a {
	case "users":
		client.GetUsers()
	case "message":
		created_at := int32(time.Now().Unix())
		m := cl.Message{LoginFrom: *from, LoginTo: *to, CreatedAt: created_at, Body: *body}
		client.SendMessage(&m)
	case "subscribe":
		client.Subscribe(*login)
	case "":
		log.Printf("Choose action")
	}
}
