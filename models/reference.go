package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Reference struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      primitive.ObjectID `bson:"user_id"`
	FullName    string             `bson:"full_name"`
	Company     string             `bson:"company"`
	Position    string             `bson:"position"`
	Email       string             `bson:"email"`
	PhoneNumber string             `bson:"phone_number"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}
