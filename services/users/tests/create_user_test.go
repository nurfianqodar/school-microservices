package user_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	pbusers "github.com/nurfianqodar/school-microservices/services/users/pb/users/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func createService() pbusers.UserServiceClient {
	host, ok := os.LookupEnv("HOST")
	if !ok {
		log.Fatal("HOST not found")
	}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatal("PORT not found")
	}
	addr := fmt.Sprintf("%s:%s", host, port)

	client, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("unable to connect to client: %s", err.Error())
	}
	service := pbusers.NewUserServiceClient(client)

	return service
}

func TestCreateUser(t *testing.T) {
	service := createService()

	t.Run("Should success create user", func(t *testing.T) {
		ctx := context.TODO()
		req := &pbusers.CreateOneUserRequest{
			Email:    "dummy@email.com",
			Password: "secretpassword",
			Role:     pbusers.UserRole_Student,
		}

		res, err := service.CreateOneUser(ctx, req)
		if err != nil {
			t.Fail()
		}

		// Delete if not error
		if err == nil && res != nil {
			t.Log(res)
			_, _ = service.DeleteHardOneUser(ctx, &pbusers.DeleteHardOneUserRequest{Id: res.Id})
		}
	})

	t.Run("Should validation error", func(t *testing.T) {
		ctx := context.TODO()
		req := &pbusers.CreateOneUserRequest{
			Email:    "dummyemail.com",
			Password: "sword",
			Role:     pbusers.UserRole_Student,
		}

		_, err := service.CreateOneUser(ctx, req)
		if err == nil {
			t.Fail()
		}

		st := status.Convert(err)
		t.Log(st.Code().String())
		if len(st.Details()) == 0 {
			t.Fail()
		}
		if st.Code() != codes.InvalidArgument {
			t.Fail()
		}
	})

	t.Run("Should conflict error", func(t *testing.T) {
		ctx := context.TODO()
		req := &pbusers.CreateOneUserRequest{
			Email:    "dummy@email.com",
			Password: "secretpassword",
			Role:     pbusers.UserRole_Student,
		}

		res, err := service.CreateOneUser(ctx, req)
		if err != nil {
			t.Fail()
		}
		defer func() {
			if err == nil && res != nil {
				t.Log(res)
				_, _ = service.DeleteHardOneUser(ctx, &pbusers.DeleteHardOneUserRequest{Id: res.Id})
			}
		}()

		_, err2 := service.CreateOneUser(ctx, req)
		if err2 == nil {
			t.Fail()
		}

		st := status.Convert(err2)
		t.Log(st.Code().String())
		if st.Code() != codes.AlreadyExists {
			t.Fail()
		}
	})
}

func BenchmarkCreateUser(b *testing.B) {
	service := createService()

	for b.Loop() {

		ctx := context.TODO()

		req := &pbusers.CreateOneUserRequest{
			Email:    fmt.Sprintf("user%s@email.com", uuid.NewString()),
			Password: "securepassword",
			Role:     pbusers.UserRole_Student,
		}

		res, err := service.CreateOneUser(ctx, req)
		if err != nil {
			b.Errorf("unexpected error: %v", err)
		}

		// Cleanup
		if res != nil {
			_, _ = service.DeleteHardOneUser(ctx, &pbusers.DeleteHardOneUserRequest{Id: res.Id})
		}

	}
}
