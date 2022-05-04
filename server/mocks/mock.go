package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"chat_go/storage"

	pb "chat_go/chat_protos"
)

type ChatServerMock struct {
	mock.Mock
	storage.StorageMock
}

func (s ChatServerMock) GetUsers(ctx context.Context, in *pb.GetUsersRequest, opts ...grpc.CallOption) (*pb.GetUsersReply, error) {
	args := s.Called(ctx, in)
	return args.Get(0).(*pb.GetUsersReply), args.Error(1)
}

func (s ChatServerMock) SendMessage(ctx context.Context, in *pb.SendMessageRequest, opts ...grpc.CallOption) (*pb.SendMessageReply, error) {
	args := s.Called(ctx, in)
	return args.Get(0).(*pb.SendMessageReply), args.Error(1)
}

func (c ChatServerMock) Subscribe(ctx context.Context, in *pb.SubscribeRequest, opts ...grpc.CallOption) (pb.Chat_SubscribeClient, error) {
	args := c.Called(ctx, in)
	return args.Get(0).(pb.Chat_SubscribeClient), args.Error(1)
}

