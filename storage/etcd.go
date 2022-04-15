package storage

import (
	"context"
	"encoding/json"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const(
	USER_PREFIX = "user."
	MESSAGE_PREFIX = "message."
)

// Base for creating storages. 
type Storage struct {
	Endpoints []string

	storage clientv3.Client
}

// EtcdStorage creates new Storage using etcd client/v3. 
func EtcdStorage(endpoints []string, dialTimeout time.Duration) (*Storage, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	return &Storage{Endpoints: endpoints, storage: *cli}, nil
}

// Close closes connection with etcd. 
func (s *Storage) Close() error {
	return s.storage.Close()
}

// CreateUser saves user object into etcd using user key. 
func (s *Storage) CreateUser(u User) error {
	ctx := context.Background()
	k :=  USER_PREFIX + u.Login
	v, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}
	s.storage.Put(ctx, k, string(v))
	return nil
}

// GetUsers returns list of users. 
func (s *Storage) GetUsers() (users []User, err error) {
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
func (s *Storage) CreateMessage(m Message) error {
	ctx := context.Background()
	k := MESSAGE_PREFIX + m.LoginTo + m.LoginFrom + string(m.CreatedAt)
	v, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	s.storage.Put(ctx, k, string(v))

	return nil
}

// GetMessages returns list of messages for specific user. 
func (s *Storage) GetMessages(login string) (messages []Message, err error) {
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
func (s *Storage) DeleteMessage(m Message) {
	ctx := context.Background()
	k := MESSAGE_PREFIX + m.LoginTo + m.LoginFrom + string(m.CreatedAt)
	s.storage.Delete(ctx, k)
}
