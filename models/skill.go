package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Skill struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Name      string             `bson:"name"`
	Rate      uint8              `bson:"rate"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}
