package svc

import (
	"context"
	"log"
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
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	errUserNotFound = status.Error(codes.NotFound, "user not found")
)

// Service server
type service struct {
	pbusers.UnimplementedUserServiceServer
	q *db.Queries
}

// Constructor
func New(q *db.Queries) pbusers.UserServiceServer {
	return &service{
		q: q,
	}
}

// Service Implementation
func (s *service) CreateOneUser(
	ctx context.Context,
	req *pbusers.CreateOneUserRequest,
) (*pbusers.CreateOneUserResponse, error) {
	// Validate request
	if err := v.Validate.Struct(req); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			return nil, errs.ConvertValidationError(validationErrs, v.Trans)
		} else {
			log.Printf("error: failed to validate data. %s", err.Error())
			return nil, errs.ErrInternalServer
		}
	}

	// Check email avaliable
	countEmail, err := s.q.CountEmailUser(ctx, req.Email)
	if err != nil {
		log.Printf("error: failed to count email. %s", err.Error())
		return nil, errs.ErrInternalServer
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
		return nil, errs.ErrInternalServer
	}

	// -- Convert role
	role, err := convertRoleProtoToRoleDb(req.Role)
	if err != nil {
		return nil, err
	}

	// -- Create db args instance
	dbArgs := &db.CreateOneUserParams{
		ID:           newUUID,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		Role:         role,
	}

	// -- execute query
	result, err := s.q.CreateOneUser(ctx, dbArgs)
	if err != nil {
		log.Printf("error: failed to insert new user. %s\n", err.Error())
		return nil, errs.ErrInternalServer
	}

	return &pbusers.CreateOneUserResponse{
		Id: result.String(),
	}, nil
}

func (s *service) DeleteHardOneUser(
	ctx context.Context,
	req *pbusers.DeleteHardOneUserRequest,
) (*pbusers.DeleteHardOneUserResponse, error) {
	// Parse uuid
	reqUUID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid uuid")
	}

	// count user by id
	count, err := s.q.CountIDUser(ctx, reqUUID)
	if err != nil {
		log.Printf("error: failed to count user by id: %e\n", err)
		return nil, errs.ErrInternalServer
	}
	if count == 0 {
		return nil, errUserNotFound
	}

	// Delete user
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
	reqUUID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid uuid")
	}

	// count user by id
	count, err := s.q.CountIDUser(ctx, reqUUID)
	if err != nil {
		log.Printf("error: failed to count user by id: %e\n", err)
		return nil, errs.ErrInternalServer
	}
	if count == 0 {
		return nil, errUserNotFound
	}

	deletedID, err := s.q.DeleteSoftOneUser(ctx, reqUUID)
	// Delete user
	if err != nil {
		log.Printf("error: failed to soft delete user by id: %e\n", err)
		return nil, errs.ErrInternalServer
	}

	return &pbusers.DeleteSoftOneUserResponse{
		Id: deletedID.String(),
	}, nil
}

func (s *service) GetManyUser(
	ctx context.Context,
	req *pbusers.GetManyUserRequest,
) (*pbusers.GetManyUserResponse, error) {
	res, err := s.q.GetManyUser(ctx, &db.GetManyUserParams{
		Limit:  int32(req.Limit),
		Offset: int32(req.Offset),
	})
	if err != nil {
		return nil, errs.ErrInternalServer
	}
	users := make([]*pbusers.UserSummary, 0, len(res))
	log.Printf("found %d users\n", len(users))
	for _, user := range res {
		pbRole := convertRoleDbToRoleProto(user.Role)
		users = append(users, &pbusers.UserSummary{
			Id:    user.ID.String(),
			Email: user.Email,
			Role:  pbRole,
		})
	}

	return &pbusers.GetManyUserResponse{
		Users: users,
	}, nil
}

func (s *service) GetOneUser(
	ctx context.Context,
	req *pbusers.GetOneUserRequest,
) (*pbusers.GetOneUserResponse, error) {
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
	user, err := s.q.GetOneUser(ctx, reqUUID)
	if err != nil {
		log.Printf("error: failed to get one user by id: %e\n", err)
		return nil, errs.ErrInternalServer
	}
	return &pbusers.GetOneUserResponse{
		Id:        user.ID.String(),
		Email:     user.Email,
		Role:      convertRoleDbToRoleProto(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt.Time),
	}, nil
}

