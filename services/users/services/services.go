package svc

import (
	"context"

	"github.com/nurfianqodar/school-microservices/services/users/db"
	pbusers "github.com/nurfianqodar/school-microservices/services/users/pb/users/v1"
)

type service struct {
	q *db.Queries
	pbusers.UnimplementedUserServiceServer
}

func New() pbusers.UserServiceServer {
	return &service{}
}

func (s *service) CreateOneUser(ctx context.Context, req *pbusers.CreateOneUserRequest) (*pbusers.CreateOneUserResponse, error) {
	panic("not implemented")
}
