package storage

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"io"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	USER_PREFIX    = "user."
	MESSAGE_PREFIX = "message."
)

// Base interface for creating storages. All storages need to provide methods
// for creating users, getting all users, creating messages,
// getting all messages per user, removing specific message for specific user.
type Storage interface {
	// CreateUser saves user in storage.
	CreateUser(u User) error

	// GetUsers returns users list from storage.
	GetUsers() (users []User, err error)

	// CreateMessage saves messages in storage.
	CreateMessage(m Message) error

	// GetMessages retrieves user's login and returns list of messages from storage.
	GetMessages(login string) (messages []Message, err error)

	// DeleteMessage deletes user-read messages.
	DeleteMessage(m Message) error
}

// Base for creating etcd storages.
type EtcdStorage struct {
	Endpoints []string

	storage EtcdClient
}

var _ Storage = &EtcdStorage{}

// Base interface for creating etcd client. 
type EtcdClient interface{
	clientv3.KV
	io.Closer
}

// NewEtcdStorage creates new Storage using etcd client/v3.
func NewEtcdStorage(endpoints []string, dialTimeout time.Duration) (*EtcdStorage, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	return &EtcdStorage{Endpoints: endpoints, storage: cli}, nil
}

// Close closes connection with etcd.
func (s *EtcdStorage) Close() error {
	return s.storage.Close()
}

// CreateUser saves user object into etcd using user key.
func (s *EtcdStorage) CreateUser(u User) error {
	ctx := context.Background()
	k := USER_PREFIX + u.Login
	v, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}
	_, err = s.storage.Put(ctx, k, string(v))
	return err
}

// GetUsers returns list of users.
func (s *EtcdStorage) GetUsers() (users []User, err error) {
	ctx := context.Background()
	opts := []clientv3.OpOption{
		clientv3.WithPrefix(),
	}
	res, err := s.storage.Get(ctx, USER_PREFIX, opts...)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range res.Kvs {
		var u User
		err := json.Unmarshal(item.Value, &u)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, u)
	}
	return users, nil
}

// CreateMessage saves message object into etcd using message key.
// Message key includes user login and timestamp created_at to be unique.
func (s *EtcdStorage) CreateMessage(m Message) error {
	ctx := context.Background()
	k := MESSAGE_PREFIX + m.LoginTo + m.LoginFrom + string(m.CreatedAt)
	v, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	_, err = s.storage.Put(ctx, k, string(v))
	return err
}

// GetMessages returns list of messages for specific user.
func (s *EtcdStorage) GetMessages(login string) (messages []Message, err error) {
	ctx := context.Background()
	opts := []clientv3.OpOption{
		clientv3.WithPrefix(),
	}
	k := MESSAGE_PREFIX + login
	res, err := s.storage.Get(ctx, k, opts...)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range res.Kvs {
		var m Message
		err := json.Unmarshal(item.Value, &m)
		if err != nil {
			log.Fatal(err)
		}
		messages = append(messages, m)
	}
	return messages, nil
}

// DeleteMessage deletes message from storage.
func (s *EtcdStorage) DeleteMessage(m Message) error {
	ctx := context.Background()
	k := MESSAGE_PREFIX + m.LoginTo + m.LoginFrom + string(m.CreatedAt)
	_, err := s.storage.Delete(ctx, k)
	return err
}
