package chatclient

import (
	"context"
	"testing"

	pb "chat_go/chat_protos"

	st "chat_go/storage"

	"chat_go/client/chatclient/mocks"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

)

func TestUnitNewClient(t *testing.T) {
	assert := assert.New(t)

	endpoint := "chat.endpoint"
	_, err := NewClient(endpoint)

	assert.Nil(err)
}

// Tests 'GetUsers' method. 
func TestGetUsers(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	mockClientConn := &grpc.ClientConn{}
	mockChatClient := mocks.ChatClientMock{}

	mockGetUsersReply := &pb.GetUsersReply{
		Users: []*pb.User{
			{Login: "user1", FullName: "u_user1"},
			{Login: "user2", FullName: "u_user2"},
		},
	}
	mockGetUsersRequest := &pb.GetUsersRequest{}
	mockChatClient.On("GetUsers", ctx, mockGetUsersRequest).Return(mockGetUsersReply, nil)

	expectedUsers :=  []*pb.User{
			{Login: "user1", FullName: "u_user1"},
			{Login: "user2", FullName: "u_user2"},
		}

	client := &Client{
		Endpoint: "123",
		conn:     mockClientConn,
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
	ctx := context.Background()

	mockClientConn := &grpc.ClientConn{}
	mockChatClient := mocks.ChatClientMock{}

	m := &st.Message{
		LoginFrom: "userA",
		LoginTo: "userB",
		Body: "Hello!",
	}

	mockSendMessageRequest := &pb.SendMessageRequest{
		Message: &pb.Message{
			LoginFrom: m.LoginFrom,
			LoginTo: m.LoginTo,
			Body: m.Body,
		},
	}

	mockSendMessageReply := &pb.SendMessageReply{
		Status: m.LoginTo + " received message from " + m.LoginFrom,
	}

	mockChatClient.On("SendMessage", ctx, mockSendMessageRequest).Return(mockSendMessageReply, nil)

	expectedStatus :=  "userB received message from userA"

	client := &Client{
		Endpoint: "123",
		conn:     mockClientConn,
		client:   mockChatClient,
	}
	status, err := client.SendMessage(m)

	assert.Nil(err)
	assert.Equal(expectedStatus, status, "Message status should be equal.")
	mockChatClient.AssertExpectations(t)
}
