package main

import (
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/maximegorov13/go-grpc-example/pkg/api/example"
)

func main() {
	server := grpc.NewServer()

	service := ExampleService{}

	example.RegisterExampleServer(server, service)

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatal(err)
	}

	if err = server.Serve(lis); err != io.EOF {
		log.Fatal(err)
	}
}

type ExampleService struct {
	example.UnimplementedExampleServer
}
