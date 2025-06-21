package svc

import (
	"github.com/nurfianqodar/school-microservices/services/users/db"
	pbusers "github.com/nurfianqodar/school-microservices/services/users/pb/users/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func convertRole(r pbusers.UserRole) (db.UserRole, error) {
	switch r {
	case pbusers.UserRole_Unspecified:
		return "", status.Error(codes.InvalidArgument, "invalid user role")
	case pbusers.UserRole_Teacher:
		return db.UserRoleTeacher, nil
	case pbusers.UserRole_Staff:
		return db.UserRoleStaff, nil
	case pbusers.UserRole_Student:
		return db.UserRoleStudent, nil
	case pbusers.UserRole_Parent:
		return db.UserRoleParent, nil
	default:
		return "", status.Error(codes.InvalidArgument, "invalid user role")
	}
}
