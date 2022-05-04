package main

import (
	"context"
	"fmt"

	"reflect"
	"testing"
	"time"

	pb "chat_go/chat_protos"

	st "chat_go/storage"

	"github.com/stretchr/testify/assert"
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
