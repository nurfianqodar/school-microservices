package tests

import (
	"context"
	"fmt"
	"os"
	"testing"

	pbusers "github.com/nurfianqodar/school-microservices/services/users/pb/users/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestCreateUser(t *testing.T) {
	host, ok := os.LookupEnv("HOST")
	if !ok {
		t.Log("HOST environment variable was not set")
		t.FailNow()
	}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		t.Log("PORT environment variable was not set")
		t.FailNow()
	}
	addr := fmt.Sprintf("%s:%s", host, port)

	client, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	service := pbusers.NewUserServiceClient(client)
	res, err := service.CreateOneUser(context.TODO(), &pbusers.CreateOneUserRequest{
		Email:    "dummy@gmail.com",
		Password: "secretpassword",
		Role:     pbusers.UserRole_Student,
	})
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(res)

	_, _ = service.DeleteHardOneUser(context.TODO(), &pbusers.DeleteHardOneUserRequest{Id: res.Id})
}
