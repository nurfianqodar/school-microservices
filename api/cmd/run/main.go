package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nurfianqodar/school-microservices/api/handlers"
	pbusers "github.com/nurfianqodar/school-microservices/services/users/pb/users/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Create router
	r := http.NewServeMux()

	// Create user service and handler
	userSvcHost, ok := os.LookupEnv("USER_SERVICE_HOST")
	if !ok {
		log.Fatal("USER_SERVICE_HOST variable was not set")
	}
	userSvcPort, ok := os.LookupEnv("USER_SERVICE_PORT")
	if !ok {
		log.Fatal("USER_SERVICE_PORT variable was not set")
	}
	userSvcAddr := fmt.Sprintf("%s:%s", userSvcHost, userSvcPort)
	userServiceClient, err := grpc.NewClient(userSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error: unable connecting to user service client. %e\n", err)
	}
	defer func() {
		if userServiceClient != nil {
			if err := userServiceClient.Close(); err != nil {
				log.Printf("error: failed to close user service client. %e", err)
			}
		}
	}()
	userSvc := pbusers.NewUserServiceClient(userServiceClient)
	userHandler := handlers.NewUserHandler(userSvc)
	userHandler.RegisterRouter(r)

	// Run http server
	host, ok := os.LookupEnv("HOST")
	if !ok {
		log.Fatal("HOST variable was not set")
	}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatal("PORT variable was not set")
	}
	addr := fmt.Sprintf("%s:%s", host, port)

	log.Printf("server listening on %s\n", addr)
	if err = http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
