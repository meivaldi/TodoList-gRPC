package main

import (
	"log"
	"net"

	"github.com/meivaldi/TodoList-gRPC/todolist/todolistpb"
	"google.golang.org/grpc"
)

type server struct {
	todolistpb.UnimplementedTodoListServiceServer
}

func main() {
	lis, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	todolistpb.RegisterTodoListServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
