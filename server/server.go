package server

import (
	"context"
	"log"
	"sync"

	"github.com/gofrs/uuid"

	pbExample "github.com/johanbrandhorst/grpc-gateway-boilerplate/proto"
)

// Backend implements the protobuf interface
type Backend struct {
	mu    *sync.RWMutex
	users []*pbExample.User
}

// New initializes a new Backend struct.
func New() *Backend {
	return &Backend{
		mu: &sync.RWMutex{},
	}
}

// AddUser adds a user to the in-memory store.
func (b *Backend) AddUser(ctx context.Context, _ *pbExample.AddUserRequest) (*pbExample.User, error) {
	log.Println("add user")
	b.mu.Lock()
	defer b.mu.Unlock()

	user := &pbExample.User{
		Id: uuid.Must(uuid.NewV4()).String(),
	}
	b.users = append(b.users, user)

	return user, nil
}

// ListUsers lists all users in the store.
func (b *Backend) ListUsers(_ *pbExample.ListUsersRequest, srv pbExample.UserService_ListUsersServer) error {
	log.Println("list users")
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, user := range b.users {
		err := srv.Send(user)
		if err != nil {
			return err
		}
	}

	return nil
}

// Ping ping returns pong
func (b *Backend) Ping(ctx context.Context, _ *pbExample.PingRequest) (*pbExample.PongResponse, error) {
	log.Println("ping")

	b.mu.Lock()
	defer b.mu.Unlock()

	return &pbExample.PongResponse{
		Message: "pong",
	}, nil
}
