package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"context"
	"os/signal"
	votingpb "github.com/guilherme-brandao/crypto-vote-platform/proto"
	"github.com/guilherme-brandao/crypto-vote-platform/database"
	"github.com/guilherme-brandao/crypto-vote-platform/services"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

var collection *mongo.Collection
var votingService votingpb.VotingServiceServer
var mongoCtx context.Context

func main() {

	database.Init()

	mongoCtx = database.GetContext()
	collection = database.GetCollection("ranking")
	votingService = services.New(collection, mongoCtx)

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Starting server on port :50051...")

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Unable to listen on port :50051: %v", err)
	}

	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)
	votingpb.RegisterVotingServiceServer(server, votingService)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	fmt.Println("Server succesfully started on port :50051")

	c := make(chan os.Signal)

	signal.Notify(c, os.Interrupt)

	<-c

	fmt.Println("\nStopping the server...")
	server.Stop()
	listener.Close()
	database.CloseDBConnection()
	fmt.Println("Done.")

}
