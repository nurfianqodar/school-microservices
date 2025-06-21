package svc

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/nurfianqodar/school-microservices/services/users/db"
	pbusers "github.com/nurfianqodar/school-microservices/services/users/pb/users/v1"
	"github.com/nurfianqodar/school-microservices/services/users/utils/token"
	v "github.com/nurfianqodar/school-microservices/services/users/utils/validation"
	"github.com/nurfianqodar/school-microservices/utils/errs"
	"github.com/nurfianqodar/school-microservices/utils/hasher"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errUserNotFound = status.Error(codes.NotFound, "user not found")
)

type service struct {
	pbusers.UnimplementedUserServiceServer
	mu sync.Mutex
	q  *db.Queries
}

func New(q *db.Queries) pbusers.UserServiceServer {
	return &service{
		q:  q,
		mu: sync.Mutex{},
	}
}

func (s *service) CreateOneUser(
	ctx context.Context,
	req *pbusers.CreateOneUserRequest,
) (*pbusers.CreateOneUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Validate request
	if err := v.Validate.Struct(req); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			return nil, errs.ConvertValidationError(validationErrs, v.Trans)
		} else {
			log.Printf("error: failed to validate data. %s", err.Error())
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	// Check email avaliable
	countEmail, err := s.q.CountEmailUser(ctx, req.Email)
	if err != nil {
		log.Printf("error: failed to count email. %s", err.Error())
		return nil, status.Error(codes.Internal, "internal server error")
	}
	if countEmail != 0 {
		return nil, status.Error(codes.AlreadyExists, "email already exist")
	}

	// Hash password
	passwordHash, err := hasher.GenerateFromPassword(req.Password, hasher.DefaultConfig)
	if err != nil {
		return nil, err
	}

	// Save to db
	// -- Generate uuid
	newUUID, err := uuid.NewV7()
	if err != nil {
		log.Printf("error: failed to generate new uuid v7. %s\n", err.Error())
		return nil, status.Error(codes.Internal, "internal server error")
	}

	// Convert role
	var role db.UserRole

	switch req.Role {
	case pbusers.UserRole_Unspecified:
		return nil, status.Error(codes.InvalidArgument, "invalid user role")
	case pbusers.UserRole_Teacher:
		role = db.UserRoleTeacher
	case pbusers.UserRole_Staff:
		role = db.UserRoleStaff
	case pbusers.UserRole_Student:
		role = db.UserRoleStudent
	case pbusers.UserRole_Parent:
		role = db.UserRoleParent
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid user role")
	}

	dbArgs := &db.CreateOneUserParams{
		ID:           newUUID,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		Role:         role,
	}

	result, err := s.q.CreateOneUser(ctx, dbArgs)
	if err != nil {
		log.Printf("error: failed to insert new user. %s\n", err.Error())
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pbusers.CreateOneUserResponse{
		Id: result.String(),
	}, nil
}

func (s *service) DeleteHardOneUser(
	ctx context.Context,
	req *pbusers.DeleteHardOneUserRequest,
) (*pbusers.DeleteHardOneUserResponse, error) {
	// count user by id
	reqUUID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid uuid")
	}

	count, err := s.q.CountIDUser(ctx, reqUUID)
	if err != nil {
		log.Printf("error: failed to count user by id: %e\n", err)
		return nil, errs.ErrInternalServer
	}
	if count == 0 {
		return nil, errUserNotFound
	}

	deletedID, err := s.q.DeleteHardOneUser(ctx, reqUUID)
	if err != nil {
		log.Printf("error: failed to hard delete user by id: %e\n", err)
		return nil, errs.ErrInternalServer
	}

	return &pbusers.DeleteHardOneUserResponse{
		Id: deletedID.String(),
	}, nil
}

func (s *service) DeleteSoftOneUser(
	ctx context.Context,
	req *pbusers.DeleteSoftOneUserRequest,
) (*pbusers.DeleteSoftOneUserResponse, error) {
	panic("not implemented")
}

func (s *service) GetManyUser(
	ctx context.Context,
	req *pbusers.GetManyUserRequest,
) (*pbusers.GetManyUserResponse, error) {
	panic("not implemented")
}

func (s *service) GetOneCredentialUserByEmail(
	ctx context.Context,
	req *pbusers.GetOneCredentialUserByEmailRequest,
) (*pbusers.GetOneCredentialUserByEmailResponse, error) {
	panic("not implemented")
}

func (s *service) GetOneUser(
	ctx context.Context,
	req *pbusers.GetOneUserRequest,
) (*pbusers.GetOneUserResponse, error) {
	panic("not implemented")
}

func (s *service) UpdateOneEmailUser(
	ctx context.Context,
	req *pbusers.UpdateOneEmailUserRequest,
) (*pbusers.UpdateOneEmailUserResponse, error) {
	panic("not implemented")
}

func (s *service) UpdateOnePasswordUser(
	ctx context.Context,
	req *pbusers.UpdateOnePasswordUserRequest,
) (*pbusers.UpdateOnePasswordUserResponse, error) {
	panic("not implemented")
}

func (s *service) UpdateOneRoleUser(
	ctx context.Context,
	req *pbusers.UpdateOneRoleUserRequest,
) (*pbusers.UpdateOneRoleUserResponse, error) {
	panic("not implemented")
}

func (s *service) LoginUser(
	ctx context.Context,
	req *pbusers.LoginUserRequest,
) (*pbusers.LoginUserResponse, error) {
	// Count email
	count, err := s.q.CountEmailUser(ctx, req.Email)
	if err != nil {
		log.Printf("error: failed to count user by email. %s", err.Error())
		return nil, status.Error(codes.Internal, "internal server error")
	}
	if count != 1 {
		return nil, errs.ErrInvalidCredential
	}

	// Get credential
	creds, err := s.q.GetOneCredentialUserByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("error: failed to get credential. %s", err.Error())
		return nil, status.Error(codes.Internal, "internal server error")
	}

	// Compare password and hash
	if err = hasher.CompareHashWithPassword(creds.PasswordHash, req.Password); err != nil {
		return nil, err
	}

	// Create access and refresh token
	// TODO fix audience
	accessToken, err := token.CreateToken(token.TokenTypeAccess, creds.ID.String(), time.Minute*30, []string{"null"})
	if err != nil {
		return nil, err
	}
	refreshToken, err := token.CreateToken(token.TokenTypeRefresh, creds.ID.String(), time.Hour*24*30, []string{"null"})
	if err != nil {
		return nil, err
	}

	return &pbusers.LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
