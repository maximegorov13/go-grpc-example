package main

import (
	"context"
	"io"
	"log"
	"math/rand"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/maximegorov13/go-grpc-example/pkg/api/example"
)

func main() {
	server := grpc.NewServer()

	service := &ExampleService{
		storage: make(map[uint64]*Post, 1),
	}

	example.RegisterExampleServer(server, service)

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatal(err)
	}

	reflection.Register(server)

	log.Println("gRPC server listen on :8082")

	if err = server.Serve(lis); err != io.EOF {
		log.Fatal(err)
	}
}

type Post struct {
	ID       uint64
	Title    string
	Content  string
	AuthorID string
}

type ExampleService struct {
	example.UnimplementedExampleServer

	storage map[uint64]*Post
	mu      sync.RWMutex
}

func (s *ExampleService) CreatePost(ctx context.Context, req *example.CreatePostRequest) (*example.CreatePostResponse, error) {

	id := rand.Uint64()
	post := &Post{
		ID:       id,
		Title:    req.GetTitle(),
		Content:  req.GetContent(),
		AuthorID: req.GetAuthorId(),
	}
	s.mu.Lock()
	s.storage[id] = post
	s.mu.Unlock()

	return &example.CreatePostResponse{
		PostId: id,
	}, nil
}

func (s *ExampleService) ListPosts(ctx context.Context, req *example.ListPostsRequest) (*example.ListPostsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPosts not implemented")
}
