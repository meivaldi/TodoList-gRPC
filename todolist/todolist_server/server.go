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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (*server) CreateTodoList(ctx context.Context, req *todolistpb.CreateTodoListRequest) (*todolistpb.CreateTodoListResponse, error) {
	todoList := req.GetTodoList()

	data := todoListItem{
		Title:       todoList.GetTitle(),
		Description: todoList.GetDescription(),
		Thumbnail:   todoList.GetThumbnail(),
		Priority:    int(todoList.GetPriority()),
		Date:        todoList.GetDate(),
	}

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal Server Error: %v\n", err),
		)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintln("Couldn't convert to OID"),
		)
	}

	return &todolistpb.CreateTodoListResponse{
		TodoList: &todolistpb.TodoList{
			Id:          oid.Hex(),
			Title:       todoList.GetTitle(),
			Description: todoList.GetDescription(),
			Thumbnail:   todoList.GetThumbnail(),
			Priority:    todoList.GetPriority(),
			Date:        todoList.GetDate(),
		},
	}, nil
}

func main() {
	//set flags so if we get an error, we can see which line cause an error
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("TodoList Service Started...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fmt.Println("Connecting MongoDB")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Couldn't connect to mongodb: %v\n", err)
	}

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

	//create custom interrupt, so code below get executed after finished/stopped
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
