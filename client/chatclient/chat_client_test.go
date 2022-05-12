package chatclient

import (
	"context"
	"io"
	// "errors"
	"log"
	"reflect"
	"testing"

	pb "chat_go/chat_protos"

	st "chat_go/storage"

	"chat_go/client/chatclient/mocks"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

// Tests 'NewClient' method.
func TestUnitNewClient(t *testing.T) {
	assert := assert.New(t)

	endpoint := "chat.endpoint"

	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	expectedClient := &Client{Endpoint: endpoint, conn: conn, client: pb.NewChatClient(conn)}
	client, err := NewClient(endpoint)

	assert.Nil(err)
	assert.Equal(reflect.TypeOf(expectedClient.client), reflect.TypeOf(client.client), "Type of client must be equal.")
	assert.Equal(expectedClient.Endpoint, client.Endpoint, "Endpoint must be equal.")
}

// Tests 'GetUsers' method.
func TestGetUsers(t *testing.T) {
	assert := assert.New(t)

	mockChatClient := mocks.ChatClientMock{}

	mockGetUsersReply := &pb.GetUsersReply{
		Users: []*pb.User{
			{Login: "user1", FullName: "u_user1"},
			{Login: "user2", FullName: "u_user2"},
		},
	}
	mockGetUsersRequest := &pb.GetUsersRequest{}
	mockChatClient.On("GetUsers", mockGetUsersRequest).Return(mockGetUsersReply, nil)

	expectedUsers := []*pb.User{
		{Login: "user1", FullName: "u_user1"},
		{Login: "user2", FullName: "u_user2"},
	}

	client := &Client{
		Endpoint: "123",
		client:   mockChatClient,
	}
	users, err := client.GetUsers()

	assert.Nil(err)
	assert.Equal(expectedUsers, users, "Lists of users should be equal.")
	mockChatClient.AssertExpectations(t)
}

// Tests 'SendMessage' method.
func TestSendMessage(t *testing.T) {
	assert := assert.New(t)

	mockChatClient := mocks.ChatClientMock{}

	m := &st.Message{
		LoginFrom: "userA",
		LoginTo:   "userB",
		Body:      "Hello!",
	}

	mockSendMessageRequest := &pb.SendMessageRequest{
		Message: &pb.Message{
			LoginFrom: m.LoginFrom,
			LoginTo:   m.LoginTo,
			Body:      m.Body,
		},
	}

	mockSendMessageReply := &pb.SendMessageReply{
		Status: m.LoginTo + " received message from " + m.LoginFrom,
	}

	mockChatClient.On("SendMessage", mockSendMessageRequest).Return(mockSendMessageReply, nil)

	expectedStatus := "userB received message from userA"

	client := &Client{
		Endpoint: "123",
		client:   mockChatClient,
	}
	status, err := client.SendMessage(m)

	assert.Nil(err)
	assert.Equal(expectedStatus, status, "Message status should be equal.")
	mockChatClient.AssertExpectations(t)
}

// Tests 'Subscribe' method.
func TestSubscribe(t *testing.T) {
	assert := assert.New(t)

	mockChatClient := mocks.ChatClientMock{}

	login := "userB"

	mockStream := &mocks.Chat_SubscribeClient{}
	messages := []*pb.Message{
		{LoginFrom: "user1", LoginTo: "userB", CreatedAt: 1234, Body: "Hello!"},
		{LoginFrom: "user2", LoginTo: "userB", CreatedAt: 4567, Body: "Hello!"},
		{LoginFrom: "user3", LoginTo: "userB", CreatedAt: 1234, Body: "Hello!"},
		{LoginFrom: "user4", LoginTo: "userB", CreatedAt: 4567, Body: "Hello!"},
	}
	for _, message := range messages{
		mockStream.On("Recv").Return(message, nil).Once()
	}

	mockStream.On("Recv").Return(nil, io.EOF).Once()

	mockSubscribeRequest := &pb.SubscribeRequest{Login: login}

	mockChatClient.On("Subscribe", context.Background(), mockSubscribeRequest).Return(mockStream, nil)

	client := &Client{
		Endpoint: "123",
		client:   mockChatClient,
	}

	channel := make(chan *pb.Message)

	go client.Subscribe(login, channel)

	var respMessages []*pb.Message

	for message := range channel {
		respMessages = append(respMessages, message)
	}	

	assert.Equal(messages, respMessages, "Lists of messages should be equal.")
	mockStream.AssertNumberOfCalls(t, "Recv", 5)
	mockChatClient.AssertExpectations(t)
}
