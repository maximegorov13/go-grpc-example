package main

import (
	"context"
	"io"
	"log"
	"math/rand"
	"net"
	"sync"

	"buf.build/go/protovalidate"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/maximegorov13/go-grpc-example/pkg/api/example"
)

func main() {
	validator, err := protovalidate.New()
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()

	service := &ExampleService{
		storage:   make(map[uint64]*Post, 1),
		validator: validator,
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

	validator protovalidate.Validator
	storage   map[uint64]*Post
	mu        sync.RWMutex
}

func (s *ExampleService) CreatePost(ctx context.Context, req *example.CreatePostRequest) (*example.CreatePostResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		log.Println(md)
	}

	if err := s.validator.Validate(req); err != nil {
		st := status.New(codes.InvalidArgument, codes.InvalidArgument.String())
		st, _ = st.WithDetails(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "request",
					Description: err.Error(),
				},
			},
		})

		return nil, st.Err()
	}

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

	header := metadata.Pairs("header-key", "val")
	grpc.SetHeader(ctx, header)
	grpc.SetTrailer(ctx, header)

	return &example.CreatePostResponse{
		PostId: id,
	}, nil
}

func (s *ExampleService) ListPosts(ctx context.Context, req *example.ListPostsRequest) (*example.ListPostsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPosts not implemented")
}
