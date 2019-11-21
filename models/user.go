package models

import (
	"github.com/vuongabc92/octocv/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Email       string             `bson:"email"`
	Password    string             `bson:"password"`
	Status      config.UserStatus  `bson:"status"`
	VerifyToken string             `bson:"verify_token"`
	VerifiedAt  time.Time          `bson:"verified_at"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}
