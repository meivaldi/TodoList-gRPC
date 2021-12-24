package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/meivaldi/TodoList-gRPC/todolist/todolistpb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type server struct {
	todolistpb.UnimplementedTodoListServiceServer
}

type todoListItem struct {
	Id          int    `bson: "id,omitempty"`
	Title       string `bson: "title"`
	Description string `bson: "desc"`
	Thumbnail   string `bson: "thumbnail"`
	Priority    int    `bson: "priority"`
	Date        string `bson: "date"`
}

var collection *mongo.Collection

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("TodoList Service Started...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fmt.Println("Connecting MongoDB")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	collection = client.Database("backend").Collection("todolist")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	todolistpb.RegisterTodoListServiceServer(s, &server{})

	go func() {
		fmt.Println("Starting server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v\n", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch
	fmt.Println("Stopping server...")
	s.Stop()
	fmt.Println("Stopping listener...")
	lis.Close()
	fmt.Println("Closing MongoDB Connection")
	client.Disconnect(context.TODO())
	fmt.Println("Finish...")
}
