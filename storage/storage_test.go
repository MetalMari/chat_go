package storage

import (
	"context"
	"encoding/json"
	"log"
	"testing"

	"chat_go/storage/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	clientv3 "go.etcd.io/etcd/client/v3"
	mvccpb "go.etcd.io/etcd/api/v3/mvccpb"
)

// Tests 'GetUsers' method. 
func TestGetUsers(t *testing.T) {
	assert := assert.New(t)

	MockEtcdClient := &mocks.EtcdClient{}

	storage:= EtcdStorage{
		Endpoints: []string{"localhost:1234"},
		storage: MockEtcdClient,
	}

	key1, err := json.Marshal("user1")
	if err != nil {
		log.Fatal(err)
	}
	value1, err := json.Marshal(User{Login: "111", FullName: "u_111"})
	if err != nil {
		log.Fatal(err)
	}

	key2, err := json.Marshal("user2")
	if err != nil {
		log.Fatal(err)
	}
	value2, err := json.Marshal(User{Login: "222", FullName: "u_222"})
	if err != nil {
		log.Fatal(err)
	}

	resp := &clientv3.GetResponse{
		Kvs: []*mvccpb.KeyValue{
			{Key: key1,
			Value: value1},
			{Key: key2,
			Value: value2},			
		},
	}

	MockEtcdClient.On("Get", context.Background(), "user.", mock.Anything).Return(resp, nil)

	users, err := storage.GetUsers()
	expectedUsers := []User{
	    {Login: "111", FullName: "u_111"},
	    {Login: "222", FullName: "u_222"},
	}

	assert.Nil(err)
	assert.Equal(expectedUsers, users, "Lists of users should be equal.")
	MockEtcdClient.AssertExpectations(t)
}

// Tests 'GetMessages' method. 
func TestGetMessages(t *testing.T) {
	assert := assert.New(t)

	MockEtcdClient := &mocks.EtcdClient{}

	storage:= EtcdStorage{
		Endpoints: []string{"localhost:1234"},
		storage: MockEtcdClient,
	}

	login := "user1"

	key1, err := json.Marshal("message1")
	if err != nil {
		log.Fatal(err)
	}
	value1, err := json.Marshal(Message{
		LoginFrom: "user2",
		LoginTo:   "user1",
		CreatedAt: 1234,
		Body:      "Hello!",
	})
	if err != nil {
		log.Fatal(err)
	}

	key2, err := json.Marshal("message2")
	if err != nil {
		log.Fatal(err)
	}
	value2, err := json.Marshal(Message{
		LoginFrom: "user3",
		LoginTo:   "user1",
		CreatedAt: 12345678,
		Body:      "Hello!",
	})
	if err != nil {
		log.Fatal(err)
	}

	resp := &clientv3.GetResponse{
		Kvs: []*mvccpb.KeyValue{
			{Key: key1,
			Value: value1},
			{Key: key2,
			Value: value2},			
		},
	}
	key := "message." + login
	MockEtcdClient.On("Get", context.Background(), key, mock.Anything).Return(resp, nil)

	expectedMessages := []Message{
		{LoginFrom: "user2", LoginTo: "user1", CreatedAt: 1234, Body: "Hello!"},
	    {LoginFrom: "user3", LoginTo: "user1", CreatedAt: 12345678, Body: "Hello!"},
	}
	messages, err := storage.GetMessages(login)

	assert.Nil(err)
	assert.Equal(expectedMessages, messages, "Lists of messages should be equal.")
	MockEtcdClient.AssertExpectations(t)
}

// Tests 'CreateUser' method. 
func TestCreateUser(t *testing.T) {
	assert := assert.New(t)

	MockEtcdClient := &mocks.EtcdClient{}

	storage:= EtcdStorage{
		Endpoints: []string{"localhost:1234"},
		storage: MockEtcdClient,
	}

	user := User{Login: "111", FullName: "u_111"}
	key := "user.111"
	value, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}
	resp := &clientv3.PutResponse{}

	MockEtcdClient.On("Put", context.Background(), key, string(value)).Return(resp, nil)

	err = storage.CreateUser(user)

	assert.Nil(err)
	MockEtcdClient.AssertExpectations(t)
}

// Tests 'CreateMessage' method. 
func TestCreateMessage(t *testing.T) {
	assert := assert.New(t)

	MockEtcdClient := &mocks.EtcdClient{}

	storage:= EtcdStorage{
		Endpoints: []string{"localhost:1234"},
		storage: MockEtcdClient,
	}

	message := Message{LoginFrom: "user2", LoginTo: "user1", CreatedAt: 1234, Body: "Hello!"}
	key := "message.user1user2" + string(message.CreatedAt)
	value, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}
	resp := &clientv3.PutResponse{}

	MockEtcdClient.On("Put", context.Background(), key, string(value)).Return(resp, nil)

	err = storage.CreateMessage(message)

	assert.Nil(err)
	MockEtcdClient.AssertExpectations(t)
}

// Tests 'DeleteMessage' method. 
func TestDeleteMessage(t *testing.T) {
	assert := assert.New(t)

	MockEtcdClient := &mocks.EtcdClient{}

	storage:= EtcdStorage{
		Endpoints: []string{"localhost:1234"},
		storage: MockEtcdClient,
	}

	message := Message{LoginFrom: "user2", LoginTo: "user1", CreatedAt: 1234, Body: "Hello!"}
	key := "message.user1user2" + string(message.CreatedAt)
	resp := &clientv3.DeleteResponse{}
	
	MockEtcdClient.On("Delete", context.Background(), key).Return(resp, nil)

	err := storage.DeleteMessage(message)

	assert.Nil(err)
	MockEtcdClient.AssertExpectations(t)
}
