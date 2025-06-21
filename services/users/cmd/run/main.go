package main

import (
	"log"
	"net"

	pbusers "github.com/nurfianqodar/school-microservices/services/users/pb/users/v1"
	svc "github.com/nurfianqodar/school-microservices/services/users/services"
	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer()
	service := svc.New()
	pbusers.RegisterUserServiceServer(server, service)

	ln, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("server listening on %s", "127.0.0.1:50051")
	if err = server.Serve(ln); err != nil {
		log.Fatalln(err)
	}
}
