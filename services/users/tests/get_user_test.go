package user_test

import (
	"context"
	"testing"

	pbusers "github.com/nurfianqodar/school-microservices/services/users/pb/users/v1"
)

func TestGetUser(t *testing.T) {
	service := createService()
	ctx := context.TODO()
	res, err := service.CreateOneUser(ctx, &pbusers.CreateOneUserRequest{
		Email:    "example@email.sch.id",
		Password: "secretpassword",
		Role:     pbusers.UserRole_Student,
	})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer service.DeleteHardOneUser(ctx, &pbusers.DeleteHardOneUserRequest{Id: res.Id})

	res2, err := service.GetOneUser(ctx, &pbusers.GetOneUserRequest{
		Id: res.Id,
	})
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	t.Log(res2)
	if res2.Id != res.Id {
		t.Fail()
	}

	if res2.Email != "example@email.sch.id" {
		t.Fail()
	}
	if res2.Role != pbusers.UserRole_Student {
		t.Fail()
	}
}
