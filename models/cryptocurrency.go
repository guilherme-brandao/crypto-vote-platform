package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Cryptocurrency struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name 	   string             `bson:"name" binding:"required"`
	Upvotes    int64              `bson:"upvotes"`
	Downvotes  int64              `bson:"downvotes"`
	Score      int64			  `bson:"score"`
}
