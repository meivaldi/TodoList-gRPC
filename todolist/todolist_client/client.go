package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/meivaldi/TodoList-gRPC/todolist/todolistpb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	c := todolistpb.NewTodoListServiceClient(conn)

	//doUnary(c)
	//doServerStreaming(c)
	doClientStreaming(c)
}

func doUnary(c todolistpb.TodoListServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &todolistpb.TodoListRequest{
		TodoList: &todolistpb.TodoList{
			Title:       "Golang gRPC",
			Description: "Belajar golang gRPC",
			Thumbnail:   "https://meivaldi.com/grpc/1.png",
			Priority:    1,
		},
	}

	res, err := c.TodoList(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling TodoList RPC: %v", err)
	}

	log.Printf("Response from TodoList RPC: %v", res)
}

func doServerStreaming(c todolistpb.TodoListServiceClient) {
	fmt.Println("Starting server streaming RPC...")

	req := &todolistpb.TodoListManyTimesRequests{
		Todolist: &todolistpb.TodoList{
			Title:       "Golang gRPC Server Stream",
			Description: "Belajar golang gRPC server stream",
			Thumbnail:   "https://meivaldi.com/grpc/2.png",
			Priority:    1,
		},
	}

	resStream, err := c.TodoListManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to get streaming result: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}

		log.Printf("Response from TodoListManyTimes: %v", msg.GetResult())
	}
}

func doClientStreaming(c todolistpb.TodoListServiceClient) {
	fmt.Println("Starting client streaming RPC...")

	requests := []*todolistpb.LongTodoListRequest{
		&todolistpb.LongTodoListRequest{
			TodoList: &todolistpb.TodoList{
				Title:       "GO",
				Description: "Golang",
				Thumbnail:   "go.png",
				Priority:    1,
			},
		},
		&todolistpb.LongTodoListRequest{
			TodoList: &todolistpb.TodoList{
				Title:       "gRPC",
				Description: "gRPC for Golang",
				Thumbnail:   "grpc.png",
				Priority:    1,
			},
		},
		&todolistpb.LongTodoListRequest{
			TodoList: &todolistpb.TodoList{
				Title:       "MongoDB",
				Description: "MongoDB document based database",
				Thumbnail:   "mongodb.png",
				Priority:    2,
			},
		},
		&todolistpb.LongTodoListRequest{
			TodoList: &todolistpb.TodoList{
				Title:       "Redis",
				Description: "Redis key-value based database",
				Thumbnail:   "redis.png",
				Priority:    2,
			},
		},
		&todolistpb.LongTodoListRequest{
			TodoList: &todolistpb.TodoList{
				Title:       "NSQ",
				Description: "Message Broker",
				Thumbnail:   "nsq.png",
				Priority:    1,
			},
		},
	}

	stream, err := c.LongTodoList(context.Background())
	if err != nil {
		log.Fatalf("Failed to stream: %v", err)
	}

	for _, req := range requests {
		fmt.Println("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v", err)
	}

	fmt.Printf("LongTodoList Response: %v\n", res)
}
