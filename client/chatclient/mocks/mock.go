package mocks

import (
	"context"
	"io"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	pb "chat_go/chat_protos"
)

// ChatClient is a mock type for the ChatClient type. 
type ChatClientMock struct {
	mock.Mock
	pb.ChatClient
}

// GetUsers provides a mock function with given fields: ctx, in, opts. 
func (c ChatClientMock) GetUsers(ctx context.Context, in *pb.GetUsersRequest, opts ...grpc.CallOption) (*pb.GetUsersReply, error) {
	args := c.Called(in)
	return args.Get(0).(*pb.GetUsersReply), args.Error(1)
}

// SendMessage provides a mock function with given fields: ctx, in, opts. 
func (c ChatClientMock) SendMessage(ctx context.Context, in *pb.SendMessageRequest, opts ...grpc.CallOption) (*pb.SendMessageReply, error) {
	args := c.Called(in)
	return args.Get(0).(*pb.SendMessageReply), args.Error(1)
}

// Subscribe provides a mock function with given fields: ctx, in, opts. 
func (c ChatClientMock) Subscribe(ctx context.Context, in *pb.SubscribeRequest, opts ...grpc.CallOption) (pb.Chat_SubscribeClient, error) {
	args := c.Called(ctx, in)
	return args.Get(0).(pb.Chat_SubscribeClient), args.Error(1)
}

type ConnectionCloserMock struct {
	io.Closer
	mock.Mock
}

func (c ConnectionCloserMock) Close() error {
	args := c.Called()
	return args.Error(0)
}
