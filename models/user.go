package models

import (
	"github.com/vuongabc92/octocv/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	Email       string             `bson:"email"`
	Password    string             `bson:"password"`
	Status      config.UserStatus  `bson:"status"`
	VerifyToken string             `bson:"verify_token"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}
