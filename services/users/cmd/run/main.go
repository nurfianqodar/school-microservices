package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/nurfianqodar/school-microservices/services/users/db"
	pbusers "github.com/nurfianqodar/school-microservices/services/users/pb/users/v1"
	svc "github.com/nurfianqodar/school-microservices/services/users/services"
	"google.golang.org/grpc"
)

func main() {
	dsn, ok := os.LookupEnv("DSN")
	if !ok {
		log.Fatalln("DSN environment variable was not set")
	}

	// Create database connection
	ctx := context.Background()
	dbConn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalln(err)
	}
	defer dbConn.Close(ctx)

	// Create queries instance
	q := db.New(dbConn)

	// Create server
	server := grpc.NewServer()
	service := svc.New(q)
	pbusers.RegisterUserServiceServer(server, service)

	// Create listener and runserver
	host, ok := os.LookupEnv("HOST")
	if !ok {
		log.Fatal("HOST environment variable was not set")
	}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatal("PORT environment variable was not set")
	}
	addr := fmt.Sprintf("%s:%s", host, port)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("server listening on %s", "127.0.0.1:50051")
	if err = server.Serve(ln); err != nil {
		log.Fatalln(err)
	}
}
