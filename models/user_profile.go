package models

import (
	"github.com/vuongabc92/octocv/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserProfile struct {
	Id          primitive.ObjectID `bson:"id"`
	UserId      primitive.ObjectID `bson:"user_id"`
	JobTitle    string             `bson:"job_title"`
	FirstName   string             `bson:"first_name"`
	LastName    string             `bson:"last_name"`
	FullName    string             `bson:"full_name"`
	AvatarImage string             `bson:"avatar_image"`
	CoverImage  string             `bson:"cover_image"`
	Birthday    time.Time          `bson:"birthday"`
	CountryID   config.ID          `bson:"country_id"`
	Address     string             `bson:"address"`
	PhoneNumber string             `bson:"phone_number"`
	Links       string             `bson:"links"`
	About       string             `bson:"about"`
	Hobbit      string             `bson:"hobbit"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}
