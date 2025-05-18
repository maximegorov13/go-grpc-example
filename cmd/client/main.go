package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

	ctx := context.Background()

	req := &example.CreatePostRequest{
		Title:    "example",
		AuthorId: "1",
		Content:  "Test",
	}

	newCtx := metadata.NewOutgoingContext(ctx, metadata.Pairs("client-header-key", "val"))

	var headers, trailers = metadata.MD{}, metadata.MD{}
	resp, err := client.CreatePost(newCtx, req,
		grpc.Header(&headers),
		grpc.Trailer(&trailers),
	)
	if err != nil {
		switch status.Code(err) {
		case codes.InvalidArgument:
			log.Println("некорректный запрос")
		default:
			log.Fatal(err)
		}

		if st, ok := status.FromError(err); ok {
			log.Println("code", st.Code(), "message", st.Message(), "details", st.Details())
		} else {
			log.Println("not grpc")
		}
	}

	log.Println("headers:", headers, "trailers:", trailers)

	log.Println(resp.GetPostId())

	bytes, err := protojson.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(bytes))
}
