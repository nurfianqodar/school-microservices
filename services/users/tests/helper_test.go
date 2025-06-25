package user_test

import (
	"fmt"
	"log"
	"os"

	pbusers "github.com/nurfianqodar/school-microservices/services/users/pb/users/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
