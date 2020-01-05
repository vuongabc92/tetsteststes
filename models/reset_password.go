package models

import (
	"github.com/vuongabc92/octocv/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ResetPassword struct {
	Id        primitive.ObjectID         `bson:"_id,omitempty"`
	Email     string                     `bson:"email"`
	Token     string                     `bson:"token"`
	Status    config.ResetPasswordStatus `bson:"status"`
	CreatedAt time.Time                  `bson:"created_at"`
	UpdatedAt time.Time                  `bson:"updated_at"`
}
