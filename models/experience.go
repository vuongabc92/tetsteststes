package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Experience struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      primitive.ObjectID `bson:"user_id"`
	CompanyName string             `bson:"company_name"`
	JobTitle    string             `bson:"job_title"`
	TimeFrom    time.Time          `bson:"time_from"`
	TimeTo      time.Time          `bson:"time_to"`
	Description string             `bson:"description"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}
