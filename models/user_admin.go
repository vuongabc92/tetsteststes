package models

import (
	"github.com/vuongabc92/octocv/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserAdmin struct {
	ID        primitive.ObjectID     `bson:"id"`
	Email     string                 `bson:"email"`
	Password  string                 `bson:"password"`
	Status    config.UserAdminStatus `bson:"status"`
	RoleID    primitive.ObjectID     `bson:"role_id"`
	CreatedAt time.Time              `bson:"created_at"`
	UpdatedAt time.Time              `bson:"updated_at"`
}