func (s *service) UpdateOneEmailUser(
	ctx context.Context,
	req *pbusers.UpdateOneEmailUserRequest,
) (*pbusers.UpdateOneEmailUserResponse, error) {
	// Parse uuid and validation
	reqUUID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid uuid")
	}
	if err = v.Validate.Struct(req); err != nil {
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			return nil, errs.ConvertValidationError(validationErr, v.Trans)
		} else {
			log.Printf("error: unable to use validator. %s", err.Error())
			return nil, errs.ErrInternalServer
		}
	}

	// Count user
	count, err := s.q.CountIDUser(ctx, reqUUID)
	if err != nil {
		log.Printf("error: failed to count user by id: %e\n", err)
		return nil, errs.ErrInternalServer
	}
	if count == 0 {
		return nil, errUserNotFound
	}

	// Update
	userID, err := s.q.UpdateOneEmailUser(ctx, &db.UpdateOneEmailUserParams{
		ID:    reqUUID,
		Email: req.Email,
	})
	if err != nil {
		log.Printf("error: failed to update user email. %s", err.Error())
		return nil, errs.ErrInternalServer
	}
	return &pbusers.UpdateOneEmailUserResponse{
		Id: userID.String(),
	}, nil
}

func (s *service) UpdateOnePasswordUser(
	ctx context.Context,
	req *pbusers.UpdateOnePasswordUserRequest,
) (*pbusers.UpdateOnePasswordUserResponse, error) {
	// Parse uuid and validation
	reqUUID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid uuid")
	}
	if err = v.Validate.Struct(req); err != nil {
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			return nil, errs.ConvertValidationError(validationErr, v.Trans)
		} else {
			log.Printf("error: unable to use validator. %s", err.Error())
			return nil, errs.ErrInternalServer
		}
	}

	// Count user
	count, err := s.q.CountIDUser(ctx, reqUUID)
	if err != nil {
		log.Printf("error: failed to count user by id: %e\n", err)
		return nil, errs.ErrInternalServer
	}
	if count == 0 {
		return nil, errUserNotFound
	}

	// Hash password
	hash, err := hasher.GenerateFromPassword(req.Password, hasher.DefaultConfig)
	if err != nil {
		log.Printf("error: failed to hash password. %e\n", err)
	}

	// Update
	userID, err := s.q.UpdateOnePasswordUser(ctx, &db.UpdateOnePasswordUserParams{
		ID:           reqUUID,
		PasswordHash: hash,
	})
	if err != nil {
		log.Printf("error: failed to update user email. %s", err.Error())
		return nil, errs.ErrInternalServer
	}
	return &pbusers.UpdateOnePasswordUserResponse{
		Id: userID.String(),
	}, nil
}

func (s *service) UpdateOneRoleUser(
	ctx context.Context,
	req *pbusers.UpdateOneRoleUserRequest,
) (*pbusers.UpdateOneRoleUserResponse, error) {
	// Parse uuid and validation
	reqUUID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid uuid")
	}
	if err = v.Validate.Struct(req); err != nil {
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			return nil, errs.ConvertValidationError(validationErr, v.Trans)
		} else {
			log.Printf("error: unable to use validator. %s", err.Error())
			return nil, errs.ErrInternalServer
		}
	}

	// Count user
	count, err := s.q.CountIDUser(ctx, reqUUID)
	if err != nil {
		log.Printf("error: failed to count user by id: %e\n", err)
		return nil, errs.ErrInternalServer
	}
	if count == 0 {
		return nil, errUserNotFound
	}

	role, err := convertRoleProtoToRoleDb(req.Role)
	if err != nil {
		return nil, err
	}

	// Update
	userID, err := s.q.UpdateOneRoleUser(ctx, &db.UpdateOneRoleUserParams{
		ID:   reqUUID,
		Role: role,
	})
	if err != nil {
		log.Printf("error: failed to update user email. %s", err.Error())
		return nil, errs.ErrInternalServer
	}
	return &pbusers.UpdateOneRoleUserResponse{
		Id: userID.String(),
	}, nil
}

func (s *service) LoginUser(
	ctx context.Context,
	req *pbusers.LoginUserRequest,
) (*pbusers.LoginUserResponse, error) {
	// Count email
	count, err := s.q.CountEmailUser(ctx, req.Email)
	if err != nil {
		log.Printf("error: failed to count user by email. %s", err.Error())
		return nil, errs.ErrInternalServer
	}
	if count != 1 {
		return nil, errs.ErrInvalidCredential
	}

	// Get credential
	creds, err := s.q.GetOneCredentialUserByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("error: failed to get credential. %s", err.Error())
		return nil, errs.ErrInternalServer
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

func (s *service) VerifyTokenUser(
	ctx context.Context,
	req *pbusers.VerifyTokenUserRequest,
) (*pbusers.VerifyTokenUserResponse, error) {
	claims, err := token.VerifyToken(req.AccessToken)
	if err != nil {
		return nil, err
	}
	return &pbusers.VerifyTokenUserResponse{
		Exp: timestamppb.New(claims.Exp.Local()),
		Iat: timestamppb.New(claims.Iat.Local()),
		Nbf: timestamppb.New(claims.Nbf.Local()),
		Iss: claims.Iss,
		Sub: claims.Sub,
		Aud: claims.Aud,
	}, nil
}
