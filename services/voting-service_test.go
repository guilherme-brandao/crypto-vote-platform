package services

import (
	"context"
	"testing"
	"log"
	"net"
	"github.com/guilherme-brandao/crypto-vote-platform/database"
	votingpb "github.com/guilherme-brandao/crypto-vote-platform/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type CryptocurrencyReturn struct {
	id         string
	name 	   string             
	upvotes    int64              
	downvotes  int64              
	score      int64			  
}

const bufSize = 1024 * 1024

var lis *bufconn.Listener

var votingService votingpb.VotingServiceServer
var mongoCtx context.Context
var collection *mongo.Collection

func init() {
	database.Init()

	mongoCtx = database.GetContext()
	collection = database.GetCollection("ranking-test-db")
	
	votingService = New(collection, mongoCtx)

    lis = bufconn.Listen(bufSize)
    s := grpc.NewServer()
    votingpb.RegisterVotingServiceServer(s, votingService)
    go func() {
        if err := s.Serve(lis); err != nil {
            log.Fatalf("Server exited with error: %v", err)
        }
    }()
}

func bufDialer(context.Context, string) (net.Conn, error) {
    return lis.Dial()
}

func TestInputCrypto(t *testing.T) {
    ctx := context.Background()
    conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
    if err != nil {
        t.Fatalf("Failed to dial bufnet: %v", err)
    }
    defer conn.Close()
	client := votingpb.NewVotingServiceClient(conn)
	
	req := votingpb.Crypto{
		Id: "1234567890",
		Name: "testCoinInput",
		Upvotes: 0,
		Downvotes: 0,
		Score: 0,
	}

    resp, err := client.CreateCrypto(ctx, &req)
    if err != nil {
        t.Fatalf("CreateCrypto failed: %v", err)
    }
	log.Printf("Response: %+v", resp.Crypto)

	assert.Equal(t, resp.Crypto.Name, req.Name)

}

func TestGetCrypto(t *testing.T) {
    ctx := context.Background()
    conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
    if err != nil {
        t.Fatalf("Failed to dial bufnet: %v", err)
    }
    defer conn.Close()
	client := votingpb.NewVotingServiceClient(conn)

	cryptoMock :=  CryptocurrencyReturn{
		id: "6056389da6cf0b94d7e8fb85",
		name: "testCoin",
		upvotes: 2,
		downvotes: 1,
		score: 1,
	}

	req := votingpb.GetCryptoReq{Id: "6056389da6cf0b94d7e8fb85"}
    resp, err := client.GetCrypto(ctx, &req)
    if err != nil {
        t.Fatalf("GetCrypto failed: %v", err)
    }
	log.Printf("Response: %+v", resp.Crypto)

	assert.Equal(t, resp.Crypto.Id, cryptoMock.id)
	assert.Equal(t, resp.Crypto.Name, cryptoMock.name)
	assert.Equal(t, resp.Crypto.Upvotes, cryptoMock.upvotes)
	assert.Equal(t, resp.Crypto.Downvotes, cryptoMock.downvotes)
	assert.Equal(t, resp.Crypto.Score, cryptoMock.score)

}
