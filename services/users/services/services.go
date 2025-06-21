package svc

import (
	"context"

	pbusers "github.com/nurfianqodar/school-microservices/services/users/pb/users/v1"
)

type service struct {
	pbusers.UnimplementedUserServiceServer
}

func New() pbusers.UserServiceServer {
	return &service{}
}

func (s *service) CreateUser(ctx context.Context, req *pbusers.CreateUserRequest) (*pbusers.CreateUserResponse, error) {
	panic("not implemented")
}
