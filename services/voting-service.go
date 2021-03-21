package services

import (
	"context"
	"fmt"

	"github.com/guilherme-brandao/crypto-vote-platform/models"
	votingpb "github.com/guilherme-brandao/crypto-vote-platform/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type votingServiceServer struct {
	collection *mongo.Collection
	mongoCtx   context.Context
}

func New(collection *mongo.Collection, mongoCtx context.Context) votingpb.VotingServiceServer {
	return &votingServiceServer{
		collection: collection,
		mongoCtx:   mongoCtx,
	}
}

func (s *votingServiceServer) CreateCrypto(ctx context.Context, req *votingpb.Crypto) (*votingpb.CreateCryptoRes, error) {

	if req.GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Verify the fields!"))
	}

	if req.GetUpvotes() > 0 || req.GetDownvotes() > 0 || req.GetUpvotes() > 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Crypto must be initialized with votes equal to zero!"))
	}

	cryptocurrency := req
	data := models.Cryptocurrency{
		Name:      cryptocurrency.GetName(),
		Upvotes:   cryptocurrency.GetUpvotes(),
		Downvotes: cryptocurrency.GetDownvotes(),
		Score:     cryptocurrency.GetScore(),
	}

	fmt.Println(data)

	result, err := s.collection.InsertOne(s.mongoCtx, data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)

	}
	oid := result.InsertedID.(primitive.ObjectID)
	cryptocurrency.Id = oid.Hex()
	return &votingpb.CreateCryptoRes{Crypto: cryptocurrency}, nil
}

func (s *votingServiceServer) GetCrypto(ctx context.Context, req *votingpb.GetCryptoReq) (*votingpb.GetCryptoRes, error) {

	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}

	result := s.collection.FindOne(ctx, bson.M{"_id": oid})
	data := models.Cryptocurrency{}
	if err := result.Decode(&data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find crypto with Object Id %s: %v", req.GetId(), err))
	}

	response := &votingpb.GetCryptoRes{
		Crypto: &votingpb.Crypto{
			Id:        oid.Hex(),
			Name:      data.Name,
			Upvotes:   data.Upvotes,
			Downvotes: data.Downvotes,
			Score:     data.Score,
		},
	}
	return response, nil
}

func (s *votingServiceServer) DeleteCrypto(ctx context.Context, req *votingpb.DeleteCryptoReq) (*votingpb.DeleteCryptoRes, error) {
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}

	_, err = s.collection.DeleteOne(ctx, bson.M{"_id": oid})

	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not delete crypto with id %s: %v", req.GetId(), err))
	}

	return &votingpb.DeleteCryptoRes{
		Success: true,
	}, nil
}

func (s *votingServiceServer) UpvoteCrypto(ctx context.Context, req *votingpb.UpvoteCryptoReq) (*votingpb.UpvoteCryptoRes, error) {

	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Could not convert the supplied crypto id to a MongoDB ObjectId: %v", err),
		)
	}

	filter := bson.M{"_id": oid}

	_, err = s.collection.UpdateOne(ctx, filter, bson.M{"$inc": bson.M{"upvotes": 1, "score": 1}}, options.Update().SetUpsert(true))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find crypto with id %s: %v", req.GetId(), err))
	}

	return &votingpb.UpvoteCryptoRes{
		Success: true,
	}, nil
}

func (s *votingServiceServer) DownvoteCrypto(ctx context.Context, req *votingpb.DownvoteCryptoReq) (*votingpb.DownvoteCryptoRes, error) {

	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Could not convert the supplied crypto id to a MongoDB ObjectId: %v", err),
		)
	}

	filter := bson.M{"_id": oid}

	_, err = s.collection.UpdateOne(ctx, filter, bson.M{"$inc": bson.M{"downvotes": 1, "score": -1}}, options.Update().SetUpsert(true))

	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find crypto with id %s: %v", req.GetId(), err))
	}

	return &votingpb.DownvoteCryptoRes{
		Success: true,
	}, nil
}

func (s *votingServiceServer) ListCryptos(req *votingpb.ListCryptosReq, stream votingpb.VotingService_ListCryptosServer) error {

	data := &models.Cryptocurrency{}

	cursor, err := s.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unknown internal error: %v", err))
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {

		err := cursor.Decode(data)

		if err != nil {
			return status.Errorf(codes.Unavailable, fmt.Sprintf("Could not decode data: %v", err))
		}

		stream.Send(&votingpb.ListCryptosRes{
			Crypto: &votingpb.Crypto{
				Id:        data.ID.Hex(),
				Name:      data.Name,
				Upvotes:   data.Upvotes,
				Downvotes: data.Downvotes,
				Score:     data.Score,
			},
		})
	}

	if err := cursor.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unkown cursor error: %v", err))
	}
	return nil
}

