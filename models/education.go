package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Education struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      primitive.ObjectID `bson:"user_id"`
	SchoolName  string             `bson:"school_name"`
	Title       string             `bson:"title"`
	FromTime    time.Time          `bson:"from_time"`
	ToTime      time.Time          `bson:"to_time"`
	Description string             `bson:"description"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}
