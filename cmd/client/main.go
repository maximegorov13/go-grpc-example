package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/maximegorov13/go-grpc-example/pkg/api/example"
)

func main() {
	conn, err := grpc.NewClient(":8082",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := example.NewExampleClient(conn)

	resp, err := client.CreatePost(context.Background(), &example.CreatePostRequest{
		Title:    "example",
		AuthorId: "",
		Content:  "Test",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp.GetPostId())

	bytes, err := protojson.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(bytes))
}
