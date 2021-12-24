package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/meivaldi/TodoList-gRPC/todolist/todolistpb"
	"google.golang.org/grpc"
)

type server struct {
	todolistpb.UnimplementedTodoListServiceServer
}

func (s *server) TodoList(ctx context.Context, req *todolistpb.TodoListRequest) (*todolistpb.TodoListResponse, error) {
	fmt.Printf("TodoList function was invoked with %v\n", req)
	title := req.GetTodoList().GetTitle()
	description := req.GetTodoList().GetDescription()
	thumbnail := req.GetTodoList().GetThumbnail()
	priority := req.GetTodoList().GetPriority()

	result := title + ", " + description + ", " + thumbnail + ", " + strconv.Itoa(int(priority))
	res := &todolistpb.TodoListResponse{
		Result: result,
	}

	return res, nil
}

func (s *server) TodoListManyTimes(req *todolistpb.TodoListManyTimesRequests, stream todolistpb.TodoListService_TodoListManyTimesServer) error {
	fmt.Printf("TodoList stream server function was invoked with %v\n", req)

	title := req.GetTodolist().GetTitle()
	description := req.GetTodolist().GetDescription()
	thumbnail := req.GetTodolist().GetThumbnail()
	priority := req.GetTodolist().GetPriority()

	for i := 1; i <= 10; i++ {
		result := strconv.Itoa(i) + ". " + title + ", " + description + ", " + thumbnail + ", " + strconv.Itoa(int(priority))
		res := &todolistpb.TodoListManyTimesResponses{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}

	return nil
}

func (s *server) LongTodoList(stream todolistpb.TodoListService_LongTodoListServer) error {
	fmt.Printf("TodoList stream client function was invoked with a streaming request\n")
	result := ""

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&todolistpb.LongTodoListResponse{
				Result: result,
			})
		}

		if err != nil {
			log.Fatalf("Failed to send streaming request: %v", err)
		}

		title := req.GetTodoList().GetTitle()
		result += title + "; "
	}
}

func (s *server) TodoListEveryone(stream todolistpb.TodoListService_TodoListEveryoneServer) error {
	fmt.Printf("TodoList stream server function was invoked with a streaming request\n")

	return nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	todolistpb.RegisterTodoListServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
