package main

import (
	"context"
	"errors"
	"fmt"

	"reflect"
	"testing"
	"time"

	pb "chat_go/chat_protos"

	st "chat_go/storage"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

// Tests 'GetUsers' method.
func TestGetUsers(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	storageMock := &st.StorageMock{}
	server := server{
		storage: storageMock,
	}

	storageUsers := []st.User{
		{Login: "111", FullName: "u_111"},
		{Login: "222", FullName: "u_222"},
	}

	returnUsers := []*pb.User{
		{Login: "111", FullName: "u_111"},
		{Login: "222", FullName: "u_222"},
	}

	storageMock.On("GetUsers").Return(storageUsers, nil)

	expectedReply := &pb.GetUsersReply{Users: returnUsers}
	mockGetUsersRequest := &pb.GetUsersRequest{}

	users, err := server.GetUsers(ctx, mockGetUsersRequest)

	assert.Nil(err)
	assert.Equal(expectedReply, users, "Lists of users should be equal.")
	storageMock.AssertExpectations(t)
}

// Tests 'SendMessage' method.
func TestSendMessage(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	storageMock := &st.StorageMock{}
	server := server{
		storage: storageMock,
	}

	message := st.Message{
		LoginFrom: "userA",
		LoginTo:   "userB",
		CreatedAt: int32(time.Now().Unix()),
		Body:      "Hello!",
	}

	storageMock.On("CreateMessage", message).Return(nil)

	mockSendMessageRequest := &pb.SendMessageRequest{
		Message: &pb.Message{
			LoginFrom: message.LoginFrom,
			LoginTo:   message.LoginTo,
			CreatedAt: message.CreatedAt,
			Body:      message.Body,
		},
	}

	expected := &pb.SendMessageReply{
		Status: message.LoginTo + " received message from " + message.LoginFrom,
	}
	result, err := server.SendMessage(ctx, mockSendMessageRequest)

	assert.Nil(err)
	assert.Equal(expected, result, "Message status should be equal.")
	storageMock.AssertExpectations(t)
}

// Tests 'fillUsers' method.
func TestFillUsers(t *testing.T) {
	storage := &st.StorageMock{}

	storage.On("GetUsers").Return(nil, nil)
	for i := 0; i < 4; i++ {
		login := fmt.Sprintf("user%d", i)
		full_name := fmt.Sprintf("user%d_user%d", i, i)
		user := st.User{Login: login, FullName: full_name}
		storage.On("CreateUser", user).Return(nil)
	}

	fillUsers(storage)

	storage.AssertExpectations(t)
}

// Tests 'createStorage' method.
func TestCreateStorage(t *testing.T) {
	assert := assert.New(t)

	wantSrorage, _ := st.NewEtcdStorage([]string{"localhost:1234"}, 5*time.Second)
	storage := createStorage("localhost", 1234)

	assert.Equal(reflect.TypeOf(wantSrorage), reflect.TypeOf(storage), "Type of storage must be equal.")
}

type mockChat_SubscribeServer struct {
	grpc.ServerStream
	Messages []*pb.Message
}

func (_m *mockChat_SubscribeServer) Send(message *pb.Message) error {
	_m.Messages = append(_m.Messages, message)
	return nil
}

// Tests 'Subscribe' method.
func TestSubscribe(t *testing.T) {
	assert := assert.New(t)
	storageMock := &st.StorageMock{}
	server := server{
		storage: storageMock,
	}

	login := "userA"
	storageMessages := []st.Message{
		{LoginFrom: "user1", LoginTo: "userB", CreatedAt: 1234, Body: "Hello!"},
		{LoginFrom: "user2", LoginTo: "userB", CreatedAt: 4567, Body: "Hello!"},
		{LoginFrom: "user3", LoginTo: "userB", CreatedAt: 1234, Body: "Hello!"},
		{LoginFrom: "user4", LoginTo: "userB", CreatedAt: 4567, Body: "Hello!"},
	}

	mockSubscribeRequest := &pb.SubscribeRequest{Login: login}
	expectedErrorMsg := "Messages ended."

	storageMock.On("GetMessages", login).Return(storageMessages, nil).Once()
	storageMock.On("GetMessages", login).Return(nil, errors.New(expectedErrorMsg)).Once()

	for _, message := range storageMessages {
		storageMock.On("DeleteMessage", message).Return(nil)
	}

	mockStream := &mockChat_SubscribeServer{}

	err := server.Subscribe(mockSubscribeRequest, mockStream)

	expectedMessages := []*pb.Message{
		{LoginFrom: "user1", LoginTo: "userB", CreatedAt: 1234, Body: "Hello!"},
		{LoginFrom: "user2", LoginTo: "userB", CreatedAt: 4567, Body: "Hello!"},
		{LoginFrom: "user3", LoginTo: "userB", CreatedAt: 1234, Body: "Hello!"},
		{LoginFrom: "user4", LoginTo: "userB", CreatedAt: 4567, Body: "Hello!"},
	}
	messages := mockStream.Messages

	assert.EqualErrorf(err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
	assert.Equal(4, len(mockStream.Messages), "Sites expected to contain 4 messages.")
	assert.Equal(expectedMessages, messages, "Lists of messages should be equal.")
	storageMock.AssertNumberOfCalls(t, "GetMessages", 2)
	storageMock.AssertNumberOfCalls(t, "DeleteMessage", 4)
}
